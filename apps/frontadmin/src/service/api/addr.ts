import { request } from '../request';

export function fetchGetAddressList(params?: Api.Address.AddressSearchParams) {
  return request<Api.Address.AddressList>({
    url: '/address/list',
    method: 'get',
    params
  });
}

export function postSaveAddress(params?: Api.Address.Address) {
  return request<Api.Address.Address>({
    url: '/address/save',
    method: 'post',
    data: params
  });
}

export function fetchGetAddressDetail(id?: number) {
  return request<Api.Address.Address>({
    url: `/address/detail?id=${id}`,
    method: 'get'
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

export function fetchGetAddressGroupDetail(id: number) {
  return request<Api.Address.AddressGroup>({
    url: `/address/group/detail?id=${id}`,
    method: 'get'
  });
}

export function postSaveAddressGroup(params?: Api.Address.AddressGroup) {
  return request<Api.Address.AddressGroup>({
    url: `/address/group/save`,
    method: 'post',
    data: params
  });
}

export function postCreateAddressGroup(params?: Api.Address.AddressGroup) {
  return request<Api.Address.AddressGroup>({
    url: `/address/group/create`,
    method: 'post',
    data: params
  });
}

export function postGenerateAddress(params?: Api.Address.AddressGenerate) {
  return request<Api.Address.AddressGenerate>({
    url: '/address/generate',
    method: 'post',
    data: params
  });
}
