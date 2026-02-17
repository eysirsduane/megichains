declare namespace Api {
  /**
   * namespace SystemManage
   *
   * backend api module: "systemManage"
   */
  namespace Evm {
    type Log = Common.CommonRecord<{
      chain: string;
      currency: Api.Common.CurrencyTypo;
      tx_hash: string;
      amount: number;
      sun: number;
      from_hex: string;
      to_hex: string;
      contract: string;
      block_timestamp: number;
      removed: boolean;
    }>;

    type LogSearchParams = CommonType.RecordNullable<
      Pick<Api.Evm.Log, 'chain' | 'currency'  | 'tx_hash' | 'from_hex' | 'to_hex' | 'id'> &
        Api.Common.CommonTimeSearchParams
    >;

    /** user list */
    type LogList = Common.PaginatingQueryRecord<Log>;

    type LogDetail = Common.CommonRecord<{
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
