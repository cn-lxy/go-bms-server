package dbutil

import (
	"fmt"
	"log"
	"testing"

	"github.com/cn-lxy/bms_go/utils"
)

func TestDbQuery(t *testing.T) {
	utils.PrintfColorStr(utils.Yellow, "TestDbQuery")
	var sqlString string = "select * from users"
	data := Query(sqlString)
	for i, d := range data {
		fmt.Printf("%d: %#v\n", i, d)
	}
}

func TestDbInsert(t *testing.T) {
	utils.PrintfColorStr(utils.Yellow, "TestDbInsert")
	sql := fmt.Sprintf("insert into `users` (name, account, password, sex, college, birthday, register) values ('%s', '%s', '%s', '%s', '%s', '%s', now())",
		"lxy",
		"2788311915",
		"123456",
		"男",
		"信通",
		"2019-10-14")
	if err := Update(sql); err != nil {
		log.Fatal(err)
	}
}

func TestDbUpdate(t *testing.T) {
	utils.PrintfColorStr(utils.Yellow, "TestDbUpdate")
	sql := fmt.Sprintf("update `users` set name = '%s' where name = '%s'",
		"lxy-2", "lxy")
	if err := Update(sql); err != nil {
		log.Fatal(err.Error())
	}
}

func TestDbDelete(t *testing.T) {
	utils.PrintfColorStr(utils.Yellow, "TestDbDelete")
	sql := fmt.Sprintf("delete from `users` where name = '%s'", "lxy-2")
	if err := Update(sql); err != nil {
		log.Fatal(err.Error())
	}
}
