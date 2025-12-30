<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { addressStatusOptions, addressTyposOptions, chainBigTyposOptions } from '@/constants/business';
import { fetchGetAddressGroupAll } from '@/service/api';
import { useForm } from '@/hooks/common/form';
import { translateOptions } from '@/utils/common';
import { $t } from '@/locales';

defineOptions({ name: 'AddressSearch' });

interface Emits {
  (e: 'search'): void;
}
const rtvalue = ref('');

const emit = defineEmits<Emits>();

const { formRef, validate, restoreValidation } = useForm();

const model = defineModel<Api.Address.AddressSearchParams>('model', { required: true });
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

const addrGroupOptions = ref<Api.Address.AddressGroup[] | undefined>();

watch(rtvalue, () => {
  if (rtvalue.value) {
    model.value.start = Number.parseInt(rtvalue.value[0], 10);
    model.value.end = Number.parseInt(rtvalue.value[1], 10);
  } else {
    model.value.start = 0;
    model.value.end = 0;
  }
});

onMounted(async () => {
  const all = await fetchGetAddressGroupAll();
  addrGroupOptions.value = all.data?.records;
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
              <ElFormItem :label="$t('page.address.common.group_id')" prop="group_id">
                <ElSelect
                  v-model="model.group_id"
                  clearable
                  :empty-values="[0]"
                  :value-on-clear="0"
                  :placeholder="$t('page.address.common.group_id')"
                >
                  <ElOption
                    v-for="(item, idx) in addrGroupOptions"
                    :key="idx"
                    :label="item.name"
                    :value="item.id"
                  ></ElOption>
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.address.common.chain')" prop="chain">
                <ElSelect
                  v-model="model.chain"
                  clearable
                  :empty-values="['']"
                  value-on-clear=""
                  :placeholder="$t('page.address.common.chain')"
                >
                  <ElOption
                    v-for="(item, idx) in translateOptions(chainBigTyposOptions)"
                    :key="idx"
                    :label="item.label"
                    :value="item.value"
                  ></ElOption>
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.address.common.typo')" prop="typo">
                <ElSelect
                  v-model="model.typo"
                  clearable
                  :empty-values="['', undefined]"
                  value-on-clear=""
                  :placeholder="$t('page.address.common.typo')"
                >
                  <ElOption
                    v-for="(item, idx) in translateOptions(addressTyposOptions)"
                    :key="idx"
                    :label="item.label"
                    :value="item.value"
                  ></ElOption>
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.order.common.status')" prop="status">
                <ElSelect
                  v-model="model.status"
                  clearable
                  :empty-values="['']"
                  value-on-clear=""
                  :placeholder="$t('page.order.common.status')"
                >
                  <ElOption
                    v-for="(item, idx) in translateOptions(addressStatusOptions)"
                    :key="idx"
                    :label="item.label"
                    :value="item.value"
                  ></ElOption>
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.address.common.address')" prop="from_address">
                <ElInput v-model="model.address" :placeholder="$t('page.address.common.address')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.address.common.address2')" prop="address2">
                <ElInput v-model="model.address2" :placeholder="$t('page.address.common.address2')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('common.timerange')" prop="start">
                <ElDatePicker
                  v-model="rtvalue"
                  type="daterange"
                  unlink-panels
                  range-separator="To"
                  :start-placeholder="$t('common.start')"
                  :end-placeholder="$t('common.end')"
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
