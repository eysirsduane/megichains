<script setup lang="ts">
import { ref, watch } from 'vue';
import { chainTyposOptions, currencyTyposOptions } from '@/constants/business';
import { findMerchantList, postOrderTestPlace } from '@/service/api';
import { translateOptions } from '@/utils/common';
import { $t } from '@/locales';

defineOptions({ name: 'OrderTestPlaceDrawer' });

const visible = defineModel<boolean>('visible', {
  default: false
});

type Model = Pick<
  Api.Order.OrderTestPlace,
  'merchant_account' | 'chain' | 'currency' | 'mode' | 'receiver' | 'seconds' | 'notify_url'
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    merchant_account: '',
    chain: '',
    currency: '',
    mode: '测试',
    receiver: '',
    seconds: 0,
    notify_url: ''
  };
}

function closeDrawer() {
  visible.value = false;
}

interface Emits {
  (e: 'saved'): void;
}

const emit = defineEmits<Emits>();

async function save() {
  const { data, error } = await postOrderTestPlace(model.value);
  if (!error) {
    window.$message?.success($t('common.saveSuccess'));
    closeDrawer();
    emit('saved');
    return data;
  }

  return data;
}

const merchants = ref<Api.Merch.MerchantDetail[]>([]);

watch(visible, async () => {
  if (visible.value) {
    const { data, error } = await findMerchantList({ current: 1, size: 20 });
    if (!error) {
      merchants.value = data.records;
    }
  }
});
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.order.detail.title')" :size="560">
    <ElForm :model="model" label-position="top">
      <ElFormItem :label="$t('page.merch.common.merchant_account')" prop="merchant_account">
        <ElSelect v-model="model.merchant_account">
          <ElOption
            v-for="item in merchants"
            :key="item.merchant_account"
            :label="item.name"
            :value="item.merchant_account"
          ></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('common.chain')" prop="id">
        <ElSelect v-model="model.chain">
          <ElOption
            v-for="(item, idx) in translateOptions(chainTyposOptions)"
            :key="idx"
            :label="item.label"
            :value="item.value"
          ></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('common.currency')" prop="name">
        <ElSelect v-model="model.currency">
          <ElOption
            v-for="(item, idx) in translateOptions(currencyTyposOptions)"
            :key="idx"
            :label="item.label"
            :value="item.value"
          ></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('common.receiver_address')" prop="receiver">
        <ElInput v-model="model.receiver" />
      </ElFormItem>
      <ElFormItem :label="$t('common.seconds')" prop="seconds">
        <ElInputNumber v-model="model.seconds" />
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
