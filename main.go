package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"github.com/rakyll/globalconf"

	"git.jtjr.com/bg_prd_grp/financelog/common"
	"git.jtjr.com/bg_prd_grp/financelog/errcode"
	//"git.jtjr.com/bg_prd_grp/rabbitmq"
	"git.jtjr.com/broccoli/sdk"
)

var (
	//服务启动监听地址和端口
	g_srvport = flag.String("srvport", "39353", "server port")

	//消息队列MQ地址
	g_mqaddr = flag.String("mqaddr", "amqp://jyb:root@172.16.1.8:5672", "mq server addr")

	//数据库地址
	g_mysqhost   = flag.String("master_db_host", "172.16.1.16:3306", "myssql host")
	g_mysqdb     = flag.String("master_db_name", "db_jyb_test", "myssql db name")
	g_mysquser   = flag.String("master_db_user", "jiayoubao", "myssql user")
	g_mysqpasswd = flag.String("master_db_passwd", "root1234", "myssql passwd")

	//redis地址
	g_redisaddr = flag.String("redis_addr", "172.16.1.13:6379", "redis mq server addr")

	//日志上报配置
	g_loghost  = flag.String("loghost", "127.0.0.1", "log host")
	g_logtitle = flag.String("logtitle", "intrapay_service", "log title")

	//微服务配置
	g_jmfpathbase = flag.String("jmfpathbase", "com.jyblife.base.bg", "jmf path base")
	g_jmfapp      = flag.String("jmfapp", "", "jmf app name")
	g_env         = flag.String("env", "sit", "env")
	g_set         = flag.String("set", "gz_shenzhen_idc", "set")
	g_verion      = flag.String("version", "1.0.0", "version")
	g_owner       = flag.String("owner", "", "server owner")
	g_group       = flag.String("group", "*", "server group")

	//外部服务接口参数
	g_out_ptb_url     = flag.String("ptb_url", "http://172.16.1.8:9729/", "获取余额支付票据服务地址")
	g_out_ptb_timeout = flag.String("ptb_timeout", "1", "超时时间 单位秒")

	g_out_txpay_url     = flag.String("txpay_url", "http://172.16.1.8:9629/", "保存交易流水")
	g_out_txpay_timeout = flag.String("txpay_timeout", "1", "超时时间 单位秒")

	g_out_ptcou_url     = flag.String("ptcou_ur", "http://172.16.1.16:9633", "获取红包消费凭证")
	g_out_ptcou_timeout = flag.String("ptcou_timeout", "1", "超时时间 单位秒")

	g_out_txno_url     = flag.String("txno_url", "http://172.16.1.8:9629/", "申请交易流水")
	g_out_txno_timeout = flag.String("txno_timeout", "1", "超时时间 单位秒")

	g_out_ptcha_url     = flag.String("ptcha_url", "http://172.16.1.8:9629/", "获取零钱消费票据")
	g_out_ptcha_timeout = flag.String("ptcha_timeout", "1", "超时时间 单位秒")

	/*******************************************【获取平安存管转账成功明细】*******************************************/
	g_out_patrans_opensrv_flag = flag.String("patrans_opensrv_flag", "0", "是否开启获取平安存管查询转账充值成功明细服务")

	g_out_patrans_suc_url     = flag.String("patrans_suc_url", "http://172.16.1.16:8209/", "平安存管查询转账充值成功明细")
	g_out_patrans_suc_timeout = flag.String("patrans_suc_timeout", "1", "超时时间 单位秒")

	g_out_patrans_exhange    = flag.String("patrans_exhange", "depository.direct", "平安存管查询转账充值成功明细-发送事件exchange")
	g_out_patrans_routingkey = flag.String("patrans_routingkey", "offline.trans", "平安存管查询转账充值成功明细-发送事件routingkey")
	g_out_patrans_prdid      = flag.String("patrans_prdid", "235", "平安存管查询转账充值成功明细-发送事件-余额转账套餐id")

	g_out_patrans_gentime_interval = flag.String("patrans_gentime_interval", "20", "平安存管查询转账充值成功明细-获取明细时间间隔")
	g_out_patrans_begin_date       = flag.String("patrans_begin_date", "", "平安存管查询转账充值成功明细-开始日期")
)

var zkr *sdk.RegistryZk
var pmutex *common.Mutex

type reqJson struct {
	ICmd     interface{} `json:"cmd"`
	JsonData string
}

type jmfJson struct {
	Service string `json:"service"`
	Env     string `json:"env"`
	Set     string `json:"set"`
	Group   string `json:"group"`
	Version string `json:"version"`
	Params  string `json:"params"`
}

func (r *reqJson) Cmd() string {
	if nil == r.ICmd {
		return "00000000"
	}
	switch reflect.TypeOf(r.ICmd).Kind() {
	case reflect.Float64, reflect.Float32:
		return fmt.Sprintf("%8.0f", r.ICmd)
	default:
		return fmt.Sprint(r.ICmd)
	}
}

func (r *reqJson) ToString() string {
	jsonstr, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSpace(string(jsonstr))
}

type respJson struct {
	errcode.ErrCode
	Data interface{} `json:"data"`
}

func (r *respJson) ToString() string {
	jsonstr, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSpace(string(jsonstr))
}

func init() {
	log.Println("======================================== begin init ========================================")
	InitSignalHandle()
	InitSDK()
	InitConfig()
	initCmd()
	initDB()
	InitRedis()
	InitLog()

	log.Println("======================================== end init ========================================")
}

