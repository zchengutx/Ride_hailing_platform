<template>
  <div class="map-container">
    <div 
      ref="mapContainer" 
      class="map-wrapper"
      :style="{ width: width, height: height }"
    ></div>
    <div class="map-controls">
      <el-button 
        type="primary" 
        size="small" 
        @click="getCurrentLocation"
        :loading="locating"
        class="location-btn"
      >
        <el-icon><Location /></el-icon>
        {{ locating ? '定位中...' : '重新定位' }}
      </el-button>
      <div class="address-info" v-if="currentAddress">
        <el-icon><LocationInformation /></el-icon>
        <span>{{ currentAddress }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Location, LocationInformation } from '@element-plus/icons-vue'

// 定义props
interface Props {
  width?: string
  height?: string
  zoom?: number
  center?: [number, number] // [lng, lat]
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  height: '200px',
  zoom: 15,
  center: () => [116.404, 39.915] // 默认北京坐标
})

// 定义emit事件
const emit = defineEmits<{
  locationUpdate: [location: { lng: number; lat: number; address: string }]
  mapClick: [event: any]
  mapError: [error: string]
}>()

// 响应式数据
const mapContainer = ref<HTMLDivElement>()
const locating = ref(false)
const currentAddress = ref('')
let map: any = null
let geolocation: any = null
let currentMarker: any = null

// 初始化地图
const initMap = () => {
  if (!mapContainer.value || !window.BMap) {
    console.error('百度地图API未加载或容器不存在')
    emit('mapError', '百度地图API未加载')
    return
  }

  try {
    // 创建地图实例
    map = new window.BMap.Map(mapContainer.value)
    
    // 监听地图加载错误事件
    map.addEventListener('error', (e: any) => {
      console.error('百度地图加载错误:', e)
      emit('mapError', '百度地图服务错误')
      ElMessage.error('地图服务异常，将使用备用定位方案')
    })
    
    // 设置中心点和缩放级别
    const point = new window.BMap.Point(props.center[0], props.center[1])
    map.centerAndZoom(point, props.zoom)
    
    // 监听地图加载完成事件
    map.addEventListener('tilesloaded', () => {
      console.log('百度地图瓦片加载完成')
    })
    
    // 启用地图功能
    map.enableScrollWheelZoom(true) // 启用滚轮缩放
    map.enableDoubleClickZoom(true) // 启用双击缩放
    map.enableKeyboard(true) // 启用键盘操作
    
    // 添加控件
    map.addControl(new window.BMap.NavigationControl()) // 平移缩放控件
    map.addControl(new window.BMap.ScaleControl()) // 比例尺控件
    
    // 地图点击事件
    map.addEventListener('click', (e: any) => {
      emit('mapClick', e)
    })
    
    // 初始化定位服务
    geolocation = new window.BMap.Geolocation()
    
    // 设置一个延时来检查地图是否正常工作
    setTimeout(() => {
      if (map && mapContainer.value) {
        // 检查地图容器是否有内容
        const mapContent = mapContainer.value.querySelector('.BMap_mask')
        if (!mapContent) {
          console.warn('地图可能未正常加载')
          // 不立即触发错误，先尝试定位
        }
      }
      // 自动获取当前位置
      getCurrentLocation()
    }, 1000)
    
    console.log('百度地图初始化成功')
  } catch (error) {
    console.error('地图初始化失败:', error)
    const errorMsg = (error as Error).message
    if (errorMsg.includes('APP') || errorMsg.includes('ak') || errorMsg.includes('key')) {
      emit('mapError', 'API密钥无效')
      ElMessage.error('地图API密钥无效，将使用备用定位方案')
    } else {
      emit('mapError', '地图初始化失败: ' + errorMsg)
      ElMessage.error('地图加载失败，将使用备用定位方案')
    }
  }
}

