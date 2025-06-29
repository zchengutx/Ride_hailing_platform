<template>
  <div class="home">
    <!-- 地图容器 -->
    <div class="map-container">
      <div id="baiduMap" class="baidu-map"></div>

      <!-- 登录按钮 -->
      <div class="login-btn" @click="router.push('/login')" v-if="!isLoggedIn">
        <span class="icon">👤</span>
        <span class="text">去登录</span>
      </div>

      <!-- 右下角定位按钮 -->
      <div class="location-btn" @click="refreshLocation">
        <span class="icon">📍</span>
      </div>

      <!-- 右上角功能按钮 -->
      <div class="function-btns">
        <div class="top-btn" @click="router.push('/wallet')">
          <span class="icon">🎫</span>
          <span class="text">优惠券</span>
        </div>
        <div class="top-btn" @click="router.push('/order')">
          <span class="icon">🚗</span>
          <span class="text">行程</span>
        </div>
      </div>
    </div>

    <!-- 底部搜索框 -->
    <div class="bottom-panel">
      <div class="location-section">
        <div class="location-item">
          <div class="location-dot start"></div>
          <div class="location-text">从 {{ startLocation }}</div>
        </div>
        <div class="location-item destination">
          <div class="location-dot end"></div>
          <input 
            v-model="destination" 
            type="text" 
            @input="handleSearch"
            placeholder="输入你的目的地"
            class="destination-input"
          />
        </div>
      </div>

      <!-- 搜索结果列表 -->
      <div v-show="searchResults.length > 0" class="search-results">
        <div 
          v-for="(result, index) in searchResults" 
          :key="index"
          class="search-result-item"
          @click="selectDestination(result)"
        >
          {{ result.title }}
        </div>
      </div>

      <!-- 打车按钮 -->
      <button 
        class="call-car-btn" 
        @click="handleCallCar"
        :disabled="!canCallCar"
      >
        立即打车
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const map = ref(null)
const startLocation = ref('定位中...')  // 修改默认文本
const destination = ref('')
const searchResults = ref([])

// 起点和终点的经纬度
const startPoint = ref({ lng: null, lat: null })
const endPoint = ref({ lng: null, lat: null })

// 是否可以打车
const canCallCar = computed(() => {
  return startPoint.value.lng && startPoint.value.lat && endPoint.value.lng && endPoint.value.lat
})

// 是否已登录
const isLoggedIn = computed(() => !!localStorage.getItem('token'))

// 初始化百度地图
const initBaiduMap = () => {
  if (!window.BMapGL) {
    console.error('百度地图API未加载')
    return
  }

  // 创建地图实例
  map.value = new window.BMapGL.Map('baiduMap', {
    enableMapClick: false,
    enableDblclickZoom: false
  })
  
  // 设置初始中心点和缩放级别（武汉市中心）
  const point = new window.BMapGL.Point(114.31, 30.52)
  map.value.centerAndZoom(point, 15)
  
  // 启用滚轮缩放
  map.value.enableScrollWheelZoom(true)
  
  // 获取当前位置
  getCurrentLocation()
}

// 获取当前位置
const getCurrentLocation = () => {
  startLocation.value = '定位中...'
  
  const geolocation = new window.BMapGL.Geolocation()
  geolocation.getCurrentPosition(function(r) {
    if (this.getStatus() == window.BMAP_STATUS_SUCCESS) {
      const point = r.point
      startPoint.value = {
        lng: point.lng,
        lat: point.lat
      }
      
      // 输出经纬度信息
      console.log('当前位置经纬度:', {
        longitude: point.lng,
        latitude: point.lat
      })
      
      map.value.centerAndZoom(point, 16)
      map.value.clearOverlays()
      
      const marker = new window.BMapGL.Marker(point)
      map.value.addOverlay(marker)
      
      const geoc = new window.BMapGL.Geocoder()
      geoc.getLocation(point, (rs) => {
        if (rs && rs.addressComponents) {
          const addComp = rs.addressComponents
          startLocation.value = addComp.street + addComp.streetNumber
          // 输出完整地址信息
          console.log('当前位置详细信息:', {
            province: addComp.province,
            city: addComp.city,
            district: addComp.district,
            street: addComp.street,
            streetNumber: addComp.streetNumber,
            longitude: point.lng,
            latitude: point.lat
          })
        } else {
          startLocation.value = '定位失败，请重试'
        }
      })
    } else {
      startLocation.value = '定位失败，请重试'
      console.error('定位失败:', this.getStatus())
    }
  })
}

