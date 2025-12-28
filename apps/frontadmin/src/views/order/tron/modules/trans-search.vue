<script setup lang="ts">
import { ref, watch } from 'vue';
import { orderStatusOptions, orderTypoOptions } from '@/constants/business';
import { useForm } from '@/hooks/common/form';
import { translateNumberOptions, translateOptions } from '@/utils/common';
import { $t } from '@/locales';

defineOptions({ name: 'DelegateSearch' });

interface Emits {
  (e: 'reset'): void;
  (e: 'search'): void;
}
const rtvalue = ref('');

const emit = defineEmits<Emits>();

const { formRef, validate, restoreValidation } = useForm();

const model = defineModel<Api.Order.DelegateOrderSearchParams>('model', { required: true });
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

const shortcuts = [
  {
    text: 'Last week',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
      return [start, end];
    }
  },
  {
    text: 'Last month',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
      return [start, end];
    }
  },
  {
    text: 'Last 3 months',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
      return [start, end];
    }
  }
];
</script>

<template>
  <ElCard class="card-wrapper">
    <ElCollapse>
      <ElCollapseItem :title="$t('common.search')" name="user-search">
        <ElForm ref="formRef" :model="model" label-position="right" :label-width="80">
          <ElRow :gutter="24">
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.transaction.common.order_id')" prop="id">
                <ElInput v-model="model.id" :placeholder="$t('page.transaction.common.order_id')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.transaction.common.transaction_id')" prop="transaction_id">
                <ElInput v-model="model.transaction_id" :placeholder="$t('page.transaction.common.transaction_id')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.transaction.common.typo')" prop="to_base58">
                <ElSelect
                  v-model="model.typo"
                  clearable
                  :empty-values="[-1, undefined]"
                  :value-on-clear="-1"
                  :placeholder="$t('page.transaction.common.typo')"
                >
                  <ElOption
                    v-for="(item, idx) in translateNumberOptions(orderTypoOptions)"
                    :key="idx"
                    :label="item.label"
                    :value="item.value"
                  ></ElOption>
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.transaction.common.status')" prop="status">
                <ElSelect
                  v-model="model.status"
                  clearable
                  :empty-values="['', undefined]"
                  value-on-clear=""
                  :placeholder="$t('page.transaction.common.status')"
                >
                  <ElOption
                    v-for="(item, idx) in translateOptions(orderStatusOptions)"
                    :key="idx"
                    :label="item.label"
                    :value="item.value"
                  ></ElOption>
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.transaction.common.from_base58')" prop="from_base58">
                <ElInput v-model="model.from_base58" :placeholder="$t('page.transaction.common.from_base58')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.transaction.common.to_base58')" prop="to_base58">
                <ElInput v-model="model.to_base58" :placeholder="$t('page.transaction.common.to_base58')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('common.timerange')" prop="start">
                <ElDatePicker
                  v-model="rtvalue"
                  type="daterange"
                  unlink-panels
                  range-separator="To"
                  :start-placeholder="$t('page.transaction.common.start')"
                  :end-placeholder="$t('page.transaction.common.end')"
                  :shortcuts="shortcuts"
                  value-format="x"
                />
              </ElFormItem>
            </ElCol>
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
