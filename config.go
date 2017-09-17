package main

import (
	"encoding/json"
	"fmt"

	"github.com/containernetworking/cni/pkg/types"
)

// IPAMConfig is used to load the options specified in the configuration file
type IPAMConfig struct {
	types.CommonArgs
	Type                 string        `json:"type"`
	LogToFile            string        `json:"logToFile"`
	IsDebugLevel         string        `json:"isDebugLevel"`
	SubnetPrefixSize     string        `json:"subnetPrefixSize"`
	Routes               []types.Route `json:"routes"`
	RancherContainerUUID types.UnmarshallableString
	IPAddress            types.UnmarshallableString
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

	if err := types.LoadArgs(args, n.IPAM); err != nil {
		return nil, fmt.Errorf("failed to parse args %s: %v", args, err)
	}

	return n.IPAM, nil
}
