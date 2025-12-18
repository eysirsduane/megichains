<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { fetchGetDelegateWithdrawal } from '@/service/api';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';

defineOptions({ name: 'DelegateWithdrawerDrawer' });

const visible = defineModel<boolean>('visible', {
  default: false
});

const targetId = defineModel<number>('targetId', {
  default: 0
});

type Model = Pick<
  Api.Transaction.DelegateWithdrawal,
  | 'user_id'
  | 'order_id'
  | 'status'
  | 'transaction_id'
  | 'from_base58'
  | 'to_base58'
  | 'un_delegated_amount'
  | 'un_delegated_sun'
  | 'created_at'
  | 'description'
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    user_id: 0,
    order_id: 0,
    status: '',
    transaction_id: '',
    from_base58: '',
    to_base58: '',
    un_delegated_amount: 0,
    un_delegated_sun: 0,
    created_at: 0,
    description: ''
  };
}

async function getDelegateWithdrawer(oid: number): Promise<Model> {
  const { data, error } = await fetchGetDelegateWithdrawal(oid);
  if (!error) {
    return data;
  }

  return {
    user_id: 0,
    order_id: 0,
    status: '',
    transaction_id: '',
    from_base58: '',
    to_base58: '',
    un_delegated_amount: 0,
    un_delegated_sun: 0,
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
    const withdrawer = await getDelegateWithdrawer(targetId.value);
    model.value = withdrawer;
  }
});
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.transaction.delegate.withdraweral.title')" :size="400">
    <ElForm :model="model" label-position="top">
      <ElFormItem :label="$t('page.transaction.common.order_id')" prop="order_id">
        <ElInput v-model="model.order_id" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.common.transaction_id')" prop="transaction_id">
        <ElInput v-model="model.transaction_id" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.common.status')" prop="status">
        <ElInput v-model="model.status" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.common.from_base58')" prop="from_base58">
        <ElInput v-model="model.from_base58" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.common.to_base58')" prop="to_base58">
        <ElInput v-model="model.to_base58" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.delegate.withdraweral.un_delegated_amount')" prop="un_delegated_amount">
        <ElInput v-model="model.un_delegated_amount" />
      </ElFormItem>
      <ElFormItem :label="$t('page.transaction.delegate.withdraweral.un_delegated_sun')" prop="un_delegated_sun">
        <ElInput v-model="model.un_delegated_sun" />
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
