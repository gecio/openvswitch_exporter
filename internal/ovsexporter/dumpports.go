package ovsexporter

import (
	"fmt"
	"strconv"

	"github.com/gecio/go-openvswitch/ovs"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &portsCollector{}

type portsCollector struct {
	ReceivedPackets    *prometheus.Desc
	ReceivedBytes      *prometheus.Desc
	ReceivedDropped    *prometheus.Desc
	ReceivedErrors     *prometheus.Desc
	TransmittedPackets *prometheus.Desc
	TransmittedBytes   *prometheus.Desc
	TransmittedDropped *prometheus.Desc
	TransmittedErrors  *prometheus.Desc

	fn      func(string) ([]*ovs.PortStats, error)
	bridges []string
}

func newPortsCollector(fn func(string) ([]*ovs.PortStats, error), bridges []string) prometheus.Collector {
	const (
		subsystem = "ports"
	)

	var (
		labels = []string{"port", "bridge"}
	)
	return &portsCollector{
		ReceivedPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "received_packets_total"),
			"Number of packets received.",
			labels, nil,
		),
		ReceivedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "received_bytes_total"),
			"Number of bytes received.",
			labels, nil,
		),
		ReceivedDropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "received_dropped_total"),
			"Number of dropped packets on rx.",
			labels, nil,
		),
		ReceivedErrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "received_errors_total"),
			"Number of error packets on rx.",
			labels, nil,
		),
		TransmittedPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transmitted_packets_total"),
			"Number of packets transmitted.",
			labels, nil,
		),
		TransmittedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transmitted_bytes_total"),
			"Number of bytes transmitted.",
			labels, nil,
		),
		TransmittedDropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transmitted_dropped_total"),
			"Number of dropped packets on tx.",
			labels, nil,
		),
		TransmittedErrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transmitted_errors_total"),
			"Number of error packets on tx.",
			labels, nil,
		),
		bridges: bridges,
		fn:      fn,
	}
}

func (c *portsCollector) Describe(ch chan<- *prometheus.Desc) {
	ports := []*prometheus.Desc{
		c.ReceivedPackets,
	}

	for _, port := range ports {
		ch <- port
	}
}

func (c *portsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, bridge := range c.bridges {
		ports, err := c.fn(bridge)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(c.ReceivedPackets, fmt.Errorf("error dumping ports: %v", err))
			continue
		}
		for _, port := range ports {
			tuples := []struct {
				desc      *prometheus.Desc
				valueType prometheus.ValueType
				value     uint64
			}{
				{
					desc:      c.ReceivedPackets,
					valueType: prometheus.CounterValue,
					value:     port.Received.Packets,
				},
				{
					desc:      c.ReceivedBytes,
					valueType: prometheus.CounterValue,
					value:     port.Received.Bytes,
				},
				{
					desc:      c.ReceivedDropped,
					valueType: prometheus.CounterValue,
					value:     port.Received.Dropped,
				},
				{
					desc:      c.ReceivedErrors,
					valueType: prometheus.CounterValue,
					value:     port.Received.Errors,
				},
				{
					desc:      c.TransmittedPackets,
					valueType: prometheus.CounterValue,
					value:     port.Transmitted.Packets,
				},
				{
					desc:      c.TransmittedBytes,
					valueType: prometheus.CounterValue,
					value:     port.Transmitted.Bytes,
				},
				{
					desc:      c.TransmittedDropped,
					valueType: prometheus.CounterValue,
					value:     port.Transmitted.Dropped,
				},
				{
					desc:      c.TransmittedErrors,
					valueType: prometheus.CounterValue,
					value:     port.Transmitted.Errors,
				},
			}
			for _, tuple := range tuples {
				ch <- prometheus.MustNewConstMetric(tuple.desc, tuple.valueType, float64(tuple.value), strconv.Itoa(port.PortID), bridge)
			}
		}
	}
}
