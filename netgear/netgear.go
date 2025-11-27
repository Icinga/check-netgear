package netgear

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const timeout = 10 * time.Second

type Netgear struct {
	sessionToken string
	client       *http.Client

	baseUrl  *url.URL
	username string
	password string

	requestHandler func(string, *url.URL) (io.Reader, error)
}

func NewNetgear(baseUrl, username, password string) (*Netgear, error) {
	u, err := url.Parse(strings.TrimSpace(baseUrl))
	if err != nil {
		return nil, err
	}
	// be sure that the base path is /api/v1 using JoinPath
	u = u.JoinPath("api", "v1")

	n := &Netgear{
		client:   &http.Client{Timeout: timeout},
		baseUrl:  u,
		username: username,
		password: password,
	}
	n.requestHandler = n.makeRequest

	return n, nil
}

func (n *Netgear) Login() error {
	loginPayload := map[string]interface{}{
		"user": map[string]string{
			"name":     n.username,
			"password": n.password,
		},
	}

	payloadBytes, err := json.Marshal(loginPayload)
	if err != nil {
		return err
	}

	loginURL := n.baseUrl.JoinPath("login")
	req, err := http.NewRequest(http.MethodPost, loginURL.String(), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse login response JSON: %w", err)
	}

	userRaw, ok := result["user"]
	if !ok {
		return fmt.Errorf("login response: missing 'user' field")
	}

	userSection, ok := userRaw.(map[string]interface{})
	if !ok {
		return fmt.Errorf("login response: 'user' is not a JSON object")
	}

	session, ok := userSection["session"]
	if !ok {
		return fmt.Errorf("login response: missing 'session' token")
	}

	n.sessionToken, ok = session.(string)
	if !ok {
		return fmt.Errorf("login response: 'session' token is not a string (got %T)", session)
	}

	return nil
}

