package process

import (
	"bufio"
	"fmt"
	"os"
)

//var userID, userPWD string

func Margin() {
	var serverAddress string
	fmt.Println("请输入服务器地址:")
	fmt.Scanln(&serverAddress)
	for {
		var userPWD string
		fmt.Println("--------------------欢迎登陆多人聊天系统--------------------")
		fmt.Println("\t\t\t1 登陆")
		fmt.Println("\t\t\t2 注册")
		fmt.Println("\t\t\t3 退出")
		fmt.Print("请选择：_\b")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			fmt.Println("--------------------------------------------------")
			fmt.Print("请输入账户名：_\b")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			var userID = scanner.Text()
			fmt.Print("请输入密码：_\b")
			fmt.Scanln(&userPWD)
			var user *UserProcess
			err := user.Login(userID, userPWD, serverAddress)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			//fmt.Printf("用户：%v  登录成功\n", userID)
			//keepOnline.ShowMenu()
		case 2:
			fmt.Println("--------------------------------------------------")
			fmt.Print("请输入账户名：_\b")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			var userID = scanner.Text()
			fmt.Print("请输入密码：_\b")
			fmt.Scanln(&userPWD)
			var user *UserProcess
			err := user.Register(userID, userPWD, serverAddress)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("用户：%v  注册成功\n", userID)
			continue
		case 3:
			fmt.Print("确定退出？(y退出/其余键继续):_\b")
			var choose string
			fmt.Scanln(&choose)
			if choose == "y" || choose == "Y" {
				return
			} else if choose == "n" {
				continue
			}
		default:
			fmt.Println("错误输入，请重新输入")
		}
	}

}
