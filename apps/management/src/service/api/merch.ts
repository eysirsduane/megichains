import { request } from '../request';

export function findMerchantList(params?: Api.Merch.MerchantSearchParams) {
  return request<Api.Merch.MerchantList>({
    url: '/merchant/list',
    method: 'get',
    params
  });
}

export function getMerchantDetail(id: number) {
  return request<Api.Merch.MerchantDetail>({
    url: `/merchant/detail/${id}`,
    method: 'get'
  });
}

export function saveMerchantDetail(params?: Api.Merch.MerchantDetail) {
  return request<Api.Merch.MerchantDetail>({
    url: `/merchant/save`,
    method: 'post',
    data: params
  });
}
