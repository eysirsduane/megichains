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

export function fetchGetAddressFundCollectList(params?: Api.Fund.AddressFundCollectListSearchParams) {
  return request<Api.Fund.AddressFundCollectList>({
    url: '/fund/collect/list',
    method: 'get',
    params
  });
}

export function postCollectAddressFund(params?: Api.Fund.AddressFundCollectCreating) {
  return request<Api.Fund.AddressFundCollectList>({
    url: '/fund/collect',
    method: 'post',
    data: params
  });
}

export function fetchGetAddressFundLogList(params?: Api.Fund.AddressFundCollectLogListSearchParams) {
  return request<Api.Fund.AddressFundCollectLogList>({
    url: '/fund/collect/log/list',
    method: 'get',
    params
  });
}

export function fetchGetAddressFundCollectLogDetail(id: number) {
  return request<Api.Fund.AddressFundCollectLog>({
    url: `/fund/collect/log/${id}`,
    method: 'get'
  });
}
