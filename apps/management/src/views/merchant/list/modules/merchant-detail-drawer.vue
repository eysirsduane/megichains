<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { getMerchantDetail, saveMerchantDetail } from '@/service/api';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';

defineOptions({ name: 'ExchangeBillDrawer' });

const visible = defineModel<boolean>('visible', {
  default: false
});

const targetId = defineModel<number>('targetId', {
  default: 0
});

const formRef = ref();

interface Emits {
  (e: 'saved'): void;
}
const emit = defineEmits<Emits>();

type Model = Pick<
  Api.Merch.MerchantDetail,
  'id' | 'merchant_account' | 'name' | 'description' | 'created_at' | 'updated_at'
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    merchant_account: '',
    name: '',
    description: '',
    updated_at: 0,
    created_at: 0
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
    if (targetId.value > 0) {
      const { data, error } = await getMerchantDetail(targetId.value);
      if (!error) {
        model.value = data;
      }
    } else {
      model.value = createDefaultModel();
    }
  }
});

function save() {
  formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      const { data, error } = await saveMerchantDetail(model.value);
      if (!error) {
        window.$message?.success($t('common.saveSuccess'));
        closeDrawer();
        emit('saved');
        return data;
      }
    } else {
      window.$message?.error($t('common.paramsInvalid'));
    }

    return null;
  });
}
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.order.detail.title')" :size="560">
    <ElForm ref="formRef" :model="model" label-position="top">
      <ElFormItem :label="$t('common.id')" prop="id">
        <ElInput v-model="model.id" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('page.order.common.merchant_account')" prop="merchant_account">
        <ElInput v-model="model.merchant_account" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('page.merch.common.name')" prop="name">
        <ElInput v-model="model.name" />
      </ElFormItem>
      <ElFormItem :label="$t('common.description')" prop="description">
        <ElInput v-model="model.description" type="textarea" :rows="7" />
      </ElFormItem>
      <ElFormItem :label="$t('common.updated_at')" prop="created_at">
        <ElInput v-model="timeHuman" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('common.created_at')" prop="created_at">
        <ElInput v-model="timeHuman" disabled />
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
