package process

import (
	"chater/ChatProject/Client/utils"
	"chater/ChatProject/message"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"net"
)

//var (
//	Connn net.Conn
//)

type UserProcess struct {
}

// 不知道sizeof可不可以判断字节数
// 登陆defer close吗？close了不是又断开链接了嘛，应该不断，在main里面close
func (this *UserProcess) Login(userID string, userPWD string, address string) (err error) { //尽量以error返回，因为这样可以返回更详细的错误，bool太笼统了
	conn, err := net.Dial("tcp", address)
	defer conn.Close() //这个函数必须阻塞，不然函数结束，这一句defer被执行，那么保持通讯的go的conn被close，就无法保护通讯
	if err != nil {
		fmt.Println("链接失败")
		return
	}
	var loginmes message.LoginMessage
	loginmes.UserID = userID
	loginmes.UserPwd = userPWD
	var mes message.Message
	mes.Type = message.LoginType
	dataput, err := json.Marshal(loginmes)
	if err != nil {
		fmt.Println("序列化失败")
		return errors.New("序列化失败")
	}
	mes.Data = string(dataput)
	var bytes [8096]byte
	var tf = utils.Transfer{conn, bytes}
	err = tf.WritePkg(&mes)
	if err != nil {
		fmt.Println("数据传输错误4")
		return
	}
	//ERROR
	resmes, err := tf.ReadPkg() //接收返回的ResLogin结构体
	if err != nil {
		fmt.Println("客户端接收返回的登陆信息结构体错误")
		return
	}
	//ERROR
	fmt.Println("成功读取到服务端发送的登陆成功与否的结构体")
	var loginInspect message.ResMessage //该结构体含有在线用户列表
	err = json.Unmarshal([]byte(resmes.Data), &loginInspect)
	if err != nil {
		fmt.Println("反序列化错误")
		return
	}
	if loginInspect.Code == "200" {
		go KeepOnline(conn)
		fmt.Printf("用户：%v  登录成功\n", userID) //这里开始是登录成功
		//保存当前用户信息
		CurUser.Conn = conn //ERROR
		CurUser.UserID = userID
		CurUser.UserPwd = userPWD
		CurUser.UserStatus = message.UserOnline
		//
		//显示在线用户ID列表
		//fmt.Println("当前在线用户ID：")
		color.Yellow("当前在线用户ID：\n")
		for _, v := range loginInspect.AllUsers {
			//fmt.Printf("用户ID：%v\n", v)
			color.Yellow("用户ID：%v\n", v)
			SharedOnlineUsersList[v] = message.UserOnline
		}
		ShowMenu()
		//显示在线用户ID列表
	} else if loginInspect.Code == "500" {
		return errors.New("此用户不存在")
	}
	return nil
}
func (this *UserProcess) Register(userID string, userPWD string, address string) (err error) { //其实好像可以加一个判断符1/2，变成一个函数
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("注册链接服务器错误")
		return
	}
	defer conn.Close()
	var regis = &message.RegisterMessage{User: message.User{UserID: userID, UserPwd: userPWD}}
	var mes = &message.Message{}
	mes.Type = message.ResgisterType
	data, err := json.Marshal(regis)

	mes.Data = string(data)
	var bytes [8096]byte
	var tf = utils.Transfer{conn, bytes}

	err = tf.WritePkg(mes)
	if err != nil {
		fmt.Println("传输Register结构体发生错误")
		return
	}
	resmes, err := tf.ReadPkg()
	if err != nil {
		fmt.Println("读取ResRegister结构体发生错误")
		return
	}
	var resRegister = &message.ResRegisterMessage{}
	err = json.Unmarshal([]byte(resmes.Data), resRegister)
	if err != nil {
		fmt.Println("反序列化失败")
		return
	}
	if resRegister.Code == "200" {
		fmt.Println("用户已存在")
		return
	}
	return
}
