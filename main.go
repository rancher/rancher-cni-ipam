package main

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/rancher/rancher-cni-ipam/ipfinder/metadata"
)

const (
	defaultPrefixSize = "/16"
)

func cmdAdd(args *skel.CmdArgs) error {
	ipamConf, err := LoadIPAMConfig(args.StdinData, args.Args)
	if err != nil {
		return err
	}

	if ipamConf.IsDebugLevel == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if ipamConf.LogToFile != "" {
		f, err := os.OpenFile(ipamConf.LogToFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil && f != nil {
			logrus.SetOutput(f)
			defer f.Close()
		}
	}

	logrus.Debugf("rancher-cni-ipam: cmdAdd: invoked")
	logrus.Debugf("rancher-cni-ipam: %s", fmt.Sprintf("args: %#v", args))
	logrus.Debugf("rancher-cni-ipam: %s", fmt.Sprintf("ipamConf: %#v", ipamConf))
	logrus.Debugf("rancher-cni-ipam: rancher UUID: %s", ipamConf.RancherContainerUUID)

	ipf, err := metadata.NewIPFinderFromMetadata()
	if err != nil {
		return err
	}
	ipString := ipf.GetIP(args.ContainerID, string(ipamConf.RancherContainerUUID))
	if ipString == "" {
		return errors.New("No IP address found")
	}

	logrus.Debugf("rancher-cni-ipam: %s", fmt.Sprintf("ip: %#v", ipString))

	var prefixSize string

	if ipamConf.SubnetPrefixSize != "" {
		prefixSize = ipamConf.SubnetPrefixSize
	} else {
		prefixSize = defaultPrefixSize
	}

	ip, ipnet, err := net.ParseCIDR(ipString + prefixSize)
	if err != nil {
		return err
	}

	r := &types.Result{
		IP4: &types.IPConfig{
			IP: net.IPNet{IP: ip, Mask: ipnet.Mask},
		},
	}

	r.IP4.Routes = append(
		ipamConf.Routes,
	)

	logrus.Infof("rancher-cni-ipam: %s", fmt.Sprintf("r: %#v", r))
	return r.Print()
}

func cmdDel(args *skel.CmdArgs) error {
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.PluginSupports("0.1.0"))
}
