package main

import (
	"chater/ChatProject/Client/process"
	"fmt"
	"github.com/fatih/color"
	"time"
)

func main() {
	fmt.Println("是否跳过开发背景？（输入1跳过，其余键不跳过）")
	var choice string
	fmt.Scanln(&choice)
	if choice != "1" {
		fmt.Println("那天，是个风和日丽的清晨…………")
		time.Sleep(time.Second)
		color.Green("\t***:冲哥，我需要你的帮助！\n")
		time.Sleep(time.Second)
		color.Yellow("我：干嘛？\n")
		time.Sleep(time.Second)
		color.Green("\t***:你知道的……那什么……我在海南……在海外，很容易接触…………咳咳……\n")
		time.Sleep(time.Second)
		color.Yellow("我：！(吃惊）\n")
		time.Sleep(time.Second)
		color.Yellow("我：难道……你终于还是迈出了那一步？！…………！\n")
		time.Sleep(time.Second)
		color.Green("\t***：sir，don‘t say like that ， I need your help\n")
		time.Sleep(time.Second)
		color.Yellow("我：不可以！群众要有群众的样子！更何况你呢！\n")
		time.Sleep(time.Second)
		color.Green("\t***：（堵住耳朵）nia~nia~nia~我听不见\n")
		time.Sleep(time.Second)
		color.Green("\t***:事成之后，封你为*交部部长\n")
		time.Sleep(time.Second)
		color.Yellow("*:!\n")
		time.Sleep(time.Second)
		color.Yellow("*:那就不奇怪了\n")
		time.Sleep(time.Second)
		color.Yellow("*:那就不奇怪了\n")
		time.Sleep(time.Second)
		color.Yellow("*:sir，this way\n")
		time.Sleep(time.Second)
		fmt.Println("这便是这个破程序的开发背景……虽然功能很简单拉垮，却编写了很久…………")
		color.Red("（以上故事纯属戏剧性编造，本人坚决拥护中国共产党的领导）\n")
		time.Sleep(time.Second)
	}
	process.Margin() //含标题显示与登陆

}