// 搜索目的地
const handleSearch = () => {
  if (!destination.value) {
    searchResults.value = []
    return
  }
  
  // 创建地址解析器实例
  const myGeo = new window.BMapGL.Geocoder()
  
  // 使用地址解析器进行搜索
  myGeo.getPoint(destination.value, function(point){
    if (point) {
      // 获取到坐标后，使用LocalSearch进行周边搜索
      const local = new window.BMapGL.LocalSearch(map.value, {
        renderOptions: {
          map: map.value,
          autoViewport: false
        },
        pageCapacity: 8,
        onSearchComplete: function(results) {
          // 清空之前的搜索结果
          searchResults.value = []
          
          try {
            if (results && Array.isArray(results)) {
              // 处理多个结果集的情况
              results.forEach(result => {
                if (result && typeof result.getCurrentNumPois === 'function') {
                  const numPois = result.getCurrentNumPois()
                  for (let i = 0; i < numPois; i++) {
                    const poi = result.getPoi(i)
                    searchResults.value.push({
                      title: poi.title,
                      address: poi.address,
                      point: poi.point
                    })
                  }
                }
              })
            } else if (results && typeof results.getCurrentNumPois === 'function') {
              // 处理单个结果集的情况
              const numPois = results.getCurrentNumPois()
              for (let i = 0; i < numPois; i++) {
                const poi = results.getPoi(i)
                searchResults.value.push({
                  title: poi.title,
                  address: poi.address,
                  point: poi.point
                })
              }
            }

            // 如果没有找到任何结果，至少添加搜索的地点
            if (searchResults.value.length === 0) {
              searchResults.value.push({
                title: destination.value,
                address: '',
                point: point
              })
            }

            console.log('搜索结果:', searchResults.value)
          } catch (error) {
            console.error('处理搜索结果时出错:', error)
            // 至少添加当前搜索的地点
            searchResults.value.push({
              title: destination.value,
              address: '',
              point: point
            })
          }
        }
      })
      
      // 在坐标点附近搜索
      local.search(destination.value)
    } else {
      console.log('未找到该地址的具体坐标')
      searchResults.value = []
    }
  }, '上海市') // 设置为上海市，因为截图显示是上海的地址
}

// 选择目的地
const selectDestination = (result) => {
  destination.value = result.title
  searchResults.value = []
  
  const point = result.point
  endPoint.value = {
    lng: point.lng,
    lat: point.lat
  }
  
  // 清除之前的标记
  map.value.clearOverlays()
  
  // 添加起点和终点标记
  if (startPoint.value.lng && startPoint.value.lat) {
    const startMarker = new window.BMapGL.Marker(new window.BMapGL.Point(startPoint.value.lng, startPoint.value.lat))
    map.value.addOverlay(startMarker)
  }
  
  const endMarker = new window.BMapGL.Marker(new window.BMapGL.Point(point.lng, point.lat))
  map.value.addOverlay(endMarker)
  
  // 调整视野以显示两个点
  const viewPoints = []
  if (startPoint.value.lng && startPoint.value.lat) {
    viewPoints.push(new window.BMapGL.Point(startPoint.value.lng, startPoint.value.lat))
  }
  viewPoints.push(new window.BMapGL.Point(point.lng, point.lat))
  
  if (viewPoints.length > 1) {
    map.value.setViewport(viewPoints)
  } else {
    map.value.centerAndZoom(viewPoints[0], 16)
  }
}

