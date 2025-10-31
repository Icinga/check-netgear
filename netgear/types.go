package netgear

// basic infop

type DeviceInfoDetails struct {
	Uptime string `json:"upTime"`
}

type Fan struct {
	Details []struct {
		Speed       float64 `json:"speed"`
		Description string  `json:"desc"`
	} `json:"details"`
}

type Sensor struct {
	Details []struct {
		Description string  `json:"desc"`
		Temperature float64 `json:"temp"`
		MaxTemp     float64 `json:"maxTemp"`
	} `json:"details"`
}

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

// ports

type PortStatisticRow struct {
	Port        int     `json:"port"`
	InTotalPkts float64 `json:"inTotalPkts"`
	InDropPkts  float64 `json:"inDropPkts"`
	InOctets    float64 `json:"inOctets"`

	OutTotalPkts float64 `json:"outTotalPkts"`
	OutDropPkts  float64 `json:"outDropPkts"`
	OutOctets    float64 `json:"outOctets"`
}

type PortStatistics struct {
	PortStatistics struct {
		Rows []PortStatisticRow `json:"rows"`
	} `json:"portStatistics"`
}

// poe

type PoePort struct {
	Port         string  `json:"port"`
	Enable       bool    `json:"enable"`
	CurrentPower float64 `json:"currentPower"`
	PowerLimit   float64 `json:"powerLimit"`
}

type PoeStatus struct {
	PoePortConfig []PoePort `json:"poePortConfig"`
}
