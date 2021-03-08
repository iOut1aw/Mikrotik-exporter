package collector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/routeros.v2/proto"
)

type wlanIFCollector struct {
	props                 []string
	descriptions          map[string]*prometheus.Desc
	propsFrequency        []string
	descriptionsFrequency map[string]*prometheus.Desc
}

func newWlanIFCollector() routerOSCollector {
	c := &wlanIFCollector{}
	c.init()
	return c
}

func (c *wlanIFCollector) init() {
	c.props = []string{"channel", "registered-clients", "noise-floor", "overall-tx-ccq"}
	labelNames := []string{"name", "address", "interface", "channel"}
	c.descriptions = make(map[string]*prometheus.Desc)
	for _, p := range c.props {
		c.descriptions[p] = descriptionForPropertyName("wlan_interface", p, labelNames)
	}
	c.propsFrequency = []string{"name", "comment", "ssid", "frequency"}
	labelNamesFrequency := []string{"name", "address", "ssid", "comment"}
	c.descriptionsFrequency = make(map[string]*prometheus.Desc)
	for _, p := range c.propsFrequency[1:] {
		c.descriptionsFrequency[p] = descriptionForPropertyName("wlan_interface", p, labelNamesFrequency)
	}

}

func (c *wlanIFCollector) describe(ch chan<- *prometheus.Desc) {
	for _, d := range c.descriptions {
		ch <- d
	}
}

func (c *wlanIFCollector) collect(ctx *collectorContext) error {
	return c.fetchInterfaceNames(ctx)
}

func (c *wlanIFCollector) fetchInterfaceNames(ctx *collectorContext) error {
	reply, err := ctx.client.Run("/interface/wireless/print", "?disabled=false", "=.proplist="+strings.Join(c.propsFrequency, ","))
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching wireless interface names")
		return err
	}

	for _, re := range reply.Re {
		c.collectForStatFrequency(re, ctx)
		err := c.collectForInterface(re.Map["name"], ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *wlanIFCollector) collectForInterface(iface string, ctx *collectorContext) error {
	reply, err := ctx.client.Run("/interface/wireless/monitor", fmt.Sprintf("=numbers=%s", iface), "=once=", "=.proplist="+strings.Join(c.props, ","))
	if err != nil {
		log.WithFields(log.Fields{
			"interface": iface,
			"device":    ctx.device.Name,
			"error":     err,
		}).Error("error fetching interface statistics")
		return err
	}

	for _, p := range c.props[1:] {
		// there's always going to be only one sentence in reply, as we
		// have to explicitly specify the interface
		c.collectMetricForProperty(p, iface, reply.Re[0], ctx)
	}

	return nil
}

func (c *wlanIFCollector) collectForStatFrequency(re *proto.Sentence, ctx *collectorContext) {
	ssid := re.Map["ssid"]
	comment := re.Map["comment"]

	for _, p := range c.propsFrequency[3:] {
		c.collectMetricForPropertyFrequency(p, ssid, comment, re, ctx)
	}
}

func (c *wlanIFCollector) collectMetricForPropertyFrequency(property, ssid, comment string, re *proto.Sentence, ctx *collectorContext) {
	desc := c.descriptionsFrequency[property]
	if value := re.Map[property]; value != "" {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.WithFields(log.Fields{
				"device":    ctx.device.Name,
				"interface": ssid,
				"property":  property,
				"value":     value,
				"error":     err,
			}).Error("error parsing interface metric value")
			return
		}
		ctx.ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, v, ctx.device.Name, ctx.device.Address, ssid, comment)
	}
}

func (c *wlanIFCollector) collectMetricForProperty(property, iface string, re *proto.Sentence, ctx *collectorContext) {
	desc := c.descriptions[property]
	channel := re.Map["channel"]
	if re.Map[property] == "" {
		return
	}
	v, err := strconv.ParseFloat(re.Map[property], 64)
	if err != nil {
		log.WithFields(log.Fields{
			"property":  property,
			"interface": iface,
			"device":    ctx.device.Name,
			"error":     err,
		}).Error("error parsing interface metric value")
		return
	}

	ctx.ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, v, ctx.device.Name, ctx.device.Address, iface, channel)
}
