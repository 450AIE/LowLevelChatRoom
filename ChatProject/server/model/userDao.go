package model

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type UserDao struct {
	pool *redis.Pool //为什么把池子放在结构体里？
}

var SharedUserDao *UserDao

func NewUserDao(pool *redis.Pool) (user *UserDao) {
	return &UserDao{pool} //为什么偏要用工厂模式创建？我感觉这样应该可以避免误操作，在别的包里面，只可以调用方法，不可以修改字段
}

// 根据用户ID判断是否存在用户
func (this *UserDao) InspectID(conn redis.Conn, id string) (err error) { //似乎可以把conn改为this.pool.Get()取一根用，再defer
	_, err = conn.Do("HGET", "User:"+id, id) //用户不存在返回err，存在nil
	if err != nil {
		if err == redis.ErrNil { //表示这个错误是因为redis中不存在所查询的内容
			return errors.New("用户不存在")
		}
		return
	}
	//我选择在redis存储原始hash，没有序列化，所以这里不反序列化
	return nil
}
func (this *UserDao) InspectPwd(conn redis.Conn, id, pwd string) (err error) { //似乎可以把conn改为this.pool.Get()取一根用，再defer
	Pwd, err := redis.String(conn.Do("HGET", "User:"+id, "pwd"))

	if err != nil {
		fmt.Println("查询用户密码失败")
		return
	}
	if Pwd == pwd {
		return nil
	}
	//我选择在redis存储原始hash，没有序列化，所以这里不反序列化
	return errors.New("用户密码错误")
}
func (this *UserDao) Login(ID string, Pwd string) (err error) { //判断有无该用户，有的话接着判断密码是否正确
	conn := this.pool.Get()
	defer conn.Close()

	err = this.InspectID(conn, ID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = this.InspectPwd(conn, ID, Pwd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
func (this *UserDao) Register(ID string, Pwd string) (err error) {
	conn := this.pool.Get()
	err = this.InspectID(conn, ID)
	if err == redis.ErrNil { //InspectID返回的err来判断有严重的逻辑错误，之后修改！
		fmt.Println("该用户已注册过，请重新输入")
		return errors.New("该用户已注册过")
	}
	//fmt.Println("临门一脚")
	_, err = conn.Do("HMSET", "User:"+ID, "id", ID, "pwd", Pwd)
	if err != nil {
		fmt.Println("HSET错误")
		return
	}
	return
}
