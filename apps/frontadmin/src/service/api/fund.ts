import { request } from '../request';

export function fetchGetAddressFundList(params?: Api.Fund.AddressFundSearchParams) {
  return request<Api.Fund.AddressFundList>({
    url: '/fund/list',
    method: 'get',
    params
  });
}

export function fetchGetAddressFundStatistics() {
  return request<Api.Fund.AddressFundStatistics>({
    url: '/fund/statistics',
    method: 'get'
  });
}
