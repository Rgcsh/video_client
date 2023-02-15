// 
// All rights reserved
// create time '2023/2/14 16:50'
//
// Usage:
//

package gui_service

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

//
//  @Description: 包裹所有错误,防止程序异常退出
//  @param window:
//
func RecoverAllPanic(window fyne.Window) {
	err := recover()
	if err != nil {
		dialog.ShowInformation("程序异常", fmt.Sprintf("%v", err), window)
	}
}
