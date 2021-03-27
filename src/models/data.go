package models

import log "github.com/sirupsen/logrus"

// userinfo定义
type User struct {
	UserId    string `json:"user_id"`
	Nickname  string `json:"nickname"`
	RoleType  uint32 `json:"role_type,string"`
	LoginTime int64  `json:"login_time,string,omitempty"`
}

// get_userlist接口定义
type UserList struct {
	Total int32  `json:"total" comment:"userinfo数组长度"`
	Users []User `json:"users" comment:"userinfo数组内容"`
}

// Response响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// http请求
type Query struct {
	UserId   string `json:"user_id"`  	//post请求使用userid，get/delete方法从url中获取userid
	PageNo   int    `json:"page_no"`
	PageSize int    `json:"page_size"`
}

// 全部配置信息
type Config struct {
	RedisConfig  RedisConfig  `yaml:"redis_config"`
	LogConfig    LogConfig    `yaml:"log_config"`
	RouterConfig RouterConfig `yaml:"router_config"`
}

// redis配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

// log配置
type LogConfig struct {
	Path  string `yaml:"path"`
	Level int    `yaml:"level"`
}

// 启动端口配置
type RouterConfig struct {
	Port string `yaml:"port"`
}

// log_Level字典
var LogLevelMap map[int]log.Level

// 日志打印级别字典，根据Albert要求，使用int-Level的映射关系
func init() {
	LogLevelMap = make(map[int]log.Level)
	LogLevelMap[5] = log.DebugLevel //DebugLevel
	LogLevelMap[4] = log.InfoLevel  // InfoLevel
	LogLevelMap[3] = log.WarnLevel  // WarnLevel
	LogLevelMap[2] = log.ErrorLevel // ErrorLevel
	LogLevelMap[1] = log.FatalLevel // ErrorLevel
}
