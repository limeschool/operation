package errors

import "github.com/limeschool/gin"

var (
	//基础相关
	ParamsError         = &gin.CustomError{Code: 100002, Msg: "参数验证失败"}
	AssignError         = &gin.CustomError{Code: 100003, Msg: "数据赋值失败"}
	DBError             = &gin.CustomError{Code: 100004, Msg: "数据库操作失败"}
	DBNotFoundError     = &gin.CustomError{Code: 100005, Msg: "未查询到指定数据"}
	UserNotFoundError   = &gin.CustomError{Code: 100006, Msg: "账号不存在"}
	UserDisableError    = &gin.CustomError{Code: 100007, Msg: "账号已被禁用"}
	PasswordError       = &gin.CustomError{Code: 100008, Msg: "账号密码错误"}
	RsaPasswordError    = &gin.CustomError{Code: 100009, Msg: "非法账号密码"}
	IpLimitLoginError   = &gin.CustomError{Code: 100010, Msg: "当前设备登陆错误次数过多,今日已被限制登陆"}
	UserLimitLoginError = &gin.CustomError{Code: 100010, Msg: "当前账号登陆错误次数过多,已被限制登陆"}
	SuperAdminEditError = &gin.CustomError{Code: 100011, Msg: "超级管理员不允许修改"}
	SuperAdminDelError  = &gin.CustomError{Code: 100012, Msg: "超级管理员不允许删除"}
	RoleDisableError    = &gin.CustomError{Code: 1000013, Msg: "账户角色已被禁用"}
	PasswordExpireError = &gin.CustomError{Code: 1000014, Msg: "登陆密码时效已过期"}
	CaptchaError        = &gin.CustomError{Code: 1000015, Msg: "验证码错误"}

	//auth相关
	NotResourcePower     = &gin.CustomError{Code: 4003, Msg: "暂无接口资源权限"}
	TokenExpiredError    = &gin.CustomError{Code: 4001, Msg: "登陆信息已过期，请重新登陆"}
	RefTokenExpiredError = &gin.CustomError{Code: 4000, Msg: "太长时间未登陆，请重新登陆"}
	DulDeviceLoginError  = &gin.CustomError{Code: 4000, Msg: "你已在其他设备登陆"}
	TokenDataError       = &gin.CustomError{Code: 4000, Msg: "token数据异常失败"}
	TokenValidateError   = &gin.CustomError{Code: 4000, Msg: "token验证失败"}
	TokenEmptyError      = &gin.CustomError{Code: 4000, Msg: "token信息不存在"}

	//menu相关
	DulMenuNameError    = &gin.CustomError{Code: 1000030, Msg: "菜单name值不能重复"}
	MenuParentIdError   = &gin.CustomError{Code: 1000031, Msg: "父菜单id值异常"}
	DeleteRootMenuError = &gin.CustomError{Code: 1000032, Msg: "不能删除根菜单"}

	//team相关
	NotAddTeamError      = &gin.CustomError{Code: 1000040, Msg: "暂无此部门的下级部门创建权限"}
	NotEditTeamError     = &gin.CustomError{Code: 1000041, Msg: "暂无此部门的修改权限"}
	NotDelTeamError      = &gin.CustomError{Code: 1000042, Msg: "暂无此部门的删除权限"}
	NotAddTeamUserError  = &gin.CustomError{Code: 1000043, Msg: "暂无此部门的人员创建权限"}
	NotEditTeamUserError = &gin.CustomError{Code: 1000044, Msg: "暂无此部门的人员修改权限"}
	NotDelTeamUserError  = &gin.CustomError{Code: 1000045, Msg: "暂无此部门的人员删除权限"}
	TeamParentIdError    = &gin.CustomError{Code: 1000046, Msg: "父部门不能为自己"}

	//role相关
	DulKeywordError     = &gin.CustomError{Code: 1000050, Msg: "角色标志符已存在"}
	DisableCurRoleError = &gin.CustomError{Code: 1000051, Msg: "不能禁用当前用户所在角色"}

	// upload相关
	InitUploadError           = &gin.CustomError{Code: 1000060, Msg: "文件上传初始化失败"}
	FileLimitMaxSizeError     = &gin.CustomError{Code: 1000061, Msg: "文件超过规定大小限制"}
	OpenFileError             = &gin.CustomError{Code: 1000062, Msg: "上传文件打开失败"}
	UploadTypeNotSupportError = &gin.CustomError{Code: 1000063, Msg: "上传文件打开失败"}
	UploadTypeError           = &gin.CustomError{Code: 1000063, Msg: "存在不允许上传的文件类型"}

	// registry
	InitDockerError         = &gin.CustomError{Code: 1000070, Msg: "初始化docker-cli 失败"}
	LoginImageRegistryError = &gin.CustomError{Code: 1000071, Msg: "登陆镜像仓库失败"}
	LoginCodeRegistryError  = &gin.CustomError{Code: 1000071, Msg: "登陆代码仓库失败"}

	// cluster
	ParseClusterConfigError = &gin.CustomError{Code: 1000080, Msg: "解析集群config失败"}
	ConnectClusterError     = &gin.CustomError{Code: 1000081, Msg: "连接集群失败"}
	GetClusterError         = &gin.CustomError{Code: 1000081, Msg: "获取集群数据失败"}
	UpdateClusterError      = &gin.CustomError{Code: 1000081, Msg: "更新集群数据失败"}
)
