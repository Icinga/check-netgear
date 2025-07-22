package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/NETWAYS/go-check"
)

var baseURL = "http://192.168.0.239/api/v1"

// LOGIN
func login(name, password string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(fmt.Sprintf("%s/login", baseURL),
		url.Values{
			"user": {fmt.Sprintf("{\"name\": {%s}, \"password\": {%s}}", name, password)},
		},
	)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	} else {
		fmt.Printf("LOGIN - There were no errors!\n")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
	return string(body), nil
}

// LOGOUT
func logout() {
	resp, err := http.Get(fmt.Sprintf("%s/logout", baseURL))
	if err != nil {
		check.ExitError(err)
	} else {
		fmt.Printf("LOGOUT - There were no errors!\n")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

// DEVICE_INFO
func device_info() {
	resp, err := http.Get(fmt.Sprintf("%s/device_info", baseURL))
	if err != nil {
		check.ExitError(err)
	} else {
		fmt.Printf("DEVICE_INFO - There were no errors!\n")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

// DEVICE_FAN
func device_fan(fanMode int) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(fmt.Sprintf("%s/device_fan", baseURL),
		url.Values{"fanMode": {fmt.Sprintf("%d", fanMode)}})
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	} else {
		fmt.Printf("DEVICE_FAN - There were no errors!\n")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
	return string(body), nil
}

func main() {

	login("admin", "password")
	device_info()
	device_fan(2)
	logout()

	// config := check.NewConfig()
	// config.Name = "check_test"
	// config.Readme = `Test Plugin`
	// config.Version = "1.0.0"
	// _ = config.FlagSet.StringP("hostname", "H", "localhost", "Hostname to check")
	// config.ParseArguments()
	// Some checking should be done here, when --help is not passed
	// check.Exitf(check.OK, "Everything is fine - answer=%d", 42)
	// Output:
	// OK - Everything is fine - answer=42
}
