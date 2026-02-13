<script setup lang="tsx">
import { reactive, ref } from 'vue';
import { ElButton } from 'element-plus';
import { useBoolean } from '@sa/hooks';
import { currencyTyposRecord, orderStatusRecord, orderTyposRecord } from '@/constants/business';
import { fetchGetOrderList } from '@/service/api';
import { defaultSearchform, useUIPaginatedTable } from '@/hooks/common/table';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';
import OrderDetailDrawer from './modules/order-detail-drawer.vue';
import OrderSearch from './modules/order-search.vue';

defineOptions({ name: 'OrderListView' });

const searchParams = reactive(getInitSearchParams());

function getInitSearchParams(): Api.Order.OrderSearchParams {
  return {
    current: 1,
    size: 20,
    start: 0,
    end: 0,
    merch_order_id: '',
    chain: '',
    transaction_id: '',
    typo: '',
    currency: '',
    from_address: '',
    to_address: '',
    status: ''
  };
}

const { columns, columnChecks, data, getData, getDataByPage, loading, mobilePagination } = useUIPaginatedTable({
  paginationProps: {
    currentPage: searchParams.current,
    pageSize: searchParams.size
  },
  api: () => fetchGetOrderList(searchParams),
  transform: response => {
    return defaultSearchform(response);
  },
  onPaginationParamsChange: params => {
    searchParams.current = params.currentPage;
    searchParams.size = params.pageSize;
  },
  columns: () => [
    { prop: 'selection', type: 'selection', width: 48 },
    { prop: 'id', type: 'id', label: $t('common.id') },
    { prop: 'merch_order_id', label: $t('page.order.common.merch_order_id'), width: 320 },
    { prop: 'chain', label: $t('page.order.common.chain') },
    {
      prop: 'typo',
      label: $t('page.order.common.typo'),
      width: 100,
      formatter: row => {
        const tagMap: Record<Api.Common.OrderTypos, UI.ThemeColor> = {
          '': 'info',
          输入: 'success'
        };

        const label = $t(orderTyposRecord[row.typo]);

        return (
          <el-tag effect="dark" round type={tagMap[row.typo]}>
            {label}
          </el-tag>
        );
      }
    },
    {
      prop: 'status',
      label: $t('page.order.common.status'),
      width: 100,
      formatter: row => {
        const tagMap: Record<Api.Common.OrderStatus, UI.ThemeColor> = {
          '': 'info',
          已创建: 'info',
          通知失败: 'danger',
          成功: 'success'
        };

        const label = $t(orderStatusRecord[row.status]);

        return (
          <el-tag effect="dark" round type={tagMap[row.status]}>
            {label}
          </el-tag>
        );
      }
    },
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
    { prop: 'received_amount', label: $t('page.order.common.received_amount'), width: 160 },
    // { prop: 'received_sun', label: $t('page.order.common.received_sun'), width: 180 },
    { prop: 'from_address', label: $t('page.order.common.from_address'), width: 400 },
    { prop: 'to_address', label: $t('page.order.common.to_address'), width: 400 },
    // { prop: 'transaction_id', label: $t('page.order.common.transaction_id'), width: 600 },
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
          <ElButton plain type="primary" size="small" onClick={() => bill(row.id)}>
            {$t('page.order.common.detail')}
          </ElButton>
        </div>
      )
    }
  ]
});

const { bool: drawerVisible, setTrue: openDrawer } = useBoolean();

const targetId = ref(0);

function bill(id: number) {
  targetId.value = id;
  openDrawer();
}

function resetSearchParams() {
  Object.assign(searchParams, getInitSearchParams());
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <OrderSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <ElCard class="card-wrapper sm:flex-1-hidden" body-class="ht50">
      <template #header>
        <div class="flex items-center justify-between">
          <p>{{ $t('page.order.common.title') }}</p>
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
      <OrderDetailDrawer v-model:visible="drawerVisible" :target-id="targetId" />
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
