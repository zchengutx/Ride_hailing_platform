<template>
  <div class="register-container">
    <div class="logo">罗小黑打车</div>
    <form class="register-form" @submit.prevent="onRegister">
      <input v-model="form.username" type="text" placeholder="手机号/用户名" required />
      <input v-model="form.password" type="password" placeholder="密码" required />
      <input v-model="form.confirm" type="password" placeholder="确认密码" required />
      <button type="submit">注册</button>
    </form>
    <div class="to-login">
      已有账号？<router-link to="/login">登录</router-link>
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()
const form = reactive({ username: '', password: '', confirm: '' })
function onRegister() {
  if (!form.username || !form.password || !form.confirm) return
  if (form.password !== form.confirm) {
    alert('两次密码不一致')
    return
  }
  // 这里只做前端演示，实际应调用后端API
  router.push('/login')
}
</script>

<style scoped>
.register-container {
  max-width: 350px;
  margin: 80px auto;
  padding: 32px 24px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.08);
  text-align: center;
}
.logo {
  font-size: 2rem;
  font-weight: bold;
  color: #007aff;
  margin-bottom: 32px;
}
.register-form input {
  width: 100%;
  margin-bottom: 18px;
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 6px;
  font-size: 1rem;
}
.register-form button {
  width: 100%;
  padding: 10px;
  background: #007aff;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
}
.to-login {
  margin-top: 18px;
  font-size: 0.95rem;
}
.to-login a {
  color: #007aff;
  text-decoration: none;
}
</style> 