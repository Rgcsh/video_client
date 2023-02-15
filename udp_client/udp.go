// 
// All rights reserved
// create time '2022/12/9 16:56'
//
// Usage:
//

package udp_client

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"net"
	"video_client/utils/constant"
)

var UDPListen *net.UDPConn

//
//  @Description: 启动udp客户端
//
func UDPClientSetUp() {
	var err error
	UDPListen, err = net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   constant.UdpServerIp,
		Port: constant.UdpServerPort,
	})
	if err != nil {
		panic(fmt.Sprintf("监听UDP服务失败,请修改,err:", err))
	}
	fmt.Println(fmt.Sprintf("连接UDP服务成功,端口为:%v", constant.UdpServerPort))
	return
}

//
//  @Description: 收尾性工作;关闭udp client
//
func UDPClose() {
	err := UDPListen.Close()
	if err != nil {
		panic(fmt.Sprintf("关闭UDPListen失败,err:%v", err))
	}
}

//
//  @Description: 接收UDP服务端的图片数据,并展示给gui(直播使用)
//  @param Img:
//
func UDPReceive(Img *canvas.Image) {
	for {
		//定义一个长度为10万的切片
		sliceData := make([]byte, 100000, 100000)
		//将读取的数据存到数组中
		n, err := UDPListen.Read(sliceData)
		if err != nil {
			fmt.Println("读取UDP数据失败,err:", err)
			continue
		}
		if n == 0 {
			continue
		}
		//将读取的视频帧数据 以图片方式展示出来
		res := &fyne.StaticResource{
			StaticName:    "test",
			StaticContent: sliceData,
		}
		Img.Resource = res
		Img.Refresh()
	}
}

//
//  @Description: 接收channel中的message 信号数据,并发送数据给服务端
//  @param MqttMessageChannel:
//
func UDPSend(MqttMessageChannel <-chan []byte) {
	for {
		message, ok := <-MqttMessageChannel
		if !ok {
			fmt.Println("发送数据到图片处理channel关闭")
			break
		}
		b, err := UDPListen.Write(message)
		if err != nil {
			panic(fmt.Sprintf("发送UDP数据失败 err:%v", err))
		}
		fmt.Println("发送数据 %v", b)
	}
}
