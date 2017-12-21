package errcode

import (
	"errors"
	"fmt"
)

/************************************************************************/
/*							错误码                                      */
/************************************************************************/
/**
@api {POST} / [ERROR_CODE]-通用错误码
@apiName COMMON_ERROR_CODE
@apiParam {int} 2101000001 [<font color=red>对外</font>]未知错误
@apiParam {int} 2101000002 [<font color=red>对外</font>]参数错误
@apiParam {int} 2101000003 [<font color=blue>对内</font>]未知命令
@apiParam {int} 2101000004 [<font color=blue>对内</font>]请求失败
@apiParam {int} 2101000005 [<font color=blue>对内</font>]外部返回错误
@apiParam {int} 2101000006 [<font color=blue>对内</font>]网络错误
@apiParam {int} 2101000007 [<font color=blue>对内</font>]命令请求记录失败
@apiParam {int} 2101000008 [<font color=blue>对内</font>]数据库操作失败
@apiParam {int} 2101000009 [<font color=blue>对内</font>]请求外部服务出错
@apiParam {int} 2101000010 [<font color=blue>对内</font>]解析外部请求返回数据失败
@apiParam {int} 2101000011 [<font color=blue>对内</font>]RPC调用其他服务返回错误数据
@apiParam {int} 2101000012 [<font color=blue>对内</font>]线下转账（平安）不存在的三方交易流水号
@apiGroup COMMON_ERROR_CODE
@apiVersion 1.0.0
*/

/************************************************************************/
/*							数据定义                                    */
/************************************************************************/

const ErrCodeBase = 1420000000
const ErrYSICodeBase = 1201000700
const (
	ERRCODE_SUCCESS = 0
	ERRCODE_UNKNOW  = ErrCodeBase + 1 //1420000001

	ERRCODE_PARAM_INVAILD     = ErrCodeBase + 2 //1420000002
	ERRCODE_CMD_INVAILD       = ErrCodeBase + 3 //1420000003
	ERRCODE_QRY_LOG_FAILED    = ErrCodeBase + 4 //1420000004
	ERRCODE_ENCODE_LOG_FAILED = ErrCodeBase + 5 //1420000005
	ERRCODE_NO_DATA_FOUND     = ErrCodeBase + 6 //1420000006

	ERRCODE_QRY_DAILY_BONUS_FAILED    = ErrYSICodeBase + 2 //1201000702
	ERRCODE_DAILY_BONUS_NO_DATA_FOUND = ErrYSICodeBase + 3 //1201000703
)

type ErrCode struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

var mapErr map[int]string

func init() {
	mapErr = map[int]string{
		ERRCODE_SUCCESS:           "ok",
		ERRCODE_UNKNOW:            "未知错误",
		ERRCODE_PARAM_INVAILD:     "参数错误",
		ERRCODE_CMD_INVAILD:       "未知命令",
		ERRCODE_QRY_LOG_FAILED:    "query logs failed.",
		ERRCODE_ENCODE_LOG_FAILED: "encoding duty_logs faild.",
		ERRCODE_NO_DATA_FOUND:     "no data found",

		ERRCODE_QRY_DAILY_BONUS_FAILED:    "database error",
		ERRCODE_DAILY_BONUS_NO_DATA_FOUND: "no data found",
	}
}

func (err *ErrCode) SetErrCode(errcode int) error {

	err.Code = fmt.Sprint(errcode)

	var ok bool

	if err.Msg, ok = mapErr[errcode]; !ok {

		err.Code = fmt.Sprint(ERRCODE_UNKNOW)

		return errors.New("unknow errcode, origin errcode :" + fmt.Sprint(errcode))
	}

	err.Msg, _ = mapErr[errcode]

	return nil
}
