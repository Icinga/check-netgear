package netgear

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
			Unit  string `json:"unit"`
		} `json:"cpu"`
		Memory []struct {
			Usage string `json:"usage"`
			Unit  string `json:"unit"`
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
