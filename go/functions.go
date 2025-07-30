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

	client := &http.Client{Timeout: timeout}
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
	defer resp.Body.Close()

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
