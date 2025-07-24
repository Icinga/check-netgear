package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NETWAYS/go-check"
)

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

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", *hostName), bytes.NewBuffer(payloadBytes))
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

	// fmt.Println("Login response:", string(body))

	var result map[string]interface{}
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

	session, ok := userSection["session"].(string)
	if !ok {
		return fmt.Errorf("login response: missing 'session' token; body: %s", string(body))
	}

	sessionToken = session
	// fmt.Println("Session token:", sessionToken)
	return nil
}

// LOGOUT
func logout() {
	resp, err := http.Get(fmt.Sprintf("%s/logout", *hostName))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

// DEVICE_INFO
func device_info() []byte {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/device_info", *hostName), nil)
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

	return body
}

// DEVICE_FAN
func device_fan() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/device_fan", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sntp_server_cfg", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/imageInfo", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/vlan_interface_route", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/special_interface_route", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ip_http", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/profile/list", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/profile", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sw_auto_lag_cfg", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/neighbor", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/swcfg_poe", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/port_statistics", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/swcfg_ports_status", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tech_support", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/device_config", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ping_test", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/trace_test", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dns_lookup", *hostName), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sw_ptp_cfg", *hostName), nil)
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

func string_percent_to_float(percents string) float64 {
	to_return, err := strconv.ParseFloat(strings.TrimSuffix(percents, "%"), 64)
	if err != nil {
		return 0
	}
	return to_return
}

func get_cpu_usage_check_level(cpuUsage float64) int {
	level := check.OK
	if cpuUsage >= 70 {
		level = check.Warning
	}
	if cpuUsage >= 90 {
		level = check.Critical
	}
	return level
}

func get_memory_usage_check_level(memoryUsage float64) int {
	level := check.OK
	if memoryUsage >= 70 {
		level = check.Warning
	}
	if memoryUsage >= 90 {
		level = check.Critical
	}
	return level
}
