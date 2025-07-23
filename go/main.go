package main

import "fmt"

var functions = []Function{
	Function{"Retrieve Switch information", device_info},
	Function{"Retrieve the current SNTP server configuration", sntp_server_cfg},
	Function{"Retrieve firmware image information", imageInfo},
	Function{"Retrieve Management interface IP Address configuration", vlan_interface_route},
	Function{"Retrieve Out-of-Band Management Port configuration", special_interface_route},
	Function{"Retrieve Switch HTTP and HTTPS configuration", ip_http},
	Function{"List all active VLANs with their associated profiles", profile_list},
	Function{"Retrieve available network profiles", profile},
	Function{"Retrieve Auto-LAG configuration", sw_auto_lag_cfg},
	Function{"Retrieve LLDP neighbor information", neighbor},
	Function{"Retrieve PoE status for a specific port or all ports", swcfg_poe},
	Function{"Retrieve port statistics", port_statistics},
	Function{"Retrieve switch port status", swcfg_ports_status},
	Function{"Download tech support file of the switch", tech_support},
	Function{"Download the current configuration file", device_config},
	Function{"Start a ping test", ping_test},
	Function{"Perform a traceroute test", trace_test},
	Function{"Perform a DNS lookup", dns_lookup},
	Function{"Retrieve the current PTP mode", sw_ptp_cfg},
	Function{"Logout", logout},
}

func check_everything() {
	for index, element := range functions {
		fmt.Printf("\n%d - %s\n", index, element.label)
		element.function()
	}
}

func main() {

	fmt.Print("Logging in..\n")
	login("admin", "password")
	fmt.Print("Device information:\n")
	device_info()
	// check_everything()

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
