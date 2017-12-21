package main

import (
	"encoding/json"
	"git.jtjr.com/bg_prd_grp/financelog/common"
	"git.jtjr.com/bg_prd_grp/financelog/errcode"
	"github.com/astaxie/beego/orm"
	"log"
	"strconv"
)

/*
type TestDB struct {
	UserID 	  int         `json:"user_id"`
	UserName  string      `json:"user_name"`
}
*/
type DutyLogs struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
	Time   string `json:"time"`
	Amount string `json:"amount"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func GetDutyLogs(params interface{}) (interface{}, int) {
	req_data := params.(*common.ParamsFinanceLogIGetDutyLogs)
	userid := req_data.UserId
	start := req_data.Start
	pagesize := req_data.PageSize
	log.Println("GetDutyLogs:userid:", userid, " start:", start, " pagesize:", pagesize)

	//查表
	var duty_logs []DutyLogs
	//var resp RespGetDutyLogs
	sql := "select '1' `type`,"
	sql += "'" + BONUS_ACCOUNT_DAILY + "'" + " `text`,"
	sql += "tx_time `time`,CONCAT('+',IFNULL(SUM(tx_amount),0)) amount,"
	sql += "IF(ttl.tx_result = '1',0,ttl.tx_result) `status`,'确认成功' msg "
	sql += "from t_tx_log ttl,t_tx_bonus ttb,t_prd_inst tpi "
	sql += "where ttl.cust_id = " + userid
	sql += " AND ttl.tx_type = '8' AND ttb.tx_no = ttl.tx_no"
	sql += " AND tpi.prd_inst_id = ttb.prd_inst_id AND tpi.prd_inst_status != '6'"
	sql += " GROUP BY DATE(tx_time) ORDER BY `time` DESC  LIMIT "
	sql += strconv.Itoa(start) + "," + strconv.Itoa(pagesize) + ";"

	log.Println(sql)

	//从缓存获取数据
	var sql_key string
	sql_key = common.Get_sql_key(sql)

	v, ok := common.RedisGet(sql_key)
	if ok {
		log.Println("GetDutyLogs:get from redis success,key:", sql_key, "value:", v)
		//v_obj := []DutyLogs{}
		//json.Unmarshal([]byte(v), &v_obj)
		//resp.Data.LogData = v_obj
		//resp.Msg = "ok"
		//return resp, 0
		return map[string]string{
			"log": v,
		}, errcode.ERRCODE_SUCCESS
	}

	log.Println("GetDutyLogs:get from redis fail,query from db,key:", sql_key)

	o := orm.NewOrm()
	o.Using("default")

	num, db_err := o.Raw(sql).QueryRows(&duty_logs)
	if db_err != nil && db_err != orm.ErrNoRows {

		log.Println("GetDutyLogs:query db falied:", "DB_ERR: ", db_err)
		//resp.Msg = "query logs failed."
		//return resp, 0
		return map[string]string{
			"msg": "查询数据库失败",
		}, errcode.ERRCODE_QRY_LOG_FAILED
	}

	var log_str string
	if num > 0 {
		//resp.Data.LogData = duty_logs
		logs_val, err := json.Marshal(&duty_logs)
		if err != nil {
			log.Println("GetDutyLogs:encoding duty_logs faild")
			//resp.Msg = "encoding duty_logs faild."
			return map[string]string{
				"msg": "格式化json数据失败",
			}, errcode.ERRCODE_ENCODE_LOG_FAILED
		} else {
			log_str = string(logs_val)
			log.Println("GetDutyLogs:encoded data ", log_str)
		}

		//加入到redis中
		if set_ok := common.RedisSet(sql_key, log_str, 0); set_ok {
			common.RedisHSet("t_tx_log", sql_key, log_str)
			if ok := common.RedisHSet("t_tx_log", sql_key, log_str); !ok {
				log.Println("GetDutyLogs:RedisHSet failed,key ", sql_key)
			}
		} else {
			log.Println("GetDutyLogs:RedisSet failed,key:", sql_key)
		}

	} else {
		//resp.Msg = "no data found"
		return map[string]string{
			"msg": "数据库没有匹配的记录",
		}, errcode.ERRCODE_NO_DATA_FOUND
	}

	//return resp, 0
	return map[string]string{
		"log": log_str,
	}, errcode.ERRCODE_SUCCESS
}
