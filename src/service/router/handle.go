package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"zego.com/userManageServer/src/contextConf"
	log "zego.com/userManageServer/src/logger"
	"zego.com/userManageServer/src/models"
	"zego.com/userManageServer/src/service/redis"
)

// 增加用户：http://loclahost:8080/service/add_user
func addHandle(context *gin.Context) {
	log.Info(nil, "add user start")

	// 从context填入user信息
	var user models.User
	if err := context.BindJSON(&user); err != nil {
		log.Error(log.Field{
			"UserId":    user.UserId,
			"Nickname":  user.Nickname,
			"RoleType":  user.RoleType,
			"LoginTime": user.LoginTime,
		}, "BindJSON failed, error:%+v", err)
		return
	}

	// 更新登录时间
	user.LoginTime = time.Now().Unix()

	// 打印原始数据
	log.Debug(log.Field{
		"user_id:":    user.UserId,
		"nickname:":   user.Nickname,
		"role_type:":  user.RoleType,
		"login_time:": user.LoginTime},
		"add usr raw data")

	// 获取user_id 和 string类型的data
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Error(log.Field{
			"UserId":    user.UserId,
			"Nickname":  user.Nickname,
			"RoleType":  user.RoleType,
			"LoginTime": user.LoginTime,
		},
			"json marshal failed! error%+v:", err)
		return
	}
	log.Debug(log.Field{
		"UserId":    user.UserId,
		"Nickname":  user.Nickname,
		"RoleType":  user.RoleType,
		"LoginTime": user.LoginTime,
	},
		"redis存储信息")

	var resp models.Response

	// 保存到redis
	if err := redis.SetData(user.UserId, string(jsonData), user.LoginTime); err != nil {
		resp.Code = contextConf.ERROR.Code
		resp.Message = contextConf.ERROR.Msg
		log.Error(nil, "redis save data failed!, error:%+v", err)
	} else {
		resp.Code = contextConf.SUCCESS.Code
		resp.Message = contextConf.SUCCESS.Msg
	}

	context.JSON(http.StatusOK, &resp)
}

// 查询用户：http://loclahost:8080/service/get_user？name=1413
func getHandle(context *gin.Context) {
	log.Info(nil, "getuser start!")

	// 获取用户ID
	userId := context.Query("query")

	// 从redis查询ID的string格式的data
	var resp models.Response
	userStr, err := redis.GetData(userId)
	if err != nil {
		log.Error(log.Field{
			"userInfo": userStr,
		},
			"get data from redis failed! error:%+v", err)
		resp.Code = contextConf.ERROR.Code
		resp.Message = contextConf.ERROR.Msg
		return
	}

	log.Debug(log.Field{
		"userInfo": userStr,
	},
		"redis raw data.")

	users := make([]models.User, 0)

	for _, userStr := range userStr {
		user := models.User{}
		if err := json.Unmarshal([]byte(userStr), &user); err != nil {
			log.Error(log.Field{
				"UserId":    user.UserId,
				"Nickname":  user.Nickname,
				"RoleType":  user.RoleType,
				"LoginTime": user.LoginTime,
			},
				"json unmarshal failed. error:%+v", err)
		}
		users = append(users, user)
	}

	resp = models.Response{
		Code:    contextConf.SUCCESS.Code,
		Message: contextConf.SUCCESS.Msg,
		Data: models.UserList{
			Total:  1,
			PageNo: 1,
			Users:  users,
		},
	}
	context.JSON(http.StatusOK, &resp)
}

// 查询用户列表：http://loclahost:8080/service/get_userlist
func getListHandle(context *gin.Context) {
	log.Info(nil, "getUserList Start!")

	pageNoStr := context.Query("page_no")
	pageSizeStr := context.Query("page_size")

	pageNo, _ := strconv.Atoi(pageNoStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	minVal := (pageNo - 1) * pageSize
	maxVal := minVal + pageSize - 1

	var resp models.Response
	usersStr, total, err := redis.GetAllData(minVal, maxVal)

	if err != nil {
		log.Error(log.Field{
			"rawData": usersStr,
		}, "redis return failed. error:%+v", err)

		resp.Code = contextConf.ERROR.Code
		resp.Message = contextConf.ERROR.Msg
		return
	}

	users := make([]models.User, 0)

	for _, userStr := range usersStr {
		user := models.User{}
		if err := json.Unmarshal([]byte(userStr), &user); err != nil {
			log.Error(log.Field{
				"UserId":    user.UserId,
				"Nickname":  user.Nickname,
				"RoleType":  user.RoleType,
				"LoginTime": user.LoginTime,
			}, "json Unmarshal failed. error:%+v", err)
			return
		}
		users = append(users, user)
	}

	resp = models.Response{
		Code:    contextConf.SUCCESS.Code,
		Message: contextConf.SUCCESS.Msg,
		Data: models.UserList{
			Total:  total,
			PageNo: int32(pageNo),
			Users:  users,
		},
	}
	context.JSON(http.StatusOK, &resp)
}

// 删除用户：http://loclahost:8080/service/del_user/id
func deleteHandle(context *gin.Context) {
	log.Info(nil, "delete user start")
	userIdStr := context.Param("id")
	var resp models.Response
	if err := redis.DelData(userIdStr); err != nil {
		log.Error(log.Field{
			"useid": userIdStr,
		}, "delete user failed! error:%+v", err)

		resp.Code = contextConf.ERROR.Code
		resp.Message = contextConf.ERROR.Msg
		return
	}

	resp.Code = contextConf.SUCCESS.Code
	resp.Message = contextConf.SUCCESS.Msg

	context.JSON(http.StatusOK, &resp)
}
