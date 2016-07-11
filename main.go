package main

import (
	"fmt"
	"net"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/rancher/rancher-cni-ipam/ipfinder/metadata"
)

func cmdAdd(args *skel.CmdArgs) error {
	ipamConf, err := LoadIPAMConfig(args.StdinData, args.Args)
	if err != nil {
		return err
	}

	if ipamConf.LogToFile != "" {
		f, err := os.OpenFile(ipamConf.LogToFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil && f != nil {
			log.SetOutput(f)
			defer f.Close()
		}
	}

	log.Debugf("rancher-cni-ipam: cmdAdd: invoked")
	log.Debugf("rancher-cni-ipam: %s", fmt.Sprintf("args: %#v", args))
	log.Debugf("rancher-cni-ipam: %s", fmt.Sprintf("ipamConf: %#v", ipamConf))

	ipf, err := metadata.NewIPFinderFromMetadata()
	if err != nil {
		return err
	}
	ipString := ipf.GetIP(args.ContainerID)

	log.Debugf("rancher-cni-ipam: %s", fmt.Sprintf("ip: %#v", ipString))

	ip, ipnet, err := net.ParseCIDR(ipString + "/16")
	if err != nil {
		return err
	}

	// TODO: if ip is NULL, return err
	r := &types.Result{
		IP4: &types.IPConfig{
			IP: net.IPNet{IP: ip, Mask: ipnet.Mask},
		},
	}

	log.Infof("rancher-cni-ipam: %s", fmt.Sprintf("r: %#v", r))
	return r.Print()
}

func cmdDel(args *skel.CmdArgs) error {
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel)
}
