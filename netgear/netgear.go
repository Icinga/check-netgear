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

func NewNetgear(baseUrl, username, password string) (*Netgear, error) {
	u, err := url.Parse(strings.TrimSpace(baseUrl))
	if err != nil {
		return nil, err
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
	var poeStatus PoeStatus
	if err := n.doRequest(http.MethodGet, "swcfg_poe", &poeStatus); err != nil {
		return nil, err
	}
	return &poeStatus, nil
}

func (n *Netgear) doRequest(method, path string, result any) error {
	return n.doRequestURL(method, n.baseUrl.JoinPath(path), result)
}

func (n *Netgear) doRequestURL(method string, u *url.URL, result any) error {
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("session", n.sessionToken)

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if result == nil {
		// ignore response body but still read it so connection can be reused
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return nil
}

func StringPercentToFloat(percents string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSuffix(percents, "%"), 64)
}
