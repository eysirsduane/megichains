import { request } from '../request';

export function findOrderList(params?: Api.Order.OrderSearchParams) {
  return request<Api.Order.OrderList>({
    url: '/order/list',
    method: 'get',
    params
  });
}

export function fetchGetOrderDetail(id: number) {
  return request<Api.Order.OrderDetail>({
    url: `/order/detail?id=${id}`,
    method: 'get'
  });
}

export function postOrderTestPlace(params?: Api.Order.OrderTestPlace) {
  return request<null>({
    url: `/order/test/place`,
    method: 'post',
    data: params
  });
}
