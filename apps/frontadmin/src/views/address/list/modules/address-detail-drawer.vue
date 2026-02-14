<script setup lang="ts">
import { ref, watch } from 'vue';
import { fetchGetAddressDetail, fetchGetAddressGroupAll, postSaveAddress } from '@/service/api';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';

defineOptions({ name: 'AddressDetailDrawer' });

const visible = defineModel<boolean>('visible', {
  default: false
});

const targetId = defineModel<number>('targetId', {
  default: 0
});

type Model = Pick<
  Api.Address.Address,
  | 'id'
  | 'group_id'
  | 'typo'
  | 'status'
  | 'chain'
  | 'address'
  | 'address2'
  | 'bsc_usdt'
  | 'bsc_usdc'
  | 'eth_usdt'
  | 'eth_usdc'
  | 'tron_usdt'
  | 'tron_usdc'
  | 'solana_usdt'
  | 'solana_usdc'
  | 'description'
  | 'updated_at'
  | 'created_at'
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    group_id: 0,
    typo: '',
    status: '',
    chain: '',
    address: '',
    address2: '',
    bsc_usdt: 0,
    bsc_usdc: 0,
    eth_usdt: 0,
    eth_usdc: 0,
    tron_usdt: 0,
    tron_usdc: 0,
    solana_usdt: 0,
    solana_usdc: 0,
    description: '',
    updated_at: 0,
    created_at: 0
  };
}

async function getAddressDetail(id: number): Promise<Model> {
  const { data, error } = await fetchGetAddressDetail(id);
  if (!error) {
    return data;
  }

  return createDefaultModel();
}

const timeHuman = (time: number) => {
  return getHumannessDateTime(time);
};

function closeDrawer() {
  visible.value = false;
}

watch(visible, async () => {
  if (visible.value) {
    if (targetId.value > 0) {
      const detail = await getAddressDetail(targetId.value);
      model.value = detail;
    } else {
      model.value = createDefaultModel();
    }
  }
});

const typos = [
  {
    value: 'IN',
    label: '入'
  },
  {
    value: 'OUT',
    label: '出'
  },
  {
    value: 'COLLECT',
    label: '归集'
  }
];
const statuses = [
  {
    value: '空闲',
    label: '空闲'
  },
  {
    value: '占用',
    label: '占用'
  },
  {
    value: '禁用',
    label: '禁用'
  }
];
interface Emits {
  (e: 'saved'): void;
}
const emit = defineEmits<Emits>();

async function save() {
  const { data, error } = await postSaveAddress(model.value);
  if (!error) {
    window.$message?.success($t('common.updateSuccess'));
    closeDrawer();
    emit('saved');
    return data;
  }

  return data;
}

const groups = ref<Api.Address.AddressGroup[]>([]);

watch(visible, async () => {
  if (visible.value) {
    const { data, error } = await fetchGetAddressGroupAll();
    if (!error) {
      groups.value = data.records;
    }
  }
});
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.address.list.title')" :size="560">
    <ElForm :model="model" label-position="top">
      <ElFormItem :label="$t('common.id')" prop="id">
        <ElInput v-model="model.id" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('page.address.common.group_id')" prop="group_id">
        <ElSelect v-model="model.group_id">
          <ElOption v-for="item in groups" :key="item.id" :label="item.name" :value="item.id"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('page.address.common.chain')" prop="chain">
        <ElInput v-model="model.chain" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('page.address.common.typo')" prop="typo">
        <ElSelect v-model="model.typo">
          <ElOption v-for="item in typos" :key="item.value" :label="item.label" :value="item.value"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('page.address.common.status')" prop="status">
        <ElSelect v-model="model.status">
          <ElOption v-for="item in statuses" :key="item.value" :label="item.label" :value="item.value"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('page.address.common.address')" prop="address">
        <ElInput v-model="model.address" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('page.address.common.address')" prop="address2">
        <ElInput v-model="model.address2" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('page.address.common.description')" prop="description">
        <ElInput v-model="model.description" type="textarea" :rows="7" />
      </ElFormItem>
      <ElFormItem :label="$t('common.updated_at')" prop="created_at">
        <ElInput :value="timeHuman(model.updated_at)" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('common.created_at')" prop="created_at">
        <ElInput :value="timeHuman(model.created_at)" disabled />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElSpace :size="16">
        <ElButton type="primary" @click="save">{{ $t('common.confirm') }}</ElButton>
        <ElButton @click="closeDrawer">{{ $t('common.cancel') }}</ElButton>
      </ElSpace>
    </template>
  </ElDrawer>
</template>

<style scoped></style>
