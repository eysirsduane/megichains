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

export function fetchGetAddressFundLogList(params?: Api.Fund.AddressFundCollectLogSearchParams) {
  return request<Api.Fund.AddressFundCollectLogList>({
    url: '/fund/collect/log/list',
    method: 'get',
    params
  });
}

export function postCollectAddressFund(params?: Api.Fund.AddressFundCollect) {
  return request<Api.Fund.AddressFundCollectLogList>({
    url: '/fund/collect',
    method: 'post',
    params
  });
}
