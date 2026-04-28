package models

// 定义注册时的请求参数结构体
// 使用binding 标签能在绑定参数的时候调用validator执行验证
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Repassword string `json:"repassword" binding:"required,eqfield=Password"`
}
