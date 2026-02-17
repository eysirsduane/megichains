declare namespace Api {
  /**
   * namespace SystemManage
   *
   * backend api module: "systemManage"
   */
  namespace Merch {
    type Merchant = Common.CommonRecord<{
      merchant_account: string;
      name: string;
      description: string;
    }>;

    type MerchantSearchParams = CommonType.RecordNullable<
      Pick<Api.Merch.Merchant, 'id' | 'merchant_account'> & Api.Common.CommonTimeSearchParams
    >;

    /** user list */
    type MerchantList = Common.PaginatingQueryRecord<Merchant>;

    type MerchantDetail = Common.CommonRecord<{
      merchant_account: string;
      name: string;
      description: string;
    }>;
  }
}
