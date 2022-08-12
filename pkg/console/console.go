package console

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"
)

// Success 打印一条成功消息，绿色输出
func Success(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	colorOut(msg, "green")
}

// Error 打印一条报错消息，红色输出
func Error(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	colorOut(msg, "red")
}

// Warning 打印一条提示消息，黄色输出
func Warning(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	colorOut(msg, "yellow")
}

// Exit 打印一条报错消息，并退出 os.Exit(1)
func Exit(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	Error(msg)
	os.Exit(1)
}

// ExitIf 语法糖，自带 err != nil 判断
func ExitIf(err error) {
	if err != nil {
		Exit(err.Error())
	}
}

// colorOut 内部使用，设置高亮颜色
func colorOut(message, color string) {
	fmt.Fprintln(os.Stdout, ansi.Color(message, color))
}
