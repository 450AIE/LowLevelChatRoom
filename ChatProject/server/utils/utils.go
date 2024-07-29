package utils

import (
	"chater/ChatProject/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte //这样就避免每次都插创建,反正Read和Write都是覆盖，不用担心数据重复
}

func (tf *Transfer) WritePkg(mes *message.Message) (err error) { //传递没有序列化的message，但是里面的date已经序列化了
	datasend, err := json.Marshal(*mes)
	if err != nil {
		fmt.Println("序列化失败")
		return
	}
	var sendlen uint32
	sendlen = uint32(len(datasend))
	//fmt.Println("预发送长度sendlen=" + strconv.Itoa(int(sendlen)))
	//var bytes [4]byte
	binary.BigEndian.PutUint32(tf.Buf[:4], sendlen)
	n, err := tf.Conn.Write(tf.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("数据传输错误1")
		return
	} else {
		//fmt.Println("数据传输成功1")
	}
	n2, err := tf.Conn.Write(datasend)
	//fmt.Println("发送长度n2=" + strconv.Itoa(int(n2)))
	if uint32(n2) != sendlen || err != nil {
		fmt.Println("数据传输错误3")
		return
	} else {
		//fmt.Println("数据传输成功3")
	}
	return nil
}
func (tf *Transfer) ReadPkg() (mes message.Message, err error) {
	//buf := make([]byte, 8096)
	fmt.Println("正在读取客户端的信息.....")
	n, err := tf.Conn.Read(tf.Buf[:4]) //先接受包的长度，判断有无丢包,原本read发现用户已经退出，则error为io.EOF
	//fmt.Println("预读取长度Pkglen=" + strconv.Itoa(int(n)))
	if err == io.EOF { //只有阻塞在这边，且客户端退出时，才会EOF
		return
	}
	if n != 4 || err != nil {
		err = errors.New("读取信息错误1") //无法区分是读取错误还是用户退出
		return
	} else {
		//fmt.Println("读取成功1")
	}
	var Pkglen uint32
	Pkglen = binary.BigEndian.Uint32(tf.Buf[:4]) //把buf[:4]转换为uint32,Pkglen里面是一些奇怪的东西，无法理解
	//fmt.Println("读取长度Pkglen=" + strconv.Itoa(int(Pkglen)))
	n2, err := tf.Conn.Read(tf.Buf[:Pkglen]) //正确后再读取
	if uint32(n2) != Pkglen || err != nil {
		errors.New("读取信息错误2")
		return
	}
	err = json.Unmarshal(tf.Buf[:Pkglen], &mes) //反序列化
	if err != nil {
		errors.New("反序列化失败")
		return
	}
	err = nil
	return
}
