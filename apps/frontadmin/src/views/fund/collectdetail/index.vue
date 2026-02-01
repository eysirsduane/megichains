<script setup lang="tsx">
import { reactive } from 'vue';
import { fetchGetAddressFundList } from '@/service/api';
import { defaultSearchform, useUIPaginatedTable } from '@/hooks/common/table';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';
import AddressFundSearch from './modules/fund-collect-search.vue';
// import AddressFundDetailDrawer from './modules/address-fund-detail-drawer.vue';

defineOptions({ name: 'AddressList' });

const searchParams = reactive(getInitSearchParams());

function getInitSearchParams(): Api.Address.AddressSearchParams {
  return {
    current: 1,
    size: 20,
    start: 0,
    end: 0,
    address: '',
    chain: ''
  };
}

const { columns, columnChecks, data, getData, getDataByPage, loading, mobilePagination } = useUIPaginatedTable({
  paginationProps: {
    currentPage: searchParams.current,
    pageSize: searchParams.size
  },
  api: () => fetchGetAddressFundList(searchParams),
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
    { prop: 'chain', label: $t('page.fund.common.chain'), width: 80 },
    { prop: 'address', label: $t('page.fund.common.address'), width: 400 },
    { prop: 'tron_usdt', label: $t('page.fund.common.tron_usdt'), width: 200 },
    { prop: 'tron_usdc', label: $t('page.fund.common.tron_usdc'), width: 200 },
    { prop: 'bsc_usdt', label: $t('page.fund.common.bsc_usdt'), width: 200 },
    { prop: 'bsc_usdc', label: $t('page.fund.common.bsc_usdc'), width: 200 },
    { prop: 'eth_usdt', label: $t('page.fund.common.eth_usdt'), width: 200 },
    { prop: 'eth_usdc', label: $t('page.fund.common.eth_usdc'), width: 200 },
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
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <AddressFundSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <ElCard class="card-wrapper sm:flex-1-hidden" body-class="ht50">
      <template #header>
        <div class="flex items-center justify-between">
          <p>{{ $t('page.address.list.title') }}</p>
          <TableHeaderOperation
            v-model:columns="columnChecks"
            :disabled-delete="true"
            :disabled-add="true"
            :loading="loading"
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
      <!-- <AddressFundDetailDrawer v-model:visible="drawerVisible" :target-id="targetId" @saved="getDataByPage" /> -->
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
