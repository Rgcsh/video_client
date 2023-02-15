// (C) Guangcai Ren <rgc@bvrft.com>
// All rights reserved
// create time '2023/2/14 15:51'
//
// Usage:
//

package constant

import "net"

//video_server服务的地址
const VideoServer string = "http://127.0.0.1:7069"

//udp服务器的ip地址
var UdpServerIp = net.IPv4(0, 0, 0, 0)

//udp服务器的端口
const UdpServerPort int = 9090
