package influxdbAlert

import (
	"encoding/json"
	"fmt"
	"github.com/bmbbms/alert-rlist/alert"
	"github.com/bmbbms/alert-rlist/influxdbClient"
	"github.com/gin-gonic/gin"
	_ "github.com/influxdata/influxdb1-client"
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"net/http"
	"net/url"
	"time"
)

func GetAlertInfo(context *gin.Context) {
	getA06Info()
	context.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func getA06Info() error {
	url := &url.URL{
		Scheme: "https",
		Host:   "ts-wz95nqj644jxdfyl1.influxdata.rds.aliyuncs.com:3242",
	}
	config, err := influxdbClient.BuildConfig(url)
	client, err := influxdbClient.NewClient(config)
	defer client.Close()
	if err != nil {
		return err
	}

	if client == nil {
		panic("获取client失败")
	}

	query := influxdb.Query{
		Command: "SELECT count(chn_retcode) as A6_cnt FROM \"monitor_trans_rlist_exception_details\" WHERE " +
			"     (\"chn_retcode\" = 'A6' or  \"chn_retcode\" = '00A6' ) AND time>now() -1m group by time(1m) fill(0);",
		Database: "monitor",
	}
	response, err := client.Query(query)

	if err != nil {
		fmt.Println("Error Query data:", err.Error())
		return err
	}
	if response.Error() != nil {
		fmt.Println("Error parsing data:", response.Error().Error())
		return err
	}
	for _, result := range response.Results {
		if len(result.Series) == 0 {
			return nil
		}
		for _, series := range result.Series {
			for _, row := range series.Values {
				value, _ := row[1].(json.Number).Int64()
				if value > 10 {
					sendAlert(time.Now(), client, value)
				}
			}
		}
	}
	fmt.Printf("%#v", response)

	return nil
}

func sendAlert(now time.Time, client influxdb.Client, value int64) {
	query := influxdb.Query{
		Command: "SELECT count(begin_time)  FROM \"monitor_alert_history\" WHERE " +
			"(\"status\" = '1' ) AND time>now() -10m;",
		Database: "monitor",
	}
	response, err := client.Query(query)
	if err != nil {
		fmt.Println("Error Query data:", err.Error())
		return
	}
	if response.Error() != nil {
		fmt.Println("Error parsing data:", response.Error().Error())
		return
	}
	for _, result := range response.Results {
		if len(result.Series) == 0 {

			sendWxAlert("代付A6状态码告警", 10, value, now)
			point, _ := influxdb.NewPoint("monitor_alert_history", map[string]string{"alert_type": "1", "status": "1"},
				map[string]interface{}{"begin_time": time.Now()})

			batchPoints, _ := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
				Precision: "s",
				Database:  "monitor",
			})
			batchPoints.AddPoint(point)
			client.Write(batchPoints)

		}

	}

}

func sendWxAlert(s string, throld, value int64, now time.Time) {
	alertMsg := &alert.AlertMsg{
		Alert_name: s,
		Throld:     throld,
		Value:      value,
		Time:       now,
		ExtUrl: "http://172.20.20.40:3030/d/Ed_MWpl4k/dai-fu-cuo-wu-xiang-qing?" +
			"orgId=1&from=now-5m&to=now&var-idc=.*&var-ch_rottype=.*&var-status=.*&var-retcode=All",
	}
	notify := alert.NewWxNotify("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=e845e359-dda5-4731-811e-b1d09a0d8c8f")
	notify.Send(alertMsg)

}
