<script setup lang="ts">
import { ref, watch } from 'vue';
import { useForm } from '@/hooks/common/form';
import { $t } from '@/locales';

defineOptions({ name: 'MerchantSearch' });

interface Emits {
  (e: 'search'): void;
}
const rtvalue = ref('');

const emit = defineEmits<Emits>();

const { formRef, validate, restoreValidation } = useForm();

const model = defineModel<Api.Merch.MerchantSearchParams>('model', { required: true });
const initialParams = { ...model.value };

async function reset() {
  await restoreValidation();
  Object.assign(model.value, initialParams);
  model.value.start = 0;
  model.value.end = 0;
  rtvalue.value = '';
}

async function search() {
  await validate();
  emit('search');
}

watch(rtvalue, () => {
  if (rtvalue.value) {
    model.value.start = Number.parseInt(rtvalue.value[0], 10);
    model.value.end = Number.parseInt(rtvalue.value[1], 10);
  } else {
    model.value.start = 0;
    model.value.end = 0;
  }
});
</script>

<template>
  <ElCard class="card-wrapper">
    <ElCollapse>
      <ElCollapseItem :title="$t('common.search')" name="user-search">
        <ElForm ref="formRef" :model="model" label-position="right" :label-width="80">
          <ElRow :gutter="24">
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('common.id')" prop="id">
                <ElInput v-model="model.id" :placeholder="$t('common.id')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.merch.common.merchant_account')" prop="merchant_account">
                <ElInput v-model="model.merchant_account" :placeholder="$t('page.merch.common.merchant_account')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12"></ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElSpace class="w-full justify-end" alignment="end">
                <ElButton @click="reset">
                  <template #icon>
                    <icon-ic-round-refresh class="text-icon" />
                  </template>
                  {{ $t('common.reset') }}
                </ElButton>
                <ElButton type="primary" plain @click="search">
                  <template #icon>
                    <icon-ic-round-search class="text-icon" />
                  </template>
                  {{ $t('common.search') }}
                </ElButton>
              </ElSpace>
            </ElCol>
          </ElRow>
        </ElForm>
      </ElCollapseItem>
    </ElCollapse>
  </ElCard>
</template>

<style scoped></style>
