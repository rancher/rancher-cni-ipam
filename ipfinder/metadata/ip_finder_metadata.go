package metadata

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher-metadata/metadata"
)

const (
	metadataURLTemplate = "http://%v/2015-12-19"
	multiplier          = 60
	emptyIPAddress      = ""

	// DefaultMetadataAddress specifies the default value to use if nothing is specified
	DefaultMetadataAddress = "169.254.169.250"
)

// IPFinderFromMetadata is used to hold information related to
// Metadata client and other stuff.
type IPFinderFromMetadata struct {
	m *metadata.Client
}

// NewIPFinderFromMetadata returns a new instance of the IPFinderFromMetadata
func NewIPFinderFromMetadata(metadataAddress string) (*IPFinderFromMetadata, error) {
	if metadataAddress == "" {
		metadataAddress = DefaultMetadataAddress
	}
	metadataURL := fmt.Sprintf(metadataURLTemplate, metadataAddress)
	m, err := metadata.NewClientAndWait(metadataURL)
	if err != nil {
		return nil, err
	}
	return &IPFinderFromMetadata{m}, nil
}

// GetIP returns the IP address for the given container id, return an empty string
// if not found
func (ipf *IPFinderFromMetadata) GetIP(cid, rancherid string) string {
	for i := 0; i < multiplier; i++ {
		containers, err := ipf.m.GetContainers()
		if err != nil {
			logrus.Errorf("rancher-cni-ipam: Error getting metadata containers: %v", err)
			return emptyIPAddress
		}

		for _, container := range containers {
			if container.ExternalId == cid && container.PrimaryIp != "" {
				logrus.Infof("rancher-cni-ipam: got ip: %v", container.PrimaryIp)
				return container.PrimaryIp
			}
			if rancherid != "" && container.UUID == rancherid && container.PrimaryIp != "" {
				logrus.Infof("rancher-cni-ipam: got ip from rancherid: %v", container.PrimaryIp)
				return container.PrimaryIp
			}
		}
		logrus.Infof("Waiting to find IP for container: %s, %s", cid, rancherid)
		time.Sleep(500 * time.Millisecond)
	}
	logrus.Infof("ip not found for cid: %v", cid)
	return emptyIPAddress
}
