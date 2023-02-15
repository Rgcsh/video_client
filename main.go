package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/flopp/go-findfont"
	"github.com/goki/freetype/truetype"
	"os"
	"time"
	"video_client/gui_service"
	"video_client/udp_client"
)

//
//  @Description: main函数的defer操作,负责 udp服务关闭 及 通知udp服务端 停止发送视频帧 的信号
//  @param MqttMessageChannel:
//
func MainDefer(MqttMessageChannel chan []byte) {
	defer udp_client.UDPClose()
	defer func() {
		fmt.Println("程序退出前,发送停止直播信号")
		//	连接udp服务,接收视频帧数据,并播放
		MqttMessageChannel <- []byte("clientStop")
		time.Sleep(1 * time.Second)
	}()
}

//
//  @Description: 初始化,功能为 让fyne框架支持中文显示
//
func init() {
	fontPath, err := findfont.Find("Songti.ttc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found 'arial.ttf' in '%s'\n", fontPath)

	// load the font with the freetype library
	// 原作者使用的ioutil.ReadFile已经弃用
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}
	_, err = truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	_ = os.Setenv("FYNE_FONT", fontPath)
}

//
//  @Description: UDP相关配置运行
//  @param MqttMessageChannel:
//  @param img:
//
func SetUp(MqttMessageChannel chan []byte, img *canvas.Image) {
	//udp客户端启动并接收数据
	udp_client.UDPClientSetUp()
	go udp_client.UDPSend(MqttMessageChannel)
	go udp_client.UDPReceive(img)
}

//
//  @Description:
//
func main() {
	//GUI初始化
	myApp := app.New()
	myWindow := myApp.NewWindow("Video Client")

	//第一行 视频播放窗口
	img, one := gui_service.FirstLineButton()

	//第二行 直播/停止直播 按钮
	MqttMessageChannel := make(chan []byte, 20)
	SetUp(MqttMessageChannel, img)
	defer MainDefer(MqttMessageChannel)
	two := gui_service.SecondLineButton(MqttMessageChannel)

	//第三行 输入框 和 筛选按钮 放在一个子 水平布局中
	startTimeEntry, endTimeEntry, three := gui_service.ThirdLineButton()

	//第四行,功能是 选择文件夹及下载按钮
	four := gui_service.FourthLineButton(myWindow, startTimeEntry, endTimeEntry)

	//将四行放到一个window中
	all := container.NewVBox(one, two, three, four)
	myWindow.SetContent(all)

	//fyne gui框架运行及展示
	myWindow.ShowAndRun()
}
