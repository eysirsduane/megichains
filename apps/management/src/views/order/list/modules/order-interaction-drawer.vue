<script setup lang="ts">
import { ref, watch } from 'vue';
import { getOrderInteraction } from '@/service/api';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';

defineOptions({ name: 'OrderTestPlaceDrawer' });

const visible = defineModel<boolean>('visible', {
  default: false
});

const targetId = defineModel<number>('targetId', {
  default: 0
});

type Model = Pick<
  Api.Order.OrderInteraction,
  | 'id'
  | 'merchant_order_id'
  | 'place_request'
  | 'place_request_timestamp'
  | 'place_response'
  | 'place_response_timestamp'
  | 'notify_request'
  | 'notify_request_timestamp'
  | 'notify_response'
  | 'notify_response_timestamp'
  | 'description'
  | 'created_at'
  | 'updated_at'
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    merchant_order_id: 0,
    place_request: '',
    place_request_timestamp: 0,
    place_response: '',
    place_response_timestamp: 0,
    notify_request: '',
    notify_request_timestamp: 0,
    notify_response: '',
    notify_response_timestamp: 0,
    description: '',
    updated_at: 0,
    created_at: 0
  };
}

function closeDrawer() {
  visible.value = false;
}

watch(visible, async () => {
  if (visible.value) {
    if (targetId.value > 0) {
      const { data, error } = await getOrderInteraction(targetId.value);
      if (!error) {
        model.value = data;
      }
    }
  }
});
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.order.common.interaction')" :size="560">
    <ElForm :model="model" label-position="top">
      <ElFormItem :label="$t('page.order.common.merchant_order_id')" prop="merchant_order_id">
        <ElInput v-model="model.merchant_order_id" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.place_request')" prop="place_request">
        <ElInput v-model="model.place_request" type="textarea" :rows="7" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.place_request_timestamp')" prop="place_request_timestamp">
        <ElInput :value="getHumannessDateTime(model.place_request_timestamp)" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.place_response')" prop="place_response">
        <ElInput v-model="model.place_response" type="textarea" :rows="7" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.place_response_timestamp')" prop="place_response_timestamp">
        <ElInput :value="getHumannessDateTime(model.place_response_timestamp)" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.notify_request')" prop="notify_request">
        <ElInput v-model="model.notify_request" type="textarea" :rows="7" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.notify_request_timestamp')" prop="notify_request_timestamp">
        <ElInput :value="getHumannessDateTime(model.notify_request_timestamp)" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.notify_response')" prop="notify_response">
        <ElInput v-model="model.notify_response" type="textarea" :rows="7" />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.notify_response_timestamp')" prop="notify_response_timestamp">
        <ElInput :value="getHumannessDateTime(model.notify_response_timestamp)" />
      </ElFormItem>
      <ElFormItem :label="$t('common.description')" prop="description">
        <ElInput v-model="model.description" type="textarea" :rows="7" />
      </ElFormItem>
      <ElFormItem :label="$t('common.updated_at')" prop="updated_at">
        <ElInput :value="getHumannessDateTime(model.updated_at)" />
      </ElFormItem>
      <ElFormItem :label="$t('common.created_at')" prop="created_at">
        <ElInput :value="getHumannessDateTime(model.created_at)" />
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
