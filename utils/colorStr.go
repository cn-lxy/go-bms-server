package utils

import "fmt"

// 打印有颜色的字符
const (
	GREEN = 32
)

func PrintfColorStr(color int, str string) {
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", color, str)
}
