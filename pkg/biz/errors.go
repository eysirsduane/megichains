package biz

var (
	//系统基本错误码
	CodeSuccess     = NewSpecificError(0, "请求成功")
	CodeServerError = NewSpecificError(500, "服务器异常")
	CodeParamsEmpty = NewSpecificError(407, "参数为空")
	CodeParamsError = NewSpecificError(408, "参数错误")
	DatabaseError   = NewSpecificError(600, "数据库错误")
	CodeUnknown     = NewSpecificError(900, "未知错误")

	ChainClientUpToMaxCount  = NewSpecificError(901, "链客户端已达到数量上西安")
	ContractAddressNotFound  = NewSpecificError(902, "不支持的币种")
	AlreadyListenThisAddress = NewSpecificError(903, "该地址已经被监听了")
)

var (
	RangeConfigGetFailed    = NewSpecificError(1001, "范围配置获取失败")
	RangeConfigValueInvalid = NewSpecificError(1002, "范围配置值无效")
)

var (
	AuthFailed = NewSpecificError(2001, "认证失败")
)

var (
	UserNotFound                    = NewSpecificError(3001, "用户不存在")
	UserUsernameOrPasswordIncorrect = NewSpecificError(3002, "用户名或密码不正确")
	UserLoginGenerateTokenFailed    = NewSpecificError(3003, "生成token失败")
	GenerateHashPasswordFailed      = NewSpecificError(3004, "生成密码失败")
	UserCreateFailed                = NewSpecificError(3005, "创建用户失败")
)

var (
	AddressEditFailed   = NewSpecificError(4001, "更新地址失败")
	AddressCreateFailed = NewSpecificError(4002, "地址创建失败")
	AddressFindFailed   = NewSpecificError(4003, "地址查询失败")
	AddressCountFailed  = NewSpecificError(4004, "地址总数统计失败")
)

var (
	OrderFindFailed  = NewSpecificError(5001, "订单查询失败")
	OrderCountFailed = NewSpecificError(5002, "订单总数统计失败")
)

var (
	EvmOrderSaveFailed   = NewSpecificError(6001, "Evm订单保存失败")
	EvmLogSaveFailed     = NewSpecificError(6002, "Evm交易保存失败")
	EvmOrderNotifyFailed = NewSpecificError(6003, "Evm订单保存失败")
	EvmLogFindFailed     = NewSpecificError(6004, "Evm订单保存失败")
	EvmLogCountFailed    = NewSpecificError(6005, "Evm订单保存失败")
)

var (
	SolanaOrderSaveFailed = NewSpecificError(7001, "Solana订单保存失败")
)

var (
	TronOrderSaveFailed        = NewSpecificError(8001, "Tron订单保存失败")
	TronTransactionSaveFailed  = NewSpecificError(8002, "Tron交易保存失败")
	TronOrderNotifyFailed      = NewSpecificError(8003, "Tron订单通知失败")
	TronTransactionFindFailed  = NewSpecificError(8004, "Tron订单通知失败")
	TronTransactionCountFailed = NewSpecificError(8005, "Tron订单总数统计失败")
)
