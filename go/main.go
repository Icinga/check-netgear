package main

import (
	"encoding/json"
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/pflag"
	"slices"
)

var allowedModes = map[string]bool{
	"basic": true,
	"short": true,
}
var showcpu *bool
var showram *bool
var showtemp *bool
var showfans *bool

var mode *[]string
var hostName *string

func main() {
	// Arguments:
	// --nocpu - Do NOT display the CPU stat
	// --noram - Do NOT display the RAM stat
	// --notemp - Do NOT display the Temperature stat
	// --nofan - Do NOT display the Fans stat
	// -H - ip
	// -u - user
	// -p - password

	showcpu = pflag.Bool("nocpu", false, "Hide the CPU")
	showram = pflag.Bool("noram", false, "Hide the RAM")
	showtemp = pflag.Bool("notemp", false, "Hide the Temperature")
	showfans = pflag.Bool("nofans", false, "Hide the Fans")

	mode = pflag.StringSlice("mode", []string{}, "Modes to enable {basic|short}")
	hostName = pflag.StringP("hostname", "H", "http://192.168.0.239", "Hostname to use")
	username := pflag.StringP("username", "u", "", "Username to use for authentication")
	password := pflag.StringP("password", "p", "", "Password to use for authentication")

	help := pflag.BoolP("help", "h", false, "Show this help")

	pflag.Parse()

	if *help || *username == "" || *password == "" {
		pflag.Usage()
		return
	}

	if len(*mode) == 0 { //We want the basic to be the default one
		*mode = append(*mode, "basic")
	}

	*showcpu = !*showcpu
	*showram = !*showram
	*showtemp = !*showtemp
	*showfans = !*showfans

	*hostName += "/api/v1"

	err := login(*username, *password)
	if err != nil {
		fmt.Println("Error while trying to login")
		fmt.Println(err)
		return
	}

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
		if *showcpu {
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
		if *showram {
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
		if *showtemp {
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
		if *showfans {
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
}
