<script setup lang="tsx">
import { reactive, ref } from 'vue';
import { ElButton } from 'element-plus';
import { useBoolean } from '@sa/hooks';
import { delegateTyposRecord, orderStatusRecord } from '@/constants/business';
import { fetchGetDelegateOrderList } from '@/service/api';
import { defaultTransform, useUIPaginatedTable } from '@/hooks/common/table';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';
import DelegateBillDrawer from './modules/delegate-bill-drawer.vue';
import DelegateSearch from './modules/delegate-search.vue';
import DelegateWithdraweralDrawer from './modules/delegate-withdraweral-drawer.vue';

defineOptions({ name: 'DelegateView' });

const searchParams = reactive(getInitSearchParams());

function getInitSearchParams(): Api.Transaction.DelegateOrderSearchParams {
  return {
    current: 1,
    size: 20,
    status: '',
    transaction_id: '',
    typo: -1,
    currency: '',
    from_base58: '',
    to_base58: ''
  };
}

const { columns, columnChecks, data, getData, getDataByPage, loading, mobilePagination } = useUIPaginatedTable({
  paginationProps: {
    currentPage: searchParams.current,
    pageSize: searchParams.size
  },
  api: () => fetchGetDelegateOrderList(searchParams),
  transform: response => {
    return defaultTransform(response);
  },
  onPaginationParamsChange: params => {
    searchParams.current = params.currentPage;
    searchParams.size = params.pageSize;
  },
  columns: () => [
    { prop: 'selection', type: 'selection', width: 48 },
    { prop: 'id', type: 'id', label: $t('common.id') },
    { prop: 'user_id', label: $t('page.transaction.common.user_id') },
    {
      prop: 'typo',
      label: $t('page.transaction.common.typo'),
      formatter: row => {
        const tagMap: Record<Api.Common.DelegateTypos, UI.ThemeColor> = {
          '-1': 'danger',
          0: 'success',
          1: 'warning'
        };

        const label = $t(delegateTyposRecord[row.typo]);

        return (
          <el-tag effect="dark" round type={tagMap[row.typo]}>
            {label}
          </el-tag>
        );
      }
    },
    {
      prop: 'status',
      label: $t('page.transaction.common.status'),
      formatter: row => {
        const tagMap: Record<Api.Common.OrderStatus, UI.ThemeColor> = {
          '': 'info',
          已创建: 'info',
          已挂起: 'info',
          已过期: 'warning',
          已取消: 'info',
          已委托: 'primary',
          回收失败: 'danger',
          错误: 'danger',
          已完成: 'success'
        };

        const label = $t(orderStatusRecord[row.status]);

        return (
          <el-tag effect="dark" round type={tagMap[row.status]}>
            {label}
          </el-tag>
        );
      }
    },
    { prop: 'currency', label: $t('page.transaction.common.currency') },
    { prop: 'received_amount', label: $t('page.transaction.delegate.order.received_amount') },
    { prop: 'received_sun', label: $t('page.transaction.delegate.order.received_sun'), width: 180 },
    { prop: 'delegate_amount', label: $t('page.transaction.delegate.order.delegate_amount') },
    { prop: 'delegate_sun', label: $t('page.transaction.delegate.order.delegate_sun'), width: 180 },
    {
      prop: 'withdraw_time',
      label: $t('page.transaction.delegate.order.withdraw_time'),
      width: 180,
      formatter: row => {
        return getHumannessDateTime(row.withdraw_time);
      }
    },
    { prop: 'from_base58', label: $t('page.transaction.common.from_base58'), width: 320 },
    { prop: 'to_base58', label: $t('page.transaction.common.to_base58'), width: 320 },
    { prop: 'transaction_id', label: $t('page.transaction.common.transaction_id'), width: 570 },
    // { prop: 'description', label: $t('page.transaction.common.description'), width: 150 },
    {
      prop: 'created_at',
      label: $t('page.transaction.common.created_at'),
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
      width: 160,
      formatter: row => (
        <div class="flex-center">
          <ElButton type="primary" plain size="small" onClick={() => bill(row.id)}>
            {$t('page.transaction.common.bill')}
          </ElButton>
          <ElButton type="primary" plain size="small" onClick={() => withdraweral(row.id)}>
            {$t('page.transaction.common.withdraweral')}
          </ElButton>
        </div>
      )
    }
  ]
});

const { bool: drawerVisible, setTrue: openDrawer } = useBoolean();

const openWithdraweral = ref(false);
const targetId = ref(0);

function bill(id: number) {
  targetId.value = id;
  openDrawer();
}

function withdraweral(id: number) {
  targetId.value = id;
  openWithdraweral.value = true;
}

function resetSearchParams() {
  Object.assign(searchParams, getInitSearchParams());
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <DelegateSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <ElCard class="card-wrapper sm:flex-1-hidden" body-class="ht50">
      <template #header>
        <div class="flex items-center justify-between">
          <p>{{ $t('page.transaction.exchange.order.title') }}</p>
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
      <DelegateBillDrawer v-model:visible="drawerVisible" :target-id="targetId" />
      <DelegateWithdraweralDrawer v-model:visible="openWithdraweral" :target-id="targetId" />
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
