package process

import (
	"chater/ChatProject/message"
	"chater/ChatProject/server/utils"
	"fmt"
)

// 要做离线留言的话，要把mes存在数据库内，做好标记，等用户上线再发送
func (this *UserProcessor) ShortMessageProcess(mes *message.Message) {
	var bytes [8096]byte
	for id, userprocess := range SharedUserMgr.AllOnlineUsers {
		var tf = &utils.Transfer{userprocess.Conn, bytes}
		err := tf.WritePkg(mes)
		if err != nil {
			fmt.Printf("%v用户向%v用户发送消息失败", this.ID, id)
		}
	}
}
