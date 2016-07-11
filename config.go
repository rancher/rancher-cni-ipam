package main

import (
	"encoding/json"
	"fmt"
)

// IPAMConfig is used to load the options specified in the configuration file
type IPAMConfig struct {
	Type      string `json:"type"`
	LogToFile string `json:"logToFile"`
}

// Net loads the options of the CNI network configuration file
type Net struct {
	Name string      `json:"name"`
	IPAM *IPAMConfig `json:"ipam"`
}

// LoadIPAMConfig loads the IPAM configuration from the given bytes
func LoadIPAMConfig(bytes []byte, args string) (*IPAMConfig, error) {
	n := Net{}
	if err := json.Unmarshal(bytes, &n); err != nil {
		return nil, fmt.Errorf("failed to load netconf: %v", err)
	}

	if n.IPAM == nil {
		return nil, fmt.Errorf("IPAM config missing 'ipam' key")
	}

	return n.IPAM, nil
}
