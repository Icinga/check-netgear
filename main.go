package main

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/pflag"
)

// We need this variable to be global, because it is used in functions.go
var hostName *string

func main() {
	// Get flags from the cli
	hidecpu := pflag.Bool("nocpu", false, "Hide the CPU info")
	hidemem := pflag.Bool("nomem", false, "Hide the RAM info")
	hidetemp := pflag.Bool("notemp", false, "Hide the Temperature info")
	hidefans := pflag.Bool("nofans", false, "Hide the Fans info")

	// Get arguments from the cli
	mode := pflag.StringSlice("mode", []string{"basic"}, "Output modes to enable {basic|short}")
	hostName = pflag.StringP("hostname", "H", "http://192.168.112.7", "Hostname to use")
	username := pflag.StringP("username", "u", "", "Username to use for authentication")
	password := pflag.StringP("password", "p", "", "Password to use for authentication")

	// Warning & Critical values for every metric
	CPU_WARN := pflag.Float64("cpu-warning", 50, "Provide the minimum for CPU usage warning")
	CPU_CRIT := pflag.Float64("cpu-critical", 90, "Provide the minimum for CPU usage critical")
	MEM_WARN := pflag.Float64("mem-warning", 50, "Provide the minimum for RAM usage warning")
	MEM_CRIT := pflag.Float64("mem-critical", 90, "Provide the minimum for RAM usage critical")
	FAN_WARN := pflag.Float64("fan-warning", 3000, "Provide the minimum for Fan speeds warning")

	TEMP_WARN := pflag.Float64("temp-warning", 50, "Provide the minimum for Temperature warning")
	TEMP_CRIT := pflag.Float64("temp-critical", 70, "Provide the minimum for Temperature critical")
	STATS_WARN := pflag.Float64("stats-warning", 5, "Provide the minimum for Port statistics warning")
	STATS_CRIT := pflag.Float64("stats-critical", 20, "Provide the minimum for Port statistics critical")

	// Ports to check if the mode 'ports' is present
	portsToCheck := pflag.IntSlice("port", []int{1, 2, 3, 4, 5, 6, 7, 8}, "Ports to check")

	help := pflag.BoolP("help", "h", false, "Show this help")

	// Parse all the flags for later usage
	pflag.Parse()

	// Display help if the -h is given, or no username / password is provided
	if *help || *username == "" || *password == "" {
		pflag.Usage()
		return
	}

	*hostName += "/api/v1"

	// Trying to log in
	err := login(*username, *password)
	if err != nil { // Error is present, display it
		fmt.Printf("Error while trying to login: %v\n", err)
		return
	}

	// Now everything is fine and the token is saved in a global variable in functions.go
	// We can proceed with checking for every mode checked and doing corresponding stuff

	//Basic output
	if slices.Contains(*mode, "basic") {
		var data map[string]any
		inputData := device_info()
		err = json.Unmarshal(inputData, &data)
		if err != nil {
			panic(err)
		}
		deviceInfo := data["deviceInfo"].(map[string]any)

		// Basic info
		upTime := deviceInfo["details"].([]any)[0].(map[string]any)["upTime"]
		cpuUsage := string_percent_to_float(deviceInfo["cpu"].([]any)[0].(map[string]any)["usage"].(string))
		memoryUsage := string_percent_to_float(deviceInfo["memory"].([]any)[0].(map[string]any)["usage"].(string))

		// Fans details
		fanDetails := deviceInfo["fan"].([]any)[0].(map[string]any)["details"].([]any)[0].(map[string]any)
		fanName := fanDetails["desc"].(string)
		fanSpeed := fanDetails["speed"].(float64)

		// Temperature details
		sensorDetails := deviceInfo["sensor"].([]any)[0].(map[string]any)["details"].([]any)

		// worstStatus is needed for 'inheriting' the worst status from the lower levels to the top
		worstStatus := check.OK

		// Create result container
		o := result.Overall{}

		// CPU check
		if !*hidecpu {
			cpuStatus := check.OK
			if cpuUsage >= *CPU_CRIT {
				cpuStatus = check.Critical
				worstStatus = check.Critical
			} else if cpuUsage >= *CPU_WARN {
				cpuStatus = check.Warning
				worstStatus = check.Warning
			}
			cpuCheck := result.PartialResult{
				Output: fmt.Sprintf("CPU Usage: %.2f%%", cpuUsage),
			}
			err := cpuCheck.SetState(cpuStatus)
			if err != nil {
				cpuCheck.SetState(check.Unknown)
			}
			cpuCheck.Perfdata.Add(&perfdata.Perfdata{
				Label: "CPU",
				Value: cpuUsage,
				Min:   0,
				Max:   100,
			})
			o.AddSubcheck(cpuCheck)
		}

		// Memory check
		if !*hidemem {
			memoryStatus := check.OK
			if memoryUsage >= *MEM_CRIT {
				memoryStatus = check.Critical
				if worstStatus != check.Critical {
					worstStatus = check.Critical
				}
			} else if memoryUsage >= *MEM_WARN {
				memoryStatus = check.Warning
				if worstStatus < check.Warning {
					worstStatus = check.Warning
				}
			}
			memoryCheck := result.PartialResult{
				Output: fmt.Sprintf("RAM Usage: %.2f%%", memoryUsage),
			}
			err := memoryCheck.SetState(memoryStatus)
			if err != nil {
				memoryCheck.SetState(check.Unknown)
			}
			memoryCheck.Perfdata.Add(&perfdata.Perfdata{
				Label: "RAM",
				Value: memoryUsage,
				Min:   0,
				Max:   100,
			})
			o.AddSubcheck(memoryCheck)
		}

		// Temperature checks
		if !*hidetemp {
			temperatureCheck := result.PartialResult{Output: "Temperature"}
			worstTempStatus := check.OK
			for _, sensor := range sensorDetails {
				desc := sensor.(map[string]any)["desc"].(string)
				temp := sensor.(map[string]any)["temp"].(float64)

				status := check.OK
				if temp >= *TEMP_CRIT {
					status = check.Critical
					if worstStatus != check.Critical {
						worstStatus = check.Critical
					}
				} else if temp >= *TEMP_WARN {
					status = check.Warning
					if worstStatus < check.Warning {
						worstStatus = check.Warning
					}
				}
				if status > worstTempStatus {
					worstTempStatus = status
				}

				sub := result.PartialResult{
					Output: fmt.Sprintf("%s: %.1f°C", desc, temp),
				}
				err := sub.SetState(status)
				if err != nil {
					sub.SetState(check.Unknown)
				}
				sub.Perfdata.Add(&perfdata.Perfdata{
					Label: desc,
					Value: temp,
					Min:   0,
				})
				temperatureCheck.AddSubcheck(sub)
			}
			err := temperatureCheck.SetState(worstTempStatus)
			if err != nil {
				temperatureCheck.SetState(check.Unknown)
			}
			o.AddSubcheck(temperatureCheck)
		}

		// Fan checks
		if !*hidefans {
			fansCheck := result.PartialResult{Output: "Fans"}
			fanStatus := check.OK
			if fanSpeed > *FAN_WARN {
				fanStatus = check.Warning
				if worstStatus < check.Warning {
					worstStatus = check.Warning
				}
			}
			fanSub := result.PartialResult{
				Output: fmt.Sprintf("%s: %.0f RPM", fanName, fanSpeed),
			}
			err := fanSub.SetState(fanStatus)
			if err != nil {
				fanSub.SetState(check.Unknown)
			}
			fanSub.Perfdata.Add(&perfdata.Perfdata{
				Label: "Fans speed",
				Value: fanSpeed,
				Min:   0,
			})
			fansCheck.AddSubcheck(fanSub)
			err = fansCheck.SetState(fanStatus)
			if err != nil {
				fansCheck.SetState(check.Unknown)
			}
			o.AddSubcheck(fansCheck)
		}

		o.Add(worstStatus, fmt.Sprintf("Device Info: Uptime - %v", upTime))

		// Output result
		fmt.Println(o.GetOutput())
	}

	if slices.Contains(*mode, "ports") {
		var dataIn map[string]any
		portsIn := port_statistics("inbound")
		err = json.Unmarshal(portsIn, &dataIn)
		if err != nil {
			panic(err)
		}
		portsInInfo := dataIn["portStatistics"].(map[string]any)
		portRowsIn := portsInInfo["rows"].([]interface{})

		var dataOut map[string]any
		portsOut := port_statistics("outbound")
		err = json.Unmarshal(portsOut, &dataOut)
		if err != nil {
			panic(err)
		}
		portsOutInfo := dataOut["portStatistics"].(map[string]any)
		portRowsOut := portsOutInfo["rows"].([]interface{})
		/*for _, portRow := range portRows {
			fmt.Printf("Port %v: %v\n", portRow.(map[string]any)["port"], portRow.(map[string]any))
		}*/

		//STATUSES CALC
		worstStatus := check.OK

		// Create result container
		o := result.Overall{}
		// o.Add(check.OK, fmt.Sprintf("Device Info: Uptime - %v", upTime))

		// Ports checks
		for index, _ := range portRowsIn {
			portNumber := portRowsIn[index].(map[string]any)["port"].(float64)
			if slices.Contains(*portsToCheck, int(portNumber)) {
				portsCheck := result.PartialResult{Output: fmt.Sprintf("Port %v", portNumber)}
				worstPortsStatus := check.OK

				// inTotalPkts - Total IN packets
				if true {
					inTotalPkts := portRowsIn[index].(map[string]any)["inTotalPkts"].(float64)
					inDropPkts := portRowsIn[index].(map[string]any)["inDropPkts"].(float64)
					inOctets := portRowsIn[index].(map[string]any)["inOctets"].(float64)
					packetLossPercentage := 0.0
					if inTotalPkts > 0 {
						packetLossPercentage = inDropPkts / inTotalPkts * 100
					}
					status := check.OK
					if packetLossPercentage >= *STATS_CRIT { // Check for critical in dropped packets percentage
						status = check.Critical
						if worstStatus != check.Critical {
							worstStatus = check.Critical
						}
					} else if packetLossPercentage >= *STATS_WARN { // Check for warning in dropped packets percentage
						status = check.Warning
						if worstStatus < check.Warning {
							worstStatus = check.Warning
						}
					}
					if status > worstPortsStatus {
						worstPortsStatus = status
					}
					subInTotalPkts := result.PartialResult{
						Output: fmt.Sprintf("Total IN: %v; Packet loss: %.2f%%", human_bytes(uint64(inOctets)), packetLossPercentage),
					}
					err := subInTotalPkts.SetState(status)
					if err != nil {
						subInTotalPkts.SetState(check.Unknown)
					}
					subInTotalPkts.Perfdata.Add(&perfdata.Perfdata{
						Label: fmt.Sprintf("port %v in packet loss", portNumber),
						Value: packetLossPercentage,
						Min:   0,
						Max:   100,
					})
					portsCheck.AddSubcheck(subInTotalPkts)
				}

				// outTotalPkts - Total OUT packets
				if true {
					outTotalPkts := portRowsOut[index].(map[string]any)["outTotalPkts"].(float64)
					outDropPkts := portRowsOut[index].(map[string]any)["outDropPkts"].(float64)
					outOctets := portRowsOut[index].(map[string]any)["outOctets"].(float64)
					packetLossPercentage := 0.0
					if outTotalPkts > 0 {
						packetLossPercentage = outDropPkts / outTotalPkts * 100
					}
					status := check.OK
					if packetLossPercentage >= *STATS_CRIT { // Check for critical in dropped packets percentage
						status = check.Critical
						if worstStatus != check.Critical {
							worstStatus = check.Critical
						}
					} else if packetLossPercentage >= *STATS_WARN { // Check for warning in dropped packets percentage
						status = check.Warning
						if worstStatus < check.Warning {
							worstStatus = check.Warning
						}
					}
					if status > worstPortsStatus {
						worstPortsStatus = status
					}
					subOutTotalPkts := result.PartialResult{
						Output: fmt.Sprintf("Total OUT: %v; Packet loss: %.2f%%", human_bytes(uint64(outOctets)), packetLossPercentage),
					}
					err := subOutTotalPkts.SetState(status)
					if err != nil {
						subOutTotalPkts.SetState(check.Unknown)
					}
					subOutTotalPkts.Perfdata.Add(&perfdata.Perfdata{
						Label: fmt.Sprintf("port %v out packet loss", portNumber),
						Value: packetLossPercentage,
						Min:   0,
						Max:   100,
					})
					portsCheck.AddSubcheck(subOutTotalPkts)
				}

				err := portsCheck.SetState(worstPortsStatus)
				if err != nil {
					portsCheck.SetState(check.Unknown)
				}
				o.AddSubcheck(portsCheck)
			}
		}

		o.Add(worstStatus, "Ports Statistics")

		// Output result
		fmt.Println(o.GetOutput())
	}

	logout()
}
