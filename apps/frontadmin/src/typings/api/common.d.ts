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
    type DelegateTypos = -1 | 0 | 1;
    type ExchangeTypos = '' | 'USDT2TRX';
    type OrderStatus = '' | '已创建' | '已挂起' | '已过期' | '已取消' | '已委托' | '回收失败' | '错误' | '已完成';

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
      updateTime: string;
      /** record status */
      status: OrderStatus;
    } & T;
  }
}
