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
func login(name, password string) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(fmt.Sprintf("%s/login", baseURL),
		url.Values{
			"user": {fmt.Sprintf("{\"name\": {%s}, \"password\": {%s}}", name, password)},
		},
	)
	if err != nil {
		fmt.Printf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
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
	resp, err := http.Get(fmt.Sprintf("%s/device_info", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

// DEVICE_FAN
func device_fan(fanMode int) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(fmt.Sprintf("%s/device_fan", baseURL),
		url.Values{"fanMode": {fmt.Sprintf("%d", fanMode)}})
	if err != nil {
		fmt.Printf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

// SNTP_SERVER_CFG
func sntp_server_cfg() {
	resp, err := http.Get(fmt.Sprintf("%s/sntp_server_cfg", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func imageInfo() {
	resp, err := http.Get(fmt.Sprintf("%s/imageInfo", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func vlan_interface_route() {
	resp, err := http.Get(fmt.Sprintf("%s/vlan_interface_route", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func special_interface_route() {
	resp, err := http.Get(fmt.Sprintf("%s/special_interface_route", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func ip_http() {
	resp, err := http.Get(fmt.Sprintf("%s/ip_http", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func profile_list() {
	resp, err := http.Get(fmt.Sprintf("%s/profile/list", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func profile() {
	resp, err := http.Get(fmt.Sprintf("%s/profile", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func sw_auto_lag_cfg() {
	resp, err := http.Get(fmt.Sprintf("%s/sw_auto_lag_cfg", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func neighbor() {
	resp, err := http.Get(fmt.Sprintf("%s/neighbor", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func swcfg_poe() {
	resp, err := http.Get(fmt.Sprintf("%s/swcfg_poe", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func port_statistics() {
	resp, err := http.Get(fmt.Sprintf("%s/port_statistics", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func swcfg_ports_status() {
	resp, err := http.Get(fmt.Sprintf("%s/swcfg_ports_status", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func tech_support() {
	resp, err := http.Get(fmt.Sprintf("%s/tech_support", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func device_config() {
	resp, err := http.Get(fmt.Sprintf("%s/device_config", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func ping_test() {
	resp, err := http.Get(fmt.Sprintf("%s/ping_test", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func trace_test() {
	resp, err := http.Get(fmt.Sprintf("%s/trace_test", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func dns_lookup() {
	resp, err := http.Get(fmt.Sprintf("%s/dns_lookup", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func sw_ptp_cfg() {
	resp, err := http.Get(fmt.Sprintf("%s/sw_ptp_cfg", baseURL))
	if err != nil {
		check.ExitError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}

func main() {

	login("admin", "password")
	device_info()
	device_fan(2)
	sntp_server_cfg()
	imageInfo()
	vlan_interface_route()
	special_interface_route()
	ip_http()
	profile_list()
	profile()
	sw_auto_lag_cfg()
	neighbor()
	swcfg_poe()
	port_statistics()
	swcfg_poe()
	swcfg_ports_status()
	tech_support()
	device_config()
	ping_test()
	trace_test()
	dns_lookup()
	sw_ptp_cfg()
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
