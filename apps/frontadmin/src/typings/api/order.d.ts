declare namespace Api {
  /**
   * namespace SystemManage
   *
   * backend api module: "systemManage"
   */
  namespace Order {
    type Order = Common.CommonRecord<{
      merch_order_id: string;
      transaction_id: string;
      typo: Api.Common.OrderTypos;
      currency: Api.Common.CurrencyTypos;
      status: Api.Common.OrderStatus;
      chain: string;
      received_amount: number;
      received_sun: number;
      from_address: string;
      to_address: string;
      description: string;
    }>;

    type OrderSearchParams = CommonType.RecordNullable<
      Pick<
        Api.Order.Order,
        | 'merch_order_id'
        | 'chain'
        | 'transaction_id'
        | 'typo'
        | 'currency'
        | 'from_address'
        | 'to_address'
        | 'status'
        | 'id'
      > &
        Api.Common.CommonTimeSearchParams
    >;

    /** user list */
    type OrderList = Common.PaginatingQueryRecord<Order>;

    type OrderDetail = Common.CommonRecord<{
      log_id: number;
      merch_order_id: string;
      transaction_id: string;
      typo: Api.Common.OrderTypos;
      status: Api.Common.OrderStatus;
      currency: string;
      chain: string;
      received_amount: number;
      received_sun: number;
      from_address: string;
      to_address: string;
      description: string;
    }>;
  }
}
