package metadata

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher-metadata/metadata"
	"time"
)

const (
	metadataURL = "http://rancher-metadata/latest"
	empty       = ""
)

// IPFinderFromMetadata is used to hold information related to
// Metadata client and other stuff.
type IPFinderFromMetadata struct {
	m *metadata.Client
}

// NewIPFinderFromMetadata returns a new instance of the IPFinderFromMetadata
func NewIPFinderFromMetadata() (*IPFinderFromMetadata, error) {
	m, err := metadata.NewClientAndWait(metadataURL)
	if err != nil {
		return nil, err
	}
	return &IPFinderFromMetadata{m}, nil
}

// GetIP returns the IP address for the given container id, return an empty string
// if not found
func (ipf *IPFinderFromMetadata) GetIP(cid string) string {

	for i := 0; i < 600; i++ {
		containers, err := ipf.m.GetContainers()
		if err != nil {
			log.Errorf("rancher-cni-ipam: Error getting metadata containers: %v", err)
			return empty
		}

		for _, container := range containers {
			if container.ExternalId == cid {
				log.Infof("rancher-cni-ipam: got ip: %v", container.PrimaryIp)
				return container.PrimaryIp
			}
		}
		log.Infof("Waiting to find IP for container: %s", cid)
		time.Sleep(500 * time.Millisecond)
	}
	log.Infof("ip not found for cid: %v", cid)
	return empty
}
