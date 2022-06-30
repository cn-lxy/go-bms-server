package utils

import "fmt"

// 打印有颜色的字符
const (
	Green  = 32
	Yellow = 33
	Red    = 31
)

func PrintfColorStr(color int, str string) {
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", color, str)
}
