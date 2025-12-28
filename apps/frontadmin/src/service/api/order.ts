import { request } from '../request';

export function fetchGetOrderList(params?: Api.Order.OrderSearchParams) {
  return request<Api.Order.OrderList>({
    url: '/order/list',
    method: 'get',
    params
  });
}

export function fetchGetOrderDetail(id: number) {
  return request<Api.Order.OrderDetail>({
    url: `/order/get?id=${id}`,
    method: 'get'
  });
}
