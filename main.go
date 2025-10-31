package main

import (
	"encoding/json"
	"fmt"
	"main/netgear"
	"os"
	"slices"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/pflag"
)

func main() {
	hidecpu := pflag.Bool("nocpu", false, "Hide the CPU info")
	hidemem := pflag.Bool("nomem", false, "Hide the RAM info")
	hidetemp := pflag.Bool("notemp", false, "Hide the Temperature info")
	hidefans := pflag.Bool("nofans", false, "Hide the Fans info")

	mode := pflag.StringSlice("mode", []string{"basic"}, "Output modes to enable {basic|ports|poe|all}")
	hostName := pflag.StringP("hostname", "H", "http://192.168.112.7", "Hostname to use")
	username := pflag.StringP("username", "u", "", "Username for authentication")
	password := pflag.StringP("password", "p", "", "Password for authentication")

	// Thresholds
	CPU_WARN := pflag.Float64("cpu-warning", 50, "CPU usage warning threshold")
	CPU_CRIT := pflag.Float64("cpu-critical", 90, "CPU usage critical threshold")
	MEM_WARN := pflag.Float64("mem-warning", 50, "RAM usage warning threshold")
	MEM_CRIT := pflag.Float64("mem-critical", 90, "RAM usage critical threshold")
	FAN_WARN := pflag.Float64("fan-warning", 3000, "Fan speed warning threshold")
	TEMP_WARN := pflag.Float64("temp-warning", 50, "Temperature warning threshold")
	TEMP_CRIT := pflag.Float64("temp-critical", 70, "Temperature critical threshold")
	STATS_WARN := pflag.Float64("stats-warning", 5, "Port stats warning threshold")
	STATS_CRIT := pflag.Float64("stats-critical", 20, "Port stats critical threshold")

	portsToCheck := pflag.IntSlice("port", []int{1, 2, 3, 4, 5, 6, 7, 8}, "Ports to check")
	help := pflag.BoolP("help", "h", false, "Show this help")

	pflag.Parse()

	if *help || *username == "" || *password == "" {
		pflag.Usage()
		return
	}

	n := netgear.NewNetgear(*hostName, *username, *password)
	if err := n.Login(); err != nil {
		fmt.Printf("Error while trying to login: %v\n", err)
		return
	}
	defer func() { _ = n.Logout() }()

	if slices.Contains(*mode, "all") {
		*mode = append(*mode, "basic", "ports", "poe")
	}

	worstStatus := check.OK

	// Basic check
	if slices.Contains(*mode, "basic") {
		var deviceInfo netgear.DeviceInfo
		inputData, err := n.DeviceInfo()
		if err != nil {
			fmt.Printf("Error retrieving device info: %v\n", err)
			os.Exit(int(check.Unknown))
		}
		if err := json.Unmarshal(inputData, &deviceInfo); err != nil {
			fmt.Printf("Failed to parse JSON: %v\n", err)
			os.Exit(int(check.Unknown))
		}

		upTime := deviceInfo.DeviceInfo.Details[0].Uptime
		cpuUsage := netgear.StringPercentToFloat(deviceInfo.DeviceInfo.Cpu[0].Usage)
		memUsage := netgear.StringPercentToFloat(deviceInfo.DeviceInfo.Memory[0].Usage)
		fan := deviceInfo.DeviceInfo.Fan[0].Details[0]
		fanName := fan.Description
		fanSpeed := fan.Speed
		sensorDetails := deviceInfo.DeviceInfo.Sensor[0].Details

		o := result.Overall{}

		// CPU
		if !*hidecpu {
			cpuStatus := statusByThreshold(cpuUsage, *CPU_WARN, *CPU_CRIT)
			worstStatus = maxStatus(worstStatus, cpuStatus)
			cpuCheck := result.PartialResult{
				Output: fmt.Sprintf("CPU Usage: %.2f%%", cpuUsage),
			}
			_ = cpuCheck.SetState(cpuStatus)
			cpuCheck.Perfdata.Add(&perfdata.Perfdata{Label: "CPU", Value: cpuUsage, Min: 0, Max: 100})
			o.AddSubcheck(cpuCheck)
		}

		// Memory
		if !*hidemem {
			memStatus := statusByThreshold(memUsage, *MEM_WARN, *MEM_CRIT)
			worstStatus = maxStatus(worstStatus, memStatus)
			memCheck := result.PartialResult{
				Output: fmt.Sprintf("RAM Usage: %.2f%%", memUsage),
			}
			_ = memCheck.SetState(memStatus)
			memCheck.Perfdata.Add(&perfdata.Perfdata{Label: "RAM", Value: memUsage, Min: 0, Max: 100})
			o.AddSubcheck(memCheck)
		}

		// Temperature
		if !*hidetemp {
			tempCheck := result.PartialResult{Output: "Temperature"}
			worstTempStatus := check.OK
			for _, s := range sensorDetails {
				desc := s.Description
				temp := s.Temperature
				status := statusByThreshold(temp, *TEMP_WARN, *TEMP_CRIT)
				worstTempStatus = maxStatus(worstTempStatus, status)
				worstStatus = maxStatus(worstStatus, status)

				sub := result.PartialResult{
					Output: fmt.Sprintf("%s: %.1f°C", desc, temp),
				}
				_ = sub.SetState(status)
				sub.Perfdata.Add(&perfdata.Perfdata{Label: desc, Value: temp, Min: 0})
				tempCheck.AddSubcheck(sub)
			}
			_ = tempCheck.SetState(worstTempStatus)
			o.AddSubcheck(tempCheck)
		}
		// Fans
		if !*hidefans {
			fanStatus := check.OK
			if fanSpeed > *FAN_WARN {
				fanStatus = check.Warning
				worstStatus = maxStatus(worstStatus, fanStatus)
			}
			fansCheck := result.PartialResult{Output: "Fans"}
			fanSub := result.PartialResult{
				Output: fmt.Sprintf("%s: %.0f RPM", fanName, fanSpeed),
			}
			_ = fanSub.SetState(fanStatus)
			fanSub.Perfdata.Add(&perfdata.Perfdata{
				Label: "Fan Speed", Value: fanSpeed, Min: 0,
			})
			fansCheck.AddSubcheck(fanSub)
			_ = fansCheck.SetState(fanStatus)
			o.AddSubcheck(fansCheck)
		}

		o.Add(worstStatus, fmt.Sprintf("Device Info: Uptime - %v", upTime))
		fmt.Println(o.GetOutput())
	}

	// ports
	if slices.Contains(*mode, "ports") {
		var inStats, outStats netgear.PortStatistics
		portsIn, _ := n.PortStatistics("inbound")
		portsOut, _ := n.PortStatistics("outbound")
		_ = json.Unmarshal(portsIn, &inStats)
		_ = json.Unmarshal(portsOut, &outStats)

		inRows := inStats.PortStatistics.Rows
		outRows := outStats.PortStatistics.Rows

		o := result.Overall{}
		worstPortsStatus := check.OK

		for i := range inRows {
			in := inRows[i]
			out := outRows[i]
			portNumber := in.Port

			if slices.Contains(*portsToCheck, portNumber) {
				portCheck := result.PartialResult{Output: fmt.Sprintf("Port %v", portNumber)}
				inLoss := lossPercent(in.InDropPkts, in.InTotalPkts)
				outLoss := lossPercent(out.OutDropPkts, out.OutTotalPkts)

				inStatus := statusByThreshold(inLoss, *STATS_WARN, *STATS_CRIT)
				outStatus := statusByThreshold(outLoss, *STATS_WARN, *STATS_CRIT)
				portStatus := maxStatus(inStatus, outStatus)
				worstPortsStatus = maxStatus(worstPortsStatus, portStatus)
				worstStatus = maxStatus(worstStatus, portStatus)

				addPerfSubcheck := func(label string, loss float64, status int) {
					sub := result.PartialResult{
						Output: fmt.Sprintf("%s: %.2f%% loss", label, loss),
					}
					_ = sub.SetState(status)
					sub.Perfdata.Add(&perfdata.Perfdata{
						Label: fmt.Sprintf("port %v %s loss", portNumber, label),
						Value: loss, Min: 0, Max: 100,
					})
					portCheck.AddSubcheck(sub)
				}

				addPerfSubcheck("IN", inLoss, inStatus)
				addPerfSubcheck("OUT", outLoss, outStatus)

				_ = portCheck.SetState(portStatus)
				o.AddSubcheck(portCheck)
			}
		}

		o.Add(worstPortsStatus, "Ports Statistics")
		fmt.Println(o.GetOutput())
	}

	// poe stuff
	if slices.Contains(*mode, "poe") {
		o := result.Overall{}
		worstPoeStatus := check.OK
		var poeStatus netgear.PoeStatus
		inputData, _ := n.PoeStatus()
		_ = json.Unmarshal(inputData, &poeStatus)

		for _, port := range poeStatus.PoePortConfig {
			state := "disabled"
			if port.Enable {
				state = "enabled"
			}

			status := check.OK
			if port.CurrentPower > port.PowerLimit {
				status = check.Critical
			} else if port.CurrentPower == port.PowerLimit {
				status = check.Warning
			}

			worstPoeStatus = maxStatus(worstPoeStatus, status)
			worstStatus = maxStatus(worstStatus, status)

			poeCheck := result.PartialResult{
				Output: fmt.Sprintf(
					"Port %v is %v. Current power: %.2f/%.2fV",
					port.Port, state, port.CurrentPower/1000, port.PowerLimit/1000,
				),
			}
			_ = poeCheck.SetState(status)
			poeCheck.Perfdata.Add(&perfdata.Perfdata{
				Label: fmt.Sprintf("port %v power", port.Port),
				Value: port.CurrentPower, Min: 0, Max: port.PowerLimit,
			})
			o.AddSubcheck(poeCheck)
		}

		o.Add(worstPoeStatus, "Power over Ethernet Statistics")
		fmt.Println(o.GetOutput())
	}

	os.Exit(int(worstStatus))
}

// util
func statusByThreshold(value, warn, crit float64) int {
	switch {
	case value >= crit:
		return check.Critical
	case value >= warn:
		return check.Warning
	default:
		return check.OK
	}
}

func maxStatus(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func lossPercent(drop, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return drop / total * 100
}
