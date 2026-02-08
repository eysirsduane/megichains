/**
 * Namespace Api
 *
 * All backend api type
 */
declare namespace Api {
  namespace Common {
    /** common params of paginating */
    interface PaginatingCommonParams {
      /** current current number */
      current: number;
      /** current size */
      size: number;
      /** total count */
      total: number;
    }

    interface PaginatingTimeCommonParams {
      /** current current number */
      current: number;
      /** current size */
      size: number;
      /** total count */
      total: number;
      start: number;
      end: number;
    }

    /** common params of paginating query list data */
    interface PaginatingQueryRecord<T = any> extends PaginatingCommonParams {
      records: T[];
    }

    interface QueryRecord<T = any> {
      records: T[];
    }

    /** common search params of table */
    type CommonSearchParams = Pick<Common.PaginatingCommonParams, 'current' | 'size'>;
    type CommonTimeSearchParams = Pick<Common.PaginatingTimeCommonParams, 'current' | 'size' | 'start' | 'end'>;

    /** common record */
    type CommonRecord<T = any> = {
      /** record id */
      id: number;
      /** record creator */
      created_at: number;
      /** record create time */
      /** record update time */
      updated_at: number;
      /** record status */
    } & T;

    type CommonPureRecord<T = any> = {
      /** record id */
      id: number;
    } & T;

    type EnableStatus = '1' | '2';
    type OrderTypos = '' | '输入';
    type CurrencyTypos = '' | 'USDT' | 'USDC';
    type ChainTypos = '' | 'TRON' | 'ETH' | 'BSC' | 'SOLANA';
    type ChainBigTypos = '' | 'TRON' | 'EVM' | 'SOLANA';
    type OrderStatus = '' | '已创建' | '通知失败' | '成功';
    type CollectStatus = '' | '已创建' | '处理中' | '部分成功' | '成功' | '失败';
    type CollectLogStatus = '' | '已创建' | '成功' | '失败';
    type AddressTypos = '' | 'IN' | 'OUT' | 'COLLECT';
    type AddressStatus = '' | '禁用' | '空闲' | '占用';
    type AddressGroupStatus = '' | '禁用' | '开放';
    type TronContractAddresses = '' | 'TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj' | 'TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf';
  }
}
