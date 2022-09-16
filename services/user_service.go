package services

import (
	"app/libs/redislib"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

const UserListKey string = "UserListKey"
const LoginUserListKey string = "LoginUserListKey"

type UserService struct {
}

type User struct {
	Username string `username:"string"`
	Password string `password:"string"`
}

func (userService *UserService) Register(ctx *gin.Context) interface{} {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	redisClient := redislib.GetRedisClient()
	var users map[string]User
	users = make(map[string]User) //在使用时没有进行初始化map，导致使用时失败，或者直接声明时，使用
	userListJson, _ := redisClient.Get(UserListKey).Result()
	_ = json.Unmarshal([]byte(userListJson), &users)
	// 获取 map 中某个 key 是否存在的语法。如果 ok 是 true，表示 key 存在，key 对应的值就是 value ，反之表示 key 不存在。
	_, ok := users[username]
	if !ok {
		users[username] = User{
			Username: username,
			Password: password,
		}
		newUserListJson, _ := json.Marshal(users)
		redisClient.Set(UserListKey, newUserListJson, 86400*time.Second)
		return "注册成功"
	} else {
		return "用户已注册"
	}
}

func (userService *UserService) SignIn(ctx *gin.Context) interface{} {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	redisClient := redislib.GetRedisClient()
	var users map[string]User
	users = make(map[string]User) //在使用时没有进行初始化map，导致使用时失败，或者直接声明时，使用
	userListJson, _ := redisClient.Get(UserListKey).Result()
	_ = json.Unmarshal([]byte(userListJson), &users)
	// 获取 map 中某个 key 是否存在的语法。如果 ok 是 true，表示 key 存在，key 对应的值就是 value ，反之表示 key 不存在。
	user, ok := users[username]
	if !ok {
		return "用户不存在"
	} else {
		if password != user.Password {
			return "密码错误"
		}
		var loginUsers map[string]User
		loginUsers = make(map[string]User) //在使用时没有进行初始化map，导致使用时失败，或者直接声明时，使用
		loginUserListJson, _ := redisClient.Get(LoginUserListKey).Result()
		_ = json.Unmarshal([]byte(loginUserListJson), &loginUsers)
		_, hasLogin := loginUsers[username]
		if hasLogin {
			return "用户已登录"
		}
		loginUsers[username] = User{
			Username: username,
			Password: password,
		}
		newLoginUserListJson, _ := json.Marshal(loginUsers)
		redisClient.Set(LoginUserListKey, newLoginUserListJson, 86400*time.Second)
		return "登录成功"
	}
}

func (userService *UserService) SignOut(ctx *gin.Context) interface{} {
	username := ctx.PostForm("username")
	redisClient := redislib.GetRedisClient()
	var loginUsers map[string]User
	loginUsers = make(map[string]User) //在使用时没有进行初始化map，导致使用时失败，或者直接声明时，使用
	loginUserListJson, _ := redisClient.Get(LoginUserListKey).Result()
	_ = json.Unmarshal([]byte(loginUserListJson), &loginUsers)
	_, ok := loginUsers[username]
	if !ok {
		return "用户未登录"
	} else {
		delete(loginUsers, username)
		newLoginUserListJson, _ := json.Marshal(loginUsers)
		redisClient.Set(LoginUserListKey, newLoginUserListJson, 86400*time.Second)
		return "退出登录成功"
	}
}

func (userService *UserService) List(ctx *gin.Context) interface{} {
	redisClient := redislib.GetRedisClient()
	var users map[string]User
	users = make(map[string]User) //在使用时没有进行初始化map，导致使用时失败，或者直接声明时，使用
	userListJson, _ := redisClient.Get(UserListKey).Result()
	_ = json.Unmarshal([]byte(userListJson), &users)
	return users
}
