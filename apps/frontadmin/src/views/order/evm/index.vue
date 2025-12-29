<script setup lang="tsx">
import { reactive } from 'vue';
import { currencyTyposRecord } from '@/constants/business';
import { fetchGetLogList } from '@/service/api';
import { defaultSearchform, useUIPaginatedTable } from '@/hooks/common/table';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';
import TransSearch from './modules/log-search.vue';

defineOptions({ name: 'TransSearch' });

const searchParams = reactive(getInitSearchParams());

function getInitSearchParams(): Api.Evm.LogSearchParams {
  return {
    current: 1,
    size: 20,
    id: 0,
    chain: '',
    currency: '',
    tx_hash: '',
    from_hex: '',
    to_hex: '',
    start: 0,
    end: 0
  };
}

const { columns, columnChecks, data, getData, getDataByPage, loading, mobilePagination } = useUIPaginatedTable({
  paginationProps: {
    currentPage: searchParams.current,
    pageSize: searchParams.size
  },
  api: () => fetchGetLogList(searchParams),
  transform: response => {
    return defaultSearchform(response);
  },
  onPaginationParamsChange: params => {
    searchParams.current = params.currentPage;
    searchParams.size = params.pageSize;
  },
  columns: () => [
    // { prop: 'selection', type: 'selection', width: 48 },
    { prop: 'id', type: 'id', label: $t('common.id') },
    { prop: 'chain', label: $t('page.order.common.chain') },
    {
      prop: 'currency',
      label: $t('page.order.common.currency'),
      width: 100,
      formatter: row => {
        const tagMap: Record<Api.Common.CurrencyTypos, UI.ThemeColor> = {
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
    { prop: 'amount', label: $t('page.order.common.amount'), width: 160 },
    { prop: 'sun', label: $t('page.order.common.sun'), width: 180 },
    { prop: 'contract', label: $t('page.order.common.contract'), width: 380 },
    { prop: 'tx_hash', label: $t('page.order.common.transaction_id'), width: 560 },
    { prop: 'from_hex', label: $t('page.order.common.from_address'), width: 380 },
    { prop: 'to_hex', label: $t('page.order.common.to_address'), width: 380 },
    { prop: 'block_timestamp', label: $t('page.order.common.block_timestamp'), width: 150 },
    { prop: 'removed', label: $t('page.order.common.removed'), width: 80 },
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
    // {
    //   prop: 'operate',
    //   fixed: true,
    //   label: $t('common.operate'),
    //   align: 'center',
    //   width: 160,
    //   formatter: row => (
    //     <div class="flex-center">
    //       <ElButton type="primary" plain size="small" onClick={() => bill(row.id)}>
    //         {$t('page.transaction.common.bill')}
    //       </ElButton>
    //       <ElButton type="primary" plain size="small" onClick={() => withdraweral(row.id)}>
    //         {$t('page.transaction.common.withdraweral')}
    //       </ElButton>
    //     </div>
    //   )
    // }
  ]
});

function resetSearchParams() {
  Object.assign(searchParams, getInitSearchParams());
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <TransSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <ElCard class="card-wrapper sm:flex-1-hidden" body-class="ht50">
      <template #header>
        <div class="flex items-center justify-between">
          <p>{{ $t('page.order.tron.title') }}</p>
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