// 处理打车请求
const handleCallCar = async () => {
  if (!canCallCar.value) {
    alert('请先选择目的地')
    return
  }
  
  try {
    // TODO: 调用后端打车接口
    const orderData = {
      startLocation: {
        lng: startPoint.value.lng,
        lat: startPoint.value.lat,
        address: startLocation.value
      },
      endLocation: {
        lng: endPoint.value.lng,
        lat: endPoint.value.lat,
        address: destination.value
      }
    }
    console.log('打车请求数据:', orderData)
    // 这里添加调用后端API的代码
    
    // 跳转到订单页面
    router.push('/order')
  } catch (error) {
    console.error('打车失败:', error)
    alert('打车失败，请重试')
  }
}

// 刷新定位
const refreshLocation = () => {
  getCurrentLocation()
}

onMounted(() => {
  const checkBaiduMap = () => {
    if (window.BMapGL) {
      initBaiduMap()
    } else {
      setTimeout(checkBaiduMap, 100)
    }
  }
  checkBaiduMap()
})
</script>

<style scoped>
.home {
  height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
}

.map-container {
  flex: 1;
  position: relative;
}

.baidu-map {
  width: 100%;
  height: 100%;
}

.bottom-panel {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  padding: 20px;
  border-radius: 20px 20px 0 0;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.location-section {
  position: relative;
  margin-bottom: 15px;
}

.location-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  gap: 12px;
  background: white;
}

.location-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.location-dot.start {
  background: #00C896;
}

.location-dot.end {
  background: #FF6B35;
}

.location-text {
  font-size: 14px;
  color: #333;
  flex: 1;
}

.destination-input {
  border: none;
  outline: none;
  font-size: 14px;
  width: 100%;
  color: #333;
  background: #f5f5f5;
  padding: 8px 12px;
  border-radius: 4px;
}

.destination-input::placeholder {
  color: #999;
}

.search-results {
  position: absolute;
  bottom: 100%;
  left: 0;
  right: 0;
  background: white;
  border-radius: 12px;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
  max-height: 200px;
  overflow-y: auto;
  z-index: 1000;
}

.search-result-item {
  padding: 12px;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
}

.search-result-item:hover {
  background: #f9f9f9;
}

.call-car-btn {
  width: 100%;
  padding: 14px;
  margin-top: 16px;
  background: #00C896;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
}

.call-car-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

/* 顶部按钮样式 */
.top-buttons {
  position: absolute;
  top: 20px;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-between;
  padding: 0 20px;
  z-index: 1000;
}

.button-group {
  display: flex;
  gap: 10px;
}

.top-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.9);
  padding: 8px 12px;
  border-radius: 20px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.top-btn:hover {
  background: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.top-btn .icon {
  font-size: 18px;
}

.top-btn .text {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

/* 登录按钮样式 */
.login-btn {
  position: absolute;
  top: 20px;
  left: 20px;
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.9);
  padding: 8px 16px;
  border-radius: 20px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  z-index: 1000;
}

.login-btn:hover {
  background: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 定位按钮样式 */
.location-btn {
  position: absolute;
  bottom: 120px;  /* 位于底部搜索框上方 */
  right: 20px;
  width: 40px;
  height: 40px;
  background: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  z-index: 1000;
}

.location-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 右上角功能按钮 */
.function-btns {
  position: absolute;
  top: 20px;
  right: 20px;
  display: flex;
  gap: 10px;
  z-index: 1000;
}

.top-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.9);
  padding: 8px 12px;
  border-radius: 20px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.top-btn:hover {
  background: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.icon {
  font-size: 18px;
}

.text {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

/* 适配移动端 */
@media (max-width: 768px) {
  .top-buttons {
    padding: 0 15px;
  }

  .top-btn {
    padding: 6px 10px;
  }

  .top-btn .text {
    font-size: 12px;
  }

  .login-btn,
  .top-btn {
    padding: 6px 10px;
  }

  .text {
    font-size: 12px;
  }

  .location-btn {
    bottom: 100px;
    width: 36px;
    height: 36px;
  }

  .bottom-panel {
    padding: 15px;
  }

  .location-item {
    padding: 10px 0;
  }

  .destination-input {
    font-size: 14px;
  }

  .call-car-btn {
    padding: 12px;
    font-size: 15px;
  }
}
</style> 