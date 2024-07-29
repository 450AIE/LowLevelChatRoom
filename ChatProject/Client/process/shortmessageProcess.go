package process

import (
	"chater/ChatProject/Client/utils"
	"chater/ChatProject/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
} //像这种没有任何字段，专门用于调用特定方法的结构体，其实可以定义一个全局变量使用，节约性能，不过万一登录失败呢没有使用呢？在登录成功后的二级菜单的switch外面也可以

func (sm *SmsProcess) SendGroup(content string) (err error) {
	var mes = &message.Message{}
	mes.Type = message.ShortMesType
	var shortmes = &message.ShortMessage{Content: content, User: message.User{CurUser.UserID, "", CurUser.UserStatus}}
	data, err := json.Marshal(shortmes)
	if err != nil {
		fmt.Println("序列化失败")
		return
	}
	mes.Data = string(data)
	var bytes [8096]byte
	var tf = &utils.Transfer{CurUser.Conn, bytes}
	err = tf.WritePkg(mes)
	if err != nil {
		fmt.Println("传输短消息错误")
		return
	}
	return nil
}
