package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()                       // 创建应用程序实例
	window := app.NewWindow("Hello world") // 创建窗口，标题为"Hello Wolrd"

	window.SetContent(widget.NewLabel("Hello world!")) // 往窗口中放入一个内容为"Hello world!"的标签控件
	window.ShowAndRun()                                //展示并运行程序
}
