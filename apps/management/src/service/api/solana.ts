import { request } from '../request';

export function getSolanaTransList(params?: Api.Solana.TransSearchParams) {
  return request<Api.Solana.TransList>({
    url: '/solana/trans/list',
    method: 'get',
    params
  });
}

export function getSolanaTransDetail(id: number) {
  return request<Api.Solana.TransDetail>({
    url: `/solana/trans/detail/${id}`,
    method: 'get'
  });
}
