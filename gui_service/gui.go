// (C) Guangcai Ren <rgc@bvrft.com>
// All rights reserved
// create time '2023/2/14 16:48'
//
// Usage:
//

package gui_service

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/asmcos/requests"
	"video_client/utils/constant"
	"video_client/utils/gtime"
	"video_client/utils/requests_handler"
)

//
//  @Description: 第一行 视频播放窗口
//  @return *canvas.Image:
//  @return *fyne.Container:
//
func FirstLineButton() (*canvas.Image, *fyne.Container) {
	img := canvas.NewImageFromResource(theme.FyneLogo())
	img.FillMode = canvas.ImageFillOriginal
	first := container.NewGridWithColumns(1, img)
	return img, first
}

//
//  @Description: 第二行 直播/停止直播按钮及操作
//  @param MqttMessageChannel:
//  @param img:
//  @return *fyne.Container:
//
func SecondLineButton(MqttMessageChannel chan []byte) *fyne.Container {
	//第二行 监控指标按钮
	playButton := widget.NewButton("开始直播", func() {
		fmt.Println("点击直播按钮了")
		//	连接udp服务,接收视频帧数据,并播放
		MqttMessageChannel <- []byte("clientPlay")
	})
	playButton.Icon = theme.ConfirmIcon()
	stopButton := widget.NewButton("停止直播", func() {
		fmt.Println("点击停止直播按钮了")
		//	连接udp服务,接收视频帧数据,并播放
		MqttMessageChannel <- []byte("clientStop")
	})
	stopButton.Icon = theme.CancelIcon()
	return container.NewGridWithColumns(2, playButton, stopButton)
}

//
//  @Description: 第三行 输入框 和 筛选按钮及相关操作函数
//  @return *widget.Entry:
//  @return *widget.Entry:
//  @return *fyne.Container:
//
func ThirdLineButton() (*widget.Entry, *widget.Entry, *fyne.Container) {
	startTimeEntry := widget.NewEntry()
	startTimeEntry.SetPlaceHolder("输入开始时间")
	startTimeEntry.SetText(gtime.GetCurrentTime())
	endTimeEntry := widget.NewEntry()
	endTimeEntry.SetPlaceHolder("输入结束时间")
	three := container.NewGridWithColumns(2, startTimeEntry, endTimeEntry)
	return startTimeEntry, endTimeEntry, three
}

//
//  @Description: 第四行 选择文件夹及下载按钮 和对应函数操作
//  @param myWindow:
//  @param startTimeEntry:
//  @param endTimeEntry:
//  @return *fyne.Container:
//
func FourthLineButton(myWindow fyne.Window, startTimeEntry *widget.Entry, endTimeEntry *widget.Entry) *fyne.Container {
	var DownloadDir string
	HistoryVideoDownloadFunc := func() {
		defer RecoverAllPanic(myWindow)
		//根据开始/结束时间 查看时间范围内的 文件名列表
		resp, err := requests.PostJson(fmt.Sprintf("%v/video/query", constant.VideoServer), map[string]interface{}{"start_time": startTimeEntry.Text, "end_time": endTimeEntry.Text})
		if err != nil {
			dialog.ShowInformation("request error", fmt.Sprintf("%v", err), myWindow)
			return
		}

		respData, err := requests_handler.ResponseCheck(resp)
		if err != nil {
			dialog.ShowInformation("response check error", fmt.Sprintf("%v", err), myWindow)
			return
		}

		if respData["result"] != nil {
			files := respData["result"].([]interface{})
			//	开始下载
			for _, file := range files {
				x, err := json.Marshal(map[string]interface{}{"FileName": file})
				if err != nil {
					dialog.ShowInformation("json marshal error", fmt.Sprintf("%v", err), myWindow)
					return
				}
				//下载文件
				resp, err := requests.Get(fmt.Sprintf("%v/video/download?FileName=%v", constant.VideoServer, file), x)
				if err != nil {
					dialog.ShowInformation("download status error", fmt.Sprintf("%v", err), myWindow)
					return
				}
				err = resp.SaveFile(fmt.Sprintf("%v/%v", DownloadDir, file))
				if err != nil {
					dialog.ShowInformation("download error", fmt.Sprintf("%v", err), myWindow)
					return
				}
			}
		}
		dialog.ShowInformation("下载成功", fmt.Sprintf("请在 %v 文件夹中查看下载的文件", DownloadDir), myWindow)
	}

	filterButton := widget.NewButton("历史视频下载", HistoryVideoDownloadFunc)
	filterButton.Icon = theme.ConfirmIcon()

	DirButton := widget.NewButton("选择下载文件夹", func() {
		dialog.ShowFolderOpen(func(closer fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowInformation("open dir error", fmt.Sprintf("%v", err), myWindow)
				return
			}
			if closer != nil {
				DownloadDir = closer.Path()
				dialog.ShowInformation("下载文件夹选择成功", fmt.Sprintf("%v", DownloadDir), myWindow)
			}
		}, myWindow)
	})
	filterButton.Icon = theme.ConfirmIcon()

	return container.NewGridWithColumns(2, DirButton, filterButton)
}
