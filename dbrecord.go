package main

import (
	"fmt"
	"reflect"

	"github.com/astaxie/beego/orm"
	//_ "github.com/go-sql-driver/mysql"
)

type DBRecord struct {
	Data   interface{}
	Status int
	Result string
}

func (d *DBRecord) Create() bool {

	recordMap := make(map[string]string)

	tblName := getKeyAndValue(reflect.ValueOf(d.Data), &recordMap)

	if len(tblName) <= 0 {
		fmt.Println("log table is NULL")
		return false
	}

	fmt.Println("map:", recordMap)

	sqlK := " ("
	sqlV := " ("

	for k, v := range recordMap {
		sqlK += k + ","
		sqlV += "'" + v + "',"
	}

	sqlK += "status,result)"
	sqlV += fmt.Sprint(d.Status) + ",'" + d.Result + "')"

	sql := "insert into " + tblName + sqlK + " values " + sqlV

	fmt.Println("dbr sql:", sql)

	o := orm.NewOrm()
	if _, err := o.Raw(sql).Exec(); err != nil {
		fmt.Println("create DBRecords db err:", err)
		return false
	}
	return true

}

func getKeyAndValue(data reflect.Value, rmap *map[string]string) string {

	value := data.Elem()      //获取value对象
	theType := value.Type()   //获取value的类型
	theKind := theType.Kind() //获取类型的分类说

	var tblName string

	if theKind == reflect.Struct {
		for i := 0; i < theType.NumField(); i++ {

			if theType.Field(i).Type.Kind() == reflect.Struct {
				ts := value.Field(i).Addr()
				getKeyAndValue(ts, rmap)
				//目前只有自己划拨会有数组结构，但支持单个记录，因此做如下处理
			} else if theType.Field(i).Type.Kind() == reflect.Slice {
				ts := value.Field(i).Index(0).Addr()
				getKeyAndValue(ts, rmap)
			} else {
				tagName := theType.Field(i).Tag.Get("dbr") //json里面的key
				if len(tagName) > 0 {
					(*rmap)[tagName] = fmt.Sprint(value.Field(i).Interface())
				} else {
					tn := theType.Field(i).Tag.Get("dbr_t")
					if len(tn) > 0 {
						tblName = tn
					}
				}
			}

		}
	}
	return tblName
}
