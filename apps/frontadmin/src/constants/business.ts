import { transformNumberRecordToOption, transformRecordToOption } from '@/utils/common';

export const delegateTyposRecord: Record<Api.Common.DelegateTypos, App.I18n.I18nKey> = {
  '-1': 'page.transaction.common.all',
  0: 'page.transaction.common.bandwidth',
  1: 'page.transaction.common.energy'
};

export const exchangeTyposRecord: Record<Api.Common.ExchangeTypos, App.I18n.I18nKey> = {
  '': 'page.transaction.common.all',
  USDT2TRX: 'page.transaction.common.usdt2trx'
};

export const orderStatusRecord: Record<Api.Common.OrderStatus, App.I18n.I18nKey> = {
  '': 'page.transaction.common.all',
  已创建: 'page.transaction.common.created',
  已挂起: 'page.transaction.common.pending',
  已过期: 'page.transaction.common.expired',
  已取消: 'page.transaction.common.canceled',
  已委托: 'page.transaction.common.delegated',
  回收失败: 'page.transaction.common.withdrawfailed',
  错误: 'page.transaction.common.error',
  已完成: 'page.transaction.common.finished'
};

export const delegateTypoOptions = transformNumberRecordToOption(delegateTyposRecord);
export const exchangeTypoOptions = transformNumberRecordToOption(exchangeTyposRecord);

export const orderStatusOptions = transformRecordToOption(orderStatusRecord);
