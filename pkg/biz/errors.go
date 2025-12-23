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
	ChainClientUpToMaxCount = NewSpecificError(3001, "链客户端已达到数量上西安")
	AddressCreateFailed     = NewSpecificError(3002, "地址创建失败")
)

var (
	EvmTransactionSaveFailed      = NewSpecificError(4001, "Evm交易失败")
	EvmTransactionDeleteFailed    = NewSpecificError(4002, "Evm交易删除失败")
	EvmTransactionUpdateFailed    = NewSpecificError(4003, "Evm交易更新失败")
	EvmTransactionFindFailed      = NewSpecificError(4004, "Evm交易查询失败")
	EvmTransactionContractInvalid = NewSpecificError(4005, "Evm交易合约地址无效")
	EvmTransactionStatusInvalid   = NewSpecificError(4006, "Evm交易状态无效")
	EvmOrderSaveFailed            = NewSpecificError(4007, "Evm订单保存失败")
)

var (
	BscTransactionSaveFailed      = NewSpecificError(4001, "BSC交易失败")
	BscTransactionDeleteFailed    = NewSpecificError(4002, "BSC交易删除失败")
	BscTransactionUpdateFailed    = NewSpecificError(4003, "BSC交易更新失败")
	BscTransactionFindFailed      = NewSpecificError(4004, "BSC交易查询失败")
	BscTransactionContractInvalid = NewSpecificError(4005, "BSC交易合约地址无效")
	BscTransactionStatusInvalid   = NewSpecificError(4006, "BSC交易状态无效")
	BscOrderSaveFailed            = NewSpecificError(4007, "BSC订单保存失败")
)
