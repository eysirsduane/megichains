<script setup lang="tsx">
import { reactive, ref } from 'vue';
import { ElButton } from 'element-plus';
import { useBoolean } from '@sa/hooks';
import { collectLogStatusRecord, currencyTyposRecord } from '@/constants/business';
import { fetchGetAddressFundLogList } from '@/service/api';
import { defaultSearchform, useUIPaginatedTable } from '@/hooks/common/table';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';
import AddressFundLogSearch from './modules/fund-collectlog-search.vue';
import AddressFundCollectLogDetailDrawer from './modules/fund-collectlog-detail-drawer.vue';

defineOptions({ name: 'AddressFundCollectLogListView' });

const searchParams = reactive(getInitSearchParams());

function getInitSearchParams(): Api.Fund.AddressFundCollectLogListSearchParams {
  return {
    current: 1,
    size: 20,
    start: 0,
    end: 0,
    collect_id: 0,
    chain: '',
    currency: '',
    status: '',
    from_address: '',
    receiver_address: ''
  };
}

const { columns, columnChecks, data, getData, getDataByPage, loading, mobilePagination } = useUIPaginatedTable({
  paginationProps: {
    currentPage: searchParams.current,
    pageSize: searchParams.size
  },
  api: () => fetchGetAddressFundLogList(searchParams),
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
    { prop: 'collect_id', label: $t('page.fund.common.collect_id'), width: 100 },
    { prop: 'chain', label: $t('page.fund.common.chain'), width: 80 },
    {
      prop: 'currency',
      label: $t('page.fund.common.currency'),
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
    {
      prop: 'status',
      label: $t('page.fund.common.status'),
      width: 120,
      formatter: row => {
        const tagMap: Record<Api.Common.CollectLogStatus, UI.ThemeColor> = {
          '': 'info',
          已创建: 'info',
          成功: 'success',
          失败: 'danger'
        };

        const label = $t(collectLogStatusRecord[row.status]);
        return (
          <el-tag effect="dark" round type={tagMap[row.status]}>
            {label}
          </el-tag>
        );
      }
    },
    { prop: 'amount', label: $t('common.amount'), width: 200 },
    { prop: 'from_address', label: $t('common.from_address'), width: 400 },
    { prop: 'receiver_address', label: $t('common.receiver_address'), width: 400 },
    { prop: 'gas_used', label: $t('common.gas_used'), width: 200 },
    { prop: 'gas_price', label: $t('common.gas_price'), width: 200 },
    { prop: 'effective_gas_price', label: $t('common.effective_gas_price'), width: 200 },
    { prop: 'total_gas_fee', label: $t('common.total_gas_fee'), width: 200 },
    { prop: 'transaction_id', label: $t('common.transaction_id'), width: 600 },
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
    },
    {
      prop: 'operate',
      fixed: true,
      label: $t('common.operate'),
      align: 'center',
      width: 80,
      formatter: row => (
        <div class="flex-center">
          <ElButton plain type="primary" size="small" onClick={() => detail(row.id)}>
            {$t('page.order.common.detail')}
          </ElButton>
        </div>
      )
    }
  ]
});

const { bool: drawerVisible, setTrue: openDrawer } = useBoolean();
const targetId = ref(0);

function detail(id: number) {
  targetId.value = id;
  openDrawer();
}

function resetSearchParams() {
  Object.assign(searchParams, getInitSearchParams());
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <AddressFundLogSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <ElCard class="card-wrapper sm:flex-1-hidden" body-class="ht50">
      <template #header>
        <div class="flex items-center justify-between">
          <p>{{ $t('page.fund.collectlog.title') }}</p>
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
      <AddressFundCollectLogDetailDrawer v-model:visible="drawerVisible" :target-id="targetId" />
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
