<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { chainBigTyposOptions, chainTyposOptions } from '@/constants/business';
import { fetchGetAddressGroupAll } from '@/service/api';
import { useForm } from '@/hooks/common/form';
import { translateOptions } from '@/utils/common';
import { $t } from '@/locales';

defineOptions({ name: 'AddressFundCollectLogListSearch' });

interface Emits {
  (e: 'search'): void;
}
const rtvalue = ref('');

const emit = defineEmits<Emits>();

const { formRef, validate, restoreValidation } = useForm();

const model = defineModel<Api.Fund.AddressFundCollectListSearchParams>('model', { required: true });
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
</script>

<template>
  <ElCard class="card-wrapper">
    <ElCollapse>
      <ElCollapseItem :title="$t('common.search')" name="user-search">
        <ElForm ref="formRef" :model="model" label-position="right" :label-width="80">
          <ElRow :gutter="24">
            <ElCol :lg="6" :md="8" :sm="12">
              <ElFormItem :label="$t('page.fund.common.chain')" prop="chain">
                <ElSelect
                  v-model="model.chain"
                  clearable
                  :empty-values="['']"
                  value-on-clear=""
                  :placeholder="$t('page.fund.common.chain')"
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
              <ElFormItem :label="$t('page.fund.common.group')" prop="address_group_id">
                <ElSelect
                  v-model="model.address_group_id"
                  clearable
                  :empty-values="['']"
                  value-on-clear=""
                  :placeholder="$t('page.fund.common.group')"
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
              <ElFormItem :label="$t('page.fund.common.currency')" prop="currency">
                <ElSelect
                  v-model="model.currency"
                  clearable
                  :empty-values="['']"
                  value-on-clear=""
                  :placeholder="$t('page.fund.common.currency')"
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
              <ElFormItem :label="$t('page.fund.common.status')" prop="status">
                <ElSelect
                  v-model="model.status"
                  clearable
                  :empty-values="['']"
                  value-on-clear=""
                  :placeholder="$t('page.fund.common.status')"
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
              <ElFormItem :label="$t('page.fund.common.to_address')" prop="to_address">
                <ElInput v-model="model.to_address" :placeholder="$t('page.fund.common.to_address')" />
              </ElFormItem>
            </ElCol>
            <ElCol :lg="6" :md="8" :sm="12"></ElCol>
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
