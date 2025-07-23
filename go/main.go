package main

import (
	"encoding/json"
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
)

func main() {

	err := login("admin", "password")
	if err != nil {
		return
	}

	var data map[string]any
	inputData := device_info()
	err = json.Unmarshal(inputData, &data)
	if err != nil {
		panic(err)
	}

	upTime := data["deviceInfo"].(map[string]any)["details"].([]any)[0].(map[string]any)["upTime"]
	cpuUsage := string_percent_to_float(data["deviceInfo"].(map[string]any)["cpu"].([]any)[0].(map[string]any)["usage"].(string))
	memoryUsage := string_percent_to_float(data["deviceInfo"].(map[string]any)["memory"].([]any)[0].(map[string]any)["usage"].(string))
	fanName := data["deviceInfo"].(map[string]any)["fan"].([]any)[0].(map[string]any)["details"].([]any)[0].(map[string]any)["desc"]
	fanSpeed := data["deviceInfo"].(map[string]any)["fan"].([]any)[0].(map[string]any)["details"].([]any)[0].(map[string]any)["speed"]

	sensorData := data["deviceInfo"].(map[string]any)["sensor"].([]any)[0].(map[string]any)["details"].([]any) // For the temperatuer
	temperature1 := sensorData[0].(map[string]any)["temp"].(float64)
	temperature1Name := sensorData[0].(map[string]any)["desc"].(string)
	temperature2 := sensorData[1].(map[string]any)["temp"].(float64)
	temperature2Name := sensorData[1].(map[string]any)["desc"].(string)
	temperature3 := sensorData[2].(map[string]any)["temp"].(float64)
	temperature3Name := sensorData[2].(map[string]any)["desc"].(string)

	worstStatusCode := check.OK

	if cpuUsage >= 90 || memoryUsage >= 90 || temperature1 > 80 || temperature2 > 80 || temperature3 > 80 {
		worstStatusCode = check.Critical
	} else if cpuUsage >= 70 || memoryUsage >= 75 || temperature1 > 70 || temperature2 > 70 || temperature3 > 70 {
		worstStatusCode = check.Warning
	}

	if fanSpeed == 0 {
		worstStatusCode = check.Warning
	}

	// Generating the graph thingie

	o := result.Overall{}
	o.Add(worstStatusCode, "General Device Info")

	o.Add(worstStatusCode, fmt.Sprintf("Uptime - %v", upTime))

	//CPU Checks
	cpuCheck := result.PartialResult{
		Output: "CPU",
	}
	cpuCheck.SetState(check.OK)

	cpuSubCheck1 := result.PartialResult{
		Output: fmt.Sprintf("Usage: %v", cpuUsage),
	}
	cpuSubCheck1.SetState(get_cpu_usage_check_level(cpuUsage))
	cpuSubCheck1.Perfdata.Add(&perfdata.Perfdata{
		Label: "CPU",
		Value: cpuUsage,
		Min:   0,
		Max:   100,
	})
	cpuCheck.AddSubcheck(cpuSubCheck1)

	o.AddSubcheck(cpuCheck)

	//Memory Checks
	memoryCheck := result.PartialResult{
		Output: "RAM",
	}
	memoryCheck.SetState(check.OK)

	memorySubCheck1 := result.PartialResult{
		Output: fmt.Sprintf("Usage: %v", memoryUsage),
	}
	memorySubCheck1.SetState(get_cpu_usage_check_level(memoryUsage))
	memorySubCheck1.Perfdata.Add(&perfdata.Perfdata{
		Label: "RAM",
		Value: memoryUsage,
		Min:   0,
		Max:   100,
	})
	memoryCheck.AddSubcheck(memorySubCheck1)

	o.AddSubcheck(memoryCheck)

	// Temperature Checks
	temperatureCheck := result.PartialResult{
		Output: "Temperature",
	}
	temperatureCheck.SetState(check.Warning)

	tempSubCheck1 := result.PartialResult{
		Output: fmt.Sprintf("%v - %v", temperature1Name, temperature1),
	}
	tempSubCheck1.SetState(check.OK)
	tempSubCheck1.Perfdata.Add(&perfdata.Perfdata{
		Label: temperature1Name,
		Value: temperature1,
		Min:   0,
	})
	temperatureCheck.AddSubcheck(tempSubCheck1)

	tempSubCheck2 := result.PartialResult{
		Output: fmt.Sprintf("%v - %v", temperature2Name, temperature2),
	}
	tempSubCheck2.SetState(check.Warning)
	tempSubCheck2.Perfdata.Add(&perfdata.Perfdata{
		Label: temperature2Name,
		Value: temperature2,
		Min:   0,
	})
	temperatureCheck.AddSubcheck(tempSubCheck2)

	tempSubCheck3 := result.PartialResult{
		Output: fmt.Sprintf("%v - %v", temperature3Name, temperature3),
	}
	tempSubCheck3.SetState(check.OK)
	tempSubCheck3.Perfdata.Add(&perfdata.Perfdata{
		Label: temperature3Name,
		Value: temperature3,
		Min:   0,
	})
	temperatureCheck.AddSubcheck(tempSubCheck3)

	o.AddSubcheck(temperatureCheck)

	// Fans Checks
	fansCheck := result.PartialResult{
		Output: "Fans",
	}
	fansCheck.SetState(check.Warning)

	fansSubCheck1 := result.PartialResult{
		Output: fmt.Sprintf("%v: %v RPM", fanName, fanSpeed),
	}
	fansSubCheck1.SetState(check.Warning)
	fansSubCheck1.Perfdata.Add(&perfdata.Perfdata{
		Label: "Fans speed",
		Value: fanSpeed,
		Min:   0,
	})
	fansCheck.AddSubcheck(fansSubCheck1)

	o.AddSubcheck(fansCheck)

	// Result
	fmt.Println(o.GetOutput())

	// login("admin", "password")
	// check_everything()

	// config := check.NewConfig()
	// config.Name = "check_test"
	// config.Readme = `Test Plugin`
	// config.Version = "1.0.0"
	// _ = config.FlagSet.StringP("hostname", "H", "localhost", "Hostname to check")
	// config.ParseArguments()
	// Some checking should be done here, when --help is not passed
	// check.Exitf(check.OK, "Everything is fine - answer=%d", 42)
	// Output:
	// OK - Everything is fine - answer=42
}
