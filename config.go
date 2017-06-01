package main

import (
	"encoding/json"
	"fmt"
	"github.com/containernetworking/cni/pkg/types"
	"strings"
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
}

// Net loads the options of the CNI network configuration file
type Net struct {
	Name         string      `json:"name"`
	BridgeSubnet string      `json:"bridgeSubnet"`
	IPAM         *IPAMConfig `json:"ipam"`
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

	//If BridgeSubnet is a valid CIDR block set the PrefixSize
	i := strings.Split(n.BridgeSubnet, "/")
	if len(i) > 1 {
		n.IPAM.SubnetPrefixSize = "/" + i[1]
	} else {
		n.IPAM.SubnetPrefixSize = ""
	}

	return n.IPAM, nil
}
