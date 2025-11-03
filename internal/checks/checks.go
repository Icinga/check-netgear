package checks

import (
	"fmt"

	"main/internal/utils"
	"main/netgear"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
)

// CPU
func CheckCPU(cpuUsage float64, warn float64, crit float64) (result.PartialResult, int) {
	status := utils.StatusByThreshold(cpuUsage, warn, crit)
	partial := result.PartialResult{
		Output: fmt.Sprintf("CPU Usage: %.2f%%", cpuUsage),
	}
	_ = partial.SetState(status)
	partial.Perfdata.Add(&perfdata.Perfdata{Label: "CPU", Value: cpuUsage, Min: 0, Max: 100})
	return partial, status
}

// Memory
func CheckMemory(memUsage float64, warn float64, crit float64) (result.PartialResult, int) {
	status := utils.StatusByThreshold(memUsage, warn, crit)
	partial := result.PartialResult{
		Output: fmt.Sprintf("RAM Usage: %.2f%%", memUsage),
	}
	_ = partial.SetState(status)
	partial.Perfdata.Add(&perfdata.Perfdata{Label: "RAM", Value: memUsage, Min: 0, Max: 100})
	return partial, status
}

// Temperature
func CheckTemperature(sensors []netgear.SensorDetail, warn float64, crit float64) (result.PartialResult, int) {
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
		_ = sub.SetState(status)
		sub.Perfdata.Add(&perfdata.Perfdata{Label: s.Description, Value: s.Temperature, Min: 0})
		partial.AddSubcheck(sub)
	}

	_ = partial.SetState(worst)
	return partial, worst
}

// Fans
func CheckFans(fanName string, fanSpeed float64, warn float64) (result.PartialResult, int) {
	status := check.OK
	if fanSpeed > warn {
		status = check.Warning
	}

	partial := result.PartialResult{Output: "Fans"}
	sub := result.PartialResult{
		Output: fmt.Sprintf("%s: %.0f RPM", fanName, fanSpeed),
	}
	_ = sub.SetState(status)
	sub.Perfdata.Add(&perfdata.Perfdata{Label: "Fan Speed", Value: fanSpeed, Min: 0})
	partial.AddSubcheck(sub)
	_ = partial.SetState(status)

	return partial, status
}

// Ports
func CheckPorts(inRows, outRows []netgear.PortStatisticRow, portsToCheck []int, warn float64, crit float64) (result.PartialResult, int) {
	overall := result.PartialResult{Output: "Ports Statistics"}
	worst := check.OK

	for i := range inRows {
		in := inRows[i]
		out := outRows[i]

		if !contains(portsToCheck, in.Port) {
			continue
		}

		portCheck := result.PartialResult{Output: fmt.Sprintf("Port %v", in.Port)}

		inLoss := utils.LossPercent(in.InDropPkts, in.InTotalPkts)
		outLoss := utils.LossPercent(out.OutDropPkts, out.OutTotalPkts)

		inStatus := utils.StatusByThreshold(inLoss, warn, crit)
		outStatus := utils.StatusByThreshold(outLoss, warn, crit)

		portStatus := max(inStatus, outStatus)
		worst = max(worst, portStatus)

		addPerfSubcheck := func(label string, loss float64, status int) {
			sub := result.PartialResult{
				Output: fmt.Sprintf("%s: %.2f%% loss", label, loss),
			}
			_ = sub.SetState(status)
			sub.Perfdata.Add(&perfdata.Perfdata{
				Label: fmt.Sprintf("port %v %s loss", in.Port, label),
				Value: loss, Min: 0, Max: 100,
			})
			portCheck.AddSubcheck(sub)
		}

		addPerfSubcheck("IN", inLoss, inStatus)
		addPerfSubcheck("OUT", outLoss, outStatus)

		_ = portCheck.SetState(portStatus)
		overall.AddSubcheck(portCheck)
	}

	_ = overall.SetState(worst)
	return overall, worst
}

// Poe
func CheckPoe(ports []netgear.PoePort) (result.PartialResult, int) {
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
		_ = poeCheck.SetState(status)
		poeCheck.Perfdata.Add(&perfdata.Perfdata{
			Label: fmt.Sprintf("port %v power", port.Port),
			Value: port.CurrentPower, Min: 0, Max: port.PowerLimit,
		})
		partial.AddSubcheck(poeCheck)
	}

	_ = partial.SetState(worst)
	return partial, worst
}

// helper
func contains(list []int, x int) bool {
	for _, v := range list {
		if v == x {
			return true
		}
	}
	return false
}
