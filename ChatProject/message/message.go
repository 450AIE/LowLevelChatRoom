package message

// 因为以下的结构体的字段都要经常取出来用，于是写作public
// 为什么明明这些结构体十分相似却要写不同的，是因为万一要增加字段可能会大改，而且方法的绑定可以避免错误
const (
	LoginType       = "LoginMessage"
	LoginResType    = "ResMessage"
	ResgisterType   = "RegisterMessage"
	ResRegisterType = "ResRegisterMessage"
	NotifyType      = "Notify"
	ShortMesType    = "ShortMessage"
	ExitType        = "ExitMessage"
	ResExitType     = "ResExitTypeMessage"
)
const (
	UserOnline = iota
	UserOffline
)

type Message struct { //统一传递Message保证处理逻辑一致
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息数据
}
type ResMessage struct {
	Code     string   `json:"code"`  //状态码，500未找到用户，200登录成功
	Error    string   `json:"error"` //错误信息
	AllUsers []string `json:"allUsers"`
	//AllUsers map[string]*process2.UserProcessor这个似乎也不可以，客户端拿到的是指针，在客户端的电脑上，这个指针指向的又不是我们需要的数据
	//AllUsers map[ID]Name 但是我的ID就是Name，怎么办
	//string为用户ID，后面是UserProcessor结构体是因为该结构体含有链接Conn，
	//这个map里相当于有所有用户的Conn，如果Conn被客户端返回接收到了，我们可以知道所有用户的在线情况，以及实现选择某几个链接私聊
	//，某几个链接群聊的功能
}
type LoginMessage struct {
	UserID  string `json:"userID"`
	UserPwd string `json:"userPwd"`
}

type RegisterMessage struct {
	User User `json:"user"`
}

type ResRegisterMessage struct {
	Code  string `json:"code"` //状态码，400表示用户已存在，200注册成功
	Error string `json:"error"`
}

type NotifyUserStatus struct {
	UserID     string `json:"userID"`
	UserStatus int    `json:"userStatus"`
}

type ShortMessage struct { //最好区分为这个和一个返回的Res类型
	Content string `json:"content"`
	User    `json:"user"`
}

type ExitMessage struct {
	UserID string `json:"userID"`
}
type ResExitMessage struct {
	UserID string `json:"userID"`
}
