package common

type ParamsFinanceLogIGetYesterday struct {
	//Interface Params: IntraPayICPayChange
	UserId string `json:"userid"`
}

type ParamsFinanceLogIGetDutyLogs struct {
	//Interface Params: IntraPayICPayChange
	UserId   string `json:"userid"`
	Start    int    `json:"start"`
	PageSize int    `json:"pagesize"`
}

//ParamsRegistEnd
