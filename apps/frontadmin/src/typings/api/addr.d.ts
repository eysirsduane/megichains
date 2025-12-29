declare namespace Api {
  /**
   * namespace SystemManage
   *
   * backend api module: "systemManage"
   */
  namespace Address {
    type Address = Common.CommonRecord<{
      group_id: number;
      chain: string;
      typo: Api.Common.AddressTypos;
      status: string;
      address: string;
      address2: string;
      description: string;
    }>;

    type AddressSearchParams = CommonType.RecordNullable<
      Pick<Api.Address.Address, 'address' | 'address2' | 'group_id' | 'chain' | 'typo' | 'status' | 'id'> &
        Api.Common.CommonTimeSearchParams
    >;

    /** user list */
    type AddressList = Common.PaginatingQueryRecord<Address>;
  }
}
