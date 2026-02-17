<script setup lang="tsx">
import { onMounted, reactive, ref } from 'vue';
import { useBoolean } from '@sa/hooks';
import { collectStatusRecord, currencyTyposRecord } from '@/constants/business';
import { fetchGetAddressFundCollectList, fetchGetAddressGroupAll } from '@/service/api';
import { defaultSearchform, useUIPaginatedTable } from '@/hooks/common/table';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';
import AddressFundSearch from './modules/fund-collect-search.vue';
import AddressFundCollectDetailDrawer from './modules/fund-collect-detail-drawer.vue';

defineOptions({ name: 'AddressFundCollectLogList' });

const searchParams = reactive(getInitSearchParams());

function getInitSearchParams(): Api.Fund.AddressFundCollectListSearchParams {
  return {
    current: 1,
    size: 20,
    start: 0,
    end: 0,
    to_address: '',
    chain: '',
    status: '',
    currency: '',
    address_group_id: 0
  };
}

const { columns, columnChecks, data, getData, getDataByPage, loading, mobilePagination } = useUIPaginatedTable({
  paginationProps: {
    currentPage: searchParams.current,
    pageSize: searchParams.size
  },
  api: () => fetchGetAddressFundCollectList(searchParams),
  transform: response => {
    return defaultSearchform(response);
  },
  onPaginationParamsChange: params => {
    searchParams.current = params.currentPage;
    searchParams.size = params.pageSize;
  },
  columns: () => [
    // { prop: 'selection', type: 'selection', width: 48 },
    { prop: 'id', type: 'id', label: $t('common.id'), width: 100 },
    { prop: 'address_group_name', label: $t('page.fund.common.group'), width: 160 },
    { prop: 'chain', label: $t('page.fund.common.chain'), width: 80 },
    {
      prop: 'currency',
      label: $t('page.fund.common.currency'),
      width: 100,
      formatter: row => {
        const tagMap: Record<Api.Common.CurrencyTypo, UI.ThemeColor> = {
          '': 'info',
          USDT: 'info',
          USDC: 'success'
        };

        const label = $t(currencyTyposRecord[row.currency]);
        return (
          <el-tag effect="dark" round type={tagMap[row.currency]}>
            {label}
          </el-tag>
        );
      }
    },
    {
      prop: 'status',
      label: $t('page.fund.common.status'),
      width: 110,
      formatter: row => {
        const tagMap: Record<Api.Common.CollectStatus, UI.ThemeColor> = {
          '': 'info',
          已创建: 'info',
          处理中: 'info',
          部分成功: 'warning',
          成功: 'success',
          失败: 'danger'
        };

        const label = $t(collectStatusRecord[row.status]);
        return (
          <el-tag effect="dark" round type={tagMap[row.status]}>
            {label}
          </el-tag>
        );
      }
    },
    { prop: 'username', label: $t('common.username'), width: 160 },
    { prop: 'receiver_address', label: $t('page.fund.common.to_address'), width: 400 },
    { prop: 'success_amount', label: $t('page.fund.common.success_amount'), width: 200 },
    { prop: 'total_count', label: $t('page.fund.common.total_count'), width: 100 },
    { prop: 'success_count', label: $t('page.fund.common.success_count'), width: 100 },
    { prop: 'total_gas_fee', label: $t('page.fund.common.total_gas_fee'), width: 200 },
    { prop: 'total_gas_fee_currency', label: $t('page.fund.common.total_gas_fee_currency'), width: 180 },
    { prop: 'amount_min', label: $t('page.fund.common.collect_amount_min'), width: 120 },
    { prop: 'fee_max', label: $t('page.fund.common.fee_max'), width: 120 },
    { prop: 'description', label: $t('page.fund.common.description'), width: 200 },
    {
      prop: 'updated_at',
      label: $t('common.updated_at'),
      width: 180,
      formatter: row => {
        return getHumannessDateTime(row.updated_at);
      }
    },
    {
      prop: 'created_at',
      label: $t('common.created_at'),
      width: 180,
      formatter: row => {
        return getHumannessDateTime(row.created_at);
      }
    }
  ]
});

function resetSearchParams() {
  Object.assign(searchParams, getInitSearchParams());
}

const targetId = ref(0);

const { bool: drawerVisible, setTrue: openDrawer } = useBoolean();

function add() {
  targetId.value = 0;
  openDrawer();
}

const addrGroupOptions = ref<Api.Address.AddressGroup[] | undefined>();

onMounted(async () => {
  const all = await fetchGetAddressGroupAll();
  addrGroupOptions.value = all.data?.records;
});
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <AddressFundSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <ElCard class="card-wrapper sm:flex-1-hidden" body-class="ht50">
      <template #header>
        <div class="flex items-center justify-between">
          <p>{{ $t('page.fund.collect.title') }}</p>
          <TableHeaderOperation
            v-model:columns="columnChecks"
            :disabled-delete="true"
            :loading="loading"
            @add="add"
            @refresh="getData"
          />
        </div>
      </template>
      <div class="h-[calc(100%-50px)]">
        <ElTable
          v-loading="loading"
          height="100%"
          class="sm:h-full"
          :data="data"
          row-key="id"
          :border="true"
          highlight-current-row
        >
          <ElTableColumn v-for="col in columns" :key="col.prop" v-bind="col" />
        </ElTable>
      </div>
      <div class="mt-20px flex justify-end">
        <ElPagination
          v-if="mobilePagination.total"
          layout="total,prev,pager,next,sizes"
          v-bind="mobilePagination"
          @current-change="mobilePagination['current-change']"
          @size-change="mobilePagination['size-change']"
        />
      </div>
      <AddressFundCollectDetailDrawer v-model:visible="drawerVisible" :target-id="targetId" @saved="getDataByPage" />
    </ElCard>
  </div>
</template>

<style lang="scss" scoped>
:deep(.el-card) {
  .ht50 {
    height: calc(100% - 50px);
  }
}
</style>