// 获取当前位置
const getCurrentLocation = () => {
  if (!geolocation || !map) {
    ElMessage.error('定位服务未初始化')
    return
  }
  
  locating.value = true
  
  geolocation.getCurrentPosition((r: any) => {
    if (geolocation.getStatus() === window.BMAP_STATUS_SUCCESS) {
      const point = new window.BMap.Point(r.point.lng, r.point.lat)
      
      // 更新地图中心
      map.panTo(point)
      
      // 清除之前的标记
      if (currentMarker) {
        map.removeOverlay(currentMarker)
      }
      
      // 添加当前位置标记
      currentMarker = new window.BMap.Marker(point)
      map.addOverlay(currentMarker)
      
      // 创建信息窗口
      const infoWindow = new window.BMap.InfoWindow('您的当前位置', {
        width: 200,
        height: 50
      })
      currentMarker.addEventListener('click', () => {
        map.openInfoWindow(infoWindow, point)
      })
      
      // 逆地理编码获取地址
      const geocoder = new window.BMap.Geocoder()
      geocoder.getLocation(point, (result: any) => {
        if (result) {
          currentAddress.value = result.address
          emit('locationUpdate', {
            lng: r.point.lng,
            lat: r.point.lat,
            address: result.address
          })
          ElMessage.success('定位成功')
        } else {
          ElMessage.warning('获取地址信息失败')
        }
      })
      
      locating.value = false
    } else {
      locating.value = false
      const errorMsg = getLocationErrorMessage(geolocation.getStatus())
      ElMessage.error('定位失败：' + errorMsg)
    }
  }, {
    enableHighAccuracy: true,
    timeout: 15000,
    maximumAge: 60000
  })
}

// 获取定位错误信息
const getLocationErrorMessage = (status: number) => {
  const errorMessages: { [key: number]: string } = {
    1: '定位权限被拒绝',
    2: '网络异常或定位失败',
    3: '定位超时',
    4: '定位结果异常',
    5: '定位不支持',
    6: '定位服务异常'
  }
  return errorMessages[status] || '未知错误'
}

// 添加标记点
const addMarker = (lng: number, lat: number, title?: string, icon?: string) => {
  if (!map) return
  
  const point = new window.BMap.Point(lng, lat)
  const marker = new window.BMap.Marker(point)
  
  if (icon) {
    const customIcon = new window.BMap.Icon(icon, new window.BMap.Size(25, 25))
    marker.setIcon(customIcon)
  }
  
  if (title) {
    const label = new window.BMap.Label(title, { offset: new window.BMap.Size(20, -10) })
    marker.setLabel(label)
  }
  
  map.addOverlay(marker)
  return marker
}

// 设置地图中心
const setCenter = (lng: number, lat: number) => {
  if (!map) return
  const point = new window.BMap.Point(lng, lat)
  map.panTo(point)
}

// 暴露方法给父组件
defineExpose({
  addMarker,
  setCenter,
  getCurrentLocation,
  map: () => map
})

onMounted(() => {
  // 监听百度地图API加载错误
  window.addEventListener('error', (e) => {
    if (e.target && (e.target as any).src && (e.target as any).src.includes('api.map.baidu.com')) {
      console.error('百度地图API加载失败:', e)
      emit('mapError', '百度地图API加载失败')
      ElMessage.error('地图服务不可用，将使用备用定位方案')
    }
  })

  // 检查百度地图API是否加载
  if (window.BMap) {
    initMap()
  } else {
    // 如果API还没加载，等待加载完成
    let retryCount = 0
    const maxRetries = 30 // 最多重试3秒
    
    const checkBMap = () => {
      if (window.BMap) {
        // 检查API是否有效（尝试创建一个简单的对象）
        try {
          new window.BMap.Map(document.createElement('div'))
          initMap()
        } catch (error) {
          console.error('百度地图API无效:', error)
          emit('mapError', '百度地图API密钥无效')
          ElMessage.error('地图服务配置有误，将使用备用定位方案')
        }
      } else if (retryCount < maxRetries) {
        retryCount++
        setTimeout(checkBMap, 100)
      } else {
        console.error('百度地图API加载超时')
        emit('mapError', '百度地图API加载超时')
        ElMessage.error('地图加载失败，将使用备用定位方案')
      }
    }
    checkBMap()
  }
})

onUnmounted(() => {
  if (map) {
    map.clearOverlays()
    map = null
  }
})
</script>

<style scoped>
.map-container {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.map-wrapper {
  position: relative;
}

.map-controls {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.location-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  background: rgba(255, 255, 255, 0.95);
  color: #ff7e00;
  border: 1px solid #ff7e00;
  backdrop-filter: blur(4px);
  font-size: 12px;
  padding: 8px 12px;
}

.location-btn:hover {
  background: #ff7e00;
  color: white;
}

.address-info {
  background: rgba(255, 255, 255, 0.95);
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 12px;
  color: #333;
  max-width: 200px;
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  gap: 4px;
  border: 1px solid #e5e5e5;
}

.address-info span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 百度地图控件样式调整 */
:deep(.BMap_cpyCtrl) {
  display: none !important;
}

:deep(.anchorBL) {
  display: none !important;
}
</style>

<script lang="ts">
// 全局声明百度地图类型
declare global {
  interface Window {
    BMap: any;
    BMAP_STATUS_SUCCESS: any;
  }
}
</script> 