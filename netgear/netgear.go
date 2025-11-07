package netgear

import (
	"bytes"
	"encoding/json"
	"fmt"
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
}

func NewNetgear(hostName, username, password string) (*Netgear, error) {
	u, err := url.Parse(strings.TrimSpace(hostName))
	if err != nil {
		return nil, err
	}
	if u == nil {
		u = &url.URL{}
	}
	// be sure that the base path is /api/v1 using JoinPath
	u = u.JoinPath("api", "v1")

	return &Netgear{
		client:   &http.Client{Timeout: timeout},
		baseUrl:  u,
		username: username,
		password: password,
	}, nil
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
	defer func() { _ = resp.Body.Close() }()

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
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)
	return nil
}

func (n *Netgear) DeviceInfo() (DeviceInfo, error) {
	data, err := n.doRequest(http.MethodGet, "device_info", nil)
	if err != nil {
		return DeviceInfo{}, err
	}
	var di DeviceInfo
	if err := json.Unmarshal(data, &di); err != nil {
		return DeviceInfo{}, err
	}
	return di, nil
}

func (n *Netgear) PortStatistics(statType string) ([]byte, error) {
	const defaultPageIndex = 1
	const defaultPageSize = 25

	u := n.baseUrl.JoinPath("port_statistics")
	q := u.Query()
	q.Set("type", statType)
	q.Set("indexPage", strconv.Itoa(defaultPageIndex))
	q.Set("pageSize", strconv.Itoa(defaultPageSize))
	u.RawQuery = q.Encode()

	return n.doRequestURL(http.MethodGet, u, nil)
}

func (n *Netgear) PoeStatus() ([]byte, error) {
	return n.doRequest(http.MethodGet, "swcfg_poe", nil)
}

func (n *Netgear) doRequest(method, path string, body io.Reader) ([]byte, error) {
	u := n.baseUrl.JoinPath(path)
	return n.doRequestURL(method, u, body)
}

func (n *Netgear) doRequestURL(method string, u *url.URL, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("session", n.sessionToken)

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func StringPercentToFloat(percents string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSuffix(percents, "%"), 64)
}
