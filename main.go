package main

import (
	"flag"
	"fmt"
	"github.com/icinga/check-netgear/internal/checks"
	"github.com/icinga/check-netgear/netgear"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

// So that flag supports slices
type stringSliceFlag []string

func (s *stringSliceFlag) String() string { return strings.Join(*s, ",") }
func (s *stringSliceFlag) Set(v string) error {
	*s = append(*s, v)
	return nil
}

type intSliceFlag []int

func (i *intSliceFlag) String() string {
	parts := make([]string, 0, len(*i))
	for _, v := range *i {
		parts = append(parts, strconv.Itoa(v))
	}
	return strings.Join(parts, ",")
}
func (i *intSliceFlag) Set(v string) error {
	n, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	*i = append(*i, n)
	return nil
}

func main() {
	hidecpu := flag.Bool("nocpu", false, "Hide the CPU info")
	hidemem := flag.Bool("nomem", false, "Hide the RAM info")
	hidetemp := flag.Bool("notemp", false, "Hide the Temperature info")
	hidefans := flag.Bool("nofans", false, "Hide the Fans info")

	mode := stringSliceFlag{}
	flag.Var(&mode, "mode", "Output modes to enable {basic|ports|poe|all} (repeatable)")

	baseURL := flag.String("base-url", "http://192.168.0.239", "Base URL to use")

	username := flag.String("username", "", "Username for authentication")
	password := flag.String("password", "", "Password for authentication")

	// Thresholds
	cpuWarn := flag.Float64("cpu-warning", 50, "CPU usage warning threshold")
	cpuCrit := flag.Float64("cpu-critical", 90, "CPU usage critical threshold")
	memWarn := flag.Float64("mem-warning", 50, "RAM usage warning threshold")
	memCrit := flag.Float64("mem-critical", 90, "RAM usage critical threshold")
	fanWarn := flag.Float64("fan-warning", 3000, "Fan speed warning threshold")
	tempWarn := flag.Float64("temp-warning", 50, "Temperature warning threshold")
	tempCrit := flag.Float64("temp-critical", 70, "Temperature critical threshold")
	statsWarn := flag.Float64("stats-warning", 5, "Port stats warning threshold")
	statsCrit := flag.Float64("stats-critical", 20, "Port stats critical threshold")

	portsToCheck := intSliceFlag{1, 2, 3, 4, 5, 6, 7, 8}
	flag.Var(&portsToCheck, "port", "Ports to check (repeatable)")

	help := flag.Bool("help", false, "Show this help")
	flag.BoolVar(help, "h", false, "Show this help (shorthand)")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *username == "" || *password == "" {
		fmt.Println("Both username and password are required")
		os.Exit(check.Unknown)
	}

	n, err := netgear.NewNetgear(*baseURL, *username, *password)
	if err != nil {
		fmt.Printf("URL error: %v", err)
		os.Exit(check.Unknown)
	}
	if err := n.Login(); err != nil {
		fmt.Printf("Error while trying to login: %v\n", err)
		os.Exit(check.Unknown)
	}
	defer func() { _ = n.Logout() }()

	if len(mode) == 0 {
		mode = append(mode, "basic")
	} else if slices.Contains(mode, "all") {
		mode = append(mode, "basic", "ports", "poe")
	}

	worstStatus := check.OK

	// Basic check
	if slices.Contains(mode, "basic") {
		deviceInfo, err := n.DeviceInfo()
		if err != nil {
			fmt.Printf("Error retrieving device info: %v\n", err)
			os.Exit(check.Unknown)
		}

		o := result.Overall{}

		// CPU
		if !*hidecpu {
			cpuUsage, err := netgear.StringPercentToFloat(deviceInfo.DeviceInfo.Cpu[0].Usage)
			if err != nil {
				fmt.Printf("Error parsing CPU usage: %v\n", err)
				os.Exit(check.Unknown)
			}

			cpuPartial, err := checks.CheckCPU(cpuUsage, *cpuWarn, *cpuCrit)
			if err != nil {
				errRes := result.NewPartialResult()
				errRes.Output = fmt.Sprintf("CPU check error: %v", err)
				o.AddSubcheck(errRes)
			} else {
				worstStatus = max(worstStatus, cpuPartial.GetStatus())
				o.AddSubcheck(*cpuPartial)
			}
		}

		// Memory
		if !*hidemem {
			memUsage, err := netgear.StringPercentToFloat(deviceInfo.DeviceInfo.Memory[0].Usage)
			if err != nil {
				fmt.Printf("Error parsing Memory usage: %v\n", err)
				os.Exit(check.Unknown)
			}

			memPartial, err := checks.CheckMemory(memUsage, *memWarn, *memCrit)
			if err != nil {
				errRes := result.NewPartialResult()
				errRes.Output = fmt.Sprintf("Memory check error: %v", err)
				o.AddSubcheck(errRes)
			} else {
				worstStatus = max(worstStatus, memPartial.GetStatus())
				o.AddSubcheck(*memPartial)
			}
		}

		// Temperature
		if !*hidetemp {
			sensorDetails := deviceInfo.DeviceInfo.Sensor[0].Details
			tempPartial, err := checks.CheckTemperature(sensorDetails, *tempWarn, *tempCrit)
			if err != nil {
				errRes := result.NewPartialResult()
				errRes.Output = fmt.Sprintf("Temperatuer check error: %v", err)
				o.AddSubcheck(errRes)
			} else {
				worstStatus = max(worstStatus, tempPartial.GetStatus())
				o.AddSubcheck(*tempPartial)
			}
		}

		// Fans
		if !*hidefans {
			fan := deviceInfo.DeviceInfo.Fan[0].Details[0]
			fanPartial, err := checks.CheckFans(fan.Description, fan.Speed, *fanWarn)
			if err != nil {
				errRes := result.NewPartialResult()
				errRes.Output = fmt.Sprintf("Fans check error: %v", err)
				o.AddSubcheck(errRes)
			} else {
				worstStatus = max(worstStatus, fanPartial.GetStatus())
				o.AddSubcheck(*fanPartial)
			}
		}

		upTime := deviceInfo.DeviceInfo.Details[0].Uptime
		o.Add(worstStatus, fmt.Sprintf("Device Info: Uptime - %v", upTime))
		fmt.Println(o.GetOutput())
	}

	// ports
	if slices.Contains(mode, "ports") {
		o := result.Overall{}
		portsIn, err := n.PortStatistics("inbound")
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("Inbound port check error: %v", err)
			o.AddSubcheck(errRes)
		}
		portsOut, err := n.PortStatistics("outbound")
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("Outbound port check error: %v", err)
			o.AddSubcheck(errRes)
		}

		inRows := portsIn.PortStatistics.Rows
		outRows := portsOut.PortStatistics.Rows

		portsPartial, err := checks.CheckPorts(inRows, outRows, portsToCheck, *statsWarn, *statsCrit)
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("Ports check error: %v", err)
			o.AddSubcheck(errRes)
		} else {
			worstStatus = max(worstStatus, portsPartial.GetStatus())
			o.AddSubcheck(*portsPartial)
			fmt.Println(o.GetOutput())
		}
	}

	// poe stuff
	if slices.Contains(mode, "poe") {
		o := result.Overall{}
		poeStatus, err := n.PoeStatus()
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("PoE check error: %v", err)
			o.AddSubcheck(errRes)
		}

		poePartial, err := checks.CheckPoe(poeStatus.PoePortConfig)
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("PoE check error: %v", err)
			o.AddSubcheck(errRes)
		} else {
			worstStatus = max(worstStatus, poePartial.GetStatus())
			o.AddSubcheck(*poePartial)
			fmt.Println(o.GetOutput())
		}
	}

	os.Exit(worstStatus)
}
