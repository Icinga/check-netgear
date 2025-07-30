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

var allowedModes = map[string]bool{
	"basic": true,
	"ports": true,
}

var hostName *string

func main() {
	hidecpu := pflag.Bool("nocpu", false, "Hide the CPU info")
	hideram := pflag.Bool("noram", false, "Hide the RAM info")
	hidetemp := pflag.Bool("notemp", false, "Hide the Temperature info")
	hidefans := pflag.Bool("nofans", false, "Hide the Fans info")

	mode := pflag.StringSlice("mode", []string{"basic"}, "Output modes to enable {basic|short}")
	hostName = pflag.StringP("hostname", "H", "http://192.168.112.19", "Hostname to use")
	username := pflag.StringP("username", "u", "", "Username to use for authentication")
	password := pflag.StringP("password", "p", "", "Password to use for authentication")

	portsToCheck := pflag.IntSlice("port", []int{1, 2, 3, 4, 5, 6, 7, 8}, "Ports to check")

	help := pflag.BoolP("help", "h", false, "Show this help")

	pflag.Parse()

	if *help || *username == "" || *password == "" {
		pflag.Usage()
		return
	}
	for _, m := range *mode {
		if _, ok := allowedModes[m]; !ok {
			fmt.Printf("Invalid display mode: %s\n", m)
			return
		}
	}

	*hostName += "/api/v1"

	err := login(*username, *password)
	if err != nil {
		fmt.Printf("Error while trying to login: %v\n", err)
		return
	}

	//Basic output
	if slices.Contains(*mode, "basic") {
		var data map[string]any
		inputData := device_info()
		err = json.Unmarshal(inputData, &data)
		if err != nil {
			panic(err)
		}
		deviceInfo := data["deviceInfo"].(map[string]any)

		// Parse base metrics
		upTime := deviceInfo["details"].([]any)[0].(map[string]any)["upTime"]
		cpuUsage := string_percent_to_float(deviceInfo["cpu"].([]any)[0].(map[string]any)["usage"].(string))
		memoryUsage := string_percent_to_float(deviceInfo["memory"].([]any)[0].(map[string]any)["usage"].(string))

		fanDetails := deviceInfo["fan"].([]any)[0].(map[string]any)["details"].([]any)[0].(map[string]any)
		fanName := fanDetails["desc"].(string)
		fanSpeed := fanDetails["speed"].(float64)

		sensorDetails := deviceInfo["sensor"].([]any)[0].(map[string]any)["details"].([]any)

		//STATUSES CALC
		worstStatus := check.OK

		// Create result container
		o := result.Overall{}
		// o.Add(check.OK, fmt.Sprintf("Device Info: Uptime - %v", upTime))

		// CPU check
		if !*hidecpu {
			cpuStatus := check.OK
			if cpuUsage >= 90 {
				cpuStatus = check.Critical
				worstStatus = check.Critical
			} else if cpuUsage >= 60 {
				cpuStatus = check.Warning
				worstStatus = check.Warning
			}
			cpuCheck := result.PartialResult{
				Output: fmt.Sprintf("CPU Usage: %.2f%%", cpuUsage),
			}
			cpuCheck.SetState(cpuStatus)
			cpuCheck.Perfdata.Add(&perfdata.Perfdata{
				Label: "CPU",
				Value: cpuUsage,
				Min:   0,
				Max:   100,
			})
			o.AddSubcheck(cpuCheck)
		}

		// Memory check
		if !*hideram {
			memoryStatus := check.OK
			if memoryUsage >= 90 {
				memoryStatus = check.Critical
				if worstStatus != check.Critical {
					worstStatus = check.Critical
				}
			} else if memoryUsage >= 70 {
				memoryStatus = check.Warning
				if worstStatus < check.Warning {
					worstStatus = check.Warning
				}
			}
			memoryCheck := result.PartialResult{
				Output: fmt.Sprintf("RAM Usage: %.2f%%", memoryUsage),
			}
			memoryCheck.SetState(memoryStatus)
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
				if temp >= 70 {
					status = check.Critical
					if worstStatus != check.Critical {
						worstStatus = check.Critical
					}
				} else if temp >= 50 {
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
				sub.SetState(status)
				sub.Perfdata.Add(&perfdata.Perfdata{
					Label: desc,
					Value: temp,
					Min:   0,
				})
				temperatureCheck.AddSubcheck(sub)
			}
			temperatureCheck.SetState(worstTempStatus)
			o.AddSubcheck(temperatureCheck)
		}

		// Fan checks
		if !*hidefans {
			fansCheck := result.PartialResult{Output: "Fans"}
			fanStatus := check.OK
			if fanSpeed == 0 {
				fanStatus = check.Warning
				if worstStatus < check.Warning {
					worstStatus = check.Warning
				}
			}
			fanSub := result.PartialResult{
				Output: fmt.Sprintf("%s: %.0f RPM", fanName, fanSpeed),
			}
			fanSub.SetState(fanStatus)
			fanSub.Perfdata.Add(&perfdata.Perfdata{
				Label: "Fans speed",
				Value: fanSpeed,
				Min:   0,
			})
			fansCheck.AddSubcheck(fanSub)
			fansCheck.SetState(fanStatus)
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
		portsInfo := dataIn["portStatistics"].(map[string]any)
		portRows := portsInfo["rows"].([]interface{})
		/*for _, portRow := range portRows {
			fmt.Printf("Port %v: %v\n", portRow.(map[string]any)["port"], portRow.(map[string]any))
		}*/

		//STATUSES CALC
		worstStatus := check.OK

		// Create result container
		o := result.Overall{}
		// o.Add(check.OK, fmt.Sprintf("Device Info: Uptime - %v", upTime))

		// Temperature checks
		for _, port := range portRows {
			portNumber := port.(map[string]any)["port"].(float64)
			if slices.Contains(*portsToCheck, int(portNumber)) {
				portsCheck := result.PartialResult{Output: fmt.Sprintf("Port %v", portNumber)}
				worstPortsStatus := check.OK

				// inTotalPkts - Total IN packets
				if true {
					inTotalPkts := port.(map[string]any)["inTotalPkts"].(float64)
					status := check.OK
					if inTotalPkts < -1 { // Check for critical in inTotalPkts values
						status = check.Critical
						if worstStatus != check.Critical {
							worstStatus = check.Critical
						}
					} else if portNumber < 0 { // The same as above
						status = check.Warning
						if worstStatus < check.Warning {
							worstStatus = check.Warning
						}
					}
					if status > worstPortsStatus {
						worstPortsStatus = status
					}
					subInTotalPkts := result.PartialResult{
						Output: fmt.Sprintf("InTotalPkts: %d", int(inTotalPkts)),
					}
					subInTotalPkts.SetState(status)
					subInTotalPkts.Perfdata.Add(&perfdata.Perfdata{
						Label: fmt.Sprintf("port %v", portNumber),
						Value: portNumber,
						Min:   0,
					})
					portsCheck.AddSubcheck(subInTotalPkts)
				}

				// inDropPkts - Total drop IN packets
				if true {
					inDropPkts := port.(map[string]any)["inDropPkts"].(float64)
					status := check.OK
					if inDropPkts >= 2600 {
						status = check.Critical
						if worstStatus != check.Critical {
							worstStatus = check.Critical
						}
					} else if inDropPkts >= 2400 {
						status = check.Warning
						if worstStatus < check.Warning {
							worstStatus = check.Warning
						}
					}
					if status > worstPortsStatus {
						worstPortsStatus = status
					}
					subInDropPkts := result.PartialResult{
						Output: fmt.Sprintf("InDropPkts: %v", inDropPkts),
					}
					subInDropPkts.SetState(status)
					subInDropPkts.Perfdata.Add(&perfdata.Perfdata{
						Label: fmt.Sprintf("port %v", portNumber),
						Value: portNumber,
						Min:   0,
					})
					portsCheck.AddSubcheck(subInDropPkts)
				}

				// inOctets - Total IN octets
				if true {
					inOctets := port.(map[string]any)["inOctets"].(float64)
					status := check.OK
					if inOctets < -1 { // Ckeck for critical in octets values
						status = check.Critical
						if worstStatus != check.Critical {
							worstStatus = check.Critical
						}
					} else if portNumber < 0 { // Check for warning in octets values
						status = check.Warning
						if worstStatus < check.Warning {
							worstStatus = check.Warning
						}
					}
					if status > worstPortsStatus {
						worstPortsStatus = status
					}
					subInOctets := result.PartialResult{
						Output: fmt.Sprintf("Bytes: %d", int(inOctets)),
					}
					subInOctets.SetState(status)
					subInOctets.Perfdata.Add(&perfdata.Perfdata{
						Label: fmt.Sprintf("port %v", portNumber),
						Value: portNumber,
						Min:   0,
					})
					portsCheck.AddSubcheck(subInOctets)
				}

				portsCheck.SetState(worstPortsStatus)
				o.AddSubcheck(portsCheck)
			}
		}

		o.Add(worstStatus, fmt.Sprintf("Ports Statistics"))

		// Output result
		fmt.Println(o.GetOutput())
	}
}
