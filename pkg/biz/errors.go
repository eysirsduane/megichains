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

var (
	//Delegate Order
	DelegateOrderCreateFailed      = NewSpecificError(1001, "租借订单创建失败")
	DelegateOrderDeleteFailed      = NewSpecificError(1002, "租借订单删除失败")
	DelegateOrderUpdateFailed      = NewSpecificError(1003, "租借订单更新失败")
	DelegateOrderBatchUpdateFailed = NewSpecificError(1004, "租借订单更新失败")
	DelegateOrderFindFailed        = NewSpecificError(1005, "租借订单查询失败")
	DelegateOrderExpiresFailed     = NewSpecificError(1006, "租借订单主动过期失败")
	DelegateOrderIsNil             = NewSpecificError(1006, "租借订单对象为无")
)

var (
	//Delegate Bill
	DelegateBillCreateFailed  = NewSpecificError(2001, "租借发货单创建失败")
	DelegateBillDeleteFailed  = NewSpecificError(2002, "租借发货单删除失败")
	DelegateBillUpdateFailed  = NewSpecificError(2003, "租借发货单更新失败")
	DelegateBillFindFailed    = NewSpecificError(2004, "租借发货单查询失败")
	DelegateBillExpiresFailed = NewSpecificError(2005, "租借发货单主动过期失败")
)

var (
	//Delegate TakeBack
	DelegateTakeBackCreateFailed  = NewSpecificError(20001, "租借收回单创建失败")
	DelegateTakeBackDeleteFailed  = NewSpecificError(20002, "租借收回单删除失败")
	DelegateTakeBackUpdateFailed  = NewSpecificError(20003, "租借收回单更新失败")
	DelegateTakeBackFindFailed    = NewSpecificError(20004, "租借收回单查询失败")
	DelegateTakeBackExpiresFailed = NewSpecificError(20005, "租借收回单主动过期失败")
)

var (
	//Exchange Order
	ExchangeOrderCreateFailed  = NewSpecificError(3001, "兑换订单创建失败")
	ExchangeOrderDeleteFailed  = NewSpecificError(3002, "兑换订单删除失败")
	ExchangeOrderUpdateFailed  = NewSpecificError(3003, "兑换订单更新失败")
	ExchangeOrderFindFailed    = NewSpecificError(3004, "兑换订单查询失败")
	ExchangeOrderExpiresFailed = NewSpecificError(3005, "兑换订单主动过期失败")
)

var (
	//Exchange Bill
	ExchangeBillCreateFailed  = NewSpecificError(4001, "兑换发货单订单创建失败")
	ExchangeBillDeleteFailed  = NewSpecificError(4002, "兑换发货单删除失败")
	ExchangeBillUpdateFailed  = NewSpecificError(4003, "兑换发货单更新失败")
	ExchangeBillFindFailed    = NewSpecificError(4004, "兑换发货单查询失败")
	ExchangeBillExpiresFailed = NewSpecificError(4005, "兑换发货单主动过期失败")
)

var (
	//Exchange Config
	ExchangeConfigGetFailed    = NewSpecificError(5001, "汇率配置获取失败")
	ExchangeConfigNotFound     = NewSpecificError(5002, "汇率配置未找到")
	ExchangeConfigValueInvalid = NewSpecificError(5003, "汇率配置数值不正确")
)

var (
	ExchangeTrx2UsdtRateIncorrect = NewSpecificError(6001, "TRX2USDT汇率获取失败")
)

// Auth & User
var (
	UserNotFound                    = NewSpecificError(7001, "用户不存在")
	UserUsernameOrPasswordIncorrect = NewSpecificError(7002, "用户名或密码不正确")
	UserLoginGenerateTokenFailed    = NewSpecificError(7003, "生成token失败")
	GenerateHashPasswordFailed      = NewSpecificError(7004, "生成密码失败")
	UserCreateFailed                = NewSpecificError(7005, "创建用户失败")
)
