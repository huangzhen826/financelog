package main

import (
	"git.jtjr.com/bg_prd_grp/financelog/common"
)

func initCmd() {

	//IRegistStart
	RegistCmd(CMD_GET_YESTERDAY_INTEREST, CmdHandle{cfunc: FFinanceLogIGetYesterdayInterest, cparam: common.ParamsFinanceLogIGetYesterday{}})
	RegistCmd(CMD_GET_DUTY_LOGS, CmdHandle{cfunc: FFinanceLogIGetDutyLogs, cparam: common.ParamsFinanceLogIGetDutyLogs{}})
	//IRegistEnd
}
