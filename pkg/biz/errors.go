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
	ConvertClientFailed      = NewSpecificError(904, "转换客户端失败")
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
	UserGetFailed                   = NewSpecificError(3006, "获取用户信息失败")
	UserListGetFailed               = NewSpecificError(3007, "获取用户列表失败")
	UserListCountFailed             = NewSpecificError(3008, "统计用户列表失败")
)

var (
	AddressSaveFailed        = NewSpecificError(4001, "保存地址失败")
	AddressCreateFailed      = NewSpecificError(4002, "地址创建失败")
	AddressFindFailed        = NewSpecificError(4003, "地址查询失败")
	AddressCountFailed       = NewSpecificError(4004, "地址总数统计失败")
	AddressGroupFindFailed   = NewSpecificError(4005, "获取地址分组失败")
	AddressGroupSaveFailed   = NewSpecificError(4006, "地址分组保存失败")
	AddressGroupFieldInvalid = NewSpecificError(4007, "地址分组字段非法")
	AddressFindGroupFailed   = NewSpecificError(4008, "获取地址分组失败")
)

var (
	OrderFindFailed               = NewSpecificError(5001, "订单查询失败")
	OrderCountFailed              = NewSpecificError(5002, "订单总数统计失败")
	OrderSaveFailed               = NewSpecificError(5003, "订单保存失败")
	OrderNotifyLogSaveFailed      = NewSpecificError(5004, "订单通知日志保存失败")
	OrderFindInteractionFailed    = NewSpecificError(5005, "订单交互日志查询失败")
	OrderInteractionLogSaveFailed = NewSpecificError(5006, "订单交互日志保存失败")
	OrderTestPlaceFailed          = NewSpecificError(5007, "订单测试下单失败")
	OrderAlreadyExistsSameListen  = NewSpecificError(5008, "已存在相同的监听订单")
	OrderInteractionLogGetFailed  = NewSpecificError(5009, "订单交互日志获取失败")
)

var (
	EvmOrderSaveFailed   = NewSpecificError(6001, "Evm订单保存失败")
	EvmLogSaveFailed     = NewSpecificError(6002, "Evm交易保存失败")
	EvmOrderNotifyFailed = NewSpecificError(6003, "Evm订单通知失败")
	EvmLogFindFailed     = NewSpecificError(6004, "Evm交易查询失败")
	EvmLogCountFailed    = NewSpecificError(6005, "Evm交易总数统计失败")
)

var (
	SolanaOrderSaveFailed        = NewSpecificError(7001, "Solana订单保存失败")
	SolanaTransactionFindFailed  = NewSpecificError(7002, "Solana交易查询失败")
	SolanaTransactionCountFailed = NewSpecificError(7003, "Solana交易总数统计失败")
	SolanaTransactionSaveFailed  = NewSpecificError(8001, "Solana交易保存失败")
)

var (
	TronOrderSaveFailed        = NewSpecificError(8001, "Tron订单保存失败")
	TronTransactionSaveFailed  = NewSpecificError(8002, "Tron交易保存失败")
	TronOrderNotifyFailed      = NewSpecificError(8003, "Tron订单通知失败")
	TronTransactionFindFailed  = NewSpecificError(8004, "Tron交易通知失败")
	TronTransactionCountFailed = NewSpecificError(8005, "Tron交易总数统计失败")
)

var (
	AddressFundCollectUnknownChain             = NewSpecificError(901, "地址资金归集未知的链")
	AddressFundCollectLogFindFailed            = NewSpecificError(902, "地址资金归集日志查询失败")
	AddressFundCollectLogCountFailed           = NewSpecificError(903, "地址资金归集日志统计失败")
	AddressFundCollectLogCreateFailed          = NewSpecificError(904, "地址资金归集日志创建失败")
	AddressFundCollectGetContractAddressFailed = NewSpecificError(905, "地址资金归集获取合约地址失败")
	AddressFundCollectParseABIFailed           = NewSpecificError(906, "地址资金归集解析ABI失败")
	AddressFundCollectPackTransferFailed       = NewSpecificError(907, "地址资金归集打包Transfer失败")
	AddressFundCollectSuggestGasPriceFailed    = NewSpecificError(908, "地址资金归集获取建议Gas价格失败")
	AddressFundCollectSignTxFailed             = NewSpecificError(909, "地址资金归集交易签名失败")
	AddressFundCollectSendTxFailed             = NewSpecificError(910, "地址资金归集交易发送失败")
	AddressFundCollectLogUpdateFailed          = NewSpecificError(911, "地址资金归集日志更新失败")
	AddressFundCollectReceiverAddressNotFound  = NewSpecificError(912, "地址资金归集接收地址未找到")
	AddressFundCollectFromAddressNotFound      = NewSpecificError(913, "地址资金归集发送地址未找到")
	AddressFundCollectPrivateKeyInvalid        = NewSpecificError(914, "地址资金归集私钥无效")
	addressFundCollectInsufficientBalance      = NewSpecificError(915, "地址资金归集余额不足")
	AddressFundCollectGetNonceFailed           = NewSpecificError(916, "地址资金归集获取Nonce失败")
	AddressFundCollectWaitMinedFailed          = NewSpecificError(917, "地址资金归集等待交易确认失败")
	AddressFundCollectNewErc20InstanceFailed   = NewSpecificError(918, "地址资金归集创建ERC20实例失败")
	AddressFundCollectEncodeTransferDataFailed = NewSpecificError(919, "地址资金归集编码Transfer数据失败")
	AddressFundCollectErc20TransferFailed      = NewSpecificError(920, "地址资金归集ERC20转账失败")
	AddressFundCollectSuggestGasTipCapFailed   = NewSpecificError(921, "地址资金归集获取建议Gas小费失败")
	AddressFundCollectGetHeaderFailed          = NewSpecificError(922, "地址资金归集获取区块头失败")
	AddressFundCollectGetChainIdFailed         = NewSpecificError(923, "地址资金归集获取链ID失败")
	AddressFundCollectPrivateKeyDecryptFailed  = NewSpecificError(924, "地址资金归集私钥解密失败")
	AddressFundCollectEstimateGasFailed        = NewSpecificError(925, "地址资金归集估算Gas失败")
	AddressFundCollectInvalidGasTipCap         = NewSpecificError(926, "地址资金归集无效的Gas小费")
	AddressFundCollectFeeOverLimit             = NewSpecificError(927, "地址资金归集手续费超过上限")
	AddressFundCollectFindFailed               = NewSpecificError(928, "地址资金归集列表查询失败")
	AddressFundCollectCountFailed              = NewSpecificError(929, "地址资金归集总数查询失败")
	AddressFundCollectInitClientFailed         = NewSpecificError(930, "地址资金归集初始化客户端失败")
	AddressFundCollectSaveFailed               = NewSpecificError(931, "地址资金归集保存失败")
	AddressFundCollectLogGetFailed             = NewSpecificError(932, "地址资金归集日志详情获取失败")
	AddressFundCollectGetAccountBalanceFailed  = NewSpecificError(933, "地址资金归集获取账户余额失败")
)

var (
	MerchantSaveFailed   = NewSpecificError(9001, "商户信息保存失败")
	MerchantCreateFailed = NewSpecificError(8005, "商户信息创建失败")
)
