package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"git.jtjr.com/bg_prd_grp/intrapay/common"
	"git.jtjr.com/bg_prd_grp/intrapay/errcode"
)

var cmdMap map[string]CmdHandle

type cmdFunc func(interface{}) (interface{}, int)

type CmdHandle struct {
	cfunc   cmdFunc
	cparam  interface{}
	blog    bool
	jmfName string
}

func (c *CmdHandle) Do(reqj reqJson) (interface{}, int) {

	common.Log_info("request :", fmt.Sprint(reqj))

	var tparams map[string]interface{}

	d := json.NewDecoder(strings.NewReader(reqj.JsonData))
	d.UseNumber()

	if err := d.Decode(&tparams); err != nil {
		common.Log_error(errcode.ERRCODE_PARAM_INVAILD, "Init Param err", fmt.Sprint(err))
		return nil, errcode.ERRCODE_PARAM_INVAILD
	}

	t := reflect.ValueOf(c.cparam).Type()

	result := reflect.New(t)

	if err := common.ParamsSetData(result, tparams); err != nil {
		return fmt.Sprint(err), errcode.ERRCODE_PARAM_INVAILD
	}

	funcParams := result.Interface()

	cresult, errc := c.cfunc(funcParams)

	if c.blog {

		str_result, _ := json.Marshal(cresult)

		dbr := DBRecord{
			Data:   funcParams,
			Status: errc / errcode.ErrCodeBase,
			Result: string(str_result),
		}

		if !dbr.Create() {
			common.Print("dbr failed")
		}
	}

	return cresult, errc

}

func DoCmd(reqj reqJson) (interface{}, int) {
	if cmddo, found := cmdMap[reqj.Cmd()]; found {
		return cmddo.Do(reqj)
	}
	return nil, errcode.ERRCODE_CMD_INVAILD
}

func RegistCmd(route string, chandle CmdHandle) error {

	if len(route) <= 0 {
		return errors.New("empty route key")
	}

	if cmdMap == nil {
		cmdMap = map[string]CmdHandle{}
	}

	if _, found := cmdMap[route]; found {
		return errors.New("duplicate route key")
	}
	cmdMap[route] = chandle

	return nil
}

func RegistJMFCmd(app, iname string, chandle CmdHandle) error {

	if len(iname) <= 0 {
		return errors.New("empty route key")
	}

	if jmfInteface, err := RegistJMF(app, iname, *g_srvport); err != nil {
		return err
	} else {
		if err = RegistCmd(jmfInteface, chandle); err != nil {
			//UnRegistCmd
			return err
		}
		return nil
	}

}
