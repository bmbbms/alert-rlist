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

var Configs = config.Config{
	WebhookUrl:  "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=e845e359-dda5-4731-811e-b1d09a0d8c8f",
	InfluxdbUrl: "https://ts-wz95nqj644jxdfyl1.influxdata.rds.aliyuncs.com:3242?user=amdin&pw=admin#20220818&db=monitor",
	ExtUrl:      "http://172.20.20.40:3030/d/Ed_MWpl4k/dai-fu-cuo-wu-xiang-qing?orgId=1&from=now-5m&to=now&var-idc=.*&var-ch_rottype=.*&var-status=.*&var-retcode=All",
}

func main() {
	if _, err := env.UnmarshalFromEnviron(&Configs); err != nil {
		return
	}
	Configs.AddFlags()

	flag.Parse()

	showVersion()
	engine := gin.Default()
	engine.GET("/alert", influxdbAlert.GetAlertInfo)
	engine.GET("/insert", influxdbAlert.InsertInfo)
	err := engine.Run(":8080")
	if err != nil {
		fmt.Println("启动失败:", err.Error())
		return
	}
}

func showVersion() {
	if ShowVer {
		fmt.Printf("build name:\t%s\n", BuildName)
		fmt.Printf("build ver:\t%s\n", BuildVersion)
		fmt.Printf("build time:\t%s\n", BuildTime)
		fmt.Printf("Commit ID:\t%s\n", CommitID)
		os.Exit(0)
	}
}
