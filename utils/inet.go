package utils

import (
	"fmt"
	"net"
	"strings"
	"github.com/vishvananda/netlink"
)

//GetRouteAndInterface returns the src and gateway to reach dst through the routing table. Error is set if impossible
func GetRouteAndInterface (dst string)(string, string, error){
	routes, err := netlink.RouteGet(net.ParseIP(dst))
	if err != nil{
		return "","",err
	}
	for _, r := range routes{
		if strings.Split(r.Dst.String(), "/")[0] == dst{
			return r.Src.String(), r.Gw.String(),nil
		}
	}
	return "","",fmt.Errorf("%s not found", dst)
}


//GetIPFromIfName returns the unicast global IPv4 address from an interface name
func GetIPFromIfName(name string)(string,error){
	ifaces, err := net.Interfaces()
	if err != nil {
		return "",fmt.Errorf("cannot get list of interfaces: %s", err)
	}
	for _, i := range ifaces {
		if name != i.Name{
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			return "", fmt.Errorf("unexpected error on interfaces details: %s", err)
		}
		for _, a := range addrs {
			if (net.ParseIP(strings.Split(a.String(), "/")[0]).IsGlobalUnicast()){
				return strings.Split(a.String(), "/")[0], nil
			}
		}
	}
	return "",fmt.Errorf("interface not found")
}

//GetNameFromIfIP returns the name of the interface from the unicast global IPv4 address or an error
func GetNameFromIfIP(addr string)(string,error){
	if net.ParseIP(addr) == nil{
		return "",fmt.Errorf("wrong ip address format: %s", addr)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return "",fmt.Errorf("cannot get list of interfaces: %s", err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", fmt.Errorf("unexpected error on interfaces details: %s", err)
		}
		for _, a := range addrs {
			if strings.Split(a.String(), "/")[0] == addr{
				return i.Name, nil
			}
		}
	}
	return "", fmt.Errorf("address %s not found on this system", addr)
}