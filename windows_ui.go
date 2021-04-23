package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"os"
	"path"
)

func openUI() {
	dir := path.Dir(os.Args[0])
	Run(dir, "cmd", "/c", "start", fmt.Sprintf("http://localhost:%d/", port))
}

func openHome() {
	dir := path.Dir(os.Args[0])
	Run(dir, "cmd", "/c", "start", "https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi/")
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Kaf-wifi")
	systray.SetTooltip("kaf-cli")
	mOpen := systray.AddMenuItem("打开界面", "打开浏览器界面")
	mHome := systray.AddMenuItem("打开官网", "打开kaf-cli发布界面")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出kaf-cli")
	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				openUI()
			case <-mHome.ClickedCh:
				openHome()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	// clean up here
	os.Exit(0)
}
