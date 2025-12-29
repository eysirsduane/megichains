declare namespace Api {
  /**
   * namespace SystemManage
   *
   * backend api module: "systemManage"
   */
  namespace Tron {
    type Trans = Common.CommonRecord<{
      chain: string;
      currency: Api.Common.CurrencyTypos;
      transaction_id: string;
      amount: number;
      sun: number;
      from_base58: string;
      to_base58: string;
      contract: string;
      block_timestamp: number;
    }>;

    type TransSearchParams = CommonType.RecordNullable<
      Pick<Api.Tron.Trans, 'currency' |  'transaction_id' | 'from_base58' | 'to_base58' | 'id'> &
        Api.Common.CommonTimeSearchParams
    >;

    /** user list */
    type TransList = Common.PaginatingQueryRecord<Trans>;

    type TransDetail = Common.CommonRecord<{
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
