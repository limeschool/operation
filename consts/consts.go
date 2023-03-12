package consts

const (
	Token     = "Authorization"
	DATABASE  = "operation"
	REDIS     = "redis"
	RedisLock = "redis"

	RsaPrivate   = "private"
	RsaPublic    = "public"
	BaseApi      = "baseApi"
	RedisBaseApi = "sysBaseApi"
)

const (
	ALLTEAM  = "ALLTEAM"  // 所有权限
	DOWNTEAM = "DOWNTEAM" // 部门以下权限
	CURTEAM  = "CURTEAM"  // 当前部门权限
	CUSTOM   = "CUSTOM"   // 自定义权限
)

const (
	JwtSuperAdmin = "superAdmin"
	JwtMapClaims  = "mapClaims"
)
