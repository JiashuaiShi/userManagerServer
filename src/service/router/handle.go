package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"zego.com/userManageServer/src/context-conf"
	log "zego.com/userManageServer/src/logger"
	"zego.com/userManageServer/src/models"
	"zego.com/userManageServer/src/service/redis"
)

// 增加用户：http://loclahost:8080/user/add_user
func addHandle(context *gin.Context) {
	var user models.User
	var resp models.Response
	logFields := log.Field{}

	// 解析post中body内容
	if err := context.BindJSON(&user); err != nil {
		// 返回错误响应
		resp.Code = context_conf.ERROR.Code
		resp.Message = context_conf.ERROR.Msg
		context.JSON(http.StatusOK, &resp)
		log.Error(logFields, "BindJSON failed, error:%+v", err)
		return
	}

	// 服务端填入登陆时间
	user.LoginTime = time.Now().Unix()

	// 解析成功后填入Field，用于err情况输出
	logFields["UserId"] = user.UserId
	logFields["Nickname"] = user.Nickname
	logFields["RoleType"] = user.RoleType
	logFields["LoginTime"] = user.LoginTime

	// 打印原始数据
	log.Debug(logFields, "addUsr raw data")

	// 转成json格式准备填入redis
	userInfo, err := json.Marshal(user)
	if err != nil {
		// 返回错误响应
		resp.Code = context_conf.ERROR.Code
		resp.Message = context_conf.ERROR.Msg
		context.JSON(http.StatusOK, &resp)
		log.Error(logFields, "json marshal failed! error%+v:", err)
		return
	}

	// 保存到redis， uid作为member， logintime作为score
	if err := redis.SetData(user.UserId, string(userInfo), user.LoginTime); err != nil {
		// 返回错误响应
		resp.Code = context_conf.ERROR.Code
		resp.Message = context_conf.ERROR.Msg
		context.JSON(http.StatusOK, &resp)
		log.Error(nil, "redis save data failed!, error:%+v", err)
		return
	}

	// 返回添加成功响应
	resp.Code = context_conf.SUCCESS.Code
	resp.Message = context_conf.SUCCESS.Msg
	context.JSON(http.StatusOK, &resp)
	log.Info(nil, "adduser success!")
	return
}

// 查询用户：http://loclahost:8080/user/get_user？name=ryder
func getHandle(context *gin.Context) {
	var resp models.Response
	logFields := log.Field{}

	// 使用http的get请求
	userId := context.Query("user_id")
	logFields["userid"] = userId

	// 从redis查询userInfo
	userInfo, err := redis.GetData(userId)
	if err != nil {
		// 返回错误响应
		log.Error(logFields, "get data from redis failed! error:%+v", err)
		resp.Code = context_conf.ERROR.Code
		resp.Message = context_conf.ERROR.Msg
		context.JSON(http.StatusOK, &resp)
		return
	}

	logFields["userInfo"] = userInfo
	log.Debug(logFields, "redis raw data.")

	// 查询单个用户, 对应前端分页table的展示，所以返回一个只有一个元素的数组，展示单个用户结果
	users := make([]models.User, 0)
	user := models.User{}

	// 把redis的记录反序列化为结构体，填回response
	if err := json.Unmarshal([]byte(userInfo), &user); err != nil {
		// 返回错误响应
		resp.Code = context_conf.ERROR.Code
		resp.Message = context_conf.ERROR.Msg
		context.JSON(http.StatusOK, &resp)
		log.Error(logFields, "json convert failed!, error:%+v", err)
		return
	}

	// 查询成功返回
	users = append(users, user)
	resp.Code = context_conf.SUCCESS.Code
	resp.Message = context_conf.SUCCESS.Msg
	// 单个用户userid查询，返回总个数1
	resp.Data = models.UserList{Total: 1, Users: users}
	context.JSON(http.StatusOK, &resp)
	log.Info(nil, "getUser Success!")
}

// 查询用户列表：http://loclahost:8080/user/get_userlist ，
// 使用http的post方法查询
func getListHandle(context *gin.Context) {
	var query models.Query
	var resp models.Response
	logFields := log.Field{}

	// 解析post中body内容
	if err := context.BindJSON(&query); err != nil {
		// 返回错误响应
		resp.Code = context_conf.ERROR.Code
		resp.Message = context_conf.ERROR.Msg
		context.JSON(http.StatusOK, &resp)
		log.Error(logFields, "BindJSON failed, error:%+v", err)
		return
	}

	pageNo := query.PageNo
	pageSize := query.PageSize
	logFields["pageNo"] = pageNo
	logFields["pageSize"] = pageNo

	// 设置单次请求的容量最大为50，配置在代码中
	if pageSize > 50 {
		// 返回错误响应, 暂时只配置成功和失败两种响应
		resp.Code = context_conf.ERROR.Code
		resp.Message = context_conf.ERROR.Msg
		context.JSON(http.StatusOK, &resp)
		log.Error(logFields, "getUserList failed! pageSize is over 50!")
		return
	}

	// redis中ZSET的查询起止序号
	start := (pageNo - 1) * pageSize
	stop := start + pageSize - 1
	usersInfo, total, err := redis.GetDataList(start, stop)

	// redis查询失败响应
	if err != nil {
		resp.Code = context_conf.ERROR.Code
		resp.Message = context_conf.ERROR.Msg
		context.JSON(http.StatusOK, &resp)
		log.Error(logFields, "redis return failed. error:%+v", err)
		return
	}

	// 返回[]models.User
	users := make([]models.User, 0)
	for _, userInfo := range usersInfo {
		user := models.User{}
		// 把redis的记录反序列化为结构体，填回response
		if err := json.Unmarshal([]byte(userInfo), &user); err != nil {
			resp.Code = context_conf.ERROR.Code
			resp.Message = context_conf.ERROR.Msg
			context.JSON(http.StatusOK, &resp)
			log.Error(nil, "json Unmarshal failed. error:%+v", err)
			return
		}
		users = append(users, user)
	}

	// 查询成功
	userList := models.UserList{Total: total, Users: users}
	resp.Data = userList
	resp.Code = context_conf.SUCCESS.Code
	resp.Message = context_conf.SUCCESS.Msg
	context.JSON(http.StatusOK, &resp)
	log.Info(nil, "getUserList success!")
}

// 删除用户：http://loclahost:8080/user/del_user/id
// 使用http的post方法
func deleteHandle(context *gin.Context) {
	var resp models.Response
	var query models.Query
	logField := log.Field{}

	// 解析post中body内容
	if err := context.BindJSON(&query); err != nil {
		// 返回错误响应
		resp.Code = context_conf.ERROR.Code
		resp.Message = context_conf.ERROR.Msg
		context.JSON(http.StatusOK, &resp)
		log.Error(logField, "BindJSON failed, error:%+v", err)
		return
	}

	//userId := context.Param("id")
	userId := query.UserId
	if err := redis.DelData(userId); err != nil {
		logField["userId"] = userId
		log.Error(logField, "delete user failed! error:%+v", err)
		resp.Code = context_conf.ERROR.Code
		resp.Message = context_conf.ERROR.Msg
		return
	}

	resp.Code = context_conf.SUCCESS.Code
	resp.Message = context_conf.SUCCESS.Msg
	log.Info(nil, "deleteUser success")
	context.JSON(http.StatusOK, &resp)
}
