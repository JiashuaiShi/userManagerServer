package models

import log "github.com/sirupsen/logrus"

// http接口
type User struct {
	UserId    string `json:"user_id"`
	Nickname  string `json:"nickname"`
	RoleType  uint32 `json:"role_type,string"`
	LoginTime int64  `json:"login_time,string,omitempty"`
}

type UserList struct {
	Total  int32  `json:"total" comment:"总记录数" example:"" validate:""`
	PageNo int32  `json:"page_no" comment:"当前页码" example:"" validate:""`
	Users  []User `json:"users" comment:"总记录" example:"" validate:""`
}

// http响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// http请求
type QuerySingle struct {
	UserId string `json:"user_id"`
}

type QueryList struct {
	PageNo   int    `json:"page_no"`
	PageSize int    `json:"page_size"`
	Query    string `json:"query"`
}

// 配置结构
type Config struct {
	RedisConfig RedisConfig `yaml:"redis_config"`
	LogConfig   LogConfig   `yaml:"log_config"`
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
	Level string `yaml:"level"`
}

// log Level字典
var LogLevelMap map[string]log.Level

func init() {
	LogLevelMap = make(map[string]log.Level)
	LogLevelMap["DebugLevel"] = log.DebugLevel
	LogLevelMap["InfoLevel"] = log.InfoLevel
	LogLevelMap["WarnLevel"] = log.WarnLevel
	LogLevelMap["ErrorLevel"] = log.ErrorLevel
}
