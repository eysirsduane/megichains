declare namespace Api {
  /**
   * namespace SystemManage
   *
   * backend api module: "systemManage"
   */
  namespace Fund {
    type AddressFund = Common.CommonRecord<{
      chain: string;
      address: string;
      tron_usdt: number;
      tron_usdc: number;
      bsc_usdt: number;
      bsc_usdc: number;
      eth_usdt: number;
      eth_usdc: number;
    }>;

    type AddressFundStatistics = Common.CommonRecord<{
      tron_usdt: number;
      tron_usdc: number;
      bsc_usdt: number;
      bsc_usdc: number;
      eth_usdt: number;
      eth_usdc: number;
    }>;

    type AddressFundCollect = Common.CommonRecord<{
      address_group_id: number;
      address_group_name: string;
      chain: string;
      currency: Api.Common.CurrencyTypos;
      status: Api.Common.CollectStatus;
      to_address?: string;
      amount_min: number;
      fee_max: number;
      tron_usdt?: number;
      tron_usdc?: number;
      bsc_usdt?: number;
      bsc_usdc?: number;
      eth_usdt?: number;
      eth_usdc?: number;
      description: string;
    }>;

    type AddressFundCollectLog = Common.CommonRecord<{
      collect_id: number;
      chain: string;
      currency: Api.Common.CurrencyTypos;
      status: Api.Common.CollectLogStatus;
      from_address: string;
      receiver_address: string;
      amount: number;
      transaction_id: string;
      gas_used: number;
      gas_price: number;
      effective_gas_price: number;
      total_gas_fee: number;
      description: string;
    }>;

    type AddressFundCollectCreating = Common.CommonPureRecord<{
      address_group_id: number;
      chain: string;
      currency: string;
      status: string;
      to_address?: string;
      amount_min: number;
      fee_max: number;
      secret_key: string;
      description: string;
    }>;

    type AddressFundSearchParams = CommonType.RecordNullable<
      Pick<Api.Fund.AddressFund, 'address' | 'chain' | 'id'> & Api.Common.CommonTimeSearchParams
    >;

    type AddressFundCollectListSearchParams = CommonType.RecordNullable<
      Pick<Api.Fund.AddressFundCollect, 'to_address' | 'address_group_id' | 'chain' | 'currency' | 'status'> &
        Api.Common.CommonTimeSearchParams
    >;

    type AddressFundCollectLogListSearchParams = CommonType.RecordNullable<
      Pick<
        Api.Fund.AddressFundCollectLog,
        'collect_id' | 'chain' | 'currency' | 'status' | 'from_address' | 'receiver_address'
      > &
        Api.Common.CommonTimeSearchParams
    >;

    /** user list */
    type AddressFundList = Common.PaginatingQueryRecord<AddressFund>;
    type AddressFundCollectList = Common.PaginatingQueryRecord<AddressFundCollect>;
    type AddressFundCollectLogList = Common.PaginatingQueryRecord<AddressFundCollectLog>;
  }
}
