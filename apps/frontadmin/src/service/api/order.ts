import { request } from '../request';

export function fetchGetOrderList(params?: Api.Order.OrderSearchParams) {
  return request<Api.Order.OrderList>({
    url: '/order/list',
    method: 'get',
    params
  });
}

export function fetchGetExchangeBill(id: number) {
  return request<Api.Order.OrderDetail>({
    url: `/order?id=${id}`,
    method: 'get'
  });
}
