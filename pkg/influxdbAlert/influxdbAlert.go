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
	"strconv"
	"time"
)

func InsertInfo(context *gin.Context) {
	InsertInfos()
	context.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})

}

func InsertInfos() {

}
func GetAlertInfo(context *gin.Context) {
	err := getA06Info()

	err = getRlistErrInfo()
	if err != nil {
		fmt.Println(err.Error())
	}
	context.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func getRlistErrInfo() error {
	client := connInfluxdb()
	defer func(client influxdb.Client) {
		err := client.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(client)

	query := influxdb.Query{
		//Command: "SELECT count(chn_retcode) as A6_cnt FROM \"monitor_trans_rlist_exception_details\" WHERE " +
		//	"     (\"chn_retcode\" = 'A6' or  \"chn_retcode\" = '00A6' ) AND time>now() -1m group by time(1m) fill(0);",
		Command: "SELECT count(status) as err_cnt FROM \"monitor_trans_rlist_exception_details\" WHERE " +
			"     (\"status\" = '2') AND time>now() -1m;",
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
	series := response.Results[0].Series
	if len(series) == 0 {
		return nil
	}
	row := series[0].Values[0]
	value, _ := row[1].(json.Number).Int64()
	fmt.Printf("过去一分钟代付错误数: %d", value)
	if value > 50 {
		sendAlert(time.Now(), client, value, 2, "代付失败订单数告警", 50)
	}

	return err

}
func connInfluxdb() influxdb.Client {
	influxdbUrl := &url.URL{
		Scheme: "https",
		Host:   "ts-wz95nqj644jxdfyl1.influxdata.rds.aliyuncs.com:3242",
	}

	config, err := influxdbClient.BuildConfig(influxdbUrl, "admin", "admin#20220818")
	client, err := influxdbClient.NewClient(config)
	if err != nil {
		panic("获取client异常" + err.Error())
	}

	if client == nil {
		panic("获取client失败")
	}
	return client

}
func getA06Info() error {
	client := connInfluxdb()
	defer func(client influxdb.Client) {
		err := client.Close()
		if err != nil {
			fmt.Println("close failed")
		}
	}(client)

	query := influxdb.Query{
		//Command: "SELECT count(chn_retcode) as A6_cnt FROM \"monitor_trans_rlist_exception_details\" WHERE " +
		//	"     (\"chn_retcode\" = 'A6' or  \"chn_retcode\" = '00A6' ) AND time>now() -1m group by time(1m) fill(0);",
		Command: "SELECT count(chn_retcode) as A6_cnt FROM \"monitor_trans_rlist_exception_details\" WHERE " +
			"     (\"chn_retcode\" = 'A6' or  \"chn_retcode\" = '00A6' ) AND time>now() -1m;",
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
					sendAlert(time.Now(), client, value, 1, "代付A6状态码告警", 10)
				}
			}
		}
	}
	fmt.Printf("%#v", response)

	return nil
}

func sendAlert(now time.Time, client influxdb.Client, value int64, alert_type int, alertTitle string, threshold int64) {
	query := influxdb.Query{
		Command: "SELECT count(begin_time)  FROM \"monitor_alert_history\" WHERE " +
			"(\"status\" = '1' ) AND \"alert_type\" ='" + strconv.Itoa(alert_type) + "' AND " + "time>now() -10m;",
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
	series := response.Results[0].Series
	if len(series) == 0 {
		sendWxAlert(alertTitle, threshold, value, now)
		point, _ := influxdb.NewPoint("monitor_alert_history", map[string]string{"alert_type": strconv.Itoa(alert_type), "status": "1"},
			map[string]interface{}{"begin_time": time.Now()})

		batchPoints, _ := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
			Precision: "s",
			Database:  "monitor",
		})
		batchPoints.AddPoint(point)
		err := client.Write(batchPoints)
		if err != nil {
			fmt.Println("write error", err.Error())
		}
	} else {
		fmt.Printf("\n静默状态-持续报警中:%s,阈值%d,当前值%d,告警时间%v", alertTitle, threshold, value, now)
	}

}

func sendWxAlert(s string, throld, value int64, now time.Time) {
	alertMsg := &alert.AlertMsg{
		Alert_name: s,
		Throld:     throld,
		Value:      value,
		Time:       now,
		ExtUrl: "https://kx-monitor.jlpay.com/d/Ed_MWpl4k/5Luj5LuY6ZSZ6K-v6K-m5oOF?" +
			"orgId=1&from=now-5m&to=now&var-idc=.*&var-ch_rottype=.*&var-status=.*&var-retcode=All",
	}
	//测试机器人
	//notify := alert.NewWxNotify("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=e845e359-dda5-4731-811e-b1d09a0d8c8f")
	notify := alert.NewWxNotify("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=0bd7e1aa-627f-4d6f-b6dd-afb1f988f841")
	notify.Send(alertMsg)

}
