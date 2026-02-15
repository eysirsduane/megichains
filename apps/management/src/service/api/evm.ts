import { request } from '../request';

export function fetchGetLogList(params?: Api.Evm.LogSearchParams) {
  return request<Api.Evm.LogList>({
    url: '/evm/log/list',
    method: 'get',
    params
  });
}
