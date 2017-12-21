package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"git.jtjr.com/bg_prd_grp/intrapay/errcode"
)

//发送http消息
/*  sUrl 		请求地址
*	sReqData 	请求数据
*	nTimeOut	超时时间
 */
func SendHttpPost(cmd, sUrl, sReqData string, nTimeout int) (int, string) {
	hc := &http.Client{
		Timeout: time.Duration(nTimeout) * time.Second,
	}

	//fmt.Printf("[0]url:%s, cmd: %s, req: %s\n", sUrl, cmd, sReqData)
	kstart := time.Now().UnixNano()
	resp, err := hc.Post(sUrl, "application/x-www-form-urlencoded", strings.NewReader(sReqData))
	if err != nil {
		log.Printf("[1]url:%s, cmd: %s, req: %s\n", sUrl, cmd, sReqData)
		Log_error(string(errcode.ERRCODE_OUTSRV_FALIED), "rpc call failed", "cmd", cmd, "url", sUrl, "req", sReqData, "rsp", "")
		return -1, fmt.Sprintf("get http send error: \n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[2]url:%s, cmd: %s, req: %s\n", sUrl, cmd, sReqData)
		Log_error(string(errcode.ERRCODE_OUTSRV_PARSE_RSP_FAILED), "parse json failed", "cmd", cmd, "url", sUrl, "req", sReqData, "rsp", "")
		return -2, fmt.Sprintf("parse http data error: \n", err)
	}

	//计算时间消耗
	kend := time.Now().UnixNano()
	kcost := (kend - kstart) / 1000000 // 达到毫秒级别

	res_data := string(body[0:])
	res_data = strings.TrimSpace(res_data)
	//log.Println("Url:", sUrl, "cmd:", cmd, "req:", sReqData, "rsp:", res_data, "cost:", kcost)
	Log_info("cmd", cmd, "cost", fmt.Sprintf("%.0f", float64(kcost)), "req", sReqData, "rsp", res_data)

	return 0, string(body[0:])
}

func SendHttpPostNoLog(cmd, sUrl, sReqData string, nTimeout int) (int, string) {
	hc := &http.Client{
		Timeout: time.Duration(nTimeout) * time.Second,
	}

	//fmt.Printf("[0]url:%s, cmd: %s, req: %s\n", sUrl, cmd, sReqData)
	kstart := time.Now().UnixNano()
	resp, err := hc.Post(sUrl, "application/x-www-form-urlencoded", strings.NewReader(sReqData))
	if err != nil {
		log.Printf("[1]url:%s, cmd: %s, req: %s\n", sUrl, cmd, sReqData)
		Log_error(string(errcode.ERRCODE_OUTSRV_FALIED), "rpc call failed", "cmd", cmd, "url", sUrl, "req", sReqData, "rsp", "")
		return -1, fmt.Sprintf("get http send error: \n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[2]url:%s, cmd: %s, req: %s\n", sUrl, cmd, sReqData)
		Log_error(string(errcode.ERRCODE_OUTSRV_PARSE_RSP_FAILED), "parse json failed", "cmd", cmd, "url", sUrl, "req", sReqData, "rsp", "")
		return -2, fmt.Sprintf("parse http data error: \n", err)
	}

	//计算时间消耗
	kend := time.Now().UnixNano()
	kcost := (kend - kstart) / 1000000 // 达到毫秒级别

	res_data := string(body[0:])
	res_data = strings.TrimSpace(res_data)
	//log.Println("Url:", sUrl, "cmd:", cmd, "req:", sReqData, "rsp:", res_data, "cost:", kcost)
	Log_info("cmd", cmd, "cost", fmt.Sprintf("%.0f", float64(kcost)), "req", sReqData, "rsp", res_data)

	return 0, string(body[0:])
}
