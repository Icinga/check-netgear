package utils

import "github.com/NETWAYS/go-check"

func StatusByThreshold(value, warn, crit float64) int {
	switch {
	case value >= crit:
		return check.Critical
	case value >= warn:
		return check.Warning
	default:
		return check.OK
	}
}

func LossPercent(drop, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return drop / total * 100
}
