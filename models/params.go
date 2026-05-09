package models

// 定义注册时的请求参数结构体
// 使用binding 标签能在绑定参数的时候调用validator执行验证
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Repassword string `json:"repassword" binding:"required,eqfield=Password"`
}

// 定义登录时的请求参数结构体
type ParamLogin struct {
	UserID   int64  `json:"user_id,string"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 获取投票数据的映射结构体
type ParamVote struct {
	// UserID 从当前登录的用户获取
	PostID    int64 `json:"post_id,string" binding:"required"`
	// 赞成(1) 反对(-1) 弃票(0)
	Direction int   `json:"direction,string" binding:"oneof=1 0 -1"`
}
