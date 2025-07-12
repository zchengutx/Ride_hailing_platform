<template>
  <div class="fallback-location">
    <div class="fallback-content">
      <div class="location-icon">
        <el-icon size="48" color="#ff7e00"><Location /></el-icon>
      </div>
      <h3>智能定位服务</h3>
      <p>正在为您提供精准的位置定位服务</p>
      
      <el-button 
        type="primary" 
        @click="getCurrentLocation"
        :loading="locating"
        class="location-btn"
      >
        <el-icon><Location /></el-icon>
        {{ locating ? '定位中...' : '重新定位' }}
      </el-button>
      
      <div class="manual-input" v-if="showManualInput">
        <el-divider>或者</el-divider>
        <el-input
          v-model="manualAddress"
          placeholder="请输入您的位置"
          @keyup.enter="useManualAddress"
        >
          <template #append>
            <el-button @click="useManualAddress">确认</el-button>
          </template>
        </el-input>
      </div>
      
      <el-button 
        text 
        type="primary" 
        @click="showManualInput = !showManualInput"
        class="manual-toggle"
      >
        {{ showManualInput ? '隐藏手动输入' : '手动输入位置' }}
      </el-button>
      
      <div class="current-location" v-if="currentAddress">
        <el-icon><LocationInformation /></el-icon>
        <span>当前位置：{{ currentAddress }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Location, LocationInformation } from '@element-plus/icons-vue'

// 定义emit事件
const emit = defineEmits<{
  locationUpdate: [location: { lng: number; lat: number; address: string }]
}>()

// 响应式数据
const locating = ref(false)
const currentAddress = ref('')
const showManualInput = ref(false)
const manualAddress = ref('')

// 使用浏览器原生定位API
const getCurrentLocation = () => {
  if (!navigator.geolocation) {
    ElMessage.error('您的浏览器不支持定位功能')
    showManualInput.value = true
    return
  }
  
  locating.value = true
  
  navigator.geolocation.getCurrentPosition(
    (position) => {
      const { latitude, longitude } = position.coords
      
      // 使用逆地理编码服务（这里使用一个简化的示例）
      getReverseGeocoding(latitude, longitude)
        .then((address) => {
          currentAddress.value = address
          emit('locationUpdate', {
            lng: longitude,
            lat: latitude,
            address: address
          })
          ElMessage.success('定位成功')
        })
        .catch(() => {
          // 如果逆地理编码失败，使用坐标信息
          const address = `纬度: ${latitude.toFixed(6)}, 经度: ${longitude.toFixed(6)}`
          currentAddress.value = address
          emit('locationUpdate', {
            lng: longitude,
            lat: latitude,
            address: address
          })
          ElMessage.success('定位成功')
        })
        .finally(() => {
          locating.value = false
        })
    },
    (error) => {
      locating.value = false
      let errorMessage = '定位失败'
      
      switch (error.code) {
        case error.PERMISSION_DENIED:
          errorMessage = '定位权限被拒绝，请在浏览器设置中允许定位'
          break
        case error.POSITION_UNAVAILABLE:
          errorMessage = '位置信息不可用'
          break
        case error.TIMEOUT:
          errorMessage = '定位请求超时'
          break
        default:
          errorMessage = '定位发生未知错误'
          break
      }
      
      ElMessage.error(errorMessage)
      showManualInput.value = true
    },
    {
      enableHighAccuracy: true,
      timeout: 10000,
      maximumAge: 60000
    }
  )
}

// 使用免费的逆地理编码服务
const getReverseGeocoding = async (lat: number, lng: number): Promise<string> => {
  try {
    // 使用免费的OpenStreetMap Nominatim服务进行逆地理编码
    const response = await fetch(
      `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}&accept-language=zh-CN`
    )
    
    if (response.ok) {
      const data = await response.json()
      if (data.display_name) {
        // 简化地址信息，只保留主要部分
        const address = data.display_name
        const parts = address.split(',')
        if (parts.length > 3) {
          return parts.slice(0, 3).join(',')
        }
        return address
      }
    }
  } catch (error) {
    console.warn('逆地理编码失败:', error)
  }
  
  // 如果逆地理编码失败，返回格式化的坐标地址
  return `位置坐标: ${lat.toFixed(4)}, ${lng.toFixed(4)}`
}

// 使用手动输入的地址
const useManualAddress = () => {
  if (!manualAddress.value.trim()) {
    ElMessage.warning('请输入有效的地址')
    return
  }
  
  currentAddress.value = manualAddress.value
  emit('locationUpdate', {
    lng: 116.404, // 默认经度（北京）
    lat: 39.915,  // 默认纬度（北京）
    address: manualAddress.value
  })
  
  ElMessage.success('地址设置成功')
  manualAddress.value = ''
  showManualInput.value = false
}

// 组件挂载后自动开始定位
onMounted(() => {
  // 延迟一点时间，让用户看到界面
  setTimeout(() => {
    getCurrentLocation()
  }, 500)
})
</script>

<style scoped>
.fallback-location {
  height: 200px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  padding: 20px;
  position: relative;
  overflow: hidden;
}

.fallback-location::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 20%, rgba(255, 255, 255, 0.1) 2px, transparent 2px),
    radial-gradient(circle at 80% 80%, rgba(255, 255, 255, 0.1) 2px, transparent 2px),
    radial-gradient(circle at 40% 60%, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
  background-size: 30px 30px, 40px 40px, 20px 20px;
  animation: float 20s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
}

.fallback-content {
  text-align: center;
  max-width: 300px;
  position: relative;
  z-index: 1;
  color: white;
}

.location-icon {
  margin-bottom: 16px;
}

.fallback-content h3 {
  margin: 0 0 8px 0;
  color: white;
  font-size: 18px;
  font-weight: 600;
}

.fallback-content p {
  margin: 0 0 20px 0;
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
  line-height: 1.5;
}

.location-btn {
  margin-bottom: 16px;
  background: rgba(255, 126, 0, 0.9);
  border-color: #ff7e00;
  backdrop-filter: blur(4px);
  color: white;
  font-weight: 500;
}

.location-btn:hover {
  background: rgba(255, 102, 0, 1);
  border-color: #ff6600;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(255, 126, 0, 0.3);
}

.manual-input {
  margin: 16px 0;
}

.manual-toggle {
  font-size: 12px;
  margin-bottom: 16px;
  color: rgba(255, 255, 255, 0.8);
}

.manual-toggle:hover {
  color: white;
}

.current-location {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.95);
  padding: 12px;
  border-radius: 8px;
  font-size: 14px;
  color: #333;
  margin-top: 16px;
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.current-location span {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style> 