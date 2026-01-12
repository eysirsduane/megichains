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

    type AddressFundSearchParams = CommonType.RecordNullable<
      Pick<Api.Fund.AddressFund, 'address' | 'chain' | 'id'> & Api.Common.CommonTimeSearchParams
    >;

    /** user list */
    type AddressFundList = Common.PaginatingQueryRecord<AddressFund>;
  }
}
