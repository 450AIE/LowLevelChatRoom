package main

import (
	"chater/ChatProject/Client/utils"
	"chater/ChatProject/message"
	"chater/ChatProject/server/model"
	"chater/ChatProject/server/process"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"net"
	"time"
)

var pool *redis.Pool

type Processor struct { //总控制器，从main拿到最初的数据，再判断类型，交由不同函数处理
	Conn net.Conn
}

func (p *Processor) ServerProcess(mes *message.Message) (err error) { //这里接收到的都是反序列化了的mes
	switch mes.Type {
	case message.LoginType:
		var user = &process.UserProcessor{Conn: p.Conn} //这里就有分层的思想，Processor拿到数据后，创建UserProcessor，交给他处理
		err = user.LoginProcess(mes)
		if err != nil {
			return
		}
	case message.ResgisterType:
		var user = &process.UserProcessor{Conn: p.Conn}
		err = user.RegisterProcess(mes)
		if err != nil {
			return
		}
	case message.ShortMesType:
		var user = &process.UserProcessor{Conn: p.Conn}
		user.ShortMessageProcess(mes)
	case message.ExitType:
		var user = &process.UserProcessor{Conn: p.Conn}
		user.ExitProcess(mes)
	default: //虽然应该不会发生在这里，但是预防万一
		return errors.New("错误的消息类型,无法处理...")
	}
	return nil
}
func (p *Processor) Process() (err error) {
	defer p.Conn.Close()
	for { //循环阻塞在这里，以接收来自同一个客户端的间隔发送的信息
		var bytes [8096]byte
		var tf = &utils.Transfer{p.Conn, bytes}
		mes, err := tf.ReadPkg() //mes已经被反序列化了，可以直接取出type进行下一步业务逻辑判断
		if err != nil {
			if err == io.EOF { //原本read发现用户已经退出，则error为io.EOF
				//fmt.Println("用户退出")
				return err
			} else {
				return err
			}
		}
		err = p.ServerProcess(&mes)
		if err != nil {
			fmt.Println("主控业务判断错误")
		}
		//fmt.Println("成功")
	}
} //	fmt.Printf("用户%v说了%v个字节的数据："+string(buf[:n])+"\n", conn.RemoteAddr().String(), n)

func initPool(address string, maxconn, maxactive int, maxsparetime time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxconn,
		MaxActive:   maxactive,
		IdleTimeout: maxsparetime,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}

func init() {
	initPool("localhost:6379", 50, 40, 300*time.Second) //启动就创建连接池
	model.SharedUserDao = model.NewUserDao(pool)
}
func main() { //main的作用：监听端口，拿到链接，启动协程，交给Processor处理
	listener, err := net.Listen("tcp", "127.0.0.1:50000")
	defer listener.Close()
	if err != nil {
		fmt.Println("监听失败")
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("连接失败")
			//continue
		} else {
			fmt.Printf("用户%v已链接\n", conn.LocalAddr().String())
			var p = Processor{conn}
			//go p.Process()这样下是无法定义变量err接收的，要用闭包
			go func() {
				err := p.Process()
				if err != nil {
					fmt.Println(err.Error())
					return //关闭的是协程
				}
			}()
		}
	}
}
