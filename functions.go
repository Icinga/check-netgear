package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NETWAYS/go-check"
)

var timeout = 10 * time.Second

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

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing body: %s\n", err)
		}
	}(resp.Body)

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
	if false { //Has to be used, but nobody wants it
		fmt.Printf("%v", string(body))
	}
}

// DEVICE_INFO
func device_info() []byte {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/device_info", *hostName), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing body: %s\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}

	return body
}

// statType - inbound / outbound
func port_statistics(statType string) []byte {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/port_statistics?type=%s&indexPage=1&pageSize=25", *hostName, statType), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing body: %s\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}

	return body
}

func poe_status() []byte {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/swcfg_poe", *hostName), nil)
	if err != nil {
		check.ExitError(err)
	}
	req.Header.Set("session", sessionToken)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing body: %s\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}

	return body
}

func string_percent_to_float(percents string) float64 {
	to_return, err := strconv.ParseFloat(strings.TrimSuffix(percents, "%"), 64)
	if err != nil {
		return 0
	}
	return to_return
}

func human_bytes(bytes uint64) string {
	if bytes == 0 {
		return "0 bytes"
	}

	const unit = 1024
	sizes := []string{"bytes", "KiB", "MiB", "GiB", "TiB", "PiB"}

	exp := int(math.Floor(math.Log(float64(bytes)) / math.Log(unit)))
	if exp >= len(sizes) {
		exp = len(sizes) - 1
	}

	value := float64(bytes) / math.Pow(unit, float64(exp))
	s := fmt.Sprintf("%.2f %s", value, sizes[exp])
	// Remove unnecessary decimals, like ".00"
	if strings.HasSuffix(s, ".00 "+sizes[exp]) {
		s = fmt.Sprintf("%.0f %s", value, sizes[exp])
	}
	return s
}
