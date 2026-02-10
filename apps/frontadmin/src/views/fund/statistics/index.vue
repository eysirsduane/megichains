<script setup lang="tsx">
import { onBeforeMount } from 'vue';
import { fetchGetAddressFundStatistics } from '@/service/api';
import { useEcharts } from '@/hooks/common/echarts';
import type { ECOption } from '@/hooks/common/echarts';

defineOptions({ name: 'AddressFundStatistics' });

const statistics: ECOption = {
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'cross',
      label: {
        backgroundColor: '#6a7985'
      }
    }
  },
  xAxis: {
    type: 'category',
    data: ['波场USDT', '波场USDC', '币安USDT', '币安USDC', '以太坊USDT', '以太坊USDC', '索拉纳USDT', '索拉纳USDC']
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      data: [],
      type: 'bar',
      color: '#8378ea',
      showBackground: true,
      barGap: 100,
      itemStyle: {
        borderRadius: [40, 40, 0, 0]
      },
      backgroundStyle: {
        color: 'rgba(180, 180, 180, 0.2)'
      }
    }
  ]
};

const { domRef: barRef, setOptions: setStaticsOptions } = useEcharts(() => statistics, {
  onRender() {}
});

onBeforeMount(async () => {
  const { data, error } = await fetchGetAddressFundStatistics();
  if (!error) {
    setStaticsOptions({
      series: [
        {
          data: [
            data.tron_usdt,
            data.tron_usdc,
            data.bsc_usdt,
            data.bsc_usdc,
            data.eth_usdt,
            data.eth_usdc,
            data.solana_usdt,
            data.solana_usdc
          ]
        }
      ]
    });
  }
});
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <ElCard class="card-wrapper">
      <div ref="barRef" class="h-400px" />
    </ElCard>
  </div>
</template>

<style lang="scss" scoped>
:deep(.el-card) {
  .ht50 {
    height: calc(100% - 50px);
  }
}
</style>
