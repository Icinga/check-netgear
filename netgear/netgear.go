package netgear

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const timeout = 10 * time.Second

type Netgear struct {
	sessionToken string
	client       *http.Client

	hostName string
	username string
	password string
}

func NewNetgear(hostName, username, password string) *Netgear {
	return &Netgear{
		client:   &http.Client{Timeout: timeout},
		hostName: strings.TrimSuffix(hostName, "/") + "/api/v1",
		username: username,
		password: password,
	}
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

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/login", n.hostName),
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse login response JSON: %w", err)
	}

	userRaw, ok := result["user"]
	if !ok {
		return fmt.Errorf("login response: missing 'user' field; body: %s", string(body))
	}

	userSection, ok := userRaw.(map[string]interface{})
	if !ok {
		return fmt.Errorf("login response: 'user' is not a JSON object; body: %s", string(body))
	}

	session, ok := userSection["session"]
	if !ok {
		return fmt.Errorf("login response: missing 'session' token; body: %s", string(body))
	}

	n.sessionToken, ok = session.(string)
	if !ok {
		return fmt.Errorf("login response: 'session' token is not a string; body: %v", session)
	}

	return nil
}

func (n *Netgear) Logout() error {
	resp, err := n.client.Get(fmt.Sprintf("%s/logout", n.hostName))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.ReadAll(resp.Body) // body ignored intentionally
	return nil
}

func (n *Netgear) DeviceInfo() ([]byte, error) {
	return n.doRequest("GET", "device_info", nil)
}

func (n *Netgear) PortStatistics(statType string) ([]byte, error) {
	path := fmt.Sprintf("port_statistics?type=%s&indexPage=1&pageSize=25", statType)
	return n.doRequest("GET", path, nil)
}

func (n *Netgear) PoeStatus() ([]byte, error) {
	return n.doRequest("GET", "swcfg_poe", nil)
}

func (n *Netgear) doRequest(method, path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", n.hostName, path), body)
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

func StringPercentToFloat(percents string) float64 {
	toReturn, err := strconv.ParseFloat(strings.TrimSuffix(percents, "%"), 64)
	if err != nil {
		return 0
	}
	return toReturn
}
