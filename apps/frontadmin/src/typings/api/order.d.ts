declare namespace Api {
  /**
   * namespace SystemManage
   *
   * backend api module: "systemManage"
   */
  namespace Order {
    type CommonSearchParams = Pick<Common.PaginatingCommonParams, 'current' | 'size'>;
    type CommonTimeSearchParams = Pick<Common.PaginatingTimeCommonParams, 'current' | 'size' | 'start' | 'end'>;

    type Order = Common.CommonRecord<{
      merch_order_id: string;
      transaction_id: string;
      typo: Api.Common.OrderTypos;
      currency: string;
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
        'chain' | 'transaction_id' | 'typo' | 'currency' | 'from_address' | 'to_address' | 'status' | 'id'
      > &
        CommonTimeSearchParams
    >;

    /** user list */
    type OrderList = Common.PaginatingQueryRecord<Order>;

    type OrderDetail = Common.CommonRecord<{
      order_id: number;
      transaction_id: string;
      currency: string;
      from_address: string;
      to_address: string;
      received_amount: number;
      received_sun: number;
      description: string;
    }>;
  }
}
