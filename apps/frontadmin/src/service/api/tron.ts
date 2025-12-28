import { request } from '../request';

export function fetchGetTransList(params?: Api.Tron.TransSearchParams) {
  return request<Api.Tron.TransList>({
    url: '/tron/trans/list',
    method: 'get',
    params
  });
}

export function fetchGetTransDetail(id: number) {
  return request<Api.Tron.TransDetail>({
    url: `/tron/trans/get?id=${id}`,
    method: 'get'
  });
}
