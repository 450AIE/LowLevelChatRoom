package process

import (
	"errors"
	"fmt"
)

// 因为UserMgr在整个服务器只有一个，且多数地方使用，故定义一个全局变量
var (
	SharedUserMgr *UserMgr
)

type UserMgr struct { //负责记录在线用户列表
	AllOnlineUsers map[string]*UserProcessor
}

func init() {
	SharedUserMgr = &UserMgr{
		make(map[string]*UserProcessor, 1024),
	}
}
func (this *UserMgr) AddOnlineUser(user *UserProcessor) { //这里应该不会error
	this.AllOnlineUsers[user.ID] = user //Add之前我们有从数据库检查是否有用户存在，所以不用担心ID一样导致覆盖
}
func (this *UserMgr) DelOnlineUser(user *UserProcessor) {
	delete(this.AllOnlineUsers, user.ID)
}
func (this *UserMgr) GetAllOnlineUser() map[string]*UserProcessor { //返回得到在线用户列表
	return this.AllOnlineUsers
}

func (this *UserMgr) GetOnlineUserConnByID(ID string) (user *UserProcessor, err error) { //根据用户名获取某个在线用户的链接
	v, ok := this.AllOnlineUsers[ID]
	if ok {
		fmt.Printf("%v用户存在\n")
		return v, nil
	} else {
		fmt.Printf("%v用户不存在\n")
		return nil, errors.New(ID + "用户不存在")
	}
}
