package process

import (
	"bufio"
	"chater/ChatProject/Client/model"
	"chater/ChatProject/Client/utils"
	"chater/ChatProject/message"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"net"
	"os"
	"strconv"
)

var (
	SharedOnlineUsersList map[string]int = make(map[string]int) //维护的在线用户列表
	CurUser                              = &model.CurrentUser{} //在用户登录成功后，完成对CurUser的初始化
	ShortMesProcess       *SmsProcess
)

// 让客户端也维护一个map，不然每次服务端都要把所有的在线用户数据返回，数据量太大
// 在登录成功后才保持一直通讯
func KeepOnline(conn net.Conn) { //不停地从服务端读取，以达到一直保持通讯的效果
	defer color.Red("这条语句出现，代表协程被关闭，无法起到保持通讯的效果(原因可能是：服务端关闭，用户端退出)")
	var bytes [8096]byte
	var tf = &utils.Transfer{conn, bytes}
	for {
		fmt.Printf("客户端正在等待接收服务端(%v)推送的数据...\n", conn.RemoteAddr().String())
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("接收数据错误，因为服务端没有数据传输") //如果是一次的偶然错误，会退出，不是很合适，但是如果是服务端凉了，continue的话，那么会一直死循环了，更不合适
			return
		}
		switch mes.Type {
		case message.NotifyType:
			var notify = &message.NotifyUserStatus{}
			err := json.Unmarshal([]byte(mes.Data), &notify)
			if err != nil {
				fmt.Println("反序列化错误")
				continue
			}
			//fmt.Printf("%v用户已上线\n", notify.UserID)
			color.Yellow("%v用户已上线\n", notify.UserID)
			//如果我们采用的map[XX]User结构体的话，可以先判断有没有该用户，不然每次都实例化一个结构体再赋值很耗费性能
			SharedOnlineUsersList[notify.UserID] = notify.UserStatus //有人登录就把其信息加入map维护
		case message.ShortMesType:
			var shortmes = &message.ShortMessage{}
			err := json.Unmarshal([]byte(mes.Data), shortmes)
			if err != nil {
				fmt.Println("反序列化错误")
				continue
			}
			//fmt.Printf("   %v用户对大家说：%v\n", shortmes.UserID, shortmes.Content)
			color.Green("%v用户对大家说：%v\n", shortmes.UserID, shortmes.Content)
		case message.ResExitType:
			var resExit = &message.ResExitMessage{}
			err := json.Unmarshal([]byte(mes.Data), resExit)
			if err != nil {
				fmt.Println("反序列化错误")
				continue
			}
			//fmt.Printf("%v用户已退出\n", resExit.UserID)
			color.Blue("%v用户已退出\n", resExit.UserID)
			delete(SharedOnlineUsersList, resExit.UserID)
		default:
			fmt.Println("服务端返回了未知的类型")
		}
	}
}
func ShowMenu() {
	for {
		fmt.Println("--------------------------------------------------")
		fmt.Println("\t\t1.显示在线用户列表")
		fmt.Println("\t\t2.发送消息")
		fmt.Println("\t\t3.离线留言")
		fmt.Println("\t\t4.退出系统")
		fmt.Println("请选择(1-4):_\b")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			ShowAllOnlineUsers()
		case 2:
			fmt.Println("请输入消息:")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			var content = scanner.Text()
			err := ShortMesProcess.SendGroup(content)
			if err != nil {
				fmt.Println("消息发送错误")
				continue
			}
		case 3:
			fmt.Println("好友功能都还没做出来，没实现点对点聊天，怎么可能有离线留言啊哥nia~nia~nia~")
		case 4:
			fmt.Println("正在退出系统...")
			var user *UserProcess
			user.Exit(CurUser.Conn, CurUser.UserID)
			//可以在推出之前传递给服务器消息，表示客户端将退出，便于服务端写入缓存等一系列操作
			//注意！os.Exit(0) 会立刻终止程序(哪怕是在协程中，也会终止该主线程)，并且不执行defer
		default:
			fmt.Println("输入错误，请重新输入")
		}
	}
}

func ShowAllOnlineUsers() {
	fmt.Println("当前在线用户列表如下：")
	var n = 1
	for i, _ := range SharedOnlineUsersList {
		//fmt.Printf(strconv.Itoa(n)+".用户:%v\n", i)
		color.Yellow(strconv.Itoa(n)+".用户:%v\n", i)
		n++
	}
}

// 传递conn和退出用户的ID，发送给客户端mes表示退出,正常发送则发送后立即退出,妈的为什么写在UserProcess就识别不到，写在这了就可以识别，艹，改正！！
func (this *UserProcess) Exit(conn net.Conn, ID string) {
	var exit = &message.ExitMessage{ID}
	var mes = &message.Message{}
	mes.Type = message.ExitType
	data, err := json.Marshal(exit)
	if err != nil {
		fmt.Println("序列化失败")
		return
	}
	mes.Data = string(data)
	var bytes [8096]byte
	var tf = &utils.Transfer{conn, bytes}
	err = tf.WritePkg(mes)
	if err != nil {
		fmt.Println("用户退出消息传递服务器发生错误")
		return
	}
	os.Exit(0)
}
