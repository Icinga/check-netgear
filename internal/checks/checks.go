package checks

import (
	"fmt"
	"slices"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
	"github.com/icinga/check-netgear/internal/utils"
	"github.com/icinga/check-netgear/netgear"
)

// CheckCPU creates a partialResult with the CPU information
func CheckCPU(cpuUsage float64, noPerfdata bool, warn float64, crit float64) (*result.PartialResult, error) {
	status := utils.StatusByThreshold(cpuUsage, warn, crit)
	partial := result.PartialResult{
		Output: fmt.Sprintf("CPU Usage: %.2f%%", cpuUsage),
	}
	if err := partial.SetState(status); err != nil {
		return nil, err
	}
	if !noPerfdata {
		partial.Perfdata.Add(&perfdata.Perfdata{Label: "CPU", Value: cpuUsage, Min: 0, Max: 100})
	}
	return &partial, nil
}

// CheckMemory creates a partialResult with the memory information
func CheckMemory(memUsage float64, noPerfdata bool, warn float64, crit float64) (*result.PartialResult, error) {
	status := utils.StatusByThreshold(memUsage, warn, crit)
	partial := result.PartialResult{
		Output: fmt.Sprintf("RAM Usage: %.2f%%", memUsage),
	}
	if err := partial.SetState(status); err != nil {
		return nil, err
	}
	if !noPerfdata {
		partial.Perfdata.Add(&perfdata.Perfdata{Label: "RAM", Value: memUsage, Min: 0, Max: 100})
	}
	return &partial, nil
}

// CheckTemperature creates a partialResult with the temperature information
func CheckTemperature(sensors []netgear.SensorDetail, noPerfdata bool, warn float64, crit float64) (*result.PartialResult, error) {
	partial := result.PartialResult{Output: "Temperature"}
	worst := check.OK

	for _, s := range sensors {
		status := utils.StatusByThreshold(s.Temperature, warn, crit)
		if status > worst {
			worst = status
		}

		sub := result.PartialResult{
			Output: fmt.Sprintf("%s: %.1f°C", s.Description, s.Temperature),
		}
		if err := sub.SetState(status); err != nil {
			return nil, err
		}
		if !noPerfdata {
			sub.Perfdata.Add(&perfdata.Perfdata{Label: s.Description, Value: s.Temperature, Min: 0})
		}
		partial.AddSubcheck(sub)
	}

	if err := partial.SetState(worst); err != nil {
		return nil, err
	}
	return &partial, nil
}

// CheckFans creates a partialResult with the fans information
func CheckFans(noPerfdata bool, fanName string, fanSpeed float64, warn float64, crit float64) (*result.PartialResult, error) {
	status := utils.StatusByThreshold(fanSpeed, warn, crit)

	partial := result.PartialResult{Output: "Fans"}
	sub := result.PartialResult{
		Output: fmt.Sprintf("%s: %.0f RPM", fanName, fanSpeed),
	}
	if err := sub.SetState(status); err != nil {
		return nil, err
	}
	if !noPerfdata {
		sub.Perfdata.Add(&perfdata.Perfdata{Label: "Fan Speed", Value: fanSpeed, Min: 0})
	}
	partial.AddSubcheck(sub)
	if err := partial.SetState(status); err != nil {
		return nil, err
	}

	return &partial, nil
}

// CheckPorts creates a partialResult with the port information
func CheckPorts(inRows, outRows []netgear.PortStatisticRow, portsToCheck []int, noPerfdata bool, warn float64, crit float64) (*result.PartialResult, error) {
	overall := result.PartialResult{Output: "Ports Statistics"}
	worst := check.OK

	for i := range inRows {
		in := inRows[i]
		out := outRows[i]

		if !slices.Contains(portsToCheck, in.Port) {
			continue
		}

		portCheck := result.PartialResult{Output: fmt.Sprintf("Port %v", in.Port)}

		inLoss := utils.LossPercent(in.InDropPkts, in.InTotalPkts)
		outLoss := utils.LossPercent(out.OutDropPkts, out.OutTotalPkts)

		inStatus := utils.StatusByThreshold(inLoss, warn, crit)
		outStatus := utils.StatusByThreshold(outLoss, warn, crit)

		portStatus := max(inStatus, outStatus)
		worst = max(worst, portStatus)

		addPerfSubcheck := func(label string, loss float64, status int) error {
			sub := result.PartialResult{
				Output: fmt.Sprintf("%s: %.2f%% loss", label, loss),
			}
			if err := sub.SetState(status); err != nil {
				return err
			}
			if !noPerfdata {
				sub.Perfdata.Add(&perfdata.Perfdata{
					Label: fmt.Sprintf("port %v %s loss", in.Port, label),
					Value: loss, Min: 0, Max: 100,
				})
			}
			portCheck.AddSubcheck(sub)
			return nil
		}

		if err := addPerfSubcheck("IN", inLoss, inStatus); err != nil {
			return nil, err
		}
		if err := addPerfSubcheck("OUT", outLoss, outStatus); err != nil {
			return nil, err
		}

		if err := portCheck.SetState(portStatus); err != nil {
			return nil, err
		}
		overall.AddSubcheck(portCheck)
	}

	if err := overall.SetState(worst); err != nil {
		return nil, err
	}
	return &overall, nil
}

// CheckPoe creates a partialResult with information about every port's POE status
func CheckPoe(ports []netgear.PoePort, noPerfdata bool) (*result.PartialResult, error) {
	partial := result.PartialResult{Output: "Power over Ethernet Statistics"}
	worst := check.OK

	for _, port := range ports {
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

		worst = max(worst, status)

		poeCheck := result.PartialResult{
			Output: fmt.Sprintf(
				"Port %v is %v. Current power: %.2f/%.2fV",
				port.Port, state, port.CurrentPower/1000, port.PowerLimit/1000,
			),
		}
		if err := poeCheck.SetState(status); err != nil {
			return nil, err
		}
		if !noPerfdata {
			poeCheck.Perfdata.Add(&perfdata.Perfdata{
				Label: fmt.Sprintf("port %v power", port.Port),
				Value: port.CurrentPower, Min: 0, Max: port.PowerLimit,
			})
		}
		partial.AddSubcheck(poeCheck)
	}

	if err := partial.SetState(worst); err != nil {
		return nil, err
	}
	return &partial, nil
}
