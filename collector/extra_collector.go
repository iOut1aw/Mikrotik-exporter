package collector

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/routeros.v2/proto"
)

type extraCollector struct {
	props        []string
	descriptions *prometheus.Desc
}

func (c *extraCollector) init() {
	c.props = []string{"current-firmware", "factory-firmware", "firmware-type", "model", "serial-number"}

	labelNames := []string{"name", "address", "currentfirmware", "factoryfirmware", "firmwaretype", "model", "serialnumber"}
	c.descriptions = description("extra", "metrics", "number of metrics", labelNames)

}

func newExtraCollector() routerOSCollector {
	c := &extraCollector{}
	c.init()
	return c
}

func (c *extraCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.descriptions
}

func (c *extraCollector) collect(ctx *collectorContext) error {
	stats, err := c.fetch(ctx)
	if err != nil {
		return err
	}

	for _, re := range stats {
		c.collectMetric(ctx, re)
	}

	return nil
}

func (c *extraCollector) fetch(ctx *collectorContext) ([]*proto.Sentence, error) {
	reply, err := ctx.client.Run("/system/routerboard/print", "=.proplist="+strings.Join(c.props, ","))
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching Extra metrics")
		return nil, err
	}
	return reply.Re, nil
}

func (c *extraCollector) collectMetric(ctx *collectorContext, re *proto.Sentence) {
	v := 1.0

	currentFirmware := re.Map["current-firmware"]
	factoryFirmware := re.Map["factory-firmware"]
	firmwareType := re.Map["firmware-type"]
	model := re.Map["model"]
	serialNumber := re.Map["serial-number"]

	ctx.ch <- prometheus.MustNewConstMetric(c.descriptions, prometheus.CounterValue, v, ctx.device.Name, ctx.device.Address, currentFirmware, factoryFirmware, firmwareType, model, serialNumber)
}
