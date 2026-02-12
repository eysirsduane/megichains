import { request } from '../request';

export function getTronTransList(params?: Api.Tron.TransSearchParams) {
  return request<Api.Tron.TransList>({
    url: '/tron/trans/list',
    method: 'get',
    params
  });
}

export function getTronTransDetail(id: number) {
  return request<Api.Tron.TransDetail>({
    url: `/tron/trans/detail?id=${id}`,
    method: 'get'
  });
}
