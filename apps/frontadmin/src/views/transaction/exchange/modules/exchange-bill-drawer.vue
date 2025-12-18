<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { fetchGetExchangeBill } from '@/service/api';
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
  Api.Transaction.ExchangeBill,
  | 'user_id'
  | 'order_id'
  | 'transaction_id'
  | 'currency'
  | 'from_base58'
  | 'to_base58'
  | 'exchanged_amount'
  | 'exchanged_sun'
  | 'status'
  | 'created_at'
  | 'description'
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    user_id: 0,
    order_id: 0,
    transaction_id: '',
    currency: '',
    from_base58: '',
    to_base58: '',
    exchanged_amount: 0,
    exchanged_sun: 0,
    status: '',
    created_at: 0,
    description: ''
  };
}

async function getExchangedBill(oid: number): Promise<Model> {
  const { data, error } = await fetchGetExchangeBill(oid);
  if (!error) {
    return data;
  }

  return {
    user_id: 0,
    order_id: 0,
    transaction_id: '',
    currency: '',
    from_base58: '',
    to_base58: '',
    exchanged_amount: 0,
    exchanged_sun: 0,
    status: '',
    created_at: 0,
    description: ''
  };
}

const timeHuman = computed(() => {
  return getHumannessDateTime(model.value.created_at);
});

function closeDrawer() {
  visible.value = false;
}

watch(visible, async () => {
  if (visible.value) {
    const bill = await getExchangedBill(targetId.value);
    model.value = bill;
  }
});
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.transaction.delegate.bill.title')" :size="400">
    <ElForm ref="form" :model="model" label-position="top">
      <ElFormItem :label="$t('page.transaction.common.order_id')" prop="order_id">
        <ElInput v-model="model.order_id" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.common.transaction_id')" prop="transaction_id">
        <ElInput v-model="model.transaction_id" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.common.currency')" prop="currency">
        <ElInput v-model="model.currency" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.common.from_base58')" prop="from_base58">
        <ElInput v-model="model.from_base58" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.common.to_base58')" prop="to_base58">
        <ElInput v-model="model.to_base58" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.exchange.bill.exchanged_amount')" prop="exchanged_amount">
        <ElInput v-model="model.exchanged_amount" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.exchange.bill.exchanged_sun')" prop="exchanged_sun">
        <ElInput v-model="model.exchanged_sun" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.common.created_at')" prop="created_at">
        <ElInput v-model="timeHuman" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.common.description')" prop="description">
        <ElInput v-model="model.description" type="textarea" :rows="7" />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElSpace :size="16">
        <ElButton @click="closeDrawer">{{ $t('common.cancel') }}</ElButton>
      </ElSpace>
    </template>
  </ElDrawer>
</template>

<style scoped></style>
