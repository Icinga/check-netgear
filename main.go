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

// Flags contains all command line flags that are relevant to check modes
type Flags struct {
	NoPerfdata bool

	HideCpu  bool
	HideMem  bool
	HideTemp bool
	HideFans bool

	CpuWarn  float64
	CpuCrit  float64
	MemWarn  float64
	MemCrit  float64
	TempWarn float64
	TempCrit float64
	FanWarn  float64
	FanCrit  float64
	PortWarn float64
	PortCrit float64

	PortsToCheck intSliceFlag
}

// ModeBasic contains all the basic hardware information of the switch, including CPU and RAM usage, temperature and fan
// speed
func ModeBasic(netgearSession *netgear.Netgear, worstStatus *int, flags *Flags) (*result.PartialResult, error) {
	deviceInfo, err := netgearSession.DeviceInfo()
	if err != nil {
		return nil, fmt.Errorf("error retrieving device info: %v\n", err)
	}

	upTime := deviceInfo.DeviceInfo.Details[0].Uptime
	if len(deviceInfo.DeviceInfo.Details) == 0 {
		return nil, fmt.Errorf("error retrieving device info")
	}

	o := result.PartialResult{
		Output: fmt.Sprintf("Device Info: Uptime - %v", upTime),
	}

	if !flags.HideCpu {
		if len(deviceInfo.DeviceInfo.Cpu) == 0 {
			return nil, fmt.Errorf("no CPU info for this device")
		}
		cpuUsage, err := netgear.StringPercentToFloat(deviceInfo.DeviceInfo.Cpu[0].Usage)
		if err != nil {
			return nil, fmt.Errorf("error parsing CPU usage: %v\n", err)
		}

		cpuPartial, err := checks.CheckCPU(cpuUsage, flags.NoPerfdata, flags.CpuWarn, flags.CpuCrit)
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("CPU check error: %v", err)
			err := errRes.SetState(check.Unknown)
			if err != nil {
				return nil, err
			}
			o.AddSubcheck(errRes)
		} else {
			*worstStatus = max(*worstStatus, cpuPartial.GetStatus())
			o.AddSubcheck(*cpuPartial)
		}
	}

	if !flags.HideMem {
		if len(deviceInfo.DeviceInfo.Memory) == 0 {
			return nil, fmt.Errorf("no Memory info for this device")
		}
		memUsage, err := netgear.StringPercentToFloat(deviceInfo.DeviceInfo.Memory[0].Usage)
		if err != nil {
			return nil, fmt.Errorf("error parsing Memory usage: %v\n", err)
		}

		memPartial, err := checks.CheckMemory(memUsage, flags.NoPerfdata, flags.MemWarn, flags.MemCrit)
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("Memory check error: %v", err)
			err := errRes.SetState(check.Unknown)
			if err != nil {
				return nil, err
			}
			o.AddSubcheck(errRes)
		} else {
			*worstStatus = max(*worstStatus, memPartial.GetStatus())
			o.AddSubcheck(*memPartial)
		}
	}

	if !flags.HideTemp {
		if len(deviceInfo.DeviceInfo.Sensor) == 0 {
			return nil, fmt.Errorf("no Temperature info for this device")
		}
		sensorDetails := deviceInfo.DeviceInfo.Sensor[0].Details
		tempPartial, err := checks.CheckTemperature(sensorDetails, flags.NoPerfdata, flags.TempWarn, flags.TempCrit)
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("Temperature check error: %v", err)
			err := errRes.SetState(check.Unknown)
			if err != nil {
				return nil, err
			}
			o.AddSubcheck(errRes)
		} else {
			*worstStatus = max(*worstStatus, tempPartial.GetStatus())
			o.AddSubcheck(*tempPartial)
		}
	}

	if !flags.HideFans {
		if len(deviceInfo.DeviceInfo.Fan) == 0 {
			return nil, fmt.Errorf("no Fan info for this device")
		}
		if len(deviceInfo.DeviceInfo.Fan[0].Details) == 0 {
			return nil, fmt.Errorf("no Fan details for this device")
		}
		fan := deviceInfo.DeviceInfo.Fan[0].Details[0]
		fanPartial, err := checks.CheckFans(flags.NoPerfdata, fan.Description, fan.Speed, flags.FanWarn, flags.FanCrit)
		if err != nil {
			errRes := result.NewPartialResult()
			errRes.Output = fmt.Sprintf("Fans check error: %v", err)
			err := errRes.SetState(check.Unknown)
			if err != nil {
				return nil, err
			}
			o.AddSubcheck(errRes)
		} else {
			*worstStatus = max(*worstStatus, fanPartial.GetStatus())
			o.AddSubcheck(*fanPartial)
		}
	}

	return &o, nil
}

// ModePorts monitors the network traffic on the ports and reports back the percentage of dropped packets
func ModePorts(netgearSession *netgear.Netgear, worstStatus *int, flags *Flags) (*result.PartialResult, error) {
	o := result.PartialResult{}
	portsIn, err := netgearSession.PortStatistics("inbound")
	if err != nil {
		errRes := result.NewPartialResult()
		errRes.Output = fmt.Sprintf("Inbound port check error: %v", err)
		err := errRes.SetState(check.Unknown)
		if err != nil {
			return nil, err
		}
		o.AddSubcheck(errRes)
		return &o, nil
	}
	portsOut, err := netgearSession.PortStatistics("outbound")
	if err != nil {
		errRes := result.NewPartialResult()
		errRes.Output = fmt.Sprintf("Outbound port check error: %v", err)
		err := errRes.SetState(check.Unknown)
		if err != nil {
			return nil, err
		}
		o.AddSubcheck(errRes)
		return &o, nil
	}

	inRows := portsIn.PortStatistics.Rows
	outRows := portsOut.PortStatistics.Rows

	portsPartial, err := checks.CheckPorts(inRows, outRows, flags.PortsToCheck, flags.NoPerfdata, flags.PortWarn, flags.PortCrit)
	if err != nil {
		errRes := result.NewPartialResult()
		errRes.Output = fmt.Sprintf("Ports check error: %v", err)
		err := errRes.SetState(check.Unknown)
		if err != nil {
			return nil, err
		}
		o.AddSubcheck(errRes)
		return &o, nil
	} else {
		*worstStatus = max(*worstStatus, portsPartial.GetStatus())
		o = *portsPartial
	}

	return &o, nil
}

