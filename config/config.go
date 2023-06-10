package config

import "flag"

// 配置信息
type Config struct {
	WebhookUrl  string `env:"ALERT_WEBHOOK"`
	InfluxdbUrl string `env:"ALERT_INFLUXDB_URL"`
	ExtUrl      string `env:"ALERT_EXTURL"`
}

func (c *Config) AddFlags() {
	flag.StringVar(&c.WebhookUrl, "alert-webhookurl", c.WebhookUrl, "webhookurl for alert")
	flag.StringVar(&c.InfluxdbUrl, "alert-InfluxdbUrl", c.InfluxdbUrl, "InfluxdbUrl for query rlist data")
	flag.StringVar(&c.ExtUrl, "alert-ExtUrl", c.ExtUrl, "ExtUrl for ext alert message")

}