func (n *Netgear) Logout() error {
	logoutURL := n.baseUrl.JoinPath("logout")
	resp, err := n.client.Get(logoutURL.String())
	if err != nil {
		return err
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	return nil
}

func (n *Netgear) DeviceInfo() (*DeviceInfo, error) {
	var di DeviceInfo
	if err := n.doRequest(http.MethodGet, "device_info", &di); err != nil {
		return nil, err
	}
	return &di, nil
}

func (n *Netgear) PortStatistics(statType string) (*PortStatistics, error) {
	const defaultPageIndex = 1
	const defaultPageSize = 25

	// NOTE: There should be a paginated request here!
	// However, as of 2025-11-25, the API does not respect the parameters "indexPage" and "pageSize" and always
	// returns the first page. Therefore, this request only works for switches with less than 25 ports.
	u := n.baseUrl.JoinPath("port_statistics")
	q := u.Query()
	q.Set("type", statType)
	q.Set("indexPage", strconv.Itoa(defaultPageIndex))
	q.Set("pageSize", strconv.Itoa(defaultPageSize))
	u.RawQuery = q.Encode()

	var stats PortStatistics
	if err := n.doRequestURL(http.MethodGet, u, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func (n *Netgear) PoeStatus() (*PoeStatus, error) {
	poeStatus := new(PoeStatus)
	if err := n.doRequest(http.MethodGet, "swcfg_poe", poeStatus); err != nil {
		return nil, err
	}
	return poeStatus, nil
}

// ModeBasic contains all the basic hardware information of the switch, including CPU and RAM usage, temperature and fan
// speed
func (n *Netgear) ModeBasic(flags *Flags) (*result.PartialResult, error) {
	deviceInfo, err := n.DeviceInfo()
	if err != nil {
		return nil, fmt.Errorf("error retrieving device info: %w", err)
	}

	if len(deviceInfo.DeviceInfo.Details) == 0 {
		return nil, fmt.Errorf("error retrieving device info")
	}
	upTime := deviceInfo.DeviceInfo.Details[0].Uptime

	o := result.PartialResult{
		Output: fmt.Sprintf("Device Info: Uptime - %v", upTime),
	}

	if !flags.HideCpu {
		if len(deviceInfo.DeviceInfo.Cpu) == 0 {
			return nil, fmt.Errorf("no CPU info for this device")
		}
		cpuUsage, err := StringPercentToFloat(deviceInfo.DeviceInfo.Cpu[0].Usage)
		if err != nil {
			return nil, fmt.Errorf("error parsing CPU usage: %w", err)
		}

		cpuPartial, err := CheckCPU(cpuUsage, flags.NoPerfdata, flags.CpuWarn, flags.CpuCrit)
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("CPU check error: %v", err)
			err := errRes.SetState(check.Unknown)
			if err != nil {
				return nil, err
			}
			o.AddSubcheck(errRes)
		} else {
			o.AddSubcheck(*cpuPartial)
		}
	}

	if !flags.HideMem {
		if len(deviceInfo.DeviceInfo.Memory) == 0 {
			return nil, fmt.Errorf("no Memory info for this device")
		}
		memUsage, err := StringPercentToFloat(deviceInfo.DeviceInfo.Memory[0].Usage)
		if err != nil {
			return nil, fmt.Errorf("error parsing Memory usage: %w", err)
		}

		memPartial, err := CheckMemory(memUsage, flags.NoPerfdata, flags.MemWarn, flags.MemCrit)
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("Memory check error: %v", err)
			err := errRes.SetState(check.Unknown)
			if err != nil {
				return nil, err
			}
			o.AddSubcheck(errRes)
		} else {
			o.AddSubcheck(*memPartial)
		}
	}

	if !flags.HideTemp {
		if len(deviceInfo.DeviceInfo.Sensor) == 0 {
			return nil, fmt.Errorf("no Temperature info for this device")
		}
		sensorDetails := deviceInfo.DeviceInfo.Sensor[0].Details
		tempPartial, err := CheckTemperature(sensorDetails, flags.NoPerfdata, flags.TempWarn, flags.TempCrit)
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("Temperature check error: %v", err)
			err := errRes.SetState(check.Unknown)
			if err != nil {
				return nil, err
			}
			o.AddSubcheck(errRes)
		} else {
			o.AddSubcheck(*tempPartial)
		}
	}

	if !flags.HideFans {
		if len(deviceInfo.DeviceInfo.Fan) == 0 {
			return nil, fmt.Errorf("no Fan info for this device")
		}
		if len(deviceInfo.DeviceInfo.Fan[0].Details) == 0 {
			return nil, fmt.Errorf("no Fan details for this device")
		}
		fan := deviceInfo.DeviceInfo.Fan[0].Details[0]
		fanPartial, err := CheckFans(flags.NoPerfdata, fan.Description, fan.Speed, flags.FanWarn, flags.FanCrit)
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("Fans check error: %v", err)
			err := errRes.SetState(check.Unknown)
			if err != nil {
				return nil, err
			}
			o.AddSubcheck(errRes)
		} else {
			o.AddSubcheck(*fanPartial)
		}
	}

	return &o, nil
}

// doRequestURL Performs an HTTP-Request to a given path on the previously defined host and stores the resulting json
// response in the object provided by the result parameter.
//
// Note: setting the result parameter to nil causes the parsing of the response to be skipped, the request is still
// performed.
func (n *Netgear) doRequest(method, path string, result any) error {
	return n.doRequestURL(method, n.baseUrl.JoinPath(path), result)
}

// doRequestURL Performs an HTTP-Request to a given URL and stores the resulting json response in the
// object provided by the result parameter.
//
// Note: setting the result parameter to nil causes the parsing of the response to be skipped, the request is still
// performed.
func (n *Netgear) doRequestURL(method string, u *url.URL, result any) error {
	if result == nil {
		return nil
	}

	body, err := n.requestHandler(method, u)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, body)
		if closer, ok := body.(io.Closer); ok {
			_ = closer.Close()
		}
	}()

	if err := json.NewDecoder(body).Decode(result); err != nil {
		return fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return nil
}

func (n *Netgear) makeRequest(method string, u *url.URL) (io.Reader, error) {
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("session", n.sessionToken)

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func StringPercentToFloat(percents string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSuffix(percents, "%"), 64)
}
