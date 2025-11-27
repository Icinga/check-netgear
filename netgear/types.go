package netgear

import (
	"strconv"
	"strings"
)

// DeviceInfoDetails represents individual uptime details returned by the deviceInfo endpoint
type DeviceInfoDetails struct {
	Uptime string `json:"upTime"`
}

// FanDetail represents information about an individual fan entry in the device
type FanDetail struct {
	Speed       float64 `json:"speed"`
	Description string  `json:"desc"`
}

// Fan contains an array of all fan detail entries
type Fan struct {
	Details []FanDetail `json:"details"`
}

// SensorDetail represents a single thermal sensor reading
type SensorDetail struct {
	Description string  `json:"desc"`
	Temperature float64 `json:"temp"`
	MaxTemp     float64 `json:"maxTemp"`
}

// Sensor contains an array of sensor detail entries
type Sensor struct {
	Details []SensorDetail `json:"details"`
}

// DeviceInfo contains high level device information
type DeviceInfo struct {
	DeviceInfo struct {
		Details []DeviceInfoDetails `json:"details"`
		Fan     []Fan               `json:"fan"`
		Sensor  []Sensor            `json:"sensor"`
		Cpu     []struct {
			Usage string `json:"usage"`
			Unit  int32  `json:"unit"`
		} `json:"cpu"`
		Memory []struct {
			Usage string `json:"usage"`
			Unit  int32  `json:"unit"`
		} `json:"memory"`
	} `json:"deviceInfo"`
}

// PortStatisticRow contains measured per port traffic
type PortStatisticRow struct {
	Port        int     `json:"port"`
	InTotalPkts float64 `json:"inTotalPkts"`
	InDropPkts  float64 `json:"inDropPkts"`
	InOctets    float64 `json:"inOctets"`

	OutTotalPkts float64 `json:"outTotalPkts"`
	OutDropPkts  float64 `json:"outDropPkts"`
	OutOctets    float64 `json:"outOctets"`
}

// PortStatistics contains traffic statistics for all ports
type PortStatistics struct {
	PortStatistics struct {
		Rows []PortStatisticRow `json:"rows"`
	} `json:"portStatistics"`
}

// PoePort represents power information for a single PoE enabled port
type PoePort struct {
	Port         string  `json:"port"`
	Enable       bool    `json:"enable"`
	CurrentPower float64 `json:"currentPower"`
	PowerLimit   float64 `json:"powerLimit"`
}

// PoeStatus represents PoE configuration for all PoE capable ports
type PoeStatus struct {
	PoePortConfig []PoePort `json:"poePortConfig"`
}

// So that flag supports slices
type StringSliceFlag []string

func (s *StringSliceFlag) String() string { return strings.Join(*s, ",") }

func (s *StringSliceFlag) Set(v string) error {
	*s = append(*s, v)
	return nil
}

type IntSliceFlag []int

func (i *IntSliceFlag) String() string {
	parts := make([]string, 0, len(*i))
	for _, v := range *i {
		parts = append(parts, strconv.Itoa(v))
	}
	return strings.Join(parts, ",")
}

func (i *IntSliceFlag) Set(v string) error {
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

	PortsToCheck IntSliceFlag
}
