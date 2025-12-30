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
      status: Api.Common.AddressStatus;
      address: string;
      address2: string;
      description: string;
    }>;

    type AddressGroup = Common.CommonRecord<{
      name: string;
      status: Api.Common.AddressGroupStatus;
      description: string;
    }>;

    type AddressSearchParams = CommonType.RecordNullable<
      Pick<Api.Address.Address, 'address' | 'address2' | 'group_id' | 'chain' | 'typo' | 'status' | 'id'> &
        Api.Common.CommonTimeSearchParams
    >;

    type AddressGroupSearchParams = CommonType.RecordNullable<
      Pick<Api.Address.AddressGroup, 'status' | 'id'> & Api.Common.CommonTimeSearchParams
    >;

    /** user list */
    type AddressList = Common.PaginatingQueryRecord<Address>;
    type AddressGroupAll = Common.QueryRecord<AddressGroup>;
    type AddressGroupList = Common.PaginatingQueryRecord<AddressGroup>;
  }
}
