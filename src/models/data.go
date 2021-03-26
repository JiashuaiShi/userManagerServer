package models

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

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type QuerySingle struct {
	UserId string `json:"user_id"`
}

type QueryList struct {
	PageNo   int    `json:"page_no"`
	PageSize int    `json:"page_size"`
	Query    string `json:"query"`
}
