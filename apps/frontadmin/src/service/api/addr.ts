import { request } from '../request';

export function fetchGetAddressList(params?: Api.Address.AddressSearchParams) {
  return request<Api.Address.AddressList>({
    url: '/address/list',
    method: 'get',
    params
  });
}

export function postEditAddress(params?: Api.Address.Address) {
  return request<Api.Address.AddressList>({
    url: '/address/edit',
    method: 'post',
    params
  });
}
