package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/NETWAYS/go-check"
)

var baseURL = "http://192.168.0.239/api/v1"
var sessionToken string

// LOGIN
func login(name, password string) error {
	loginPayload := map[string]interface{}{
		"user": map[string]string{
			"name":     name,
			"password": password,
		},
	}

	payloadBytes, err := json.Marshal(loginPayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", baseURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Login response:", string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	userSection, ok := result["user"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("login response: missing 'user' field or wrong type")
	}

	session, ok := userSection["session"].(string)
	if !ok {
		return fmt.Errorf("login response: missing 'session' token")
	}

	sessionToken = session
	fmt.Println("Session token:", sessionToken)
	return nil
}

// LOGOUT
func logout() {
	resp, err := http.Get(fmt.Sprintf("%s/logout", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

// DEVICE_INFO
func device_info() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/device_info", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

// DEVICE_FAN
func device_fan() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/device_fan", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

// SNTP_SERVER_CFG
func sntp_server_cfg() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sntp_server_cfg", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func imageInfo() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/imageInfo", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func vlan_interface_route() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/vlan_interface_route", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func special_interface_route() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/special_interface_route", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func ip_http() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ip_http", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func profile_list() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/profile/list", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func profile() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/profile", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func sw_auto_lag_cfg() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sw_auto_lag_cfg", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func neighbor() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/neighbor", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func swcfg_poe() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/swcfg_poe", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func port_statistics() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/port_statistics", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func swcfg_ports_status() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/swcfg_ports_status", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func tech_support() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tech_support", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func device_config() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/device_config", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func ping_test() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ping_test", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func trace_test() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/trace_test", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func dns_lookup() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dns_lookup", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

func sw_ptp_cfg() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sw_ptp_cfg", baseURL), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	fmt.Printf("%v", string(body))
}

type Function struct {
	label    string
	function func()
}
