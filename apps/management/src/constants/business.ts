import { transformRecordToOption } from '@/utils/common';

export const chainTyposRecord: Record<Api.Common.ChainTypos, App.I18n.I18nKey> = {
  '': 'common.all',
  TRON: 'common.tron',
  ETH: 'common.eth',
  BSC: 'common.bsc',
  SOLANA: 'common.solana'
};

export const chainBigTyposRecord: Record<Api.Common.ChainBigTypos, App.I18n.I18nKey> = {
  '': 'common.all',
  TRON: 'common.tron',
  EVM: 'common.evm',
  SOLANA: 'common.solana'
};

export const currencyTyposRecord: Record<Api.Common.CurrencyTypos, App.I18n.I18nKey> = {
  '': 'common.all',
  USDT: 'common.usdt',
  USDC: 'common.usdc'
};

export const orderTyposRecord: Record<Api.Common.OrderTypos, App.I18n.I18nKey> = {
  '': 'common.all',
  输入: 'common.payin'
};

export const orderModesRecord: Record<Api.Common.OrderModes, App.I18n.I18nKey> = {
  '': 'common.all',
  正式: 'common.prod',
  测试: 'common.test'
};

export const orderStatusRecord: Record<Api.Common.OrderStatus, App.I18n.I18nKey> = {
  '': 'common.all',
  已创建: 'common.created',
  超时: 'common.timeout',
  失败: 'common.failed',
  成功: 'common.success'
};

export const notifyStatusRecord: Record<Api.Common.NotifyStatus, App.I18n.I18nKey> = {
  '': 'common.all',
  未知: 'common.unknown',
  失败: 'common.failed',
  成功: 'common.success'
};

export const collectStatusRecord: Record<Api.Common.CollectStatus, App.I18n.I18nKey> = {
  '': 'common.all',
  已创建: 'common.created',
  处理中: 'common.processing',
  部分成功: 'common.partial_success',
  成功: 'common.success',
  失败: 'common.failed'
};

export const collectLogStatusRecord: Record<Api.Common.CollectLogStatus, App.I18n.I18nKey> = {
  '': 'common.all',
  已创建: 'common.created',
  成功: 'common.success',
  失败: 'common.failed'
};

export const addressTyposRecord: Record<Api.Common.AddressTypos, App.I18n.I18nKey> = {
  '': 'common.all',
  IN: 'common.in',
  OUT: 'common.out',
  COLLECT: 'common.collect'
};

export const addressStatusRecord: Record<Api.Common.AddressStatus, App.I18n.I18nKey> = {
  '': 'common.all',
  禁用: 'common.ban',
  空闲: 'common.vacant',
  占用: 'common.inuse'
};
export const addressGroupStatusRecord: Record<Api.Common.AddressGroupStatus, App.I18n.I18nKey> = {
  '': 'common.all',
  禁用: 'common.ban',
  开放: 'common.open'
};

export const orderTypoOptions = transformRecordToOption(orderTyposRecord);
export const orderModesOptions = transformRecordToOption(orderModesRecord);
export const orderStatusOptions = transformRecordToOption(orderStatusRecord);
export const currencyTyposOptions = transformRecordToOption(currencyTyposRecord);
export const chainTyposOptions = transformRecordToOption(chainTyposRecord);
export const chainBigTyposOptions = transformRecordToOption(chainBigTyposRecord);
export const addressTyposOptions = transformRecordToOption(addressTyposRecord);
export const addressStatusOptions = transformRecordToOption(addressStatusRecord);
export const addressGroupStatusOptions = transformRecordToOption(addressGroupStatusRecord);
export const collectLogStatusOptions = transformRecordToOption(collectLogStatusRecord);
