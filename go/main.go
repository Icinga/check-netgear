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

	sensorData := data["deviceInfo"].(map[string]any)["sensor"].([]any)[0].(map[string]any)["details"].([]any) // For the temperature
	var temperaturesNames []string
	for _, sensor := range sensorData {
		temperaturesNames = append(temperaturesNames, sensor.(map[string]any)["desc"].(string))
	}
	var temperatures []float64
	for _, sensor := range sensorData {
		temperatures = append(temperatures, sensor.(map[string]any)["temp"].(float64))
	}

	worstStatusCode := check.OK

	if fanSpeed == 0 {
		worstStatusCode = check.Warning
	}

	// Generating the graph thingie

	o := result.Overall{}
	o.Add(worstStatusCode, fmt.Sprintf("Uptime - %v", upTime))

	//CPU Checks
	cpuCheck := result.PartialResult{
		Output: fmt.Sprintf("CPU Usage: %v%%", cpuUsage),
	}
	cpuCheck.SetState(check.OK)
	cpuCheck.Perfdata.Add(&perfdata.Perfdata{
		Label: "CPU",
		Value: cpuUsage,
		Min:   0,
		Max:   100,
	})

	o.AddSubcheck(cpuCheck)

	//Memory Checks
	memoryCheck := result.PartialResult{
		Output: fmt.Sprintf("RAM Usage: %v%%", memoryUsage),
	}
	memoryCheck.SetState(check.OK)
	memoryCheck.Perfdata.Add(&perfdata.Perfdata{
		Label: "RAM",
		Value: memoryUsage,
		Min:   0,
		Max:   100,
	})

	o.AddSubcheck(memoryCheck)

	// Temperature Checks
	temperatureCheck := result.PartialResult{
		Output: "Temperature",
	}
	temperatureCheck.SetState(check.Warning)

	for index, _ := range temperatures {
		tempSubCheck := result.PartialResult{
			Output: fmt.Sprintf("%v: %v°С", temperaturesNames[index], temperatures[index]),
		}
		tempSubCheck.SetState(check.OK)
		tempSubCheck.Perfdata.Add(&perfdata.Perfdata{
			Label: temperaturesNames[index],
			Value: temperatures[index],
			Min:   0,
		})
		temperatureCheck.AddSubcheck(tempSubCheck)
	}

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