func main() {

	myAddr := ":" + *g_srvport

	fmt.Println(myAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := http.ListenAndServe(myAddr, mux)
	if err != nil {
		//服务启动失败
		fmt.Println("http server start failed :", err)
	}

}

func handler(rw http.ResponseWriter, req *http.Request) {

	var reqj reqJson
	var respj respJson
	var err error

	head := req.Header

	frameType := head.Get("frame-type")

	//获取开始时间
	kstart := time.Now().UnixNano()

	if strings.Contains(frameType, "JMF") {
		err = DecodeParamsJMF(req, &reqj)
	} else {
		err = DecodeParamsCmd(req, &reqj)
	}

	if err != nil {
		respj.Data = fmt.Sprint(err)
		respj.SetErrCode(errcode.ERRCODE_PARAM_INVAILD)

	} else {
		var data interface{}
		var ecode int

		data, ecode = DoCmd(reqj)

		respj.Data = data
		respj.SetErrCode(ecode)
	}

	kend := time.Now().UnixNano()
	kcost := (kend - kstart) / 1000000 // 达到毫秒级别

	log.Println("cmd:"+reqj.Cmd(), "request:"+strings.TrimSpace(reqj.ToString()), "response:"+respj.ToString(), "cost:"+strconv.FormatInt(kcost, 10))
	common.Log_profile(reqj.Cmd(), kcost, respj.ErrCode.Code, respj.ErrCode.Msg, "request", reqj.ToString(), "response", respj.ToString())
	io.WriteString(rw, respj.ToString())
}

func initDB() {

	dbcon := *g_mysquser + ":" + *g_mysqpasswd + "@tcp(" + *g_mysqhost + ")/" + *g_mysqdb + "?charset=utf8"
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dbcon)

	log.Println(dbcon)
	log.Println("init db success.", "dbcon:"+dbcon)
}

func InitLog() {

	common.Set_host(*g_loghost)
	common.Log_init(*g_logtitle)
	log.Println("init log success. ", "log_host:"+*g_loghost, "log_title:"+*g_logtitle)
}

func InitConfig() {
	if len(os.Args) < 2 {
		log.Println("please set conf file ")
		return
	}

	conf, err := globalconf.NewWithOptions(&globalconf.Options{
		Filename: os.Args[1],
	})

	if err != nil {
		log.Println("NewWithFilename ", os.Args[1], " fail :", err)
		os.Exit(0)
	}

	conf.ParseAll()
	log.Println("init config success. ", "NewWithFilename ", os.Args[1], " succ\n")
}

func RegistJMF(app, iname, port string) (string, error) {

	jmfPath := *g_jmfpathbase + "." + app + ".I" + iname

	si := sdk.ServiceInfo{
		App:       app,
		Env:       *g_env,
		Interface: jmfPath,
		Lang:      "go",
		Owner:     *g_owner,
		Set:       *g_set,
		Ver:       *g_verion,
		Group:     *g_group,
		TimeStamp: fmt.Sprint(time.Now().Unix()),
		Dynamic:   true,
	}

	reg := sdk.RegInfo{
		Scheme: "rest",
		Port:   *g_srvport,
		Path:   jmfPath,
		Params: si,
	}

	if err := zkr.Register(reg); err != nil {
		common.Print("new registryzk failed, err:", err)
		return "", err
	}

	return jmfPath, nil

}

func InitSDK() {

	var err error

	zkr, err = sdk.NewRegistryZk()

	if err != nil {
		common.Print("new registryzk failed, err:", err)
		os.Exit(0)
	}

	if nil == zkr {
		common.Print("new registryzk failed")
		os.Exit(0)
	}

}

func DecodeParamsJMF(req *http.Request, reqj *reqJson) error {

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	d := json.NewDecoder(strings.NewReader(string(body)))
	d.UseNumber()

	var jmfj jmfJson
	if err := d.Decode(&jmfj); err != nil {
		return err
	} else {
		reqj.ICmd = jmfj.Service
		reqj.JsonData = jmfj.Params
	}
	return nil
}

func DecodeParamsCmd(req *http.Request, reqj *reqJson) error {

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	d := json.NewDecoder(strings.NewReader(string(body)))
	d.UseNumber()

	if err := d.Decode(&reqj); err != nil {
		return err

	} else {
		reqj.JsonData = string(body)
	}

	return nil
}

func InitSignalHandle() {

	go func() {
		for {
			ch := make(chan os.Signal)

			signal.Notify(ch, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)
			sig := <-ch
			common.Print("Signal received:", sig, " \n")
			switch sig {
			case syscall.SIGHUP:
				common.Print("get sighup\n")
			case syscall.SIGINT:
				os.Exit(0)
			case syscall.SIGUSR1:
				common.Print("usr1\n")
			case syscall.SIGUSR2:
				os.Exit(0)
			default:
				common.Print("get sig:", sig, "\n")
			}
		}
	}()

}

func initMutex() {
	p := common.NewRedisMuxtex()

	fmt.Println("pool:", p)

	nodes := []common.Node{
		p,
	}

	pmutex = common.NewMutex(GLOBLE_PMUTEX_KEY, nodes)
}

func InitRedis() {
	common.InitRedis(*g_redisaddr)
	initMutex()

	log.Println("init redis success. ", "redis:"+*g_redisaddr)
}
