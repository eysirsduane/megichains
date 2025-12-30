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

export function fetchGetAddressGroupAll() {
  return request<Api.Address.AddressGroupAll>({
    url: '/address/group/all',
    method: 'get'
  });
}

export function fetchGetAddressGroupList(params?: Api.Address.AddressGroupSearchParams) {
  return request<Api.Address.AddressGroupList>({
    url: '/address/group/list',
    method: 'get',
    params
  });
}
