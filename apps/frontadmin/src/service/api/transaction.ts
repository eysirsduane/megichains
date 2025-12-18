import { request } from '../request';

export function fetchGetDelegateOrderList(params?: Api.Transaction.DelegateOrderSearchParams) {
  return request<Api.Transaction.DelegateOrderList>({
    url: '/delegate/order/list',
    method: 'get',
    params
  });
}

export function fetchGetDelegateBill(order_id: number) {
  return request<Api.Transaction.DelegateBill>({
    url: `/delegate/bill?order_id=${order_id}`,
    method: 'get'
  });
}

export function fetchGetDelegateWithdrawal(order_id: number) {
  return request<Api.Transaction.DelegateWithdrawal>({
    url: `/delegate/withdrawal?order_id=${order_id}`,
    method: 'get'
  });
}

export function fetchGetExchangeOrderList(params?: Api.Transaction.ExchangeOrderSearchParams) {
  return request<Api.Transaction.ExchangeOrderList>({
    url: '/exchange/order/list',
    method: 'get',
    params
  });
}

export function fetchGetExchangeBill(order_id: number) {
  return request<Api.Transaction.ExchangeBill>({
    url: `/exchange/bill?order_id=${order_id}`,
    method: 'get'
  });
}