// ModePoE checks the ports PoE state
func ModePoE(netgearSession *netgear.Netgear, worstStatus *int, flags *Flags) (*result.PartialResult, error) {
	o := result.PartialResult{}
	poeStatus, err := netgearSession.PoeStatus()
	if err != nil {
		errRes := result.NewPartialResult()
		errRes.Output = fmt.Sprintf("PoE check error: %v", err)
		err := errRes.SetState(check.Unknown)
		if err != nil {
			return nil, err
		}
		o.AddSubcheck(errRes)
		return &o, nil
	}

	poePartial, err := checks.CheckPoe(poeStatus.PoePortConfig, flags.NoPerfdata)
	if err != nil {
		errRes := result.NewPartialResult()
		errRes.Output = fmt.Sprintf("PoE check error: %v", err)
		err := errRes.SetState(check.Unknown)
		if err != nil {
			return nil, err
		}
		o.AddSubcheck(errRes)
		return &o, nil
	} else {
		*worstStatus = max(*worstStatus, poePartial.GetStatus())
		o = *poePartial
	}

	return &o, nil
}

func main() {
	flags := Flags{}

	flag.BoolVar(&flags.NoPerfdata, "noperfdata", false, "Do not output performance data")

	flag.BoolVar(&flags.HideCpu, "nocpu", false, "Hide the CPU info")
	flag.BoolVar(&flags.HideMem, "nomem", false, "Hide the RAM info")
	flag.BoolVar(&flags.HideTemp, "notemp", false, "Hide the Temperature info")
	flag.BoolVar(&flags.HideFans, "nofans", false, "Hide the Fans info")

	mode := stringSliceFlag{}
	flag.Var(&mode, "mode", "Output modes to enable {basic|ports|poe|all} (repeatable)")

	baseURL := flag.String("base-url", "http://192.168.0.239", "Base URL to use")

	username := flag.String("username", "", "Username for authentication")
	password := flag.String("password", "", "Password for authentication")

	// Thresholds
	flag.Float64Var(&flags.CpuWarn, "cpu-warning", 50, "CPU usage warning threshold")
	flag.Float64Var(&flags.CpuCrit, "cpu-critical", 90, "CPU usage critical threshold")
	flag.Float64Var(&flags.MemWarn, "mem-warning", 50, "RAM usage warning threshold")
	flag.Float64Var(&flags.MemCrit, "mem-critical", 90, "RAM usage critical threshold")
	flag.Float64Var(&flags.FanWarn, "fan-warning", 3000, "Fan speed warning threshold")
	flag.Float64Var(&flags.FanCrit, "fan-critical", 5000, "Fan speed critical threshold")
	flag.Float64Var(&flags.TempWarn, "temp-warning", 50, "Temperature warning threshold")
	flag.Float64Var(&flags.TempCrit, "temp-critical", 70, "Temperature critical threshold")
	flag.Float64Var(&flags.PortWarn, "stats-warning", 5, "Port stats warning threshold")
	flag.Float64Var(&flags.PortCrit, "stats-critical", 20, "Port stats critical threshold")

	flags.PortsToCheck = intSliceFlag{1, 2, 3, 4, 5, 6, 7, 8}
	flag.Var(&flags.PortsToCheck, "port", "Ports to check (repeatable)")

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

	netgearSession, err := netgear.NewNetgear(*baseURL, *username, *password)
	if err != nil {
		fmt.Printf("URL error: %v", err)
		os.Exit(check.Unknown)
	}
	if err := netgearSession.Login(); err != nil {
		fmt.Printf("Error while trying to login: %v\n", err)
		os.Exit(check.Unknown)
	}
	defer func() { _ = netgearSession.Logout() }()

	if len(mode) == 0 {
		mode = append(mode, "basic")
	} else if slices.Contains(mode, "all") {
		mode = append(mode, "basic", "ports", "poe")
	}

	worstStatus := check.OK
	o := result.Overall{}

	// Basic check
	if slices.Contains(mode, "basic") {
		subcheck, err := ModeBasic(netgearSession, &worstStatus, &flags)
		if err != nil {
			fmt.Print(err)
			os.Exit(check.Unknown)
		}
		o.AddSubcheck(*subcheck)
	}

	// ports
	if slices.Contains(mode, "ports") {
		subcheck, err := ModePorts(netgearSession, &worstStatus, &flags)
		if err != nil {
			fmt.Print(err)
			os.Exit(check.Unknown)
		}
		o.AddSubcheck(*subcheck)
	}

	// poe stuff
	if slices.Contains(mode, "poe") {
		subcheck, err := ModePoE(netgearSession, &worstStatus, &flags)
		if err != nil {
			fmt.Print(err)
			os.Exit(check.Unknown)
		}
		o.AddSubcheck(*subcheck)
	}

	fmt.Print(o.GetOutput())

	os.Exit(worstStatus)
}
