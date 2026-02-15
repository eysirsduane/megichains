<script setup lang="ts">
import { ref, watch } from 'vue';
import { fetchGetAddressFundCollectLogDetail } from '@/service/api';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';

defineOptions({ name: 'ExchangeBillDrawer' });

const visible = defineModel<boolean>('visible', {
  default: false
});

const targetId = defineModel<number>('targetId', {
  default: 0
});

type Model = Pick<
  Api.Fund.AddressFundCollectLog,
  | 'id'
  | 'collect_id'
  | 'chain'
  | 'currency'
  | 'status'
  | 'from_address'
  | 'receiver_address'
  | 'amount'
  | 'transaction_id'
  | 'gas_used'
  | 'gas_price'
  | 'effective_gas_price'
  | 'total_gas_fee'
  | 'description'
  | 'updated_at'
  | 'created_at'
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    collect_id: 0,
    chain: '',
    currency: '',
    status: '',
    from_address: '',
    receiver_address: '',
    amount: 0,
    transaction_id: '',
    gas_used: 0,
    gas_price: 0,
    effective_gas_price: 0,
    total_gas_fee: 0,
    description: '',
    updated_at: 0,
    created_at: 0
  };
}

const timeHuman = (time: number) => {
  return getHumannessDateTime(time);
};

function closeDrawer() {
  visible.value = false;
}

// watch(visible, async () => {
//   if (visible.value) {
//     const { data, error } = await fetchGetAddressFundCollectLogDetail(targetId.value);
//     if (!error) {
//       model.value = data;
//     }
//   }
// });

watch(targetId, async () => {
  if (targetId.value > 0) {
    const { data, error } = await fetchGetAddressFundCollectLogDetail(targetId.value);
    if (!error) {
      model.value = data;
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
      <ElFormItem :label="$t('page.fund.common.collect_id')" prop="collect_id">
        <ElInput v-model="model.collect_id" />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.chain')" prop="chain">
        <ElInput v-model="model.chain" />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.currency')" prop="currency">
        <ElInput v-model="model.currency" />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.status')" prop="status">
        <ElInput v-model="model.status" />
      </ElFormItem>
      <ElFormItem :label="$t('common.from_address')" prop="from_address">
        <ElInput v-model="model.from_address" />
      </ElFormItem>
      <ElFormItem :label="$t('common.receiver_address')" prop="receiver_address">
        <ElInput v-model="model.receiver_address" />
      </ElFormItem>
      <ElFormItem :label="$t('common.amount')" prop="amount">
        <ElInput v-model="model.amount" />
      </ElFormItem>
      <ElFormItem :label="$t('common.transaction_id')" prop="transaction_id">
        <ElInput v-model="model.transaction_id" />
      </ElFormItem>
      <ElFormItem :label="$t('common.gas_used')" prop="gas_used">
        <ElInput v-model="model.gas_used" />
      </ElFormItem>
      <ElFormItem :label="$t('common.gas_price')" prop="gas_price">
        <ElInput v-model="model.gas_price" />
      </ElFormItem>
      <ElFormItem :label="$t('common.effective_gas_price')" prop="effective_gas_price">
        <ElInput v-model="model.effective_gas_price" />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.total_gas_fee')" prop="total_gas_fee">
        <ElInput v-model="model.total_gas_fee" />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.description')" prop="description">
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
        <!-- <ElButton type="primary" @click="save">{{ $t('common.confirm') }}</ElButton> -->
        <ElButton @click="closeDrawer">{{ $t('common.cancel') }}</ElButton>
      </ElSpace>
    </template>
  </ElDrawer>
</template>

<style scoped></style>
