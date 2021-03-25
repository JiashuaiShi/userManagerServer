package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
	"zego.com/userManageServer/src/models"
	"zego.com/userManageServer/src/contextConf"
	"zego.com/userManageServer/src/service/redis"
)

// 增加用户：http://loclahost:8080/service/add_user
func addHandle(context *gin.Context) {
	// 从context填入user信息
	var user models.User
	if err := context.BindJSON(&user); err != nil {
		log.Println("BindJSON failed!")
		return
	}

	// 更新登录时间
	user.LoginTime = time.Now().Unix()

	// 打印原始数据
	fmt.Println("----新增用户开始-----")
	fmt.Println("user_id:", user.UserId, " nickname:", user.NickName, " role_type:", user.RoleType, " login_time:", user.LoginTime)

	// 获取user_id 和 string类型的data
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Println("struct2json Marshal failed!")
		return
	}
	fmt.Printf("Json数据: %+v\n", string(jsonData))

	var resp models.Response

	// 保存到redis
	if err := redis.SetData(user.UserId, string(jsonData), user.LoginTime); err != nil {
		resp.Code = contextConf.ERROR.Code
		resp.Message = contextConf.ERROR.Msg
		log.Println("save data failed!")
	} else {
		resp.Code = contextConf.SUCCESS.Code
		resp.Message = contextConf.SUCCESS.Msg
	}

	context.JSON(http.StatusOK, &resp)

	fmt.Println("----新增用户结束-----")
}

// 查询用户：http://loclahost:8080/service/get_user？name=1413
func getHandle(context *gin.Context) {
	fmt.Println("----查询用户开始-----")

	// 获取用户ID
	userId := context.Query("query")
	fmt.Println("查询的用户ID：", userId)

	var resp models.Response

	// 从redis查询ID的string格式的data
	userStr, err := redis.GetData(userId)
	if err != nil {
		log.Println(err.Error())
		resp.Code = contextConf.ERROR.Code
		resp.Message = contextConf.ERROR.Msg
		return
	}

	fmt.Printf("Json数据: %+v\n", userStr)

	users := make([]models.User, 0)

	for _, userStr := range userStr {
		user := models.User{}
		json.Unmarshal([]byte(userStr), &user)
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
	fmt.Println("----查询用户结束-----")
}

// 查询用户列表：http://loclahost:8080/service/get_userlist
func getListHandle(context *gin.Context) {
	fmt.Println("----查询用户列表开始-----")

	pageNoStr := context.Query("page_no")
	pageSizeStr := context.Query("page_size")

	pageNo, _ := strconv.Atoi(pageNoStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	minv := (pageNo - 1) * pageSize
	maxv := minv + pageSize - 1

	usersStr, total, err := redis.GetAllData(minv, maxv)

	var resp models.Response

	if err != nil {
		log.Println(err.Error())
		resp.Code = contextConf.ERROR.Code
		resp.Message = contextConf.ERROR.Msg
		return
	}

	users := make([]models.User, 0)

	for _, userStr := range usersStr {
		user := models.User{}
		json.Unmarshal([]byte(userStr), &user)
		users = append(users, user)
	}

	resp = models.Response{
		Code:    contextConf.SUCCESS.Code,
		Message: contextConf.SUCCESS.Msg,
		Data: models.UserList{
			Total:  int32(total),
			PageNo: int32(pageNo),
			Users:  users,
		},
	}
	context.JSON(http.StatusOK, &resp)
	fmt.Println("----查询用户列表结束-----")
}

// 删除用户：http://loclahost:8080/service/del_user/id
func deleteHandle(context *gin.Context) {
	fmt.Println("----删除用户列表开始-----")
	userIdStr := context.Param("id")
	//userId, _ := strconv.Atoi(userIdStr)
	var resp models.Response
	if err := redis.DelData(userIdStr); err != nil {
		log.Println(err.Error())
		resp.Code = contextConf.ERROR.Code
		resp.Message = contextConf.ERROR.Msg
		return
	}

	resp.Code = contextConf.SUCCESS.Code
	resp.Message = contextConf.SUCCESS.Msg

	context.JSON(http.StatusOK, &resp)
	fmt.Println("----查询用户列表完成-----")
}
