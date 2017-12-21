package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"errors"

	"github.com/fluent/fluent-logger-golang/fluent"
)

var g_logger *fluent.Fluent = nil
var g_log_tag = "monitor"
var g_log_tag_interface = "inf"
var g_log_tag_server_log = "log"
var g_log_tag_component = "com"
var g_log_tag_biz = "biz"

var g_host = "127.0.0.1"
var g_seq = 0
var locker sync.Mutex

func get_seq() int {
	locker.Lock()
	g_seq = g_seq + 1
	defer locker.Unlock()
	return g_seq
}

/*
func init() {
	init_logger()
}
*/
func Set_host(host string) {
	g_host = host
	init_logger()
}
func Log_init(monitor_id string) {
	g_log_tag = "monitor." + monitor_id
	g_log_tag_interface = g_log_tag + "." + g_log_tag_interface
	g_log_tag_server_log = g_log_tag + "." + g_log_tag_server_log
	g_log_tag_component = g_log_tag + "." + g_log_tag_component
	g_log_tag_biz = g_log_tag + "." + g_log_tag_biz
}
func init_logger() {
	defer catch_err("init_logger")
	var cfg = fluent.Config{FluentHost: g_host, Timeout: (300 * time.Millisecond)}
	var err error
	var tmp_logger *fluent.Fluent
	tmp_logger, err = fluent.New(cfg)
	if err != nil {
		panic(err)
	}
	set_logger(tmp_logger)
}
func set_logger(new_logger *fluent.Fluent) {
	locker.Lock()
	defer locker.Unlock()
	g_logger = new_logger //注意这里不能用 :=， 否则 g_logger就被认为是新的变量
}
func close_logger() {
	locker.Lock()
	defer locker.Unlock()
	if g_logger != nil {
		g_logger.Close()
	}
}
func post_logger(tag string, message interface{}) error {
	locker.Lock()
	defer locker.Unlock()
	if g_logger == nil {
		return errors.New("g_logger is nil")
	}
	return g_logger.Post(tag, message)
}

func Post(tag string, message interface{}) {
	defer func() { //拦截post时的未知异常
		close_logger()
		if err := recover(); err != nil {
			Print(err) //这里的err其实就是panic传入的内容
		}
	}()

	map_v, ok := message.(map[string]interface{})
	if ok {
		map_v["_seq"] = get_seq()
	}
	error := post_logger(tag, message)
	if error != nil {
		init_logger()
	}
}
func print_frame(frame_idx int, msg ...interface{}) {
	sys_info := get_sys_info(frame_idx)
	time_ns := time.Now().UnixNano()
	time_ms := time_ns / 1e3
	time_str := time.Unix(time_ns/1e9, 0).Format("2006-01-02 15:04:05") //方式比较特别，按照123456来记忆吧：01月02号 下午3点04分05秒 2006年

	fmt.Print("[", time_str, ".", fmt.Sprintf("%03d", time_ms%1000), "] [", sys_info["file"], "(", sys_info["line"], ")][PID=", sys_info["pid"], "] ", fmt.Sprintln(msg...)) //不定参数往下传递

}
func Print(msg ...interface{}) {
	print_frame(3, msg...)
}
func catch_err(tag string) {
	if err := recover(); err != nil {
		Print(tag, err)
	}
}
func parse_param(kvs ...interface{}) (data map[string]interface{}) {
	data = get_sys_info(3)
	var last_key string = ""
	for _, item := range kvs {
		if last_key == "" {
			tmp_key, ok := item.(string)
			if ok {
				last_key = tmp_key
			}
		} else {
			data[last_key] = item
			last_key = ""
		}
	}
	return
}

func Log_info(kvs ...interface{}) {
	defer catch_err("Log_info") //拦截异常
	data := parse_param(kvs...)
	data["log_level"] = "info"
	Post(g_log_tag_server_log, data)
}

