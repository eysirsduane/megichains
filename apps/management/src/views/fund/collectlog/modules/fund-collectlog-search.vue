<script setup lang="ts">
import { ref, watch } from 'vue';
import { chainTyposOptions, collectLogStatusOptions, currencyTyposOptions } from '@/constants/business';
import { useForm } from '@/hooks/common/form';
import { translateOptions } from '@/utils/common';
import { $t } from '@/locales';

defineOptions({ name: 'AddressFundCollectLogSearch' });

interface Emits {
  (e: 'search'): void;
}
const rtvalue = ref('');

const emit = defineEmits<Emits>();

const { formRef, validate, restoreValidation } = useForm();

const model = defineModel<Api.Fund.AddressFundCollectLogListSearchParams>('model', { required: true });
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
              <ElFormItem :label="$t('page.fund.common.collect_id')" prop="chain">
                <ElInput v-model="model.collect_id" :placeholder="$t('page.fund.common.collect_id')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('common.chain')" prop="chain">
                <ElSelect
                  v-model="model.chain"
                  clearable
                  :empty-values="['']"
                  value-on-clear=""
                  :placeholder="$t('common.chain')"
                >
                  <ElOption
                    v-for="(item, idx) in translateOptions(chainTyposOptions)"
                    :key="idx"
                    :label="item.label"
                    :value="item.value"
                  ></ElOption>
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('common.currency')" prop="currency">
                <ElSelect
                  v-model="model.currency"
                  clearable
                  :empty-values="['', undefined]"
                  value-on-clear=""
                  :placeholder="$t('common.currency')"
                >
                  <ElOption
                    v-for="(item, idx) in translateOptions(currencyTyposOptions)"
                    :key="idx"
                    :label="item.label"
                    :value="item.value"
                  ></ElOption>
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('common.status')" prop="status">
                <ElSelect
                  v-model="model.status"
                  clearable
                  :empty-values="['']"
                  value-on-clear=""
                  :placeholder="$t('common.status')"
                >
                  <ElOption
                    v-for="(item, idx) in translateOptions(collectLogStatusOptions)"
                    :key="idx"
                    :label="item.label"
                    :value="item.value"
                  ></ElOption>
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('common.from_address')" prop="from_address">
                <ElInput v-model="model.from_address" :placeholder="$t('common.from_address')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('common.receiver_address')" prop="receiver_address">
                <ElInput v-model="model.receiver_address" :placeholder="$t('common.from_address')" />
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
