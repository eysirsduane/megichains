import { transformRecordToOption } from '@/utils/common';

export const chainTyposRecord: Record<Api.Common.ChainTypos, App.I18n.I18nKey> = {
  '': 'common.all',
  TRON: 'common.tron',
  ETH: 'common.eth',
  BSC: 'common.bsc'
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

export const orderStatusRecord: Record<Api.Common.OrderStatus, App.I18n.I18nKey> = {
  '': 'common.all',
  已创建: 'common.created',
  通知失败: 'common.notifyfailed',
  成功: 'common.success'
};

export const orderTypoOptions = transformRecordToOption(orderTyposRecord);
export const orderStatusOptions = transformRecordToOption(orderStatusRecord);
export const currencyTyposOptions = transformRecordToOption(currencyTyposRecord);
