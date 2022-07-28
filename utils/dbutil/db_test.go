package dbutil

import (
	"fmt"
	"testing"

	"github.com/cn-lxy/bms_go/utils"
)

func TestDbQuery(t *testing.T) {
	// common test
	utils.PrintfColorStr(utils.Yellow, "TestDbQuery")
	var sql1 string = "select * from users"
	res1 := Query(sql1)
	for i, d := range res1 {
		fmt.Printf("%d: %#v\n", i, d)
	}

	// mutilative args test
	sql2 := "select * from users where id = ? and sex = ?"
	res2 := Query(sql2, 1, "男")
	if len(res2) == 0 {
		t.Fatal(fmt.Errorf("query err"))
	}
	fmt.Printf("%v\n", res2[0])

	// args slice test
	sql3 := "select * from `users` where id = ? and account = ? and sex = ?"
	args := []any{1, "2019136417", "男"}
	res3 := Query(sql3, args...)
	if len(res3) == 0 {
		t.Fatal(fmt.Errorf("query err"))
	}
	fmt.Printf("%v\n", res3[0])
}

func TestDbInsert(t *testing.T) {
	sql := "insert into `users` (name, account, password, sex, college, birthday, register) values (?, ?, ?, ?, ?, ?, now())"
	args := make([]any, 0)
	// name, account, password, sex, college, birthday
	args = append(args, "lxy", "2788", "123456", "男", "信通...", "2019-10-14")
	if err := Update(sql, args...); err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, "insert success")
}

func TestDbUpdate(t *testing.T) {
	sql := "update `users` set name = ? where name = ?"
	args := make([]any, 0)
	args = append(args, "lxy-2", "lxy")
	if err := Update(sql, args...); err != nil {
		t.Fatal(err)
	}
}

func TestDbDelete(t *testing.T) {
	sql := "delete from `users` where name = ?"
	args := make([]any, 0)
	args = append(args, "lxy-2")
	if err := Update(sql, args...); err != nil {
		t.Fatal(err)
	}
}
