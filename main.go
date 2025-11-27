package main

import (
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/icinga/check-netgear/netgear"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

// ModePorts monitors the network traffic on the ports and reports back the percentage of dropped packets
func ModePorts(netgearSession *netgear.Netgear, flags *netgear.Flags) (*result.PartialResult, error) {
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

	portsPartial, err := netgear.CheckPorts(inRows, outRows, flags.PortsToCheck, flags.NoPerfdata, flags.PortWarn, flags.PortCrit)
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
		o = *portsPartial
	}

	return &o, nil
}

// ModePoE checks the ports PoE state
func ModePoE(netgearSession *netgear.Netgear, flags *netgear.Flags) (*result.PartialResult, error) {
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

	poePartial, err := netgear.CheckPoe(poeStatus.PoePortConfig, flags.NoPerfdata)
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
		o = *poePartial
	}

	return &o, nil
}

func main() {
	flags := netgear.Flags{}

	flag.BoolVar(&flags.NoPerfdata, "noperfdata", false, "Do not output performance data")

	flag.BoolVar(&flags.HideCpu, "nocpu", false, "Hide the CPU info")
	flag.BoolVar(&flags.HideMem, "nomem", false, "Hide the RAM info")
	flag.BoolVar(&flags.HideTemp, "notemp", false, "Hide the Temperature info")
	flag.BoolVar(&flags.HideFans, "nofans", false, "Hide the Fans info")

	mode := netgear.StringSliceFlag{}
	flag.Var(&mode, "mode", "Output modes to enable {basic|ports|poe|all} (repeatable) (default: basic)")

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

	flags.PortsToCheck = netgear.IntSliceFlag{1, 2, 3, 4, 5, 6, 7, 8}
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
		subcheck, err := netgearSession.ModeBasic(&flags)
		if err != nil {
			fmt.Print(err)
			os.Exit(check.Unknown)
		}
		o.AddSubcheck(*subcheck)
		worstStatus = result.WorstState(worstStatus, subcheck.GetStatus())
	}

	// ports
	if slices.Contains(mode, "ports") {
		subcheck, err := ModePorts(netgearSession, &flags)
		if err != nil {
			fmt.Print(err)
			os.Exit(check.Unknown)
		}
		o.AddSubcheck(*subcheck)
		worstStatus = result.WorstState(worstStatus, subcheck.GetStatus())
	}

	// poe stuff
	if slices.Contains(mode, "poe") {
		subcheck, err := ModePoE(netgearSession, &flags)
		if err != nil {
			fmt.Print(err)
			os.Exit(check.Unknown)
		}
		o.AddSubcheck(*subcheck)
		worstStatus = result.WorstState(worstStatus, subcheck.GetStatus())
	}

	if len(o.PartialResults) == 0 {
		fmt.Print("No valid modes selected")
		os.Exit(check.Unknown)
	}

	fmt.Print(o.GetOutput())

	os.Exit(worstStatus)
}
