package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestInsert(t *testing.T) {
	admin := Admin{
		Account:  "2019001",
		Password: "123456",
	}

	admin.Insert(db)
}

func TestQuery(t *testing.T) {
	admin := Admin{}
	admin.Query(db, "id", 1)
	fmt.Printf("id: %d, account: %s, password: %s\n", admin.Model.ID, admin.Account, admin.Password)
}

func TestQueryByAccount(t *testing.T) {
	admin := Admin{}
	admin.Query(db, ACCOUNT, "admin")
	fmt.Printf("id: %d, account: %s, password: %s\n", admin.Model.ID, admin.Account, admin.Password)
}

func TestUpdate(t *testing.T) {
	admin := Admin{
		Account:  "2019001",
		Password: "abcdefg",
	}
	admin.Update(db, ACCOUNT)
	updateAfterAdmin := Admin{}
	updateAfterAdmin.Query(db, "id", admin.Model.ID)
	fmt.Printf(
		"id: %d, account: %s, password: %s\n",
		updateAfterAdmin.ID, updateAfterAdmin.Account, updateAfterAdmin.Password)
}

func TestDelete(t *testing.T) {
	admin := Admin{}
	admin.Query(db, "id", 1)
	admin.Delete(db)
}

func TestInsertAdmin(t *testing.T) {
	admin := Admin{
		Account:  "admin",
		Password: "LXY1019XYXYZ",
	}

	admin.Insert(db)
}

//表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

//表情转换
func UnicodeEmojiEncode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

// test print msg to terminal take custom color.
func TestMsgColor(t *testing.T) {
	fmt.Printf("\x1b[%dm hello world 32: 绿 \x1b[0m\n", 32)
	fmt.Printf("%s Rocket!\n", UnicodeEmojiDecode("🚀"))
}
