package biz

var (
	//系统基本错误码
	CodeSuccess     = NewSpecificError(0, "请求成功")
	CodeServerError = NewSpecificError(500, "服务器异常")
	CodeParamsEmpty = NewSpecificError(407, "参数为空")
	CodeParamsError = NewSpecificError(408, "参数错误")
	DatabaseError   = NewSpecificError(600, "数据库错误")
	CodeUnknown     = NewSpecificError(900, "未知错误")
)

var (
	//Trx Transfer
	TransferFailed = NewSpecificError(1000, "trx交易失败")
)

// Auth & User
var (
	UserNotFound                    = NewSpecificError(1001, "用户不存在")
	UserUsernameOrPasswordIncorrect = NewSpecificError(1002, "用户名或密码不正确")
	UserLoginGenerateTokenFailed    = NewSpecificError(1003, "生成token失败")
	GenerateHashPasswordFailed      = NewSpecificError(1004, "生成密码失败")
	UserCreateFailed                = NewSpecificError(1005, "创建用户失败")
)

var (
	RangeConfigGetFailed    = NewSpecificError(2001, "范围配置获取失败")
	RangeConfigValueInvalid = NewSpecificError(2002, "范围配置值无效")
)

var (
	ChainClientUpToMaxCount  = NewSpecificError(3001, "链客户端已达到数量上西安")
	AddressCreateFailed      = NewSpecificError(3002, "地址创建失败")
	ContractAddressNotFound  = NewSpecificError(3003, "不支持的币种")
	AlreadyListenThisAddress = NewSpecificError(3004, "该地址已经被监听了")
)

var (
	EvmOrderSaveFailed   = NewSpecificError(4001, "Evm订单保存失败")
	EvmLogSaveFailed     = NewSpecificError(4002, "Evm交易保存失败")
	EvmOrderNotifyFailed = NewSpecificError(4003, "Evm订单保存失败")
)

var (
	SolanaOrderSaveFailed = NewSpecificError(5001, "Solana订单保存失败")
)

var (
	TronOrderSaveFailed   = NewSpecificError(6001, "Tron订单保存失败")
	TronTransSaveFailed   = NewSpecificError(6002, "Tron交易保存失败")
	TronOrderNotifyFailed = NewSpecificError(6003, "Tron订单通知失败")
)
