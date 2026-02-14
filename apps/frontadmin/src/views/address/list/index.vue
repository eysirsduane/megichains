<script setup lang="tsx">
import { reactive, ref } from 'vue';
import { ElButton } from 'element-plus';
import { useBoolean } from '@sa/hooks';
import { addressStatusRecord, addressTyposRecord } from '@/constants/business';
import { fetchGetAddressList } from '@/service/api';
import { defaultSearchform, useUIPaginatedTable } from '@/hooks/common/table';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';
import AddressSearch from './modules/address-search.vue';
import AddressDetailDrawer from './modules/address-detail-drawer.vue';
import AddressGenerateDrawer from './modules/address-generate-drawer.vue';

defineOptions({ name: 'AddressList' });

const searchParams = reactive(getInitSearchParams());

function getInitSearchParams(): Api.Address.AddressSearchParams {
  return {
    current: 1,
    size: 20,
    start: 0,
    end: 0,
    address: '',
    address2: '',
    group_id: 0,
    chain: '',
    typo: ''
  };
}

const { columns, columnChecks, data, getData, getDataByPage, loading, mobilePagination } = useUIPaginatedTable({
  paginationProps: {
    currentPage: searchParams.current,
    pageSize: searchParams.size
  },
  api: () => fetchGetAddressList(searchParams),
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
    { prop: 'group_name', label: $t('page.address.common.group_name'), width: 120 },
    { prop: 'chain', label: $t('page.address.common.chain'), width: 80 },
    {
      prop: 'typo',
      label: $t('page.address.common.typo'),
      width: 80,
      formatter: row => {
        const tagMap: Record<Api.Common.AddressTypos, UI.ThemeColor> = {
          '': 'info',
          IN: 'success',
          OUT: 'warning',
          COLLECT: 'primary'
        };

        const label = $t(addressTyposRecord[row.typo]);
        return (
          <el-tag effect="dark" round type={tagMap[row.typo]}>
            {label}
          </el-tag>
        );
      }
    },
    {
      prop: 'status',
      label: $t('page.address.common.status'),
      width: 80,
      formatter: row => {
        const tagMap: Record<Api.Common.AddressStatus, UI.ThemeColor> = {
          '': 'info',
          禁用: 'danger',
          空闲: 'success',
          占用: 'warning'
        };

        const label = $t(addressStatusRecord[row.status]);
        return (
          <el-tag effect="dark" round type={tagMap[row.status]}>
            {label}
          </el-tag>
        );
      }
    },
    { prop: 'address', label: $t('page.address.common.address'), width: 400 },
    { prop: 'tron_usdt', label: $t('page.fund.common.tron_usdt'), width: 200 },
    { prop: 'tron_usdc', label: $t('page.fund.common.tron_usdc'), width: 200 },
    { prop: 'bsc_usdt', label: $t('page.fund.common.bsc_usdt'), width: 200 },
    { prop: 'bsc_usdc', label: $t('page.fund.common.bsc_usdc'), width: 200 },
    { prop: 'eth_usdt', label: $t('page.fund.common.eth_usdt'), width: 200 },
    { prop: 'eth_usdc', label: $t('page.fund.common.eth_usdc'), width: 200 },
    { prop: 'solana_usdt', label: $t('page.fund.common.solana_usdt'), width: 200 },
    { prop: 'solana_usdc', label: $t('page.fund.common.solana_usdc'), width: 200 },
    { prop: 'address2', label: $t('page.address.common.address2'), width: 400 },
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
          <ElButton type="primary" plain size="small" onClick={() => edit(row.id)}>
            {$t('common.edit')}
          </ElButton>
          {/* <ElButton type="primary" plain size="small" onClick={() => withdraweral(row.id)}>
            {$t('page.address.common.withdraweral')}
          </ElButton> */}
        </div>
      )
    }
  ]
});

function resetSearchParams() {
  Object.assign(searchParams, getInitSearchParams());
}

const targetId = ref(0);
const { bool: drawerVisible, setTrue: openDrawer } = useBoolean();
const { bool: drawerGenerateVisible, setTrue: openGenerateDrawer } = useBoolean();

function edit(id: number) {
  targetId.value = id;
  openDrawer();
}

function add() {
  openGenerateDrawer();
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <AddressSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />
    <ElCard class="card-wrapper sm:flex-1-hidden" body-class="ht50">
      <template #header>
        <div class="flex items-center justify-between">
          <p>{{ $t('page.address.list.title') }}</p>
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
      <AddressDetailDrawer v-model:visible="drawerVisible" :target-id="targetId" @saved="getDataByPage" />
      <AddressGenerateDrawer v-model:visible="drawerGenerateVisible" @saved="getDataByPage" />
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
