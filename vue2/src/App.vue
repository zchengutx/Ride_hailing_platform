<script setup>
import { ref } from 'vue';

// 响应式数据
const currentService = ref('express');
const startLocation = ref('当前位置');
const endLocation = ref('输入目的地');
const services = [
  { id: 'express', name: '快车', icon: '🚗' },
  { id: 'premium', name: '专车', icon: '🚘' },
  { id: 'taxi', name: '出租车', icon: '🚕' },
  { id: 'designated', name: '代驾', icon: '👨‍✈️' },
  { id: 'carpool', name: '拼车', icon: '👥' },
  { id: 'luxury', name: '豪华车', icon: '💎' }
];

const quickServices = [
  { id: 'reserve', name: '预约', icon: '⏰' },
  { id: 'charter', name: '包车', icon: '🚐' },
  { id: 'enterprise', name: '企业用车', icon: '🏢' },
  { id: 'intercity', name: '城际快线', icon: '🚄' }
];
</script>

<template>
  <router-view />
</template>

<script>
export default {
  name: 'App'
}
</script>

<style scoped>
#app {
  width: 100%;
  height: 100%;
}
</style>
