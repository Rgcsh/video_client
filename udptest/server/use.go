package main

import (
	"fmt"
	"net"
)

// UDP Server端
func main() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 30000,
	})
	if err != nil {
		fmt.Println("Listen failed, err: ", err)
		return
	}
	defer listen.Close()
	for {
		//定义一个长度为10万的切片
		data := make([]byte, 30000, 30000)
		n, addr, err := listen.ReadFromUDP(data) // 接收数据
		if err != nil {
			fmt.Println("read udp failed, err: ", err)
			continue
		}
		fmt.Println("data:%v addr:%v count:%v\n", string(data[:n]), addr, n)
		_, err = listen.WriteToUDP(data, addr) // 发送数据
		if err != nil {
			fmt.Println("Write to udp failed, err: ", err)
			continue
		}
	}
}
