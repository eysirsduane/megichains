<script setup lang="ts">
import { ref, watch } from 'vue';
import { fetchGetAddressGroupAll, postCollectAddressFund } from '@/service/api';
import { $t } from '@/locales';

defineOptions({ name: 'AddressFundCollectDetailDrawer' });

const visible = defineModel<boolean>('visible', {
  default: false
});

const targetId = defineModel<number>('targetId', {
  default: 0
});

type Model = Pick<
  Api.Fund.AddressFundCollect,
  'id' | 'address_group_id' | 'chain' | 'currency' | 'status' | 'amount_min' | 'fee_max' | 'secret_key' | 'description'
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    address_group_id: 0,
    chain: '',
    currency: '',
    status: '',
    amount_min: 0,
    fee_max: 0,
    secret_key: '',
    description: ''
  };
}

function closeDrawer() {
  visible.value = false;
}

const chains = [
  {
    value: 'BSC',
    label: 'BSC'
  },
  {
    value: 'ETH',
    label: 'ETH'
  },
  {
    value: 'TRON',
    label: 'TRON'
  }
];

const currencys = [
  {
    value: 'USDT',
    label: 'USDT'
  },
  {
    value: 'USDC',
    label: 'USDC'
  }
];
interface Emits {
  (e: 'saved'): void;
}
const emit = defineEmits<Emits>();

async function save() {
  const { data, error } = await postCollectAddressFund(model.value);
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
    if (targetId.value > 0) {
      // const detail = await getOrderDetail(targetId.value);
      // model.value = detail;
    } else {
      const { data, error } = await fetchGetAddressGroupAll();
      if (!error) {
        groups.value = data.records;
      }

      model.value = createDefaultModel();
    }
  }
});
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.fund.collectlog.detail')" :size="500">
    <ElForm :model="model" label-position="top">
      <ElFormItem :label="$t('common.id')" prop="id">
        <ElInput v-model="model.id" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.address_group')" prop="group_id">
        <ElSelect v-model="model.address_group_id" :empty-values="[0]" value-on-clear="">
          <ElOption v-for="item in groups" :key="item.id" :label="item.name" :value="item.id"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.chain')" prop="chain">
        <ElSelect v-model="model.chain">
          <ElOption v-for="item in chains" :key="item.value" :label="item.label" :value="item.value"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.currency')" prop="currency">
        <ElSelect v-model="model.currency">
          <ElOption v-for="item in currencys" :key="item.value" :label="item.label" :value="item.value"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.amount_min')" prop="amount_min">
        <ElInput v-model.number="model.amount_min" />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.fee_max')" prop="fee_max">
        <ElInput v-model.number="model.fee_max" />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.command_key')" prop="secret_key">
        <ElInput v-model="model.secret_key" />
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
