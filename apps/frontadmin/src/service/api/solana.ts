import { request } from '../request';

export function getSolanaTransList(params?: Api.Solana.TransSearchParams) {
  return request<Api.Tron.TransList>({
    url: '/solana/trans/list',
    method: 'get',
    params
  });
}

export function getSolanaTransDetail(id: number) {
  return request<Api.Tron.TransDetail>({
    url: `/solana/trans/detail/${id}`,
    method: 'get'
  });
}
