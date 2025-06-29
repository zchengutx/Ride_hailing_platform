<template>
  <div id="app">
    <router-view />
    <!-- 底部导航栏 -->
    <div class="bottom-nav">
      <router-link to="/" class="nav-item" active-class="active">
        <span class="icon">🏠</span>
        <span class="text">首页</span>
      </router-link>
      <router-link to="/order" class="nav-item" active-class="active">
        <span class="icon">📋</span>
        <span class="text">订单</span>
      </router-link>
      <router-link to="/user" class="nav-item" active-class="active">
        <span class="icon">👤</span>
        <span class="text">我的</span>
      </router-link>
    </div>
  </div>
</template>

<script setup>
// No script needed
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  min-height: 100vh;
  padding-bottom: 60px; /* 为底部导航留出空间 */
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: white;
  display: flex;
  justify-content: space-around;
  align-items: center;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  z-index: 1000;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-decoration: none;
  color: #666;
  padding: 8px;
  flex: 1;
}

.nav-item.active {
  color: #00C896;
}

.nav-item .icon {
  font-size: 24px;
  margin-bottom: 2px;
}

.nav-item .text {
  font-size: 12px;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  background-color: #f5f5f5;
}
</style>
