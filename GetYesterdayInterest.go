package main

import (
	"git.jtjr.com/bg_prd_grp/financelog/common"
	"git.jtjr.com/bg_prd_grp/financelog/errcode"
	"github.com/astaxie/beego/orm"
	"log"
	//"strconv"
)

/*
type TestDB struct {
	UserID 	  int         `json:"user_id"`
	UserName  string      `json:"user_name"`
}
*/
type DailyBonus struct {
	Bonus string `json:"pre_bonus"`
}

/*
type RespGetYesterday struct {
	Msg   string `json:"msg"`
	Bonus string `json:"pre_bonus"`
}
*/

func GetYesterdayInterest(params interface{}) (interface{}, int) {
	req_data := params.(*common.ParamsFinanceLogIGetYesterday)
	userid := req_data.UserId
	log.Println("GetYesterdayInterest:userid:", userid)

	//查表
	var bonus DailyBonus
	//var resp RespGetYesterday
	sql := "SELECT bonus FROM t_tx_daily_bonus where cust_id="
	sql += userid
	sql += " AND bonus_date = DATE(DATE_SUB(NOW(), INTERVAL 1 DAY));"

	log.Println(sql)

	o := orm.NewOrm()
	o.Using("default")

	db_err := o.Raw(sql).QueryRow(&bonus)
	if db_err != nil && db_err != orm.ErrNoRows {

		log.Println("GetYesterdayInterest:query t_tx_daily_bonus falied:", "DB_ERR: ", db_err)
		//resp.Msg = "query t_tx_daily_bonus failed"
		return map[string]string{
			"msg": "查询数据库失败",
		}, errcode.ERRCODE_QRY_DAILY_BONUS_FAILED
	}

	if db_err == orm.ErrNoRows {
		return map[string]string{
			"msg": "没有查询到该userid的昨日收益记录",
		}, errcode.ERRCODE_DAILY_BONUS_NO_DATA_FOUND
	}

	log.Println("[bonus]:", bonus.Bonus)

	//resp.Msg = "ok"
	//resp.Bonus = bonus.Bonus
	return map[string]string{
		"bonus": bonus.Bonus,
	}, errcode.ERRCODE_SUCCESS
}
