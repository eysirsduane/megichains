<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { fetchGetOrderDetail } from '@/service/api';
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
  Api.Order.OrderDetail,
  | 'order_no'
  | 'merchant_account'
  | 'merchant_order_no'
  | 'transaction_id'
  | 'log_id'
  | 'chain'
  | 'typo'
  | 'status'
  | 'currency'
  | 'from_address'
  | 'to_address'
  | 'received_amount'
  | 'received_sun'
  | 'updated_at'
  | 'created_at'
  | 'description'
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    order_no: '',
    merchant_account: '',
    merchant_order_no: '',
    log_id: 0,
    transaction_id: '',
    chain: '',
    typo: '',
    status: '',
    currency: '',
    from_address: '',
    to_address: '',
    received_amount: 0,
    received_sun: 0,
    updated_at: 0,
    created_at: 0,
    description: ''
  };
}

async function getOrderDetail(oid: number): Promise<Model> {
  const { data, error } = await fetchGetOrderDetail(oid);
  if (!error) {
    return data;
  }

  return createDefaultModel();
}

const timeHuman = computed(() => {
  return getHumannessDateTime(model.value.created_at);
});

function closeDrawer() {
  visible.value = false;
}

watch(visible, async () => {
  if (visible.value) {
    const detail = await getOrderDetail(targetId.value);
    model.value = detail;
  }
});
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.order.detail.title')" :size="560">
    <ElForm :model="model" label-position="top">
      <ElFormItem :label="$t('common.merchant_request')" prop="order_no">
        <ElInput v-model="model.order_no" />
      </ElFormItem>
      <ElFormItem :label="$t('common.merchant_response')" prop="merchant_account">
        <ElInput v-model="model.merchant_account" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.merchant_order_no')" prop="merchant_order_no">
        <ElInput v-model="model.merchant_order_no" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.log_id')" prop="order_id">
        <ElInput v-model="model.log_id" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.transaction_id')" prop="transaction_id">
        <ElInput v-model="model.transaction_id" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.chain')" prop="currency">
        <ElInput v-model="model.chain" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.typo')" prop="currency">
        <ElInput v-model="model.typo" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.status')" prop="currency">
        <ElInput v-model="model.status" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.currency')" prop="currency">
        <ElInput v-model="model.currency" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.from_address')" prop="from_address">
        <ElInput v-model="model.from_address" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.to_address')" prop="to_address">
        <ElInput v-model="model.to_address" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.received_amount')" prop="exchanged_amount">
        <ElInput v-model="model.received_amount" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.received_sun')" prop="exchanged_sun">
        <ElInput v-model="model.received_sun" />
      </ElFormItem>
      <ElFormItem :label="$t('common.updated_at')" prop="created_at">
        <ElInput v-model="timeHuman" />
      </ElFormItem>
      <ElFormItem :label="$t('common.created_at')" prop="created_at">
        <ElInput v-model="timeHuman" />
      </ElFormItem>
      <ElFormItem :label="$t('common.description')" prop="description">
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
