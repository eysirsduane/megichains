declare namespace Api {
  /**
   * namespace SystemManage
   *
   * backend api module: "systemManage"
   */
  namespace Transaction {
    type CommonSearchParams = Pick<Common.PaginatingCommonParams, 'current' | 'size'>;
    type CommonTimeSearchParams = Pick<Common.PaginatingTimeCommonParams, 'current' | 'size' | 'start' | 'end'>;

    type UserGender = '1' | '2';

    /** user */
    type DelegateOrder = Common.CommonRecord<{
      user_id: number;
      transaction_id: string;
      typo: Api.Common.DelegateTypos;
      currency: string;
      received_amount: number;
      received_sun: number;
      from_base58: string;
      to_base58: string;
      delegate_amount: number;
      withdraw_time: number;
      telegram: string;
      whatsapp: string;
      wechat: string;
      other: string;
      description: string;
    }>;

    type ExchangeOrder = Common.CommonRecord<{
      user_id: number;
      transaction_id: string;
      typo: Api.Common.ExchangeTypos;
      currency: string;
      received_amount: number;
      received_sun: number;
      from_base58: string;
      to_base58: string;
      exchange_amount: number;
      exchange_sun: number;
      then_rate: number;
      exchange_rate: number;
      exchange_discount: number;
      telegram: string;
      whatsapp: string;
      wechat: string;
      other: string;
      description: string;
    }>;

    /** user search params */
    type DelegateOrderSearchParams = CommonType.RecordNullable<
      Pick<
        Api.Transaction.DelegateOrder,
        'transaction_id' | 'typo' | 'currency' | 'from_base58' | 'to_base58' | 'status' | 'id'
      > &
        CommonTimeSearchParams
    >;

    type ExchangeOrderSearchParams = CommonType.RecordNullable<
      Pick<
        Api.Transaction.ExchangeOrder,
        'transaction_id' | 'typo' | 'currency' | 'from_base58' | 'to_base58' | 'status' | 'id'
      > &
        CommonTimeSearchParams
    >;

    /** user list */
    type DelegateOrderList = Common.PaginatingQueryRecord<DelegateOrder>;
    type ExchangeOrderList = Common.PaginatingQueryRecord<ExchangeOrder>;

    type DelegateBill = Common.CommonRecord<{
      user_id: number;
      order_id: number;
      transaction_id: string;
      currency: string;
      from_base58: string;
      to_base58: string;
      delegated_amount: number;
      delegated_sun: number;
      description: string;
    }>;

    type ExchangeBill = Common.CommonRecord<{
      user_id: number;
      order_id: number;
      transaction_id: string;
      currency: string;
      from_base58: string;
      to_base58: string;
      exchanged_amount: number;
      exchanged_sun: number;
      description: string;
    }>;

    type DelegateWithdrawal = Common.CommonRecord<{
      user_id: number;
      order_id: number;
      transaction_id: string;
      typo: number;
      from_base58: string;
      to_base58: string;
      un_delegated_amount: number;
      un_delegated_sun: number;
      description: string;
    }>;
  }
}
