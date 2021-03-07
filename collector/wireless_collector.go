package collector

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/routeros.v2/proto"
)

type wirelessCollector struct {
	props        []string
	descriptions map[string]*prometheus.Desc
}

func newWirelessCollector() routerOSCollector {
	c := &wirelessCollector{}
	c.init()
	return c
}

func (c *wirelessCollector) init() {
	c.props = []string{"name", "comment", "ssid", "frequency"}

	labelNames := []string{"name", "address", "ssid", "comment"}
	c.descriptions = make(map[string]*prometheus.Desc)
	for _, p := range c.props[1:] {
		c.descriptions[p] = descriptionForPropertyName("wireless", p, labelNames)
	}
}

func (c *wirelessCollector) describe(ch chan<- *prometheus.Desc) {
	for _, d := range c.descriptions {
		ch <- d
	}
}

func (c *wirelessCollector) collect(ctx *collectorContext) error {
	stats, err := c.fetch(ctx)
	if err != nil {
		return err
	}

	for _, re := range stats {
		c.collectForStat(re, ctx)
	}

	return nil
}

func (c *wirelessCollector) fetch(ctx *collectorContext) ([]*proto.Sentence, error) {
	reply, err := ctx.client.Run("/interface/wireless/print", "=.proplist="+strings.Join(c.props, ","))
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching interface metrics")
		return nil, err
	}

	return reply.Re, nil
}

func (c *wirelessCollector) collectForStat(re *proto.Sentence, ctx *collectorContext) {
	name := re.Map["name"]
	comment := re.Map["comment"]

	for _, p := range c.props[3:] {
		c.collectMetricForProperty(p, name, comment, re, ctx)
	}
}

func (c *wirelessCollector) collectMetricForProperty(property, iface, comment string, re *proto.Sentence, ctx *collectorContext) {
	desc := c.descriptions[property]
	if value := re.Map[property]; value != "" {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.WithFields(log.Fields{
				"device":    ctx.device.Name,
				"interface": iface,
				"property":  property,
				"value":     value,
				"error":     err,
			}).Error("error parsing interface metric value")
			return
		}
		ctx.ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, v, ctx.device.Name, ctx.device.Address, iface, comment)
	}
}
