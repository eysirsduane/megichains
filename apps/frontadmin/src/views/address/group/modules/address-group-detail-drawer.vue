<script setup lang="ts">
import { ref, watch } from 'vue';
import { fetchGetAddressGroupDetail, postSaveAddressGroup } from '@/service/api';
import { $t } from '@/locales';
import { getHumannessDateTime } from '@/locales/dayjs';

defineOptions({ name: 'AddressGroupEditDrawer' });

const visible = defineModel<boolean>('visible', {
  default: false
});

const targetId = defineModel<number>('targetId', {
  default: 0
});

type Model = Pick<Api.Address.AddressGroup, 'id' | 'name' | 'description' | 'status' | 'created_at' | 'updated_at'>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    name: '',
    status: '',
    description: '',
    created_at: 0,
    updated_at: 0
  };
}

async function getOrderDetail(oid: number): Promise<Model> {
  const { data, error } = await fetchGetAddressGroupDetail(oid);
  if (!error) {
    return data;
  }

  return createDefaultModel();
}

const timeHuman = (time: number) => {
  return getHumannessDateTime(time);
};

function closeDrawer() {
  visible.value = false;
}

interface Emits {
  (e: 'saved'): void;
}

const emit = defineEmits<Emits>();

async function save() {
  const { data, error } = await postSaveAddressGroup(model.value);
  if (!error) {
    window.$message?.success($t('common.updateSuccess'));
    closeDrawer();
    emit('saved');
    return data;
  }

  return data;
}

watch(visible, async () => {
  if (visible.value) {
    if (targetId.value > 0) {
      const detail = await getOrderDetail(targetId.value);
      model.value = detail;
    } else {
      model.value = createDefaultModel();
    }
  }
});

const options = [
  {
    value: '开放',
    label: '开放'
  },
  {
    value: '禁用',
    label: '禁用'
  }
];
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.order.detail.title')" :size="500">
    <ElForm :model="model" label-position="top">
      <ElFormItem :label="$t('common.id')" prop="id">
        <ElInput v-model="model.id" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('common.name')" prop="name">
        <ElInput v-model="model.name" />
      </ElFormItem>
      <ElFormItem :label="$t('page.address.common.status')" prop="status">
        <ElSelect v-model="model.status">
          <ElOption v-for="item in options" :key="item.value" :label="item.label" :value="item.value"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('common.description')" prop="description">
        <ElInput v-model="model.description" type="textarea" :rows="7" />
      </ElFormItem>
      <ElFormItem :label="$t('common.updated_at')" prop="updated_at">
        <ElInput :value="timeHuman(model.updated_at)" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('common.created_at')" prop="created_at">
        <ElInput :value="timeHuman(model.created_at)" disabled />
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
