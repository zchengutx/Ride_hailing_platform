<template>
  <div class="home-container">
    <header class="home-header">
      <div class="logo">罗小黑打车</div>
    </header>
    <div id="baidu-map" class="baidu-map"></div>
    <main class="home-main">
      <div class="welcome">欢迎来到罗小黑打车！</div>
      <div class="service-cards">
        <div class="service-card">
          <div class="service-icon">🚗</div>
          <div class="service-title">专车服务</div>
          <div class="service-desc">舒适出行，安全到达</div>
        </div>
        <!-- 可扩展更多服务 -->
      </div>
    </main>
    <footer class="home-footer">
      <div class="footer-item active">
        <span>🏠</span>
        <span>首页</span>
      </div>
      <div class="footer-item">
        <span>📝</span>
        <span>订单</span>
      </div>
      <div class="footer-item">
        <span>👤</span>
        <span>我的</span>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'

function initMap(lng, lat) {
  if (window.BMap) {
    const map = new window.BMap.Map('baidu-map')
    const point = new window.BMap.Point(lng, lat)
    map.centerAndZoom(point, 15)
    map.enableScrollWheelZoom(true)
    // 添加当前位置标记
    const marker = new window.BMap.Marker(point)
    map.addOverlay(marker)
  }
}

onMounted(async () => {
  try {
    const body = {
      origins: "31.0638,121.7975",
      destinations: "40.0799,116.6031"
    }
    const res = await fetch('/api/driving', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(body)
    })
    const data = await res.json()
    console.log('后端返回数据:', data)
    let drivingData = data.data
    if (typeof drivingData === 'string') {
      try {
        // Base64解码
        const jsonStr = atob(drivingData)
        drivingData = JSON.parse(jsonStr)
        console.log('Base64解码并解析后的data:', drivingData)
      } catch (e) {
        console.error('data字段Base64解码或JSON解析失败:', e)
      }
    }
    // 假设 drivingData 里有 lng, lat 字段
    if (drivingData && drivingData.lng && drivingData.lat) {
      initMap(drivingData.lng, drivingData.lat)
    } else {
      initMap(116.404, 39.915)
    }
  } catch (e) {
    console.error('请求后端接口异常:', e)
    initMap(116.404, 39.915)
  }
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}
.home-header {
  background: #fff;
  padding: 18px 0 12px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}
.logo {
  font-size: 1.7rem;
  font-weight: bold;
  color: #007aff;
}
.baidu-map {
  width: 100%;
  height: 220px;
  margin: 0 auto 18px auto;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  background: #eee;
}
.home-main {
  flex: 1;
  padding: 32px 16px 16px;
}
.welcome {
  font-size: 1.2rem;
  margin-bottom: 24px;
  color: #333;
}
.service-cards {
  display: flex;
  gap: 18px;
}
.service-card {
  flex: 1;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  padding: 18px 12px;
  text-align: center;
}
.service-icon {
  font-size: 2.2rem;
  margin-bottom: 8px;
}
.service-title {
  font-weight: bold;
  color: #007aff;
  margin-bottom: 4px;
}
.service-desc {
  color: #666;
  font-size: 0.95rem;
}
.home-footer {
  display: flex;
  background: #fff;
  box-shadow: 0 -2px 8px rgba(0,0,0,0.04);
}
.footer-item {
  flex: 1;
  text-align: center;
  padding: 10px 0 6px;
  color: #888;
  font-size: 1rem;
}
.footer-item.active {
  color: #007aff;
  font-weight: bold;
}
</style> 