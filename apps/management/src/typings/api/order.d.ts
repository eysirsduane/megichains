declare namespace Api {
  /**
   * namespace SystemManage
   *
   * backend api module: "systemManage"
   */
  namespace Order {
    type Order = Common.CommonRecord<{
      merchant_account: string;
      order_no: string;
      merchant_order_no: string;
      transaction_id: string;
      typo: Api.Common.OrderTypo;
      mode: Api.Common.OrderMode;
      currency: Api.Common.CurrencyTypo;
      status: Api.Common.OrderStatus;
      notify_status: Api.Common.NotifyStatus;
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
        | 'merchant_account'
        | 'order_no'
        | 'merchant_order_no'
        | 'chain'
        | 'transaction_id'
        | 'typo'
        | 'mode'
        | 'currency'
        | 'from_address'
        | 'to_address'
        | 'status'
        | 'notify_status'
        | 'id'
      > &
        Api.Common.CommonTimeSearchParams
    >;

    /** user list */
    type OrderList = Common.PaginatingQueryRecord<Order>;

    type OrderDetail = Common.CommonRecord<{
      log_id: number;
      order_no: string;
      merchant_account: string;
      merchant_order_no: string;
      transaction_id: string;
      typo: Api.Common.OrderTypo;
      status: Api.Common.OrderStatus;
      currency: string;
      chain: string;
      received_amount: number;
      received_sun: number;
      from_address: string;
      to_address: string;
      description: string;
    }>;

    type OrderTestPlace = {
      merchant_account: string;
      chain: Api.Common.ChainTypo;
      currency: Api.Common.CurrencyTypo;
      mode: Api.Common.OrderMode;
      receiver: string;
      seconds: number;
      notify_url: string;
    };

    type OrderInteraction = Common.CommonRecord<{
      id: number;
      merchant_order_id: number;
      place_request: string;
      place_request_timestamp: number;
      place_response: string;
      place_response_timestamp: number;
      notify_request: string;
      notify_request_timestamp: number;
      notify_response: string;
      notify_response_timestamp: number;
      description: string;
    }>;
  }
}
