<template>
  <div class="order-page">
    <h2>我的订单</h2>
    <div class="order-list">
      <!-- 订单列表将在这里显示 -->
      <div class="empty-state">
        暂无订单记录
      </div>
    </div>
  </div>
</template>

<script setup>
// 组件逻辑将在这里添加
</script>

<style scoped>
.order-page {
  padding: 20px;
}

h2 {
  margin-bottom: 20px;
  font-size: 20px;
  font-weight: 600;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
  color: #999;
  font-size: 14px;
}
</style> 