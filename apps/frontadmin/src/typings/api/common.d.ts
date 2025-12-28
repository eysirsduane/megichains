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

    /** common search params of table */
    type CommonSearchParams = Pick<Common.PaginatingCommonParams, 'current' | 'size'>;

    /**
     * enable status
     *
     * - "1": enabled
     * - "2": disabled
     */
    type EnableStatus = '1' | '2';
    type OrderTypos = '' | '入账';
    type CurrencyTypos = '' | 'USDT'|'USDC';
    type ChainTypos = '' | 'TRON'|'ETH'|'BSC';
    type OrderStatus = '' | '已创建' | '通知失败' | '成功' ;

    /** common record */
    type CommonRecord<T = any> = {
      /** record id */
      id: number;
      /** record creator */
      created_at: number;
      /** record create time */
      createTime: string;
      /** record updater */
      updateBy: string;
      /** record update time */
      updated_at: number;
      /** record status */
      status: OrderStatus;
    } & T;
  }
}
