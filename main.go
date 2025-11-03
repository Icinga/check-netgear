package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"main/internal/checks"
	"main/internal/utils"
	"main/netgear"
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

	var mode stringSliceFlag = []string{"basic"}
	flag.Var(&mode, "mode", "Output modes to enable {basic|ports|poe|all} (repeatable)")

	baseURL := flag.String("base-url", "http://192.168.1.1", "Base URL to use")
	flag.StringVar(baseURL, "H", "http://192.168.1.1", "Base URL to use (shorthand)")

	username := flag.String("username", "", "Username for authentication")
	flag.StringVar(username, "u", "", "Username for authentication (shorthand)")
	password := flag.String("password", "", "Password for authentication")
	flag.StringVar(password, "p", "", "Password for authentication (shorthand)")

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

	var portsToCheck intSliceFlag = []int{1, 2, 3, 4, 5, 6, 7, 8}
	flag.Var(&portsToCheck, "port", "Ports to check (repeatable)")

	help := flag.Bool("help", false, "Show this help")
	flag.BoolVar(help, "h", false, "Show this help (shorthand)")

	flag.Parse()

	if *help || *username == "" || *password == "" {
		flag.Usage()
		return
	}

	n := netgear.NewNetgear(*baseURL, *username, *password)
	if err := n.Login(); err != nil {
		fmt.Printf("Error while trying to login: %v\n", err)
		os.Exit(int(check.Unknown))
	}
	defer func() { _ = n.Logout() }()

	if slices.Contains(mode, "all") {
		mode = append(mode, "basic", "ports", "poe")
	}

	worstStatus := check.OK

	// Basic check
	if slices.Contains(mode, "basic") {
		deviceInfo, err := n.DeviceInfo()
		if err != nil {
			fmt.Printf("Error retrieving device info: %v\n", err)
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
			cpuPartial, cpuStatus := checks.CheckCPU(cpuUsage, *cpuWarn, *cpuCrit)
			worstStatus = max(worstStatus, cpuStatus)
			o.AddSubcheck(cpuPartial)
		}

		// Memory
		if !*hidemem {
			memPartial, memStatus := checks.CheckMemory(memUsage, *memWarn, *memCrit)
			worstStatus = max(worstStatus, memStatus)
			o.AddSubcheck(memPartial)
		}

		// Temperature
		if !*hidetemp {
			tempPartial, tempStatus := checks.CheckTemperature(sensorDetails, *tempWarn, *tempCrit)
			worstStatus = max(worstStatus, tempStatus)
			o.AddSubcheck(tempPartial)
		}

		// Fans
		if !*hidefans {
			fanPartial, fanStatus := checks.CheckFans(fanName, fanSpeed, *fanWarn)
			worstStatus = max(worstStatus, fanStatus)
			o.AddSubcheck(fanPartial)
		}

		o.Add(worstStatus, fmt.Sprintf("Device Info: Uptime - %v", upTime))
		fmt.Println(o.GetOutput())
	}

	// ports
	if slices.Contains(mode, "ports") {
		var inStats, outStats netgear.PortStatistics
		portsIn, _ := n.PortStatistics("inbound")
		portsOut, _ := n.PortStatistics("outbound")
		_ = json.Unmarshal(portsIn, &inStats)
		_ = json.Unmarshal(portsOut, &outStats)

		inRows := inStats.PortStatistics.Rows
		outRows := outStats.PortStatistics.Rows

		o := result.Overall{}
		portsPartial, portsStatus := checks.CheckPorts(inRows, outRows, portsToCheck, *statsWarn, *statsCrit)
		worstStatus = max(worstStatus, portsStatus)
		o.AddSubcheck(portsPartial)
		fmt.Println(o.GetOutput())
	}

	// poe stuff
	if slices.Contains(mode, "poe") {
		o := result.Overall{}
		var poeStatus netgear.PoeStatus
		inputData, _ := n.PoeStatus()
		_ = json.Unmarshal(inputData, &poeStatus)

		poePartial, poeWorst := checks.CheckPoe(poeStatus.PoePortConfig)
		worstStatus = max(worstStatus, poeWorst)
		o.AddSubcheck(poePartial)
		fmt.Println(o.GetOutput())
	}

	os.Exit(worstStatus)
}
