package model

import (
	"chater/ChatProject/message"
	"net"
)

type CurrentUser struct { //保存当前用户的链接和信息
	message.User `json:"message.User"`
	Conn         net.Conn `json:"conn"`
}
