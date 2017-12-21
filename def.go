package main

/******************** 命令字定义 ************************/
const CMD_GET_YESTERDAY_INTEREST = "12010007" // "12010007" 昨日收益
const CMD_GET_DUTY_LOGS = "142"               // "142" GetDutyLogs

/******************** 资产明细显示常量 ************************/
const BONUS_ACCOUNT_DAILY = "日收益结算"

/******************** 通用数据定义 ************************/
const GLOBLE_PMUTEX_KEY = "OFFLINETRANS_PMUTEX_KEY_f9d5625d42a022977a2ca70c69e5da6a" //redis分布式锁key
