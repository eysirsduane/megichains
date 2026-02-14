<script setup lang="ts">
import { ref, watch } from 'vue';
import { chainBigTyposOptions } from '@/constants/business';
import { fetchGetAddressGroupAll, postGenerateAddress } from '@/service/api';
import { translateOptions } from '@/utils/common';
import { $t } from '@/locales';

defineOptions({ name: 'AddressGenerateDrawer' });

const visible = defineModel<boolean>('visible', {
  default: false
});

type Model = Pick<Api.Address.AddressGenerate, 'chain' | 'group_id' | 'count'>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    chain: '',
    group_id: 0,
    count: 0
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
  const { data, error } = await postGenerateAddress(model.value);
  if (!error) {
    window.$message?.success($t('common.updateSuccess'));
    closeDrawer();
    emit('saved');
    return data;
  }

  return data;
}

const groups = ref<Api.Address.AddressGroup[]>([]);

watch(visible, async () => {
  if (visible.value) {
    const { data, error } = await fetchGetAddressGroupAll();
    if (!error) {
      groups.value = data.records;
    }
  }
});
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.address.list.title')" :size="560">
    <ElForm :model="model" label-position="top">
      <ElFormItem :label="$t('page.address.common.chain')" prop="chain">
        <ElSelect
          v-model="model.chain"
          clearable
          :empty-values="['']"
          value-on-clear=""
          :placeholder="$t('common.chain')"
        >
          <ElOption
            v-for="(item, idx) in translateOptions(chainBigTyposOptions)"
            :key="idx"
            :label="item.label"
            :value="item.value"
          ></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('common.group')" prop="group_id">
        <ElSelect
          v-model="model.group_id"
          clearable
          :empty-values="[0]"
          value-on-clear=""
          :placeholder="$t('page.fund.common.group')"
        >
          <ElOption v-for="item in groups" :key="item.id" :label="item.name" :value="item.id"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('common.count')" prop="count">
        <ElInputNumber v-model="model.count" type="number" />
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
