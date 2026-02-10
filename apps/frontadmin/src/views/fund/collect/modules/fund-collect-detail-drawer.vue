<script setup lang="ts">
import { ref, watch } from 'vue';
import { fetchGetAddressGroupAll, postCollectAddressFund } from '@/service/api';
import { $t } from '@/locales';

defineOptions({ name: 'AddressFundCollectDetailDrawer' });

const visible = defineModel<boolean>('visible', {
  default: false
});

const targetId = defineModel<number>('targetId', {
  default: 0
});

const formRef = ref();

type Model = Pick<
  Api.Fund.AddressFundCollectCreating,
  'id' | 'address_group_id' | 'chain' | 'currency' | 'status' | 'amount_min' | 'fee_max' | 'secret_key' | 'description'
>;

const model = ref(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    address_group_id: 0,
    chain: '',
    currency: '',
    status: '',
    amount_min: 0.0,
    fee_max: 0.0,
    secret_key: '',
    description: ''
  };
}

function closeDrawer() {
  visible.value = false;
}

const chains = [
  {
    value: 'BSC',
    label: 'BSC'
  },
  {
    value: 'ETH',
    label: 'ETH'
  },
  {
    value: 'TRON',
    label: 'TRON'
  },
  {
    value: 'SOLANA',
    label: 'SOLANA'
  }
];

const currencys = [
  {
    value: 'USDT',
    label: 'USDT'
  },
  {
    value: 'USDC',
    label: 'USDC'
  }
];
interface Emits {
  (e: 'saved'): void;
}
const emit = defineEmits<Emits>();

function save() {
  formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      const { data, error } = await postCollectAddressFund(model.value);
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

const groups = ref<Api.Address.AddressGroup[]>([]);

watch(visible, async () => {
  if (visible.value) {
    if (targetId.value > 0) {
      // const detail = await getOrderDetail(targetId.value);
      // model.value = detail;
    } else {
      const { data, error } = await fetchGetAddressGroupAll();
      if (!error) {
        groups.value = data.records;
      }

      model.value = createDefaultModel();
    }
  }
});

const rules = {
  address_group_id: [{ required: true, message: '请选择分组', trigger: 'blur' }],
  chain: [{ required: true, message: '请选择链', trigger: 'blur' }],
  currency: [{ required: true, message: '请选择币种', trigger: 'blur' }],
  amount_min: [{ required: true, message: '请输入最小归集金额', trigger: 'blur' }],
  fee_max: [{ required: true, message: '请输入最大转账费用', trigger: 'blur' }],
  secret_key: [{ required: true, message: '请输入口令密钥', trigger: 'blur' }]
};
</script>

<template>
  <ElDrawer v-model="visible" :title="$t('page.fund.collect.detail')" :size="560">
    <ElForm ref="formRef" :rules="rules" :model="model" label-position="top">
      <ElFormItem :label="$t('common.id')" prop="id">
        <ElInput v-model="model.id" disabled />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.address_group')" prop="address_group_id">
        <ElSelect v-model="model.address_group_id" :empty-values="[0]" value-on-clear="">
          <ElOption v-for="item in groups" :key="item.id" :label="item.name" :value="item.id"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.chain')" prop="chain">
        <ElSelect v-model="model.chain">
          <ElOption v-for="item in chains" :key="item.value" :label="item.label" :value="item.value"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.currency')" prop="currency">
        <ElSelect v-model="model.currency">
          <ElOption v-for="item in currencys" :key="item.value" :label="item.label" :value="item.value"></ElOption>
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.collect_amount_min')" prop="amount_min">
        <ElInputNumber v-model.number="model.amount_min" />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.fee_max')" prop="fee_max">
        <ElInputNumber v-model.number="model.fee_max" />
      </ElFormItem>
      <ElFormItem :label="$t('page.fund.common.command_key')" prop="secret_key">
        <ElInput v-model="model.secret_key" />
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
