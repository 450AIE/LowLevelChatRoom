package process

import (
	"chater/ChatProject/message"
	"chater/ChatProject/server/model"
	"chater/ChatProject/server/utils"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"net"
)

type UserProcessor struct {
	Conn net.Conn
	ID   string
	//老师在这里加入了一个UserID便于区分，但是我觉得没必要，因为使用这个UserProcessor结构体的UserMgr里的map也可以拿到ID。想了想，这里不加入ID，那么AddOnlineUser就要传递ID，麻烦
}

func (this *UserProcessor) LoginProcess(mes *message.Message) (err error) { //得到的mes是已经反序列化了的
	var login message.LoginMessage
	err = json.Unmarshal([]byte(mes.Data), &login) //把接收到到LoginMes反序列化提取出来
	if err != nil {
		fmt.Println("反序列化失败")
		return
	}
	this.ID = login.UserID
	var reslogin message.ResMessage //把登陆成功与否的信息传递回客户端
	var resmes message.Message      //没有实例化，为什么可以正常运行？
	resmes.Type = message.LoginResType
	//从数据库判断有无该用户存在,200有,500无
	err = model.SharedUserDao.Login(login.UserID, login.UserPwd)
	if err != nil {
		reslogin.Code = "500"
		reslogin.Error = err.Error()
	} else {
		reslogin.Code = "200"
		reslogin.AllUsers = make([]string, 0) //服务端创造了切片并且make，可以直接传递给客户端使用，客户端就无需make了
		SharedUserMgr.AddOnlineUser(this)
		this.NotifyOthers()
		for k, _ := range SharedUserMgr.AllOnlineUsers { //把所有在线用户ID列表传入切片
			reslogin.AllUsers = append(reslogin.AllUsers, k)
		}
	}
	//并填写Code和Error
	data, err := json.Marshal(reslogin)
	if err != nil {
		fmt.Println("序列化失败")
	}
	var bytes [8096]byte
	resmes.Data = string(data)
	tf := utils.Transfer{Conn: this.Conn, Buf: bytes}
	err = tf.WritePkg(&resmes)
	if err != nil {
		fmt.Println("数据传输错误5")
	}
	//我感觉，可以在一个用户登录成功后，把它的ID装在Message里发给所有人，即可让所有人知道此人登录
	return nil
}
func (this *UserProcessor) RegisterProcess(mes *message.Message) (err error) { //得到的mes是已经反序列化了的
	var register message.RegisterMessage
	err = json.Unmarshal([]byte(mes.Data), &register)
	if err != nil {
		fmt.Println("反序列化失败")
		return
	}
	var resRegister = &message.ResRegisterMessage{}
	var resmes = &message.Message{}
	resmes.Type = message.ResRegisterType
	err = model.SharedUserDao.Register(register.User.UserID, register.User.UserPwd)
	if err != nil {
		fmt.Println(err.Error())
		resRegister.Error = err.Error()
		resRegister.Code = "400"
	} else {
		resRegister.Code = "200"
	}
	data, err := json.Marshal(resRegister)
	if err != nil {
		fmt.Println("序列化失败")
		return
	}
	resmes.Data = string(data)
	var bytes [8096]byte
	var tf = utils.Transfer{this.Conn, bytes}
	err = tf.WritePkg(resmes)
	if err != nil {
		fmt.Println("错误")
	}
	return
}
func (this *UserProcessor) NotifyOthers() { //遍历在线用户ID列表，让每一个用户都通知别人自己在线
	for id, userprocess := range SharedUserMgr.AllOnlineUsers {
		if id == this.ID { //某个人登录了，会调用这个函数，但是客户端有显示了，于是不推送自己？？？？？？
			continue
		}
		userprocess.NotifyMeOnlineDetail(this.ID) //每一个在线用户都调用，都会被告知this用户上线了
	}
}
func (this *UserProcessor) NotifyMeOnlineDetail(ID string) { //调用该函数的用户，会被Write当前在线的用户列表，因为conn是在服务端，所以write是给客户端
	var mes message.Message //这里没有实例化也可以！！！！why
	var userStatue message.NotifyUserStatus
	mes.Type = message.NotifyType
	userStatue.UserStatus = message.UserOnline
	userStatue.UserID = ID //因为是通知“我”上线了，所以ID应该是最开始的调用者，即调用NotifyOthers的this,ID
	data, err := json.Marshal(userStatue)
	if err != nil {
		fmt.Println("序列化错误")
		return
	}
	mes.Data = string(data)
	var bytes [8096]byte
	var tf = &utils.Transfer{this.Conn, bytes}
	err = tf.WritePkg(&mes)
	if err != nil {
		fmt.Println("传输用户在线列表失败")
		return
	}
}
func (this *UserProcessor) ExitProcess(mes *message.Message) {
	//太困了。思路放这里：
	//1.服务端，就是这里，收到mes后，反序列化Data提取出ID，从在线用户列表中删除
	//2.把这个mes反序列化Data提取出的ID,创建一个ResExitMessage发送给每一个在线的用户
	//3.用户收到ResExitMessage后反序列化得到ID，从维护的map中删除这个用户，并打印“%v用户已退出”
	//4.记得实现带空格的输入
	//5.可以把用户发送的消息带上颜色，挺简单的
	var resmes = &message.Message{}
	resmes.Type = message.ResExitType
	var resExit = &message.ResExitMessage{}
	err := json.Unmarshal([]byte(mes.Data), resExit) //把退出用户的ID写入到resExit
	if err != nil {
		fmt.Println("Data反序列化失败")
	}
	this.ID = resExit.UserID
	SharedUserMgr.DelOnlineUser(this)
	resmes.Data = mes.Data //resExit和Exit数据一样的，所以可以直接发
	//接下来发给每个人“我”离线了
	var bytes [8096]byte
	for id, userprocess := range SharedUserMgr.AllOnlineUsers {
		var tf = &utils.Transfer{userprocess.Conn, bytes}
		err = tf.WritePkg(resmes)
		if err != nil {
			fmt.Printf("%v向%v用户通知离线失败\n", this.ID, id)
		}
	}
	color.Blue("%v用户退出", this.ID)
}
