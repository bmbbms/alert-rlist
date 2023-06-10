package main

import (
	"flag"
	"fmt"
	env "github.com/Netflix/go-env"
	"github.com/bmbbms/alert-rlist/config"
	"github.com/bmbbms/alert-rlist/pkg/influxdbAlert"
	"github.com/gin-gonic/gin"
	"os"
)

var configs = config.Config{
	WebhookUrl:  "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=e845e359-dda5-4731-811e-b1d09a0d8c8f",
	InfluxdbUrl: "https://ts-wz95nqj644jxdfyl1.influxdata.rds.aliyuncs.com:3242?user=amdin&pw=admin#20220818&db=monitor",
	ExtUrl:      "http://172.20.20.40:3030/d/Ed_MWpl4k/dai-fu-cuo-wu-xiang-qing?orgId=1&from=now-5m&to=now&var-idc=.*&var-ch_rottype=.*&var-status=.*&var-retcode=All",
}

func main() {
	if _, err := env.UnmarshalFromEnviron(&configs); err != nil {
		return
	}
	configs.AddFlags()

	flag.Parse()

	showVersion()
	engine := gin.Default()
	engine.GET("/alert", influxdbAlert.GetAlertInfo)
	engine.Run(":8080")
}

func showVersion() {
	if config.ShowVer {
		fmt.Printf("build name:\t%s\n", config.BuildName)
		fmt.Printf("build ver:\t%s\n", config.BuildVersion)
		fmt.Printf("build time:\t%s\n", config.BuildTime)
		fmt.Printf("Commit ID:\t%s\n", config.CommitID)
		os.Exit(0)
	}
}