func Log_error(errcode interface{}, info string, kvs ...interface{}) {
	defer catch_err("Log_error") //拦截异常
	data := parse_param(kvs...)
	data["log_level"] = "error"
	data["errcode"] = errcode
	data["info"] = info
	Post(g_log_tag_server_log, data)
}

func Log_debug(kvs ...interface{}) {
	defer catch_err("Log_debug") //拦截异常
	data := parse_param(kvs...)
	data["log_level"] = "debug"
	Post(g_log_tag_server_log, data)
}
func Log_fatal(info string, kvs ...interface{}) {
	defer catch_err("Log_fatal") //拦截异常
	data := parse_param(kvs...)
	data["log_level"] = "fatal"
	data["info"] = info
	Post(g_log_tag_server_log, data)
}

/**
 * 函数说明：被调方上报服务质量
 * 参数列表：
 *      cmd          当前处理的命令字
 *      cmt_time_s   从收到请求到发出响应包的耗时
 *      errcode      错误码或返回码
 *      info         错误信息
 */
func Log_profile(cmd interface{}, cmd_time_ms int64, errcode interface{}, info string, kvs ...interface{}) {
	defer catch_err("Log_profile") //拦截异常
	data := parse_param(kvs...)
	data["log_level"] = "profile"
	data["cmd"] = cmd
	data["cmd_time_s"] = cmd_time_ms
	data["errcode"] = errcode
	data["info"] = info
	Post(g_log_tag_server_log, data)
}

/**
 * 函数说明：调用方上报接口质量
 * 参数列表：
 *      call_cmd     调用的命令字
 *      callee_addr  调用方地址
 *      call_time    发起调用的时间
 *      time_span    接口耗时,单位：毫秒
 *      ret_code     返回码
 */
func Log_interface(call_cmd interface{}, callee_addr string, call_time string, time_span int, ret_code interface{}, kvs ...interface{}) {
	defer catch_err("Log_interface") //拦截异常
	data := parse_param(kvs)
	data["call_cmd"] = call_cmd
	data["callee_addr"] = callee_addr
	data["call_time"] = call_time
	data["time_span"] = time_span
	data["ret_code"] = ret_code
	data["pid"] = os.Getpid() //调用一次耗时约0.001ms
	Post(g_log_tag_interface, data)
}

/**
 * 函数说明:上报组件状态
 *      srv_cat      服务种类：mysql\redis\kafka\rabbitmq
 *      srv_type     服务类型。0：存储类，1：队列类
 *      disc_ratio   磁盘使用率。乘以100的整数。
 *      mem_ratio    内存使用率。乘以100的整数。
 *      cpu_ratio    cpu使用率。乘以100的整数。
 *      pid          服务的进程ID
 */
func Log_com(srv_cat string, srv_type int, disc_ratio int, mem_ratio int, cpu_ratio int, pid int, kvs ...interface{}) {
	defer catch_err("Log_com") //拦截异常
	data := parse_param(kvs...)
	data["srv_cat"] = srv_cat
	data["srv_type"] = srv_type
	data["disc_ratio"] = disc_ratio
	data["mem_ratio"] = mem_ratio
	data["cpu_ratio"] = cpu_ratio
	data["pid"] = pid
	Post(g_log_tag_component, data)
}

/**
 * 函数说明: 自定义上报。（携带系统字段）
 * 参数列表：
 *      postfix     tag的后缀
 */
func Log_custom(postfix string, kvs ...interface{}) {
	defer catch_err("Log_custom") //拦截异常
	data := parse_param(kvs...)
	Post(g_log_tag+"."+postfix, data)
}

func get_sys_info(base_frame int) (sys_info map[string]interface{}) {
	sys_info = make(map[string]interface{})
	funcName, file, line, _ := runtime.Caller(base_frame)
	funcNameStr := runtime.FuncForPC(funcName).Name()
	sys_info["pid"] = os.Getpid()
	sys_info["line"] = line
	sys_info["file"] = filepath.Base(file)
	sys_info["func"] = funcNameStr
	return
}
