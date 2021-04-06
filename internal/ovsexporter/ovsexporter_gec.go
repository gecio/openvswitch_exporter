package ovsexporter

import (
	"github.com/gecio/go-openvswitch/ovs"
	"github.com/digitalocean/go-openvswitch/ovsnl"
	"github.com/prometheus/client_golang/prometheus"
)

// New creates a new Prometheus collector which collects metrics using the
// input Open vSwitch generic netlink client.
func NewWithPortsCollector(c *ovsnl.Client, ovsc *ovs.Client, bridges []string) prometheus.Collector {
	return &collector{
		cs: []prometheus.Collector{
			// Additional generic netlink family collectors can be added here.
			newDatapathCollector(c.Datapath.List),
			newPortsCollector(ovsc.OpenFlow.DumpPorts, bridges),
		},
	}
}
