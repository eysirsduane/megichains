<script setup lang="tsx">
import { reactive } from 'vue';
import { addressGroupStatusRecord } from '@/constants/business';
import { fetchGetAddressGroupList } from '@/service/api';
import { defaultSearchform, useUIPaginatedTable } from '@/hooks/common/table';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';
import TransSearch from './modules/address-group-search.vue';

defineOptions({ name: 'TransSearch' });

const searchParams = reactive(getInitSearchParams());

function getInitSearchParams(): Api.Address.AddressGroupSearchParams {
  return {
    current: 1,
    size: 20,
    id: 0,
    status: ''
  };
}

const { columns, columnChecks, data, getData, getDataByPage, loading, mobilePagination } = useUIPaginatedTable({
  paginationProps: {
    currentPage: searchParams.current,
    pageSize: searchParams.size
  },
  api: () => fetchGetAddressGroupList(searchParams),
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
    { prop: 'name', label: $t('common.name'), width: 160 },
    {
      prop: 'status',
      label: $t('page.address.common.status'),
      width: 100,
      formatter: row => {
        const tagMap: Record<Api.Common.AddressGroupStatus, UI.ThemeColor> = {
          '': 'info',
          禁用: 'danger',
          开放: 'success'
        };

        const label = $t(addressGroupStatusRecord[row.status]);

        return (
          <el-tag effect="dark" round type={tagMap[row.status]}>
            {label}
          </el-tag>
        );
      }
    },
    { prop: 'description', label: $t('common.description'), width: 800 },
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
