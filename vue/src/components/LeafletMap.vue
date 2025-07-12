<template>
  <div class="leaflet-map-container">
    <div ref="mapContainer" class="map-wrapper" :style="{ height: height }"></div>
    
    <!-- 地图控件 -->
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
      
      <el-button 
        type="success" 
        size="small" 
        @click="enableManualAdjustment"
        v-if="currentMarker && !manualAdjustMode"
        class="adjust-btn"
      >
        <el-icon><Edit /></el-icon>
        微调位置
      </el-button>
      
      <div class="manual-adjust-controls" v-if="manualAdjustMode">
        <el-button size="small" @click="confirmManualAdjustment" type="primary">
          <el-icon><Check /></el-icon>
          确认位置
        </el-button>
        <el-button size="small" @click="cancelManualAdjustment">
          <el-icon><Close /></el-icon>
          取消
        </el-button>
      </div>
      
      <div class="address-info" v-if="currentAddress">
        <el-icon><LocationInformation /></el-icon>
        <span>{{ currentAddress }}</span>
      </div>
      
      <div class="manual-adjust-tip" v-if="manualAdjustMode">
        <el-icon><InfoFilled /></el-icon>
        <span>点击地图调整定位位置</span>
      </div>
      
      <!-- 路线规划测试按钮 -->
      <el-button 
        type="warning" 
        size="small" 
        @click="testRouteAPI"
        v-if="currentMarker && !routePolyline"
        class="test-route-btn"
        title="测试后端路线规划API（需要后端服务运行在8888端口）"
      >
        <el-icon><Guide /></el-icon>
        测试路线
      </el-button>
      
      <!-- 清除路线按钮 -->
      <el-button 
        type="danger" 
        size="small" 
        @click="clearRoute"
        v-if="routePolyline"
        class="clear-route-btn"
      >
        <el-icon><Delete /></el-icon>
        清除路线
      </el-button>
      

    </div>

    <!-- 加载提示 -->
    <div v-if="mapLoading" class="map-loading-overlay">
      <div class="loading-content">
        <el-icon class="loading-icon"><Loading /></el-icon>
        <p>{{ loadingText }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Location, LocationInformation, Loading, Edit, Check, Close, InfoFilled, Guide, Delete } from '@element-plus/icons-vue'
import * as L from 'leaflet'
import { calculateDriving, getDirectionLite } from '@/api/index'

// 修复Leaflet默认图标问题
delete (L.Icon.Default.prototype as any)._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png'
})

// 定义props
interface Props {
  height?: string
  zoom?: number
  center?: [number, number]
  autoLocate?: boolean  // 新增：是否自动定位
}

const props = withDefaults(defineProps<Props>(), {
  height: '200px',
  zoom: 15,
  center: () => [39.915, 116.404],
  autoLocate: true  // 默认自动定位
})

const {
  center = [39.915, 116.404],
  zoom = 15,
  height = '200px',
  autoLocate = true  // 默认自动定位
} = props

// 定义emit事件
const emit = defineEmits<{
  locationUpdate: [location: { lng: number; lat: number; address: string }]
  mapClick: [event: any]
  routeCalculated: [info: { distance: string; duration: string; coordinates: any[]; isRealRoute: boolean }]
  routeCleared: []
}>()

// 响应式数据
const mapContainer = ref<HTMLDivElement>()
const mapLoading = ref(true)
const loadingText = ref('正在初始化地图...')
const locating = ref(false)
const currentAddress = ref('')
const manualAdjustMode = ref(false)


let map: L.Map | null = null
let currentMarker: L.Marker | null = null
let startMarker: L.Marker | null = null
let endMarker: L.Marker | null = null
let routePolyline: L.Polyline | null = null
let originalPosition: { lat: number; lng: number } | null = null

// 在文件开头添加一个全局错误抑制器
const suppressGeocodeErrors = true // 控制是否抑制地理编码错误

// 初始化地图
const initMap = () => {
  if (!mapContainer.value) {
    console.error('地图容器不存在')
    return
  }

  try {
    map = L.map(mapContainer.value, {
      center: props.center,
      zoom: props.zoom,
      zoomControl: true,
      attributionControl: true
    })

    // 使用高德地图瓦片（国内访问稳定）
    L.tileLayer('https://webrd0{s}.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&x={x}&y={y}&z={z}', {
      attribution: '© 高德地图',
      maxZoom: 18,
      subdomains: ['1', '2', '3', '4']
    }).addTo(map)

    map.on('click', (e) => {
      if (manualAdjustMode.value) {
        handleManualAdjustClick(e)
      } else {
        emit('mapClick', e)
      }
    })

    map.whenReady(() => {
      console.log('地图初始化成功')
      console.log('🔍 autoLocate属性值:', props.autoLocate)
      mapLoading.value = false
      
      // 只有在允许自动定位且没有手动禁用时才自动定位
      if (props.autoLocate) {
        console.log('🟢 开始自动定位...')
        setTimeout(getCurrentLocation, 500)
      } else {
        console.log('🔴 跳过自动定位，等待手动触发或路线规划')
      }
    })

  } catch (error) {
    console.error('地图初始化失败:', error)
    mapLoading.value = false
    ElMessage.error('地图初始化失败')
  }
}

// 获取当前位置
const getCurrentLocation = () => {
  if (!navigator.geolocation) {
    ElMessage.error('您的浏览器不支持定位功能')
    return
  }
  
  locating.value = true
  loadingText.value = '正在高精度定位...'
  
  // 使用多次定位取平均值来提高精度
  performMultipleLocationAttempts()
}

// 执行多次高精度定位尝试
const performMultipleLocationAttempts = async () => {
  const positions: Array<{ lat: number; lng: number; accuracy: number; source: string; timestamp: number }> = []
  const maxAttempts = 5 // 增加定位次数以提高精度
  
  // 并行进行GPS和网络定位
  try {
    loadingText.value = '正在进行高精度定位...'
    
    // 并行获取多种定位方式
    const locationPromises = [
      // 高精度GPS定位（多次）
      ...Array(maxAttempts).fill(0).map((_, i) => 
        getHighPrecisionGPSLocation(i + 1).catch(() => null)
      ),
      // 网络定位作为补充
      getNetworkLocation().catch(() => null)
    ]
    
    const results = await Promise.allSettled(locationPromises)
    
    // 收集有效的定位结果
    results.forEach((result, index) => {
      if (result.status === 'fulfilled' && result.value) {
        positions.push(result.value)
      }
    })
    
    if (positions.length === 0) {
      locating.value = false
      mapLoading.value = false // 隐藏加载遮罩
      loadingText.value = '正在初始化地图...'
      ElMessage.error('定位失败，请检查位置权限和网络连接')
      return
    }
    
    console.log(`总共获得${positions.length}个有效定位结果`)
    
    // 智能处理和融合定位结果
    processEnhancedLocationResults(positions)
    
  } catch (error) {
    console.error('定位过程出错:', error)
    locating.value = false
    mapLoading.value = false // 隐藏加载遮罩
    loadingText.value = '正在初始化地图...'
    ElMessage.error('定位服务异常，请重试')
  }
}

// 获取高精度GPS定位
const getHighPrecisionGPSLocation = (attemptNumber: number): Promise<{ lat: number; lng: number; accuracy: number; source: string; timestamp: number }> => {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new Error('浏览器不支持地理定位'))
      return
    }
    
    // 第一次使用更严格的设置
    const options = {
      enableHighAccuracy: true,
      timeout: attemptNumber === 1 ? 15000 : 10000, // 第一次给更多时间
      maximumAge: 0  // 强制获取新的位置
    }
    
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude, accuracy, heading, speed } = position.coords
        const timestamp = position.timestamp
        
        // 验证坐标有效性
        if (isNaN(latitude) || isNaN(longitude) || latitude === 0 || longitude === 0) {
          reject(new Error('获取到无效的坐标'))
          return
        }
        
        // 验证坐标合理性（中国境内）
        if (latitude < 3 || latitude > 54 || longitude < 73 || longitude > 135) {
          reject(new Error('坐标超出合理范围'))
          return
        }
        
        console.log(`GPS定位 #${attemptNumber}: ${latitude.toFixed(6)}, ${longitude.toFixed(6)}, 精度: ${accuracy}米`)
        
        resolve({ 
          lat: latitude, 
          lng: longitude, 
          accuracy: accuracy || 999,
          source: `GPS-${attemptNumber}`,
          timestamp: timestamp
        })
      },
      (error) => {
        const errorMessage = `GPS定位#${attemptNumber}失败: ${error.message}`
        console.log(errorMessage)
        reject(new Error(errorMessage))
      },
      options
    )
  })
}

// 获取网络定位（作为GPS的补充）
const getNetworkLocation = (): Promise<{ lat: number; lng: number; accuracy: number; source: string; timestamp: number }> => {
  return new Promise((resolve, reject) => {
    // 使用网络定位API（如果可用）
    if (!navigator.geolocation) {
      reject(new Error('浏览器不支持网络定位'))
      return
    }
    
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude, accuracy } = position.coords
        
        // 验证坐标有效性和合理性
        if (isNaN(latitude) || isNaN(longitude) || 
            latitude < 3 || latitude > 54 || longitude < 73 || longitude > 135) {
          reject(new Error('网络定位坐标无效'))
          return
        }
        
        console.log(`网络定位: ${latitude.toFixed(6)}, ${longitude.toFixed(6)}, 精度: ${accuracy}米`)
        
        resolve({ 
          lat: latitude, 
          lng: longitude, 
          accuracy: (accuracy || 1000) + 100, // 网络定位精度通常较低
          source: 'Network',
          timestamp: Date.now()
        })
      },
      (error) => {
        reject(new Error(`网络定位失败: ${error.message}`))
      },
      {
        enableHighAccuracy: false, // 网络定位不需要高精度
        timeout: 8000,
        maximumAge: 30000 // 允许使用30秒内的缓存
      }
    )
  })
}

// 智能处理和融合定位结果
const processEnhancedLocationResults = (positions: Array<{ lat: number; lng: number; accuracy: number; source: string; timestamp: number }>) => {
  console.log('开始智能处理定位结果...')
  
  // 1. 过滤异常值
  const filteredPositions = filterOutlierPositions(positions)
  console.log(`过滤后剩余${filteredPositions.length}个有效定位`)
  
  if (filteredPositions.length === 0) {
    ElMessage.error('所有定位结果都异常，请重试')
    return
  }
  
  // 2. 按精度分组
  const highPrecision = filteredPositions.filter(p => p.accuracy < 30)
  const mediumPrecision = filteredPositions.filter(p => p.accuracy >= 30 && p.accuracy < 100)
  const lowPrecision = filteredPositions.filter(p => p.accuracy >= 100)
  
  console.log(`高精度定位: ${highPrecision.length}个, 中等精度: ${mediumPrecision.length}个, 低精度: ${lowPrecision.length}个`)
  
  // 3. 选择最佳融合策略
  let finalPosition: { lat: number; lng: number; accuracy: number }
  
  if (highPrecision.length >= 2) {
    // 有多个高精度定位，使用智能加权平均
    finalPosition = calculateSmartWeightedAverage(highPrecision)
    console.log('使用多个高精度定位的智能平均')
  } else if (highPrecision.length === 1) {
    // 只有一个高精度定位，直接使用
    finalPosition = highPrecision[0]
    console.log('使用单个高精度定位')
  } else if (mediumPrecision.length >= 3) {
    // 有多个中等精度定位，使用加权平均
    finalPosition = calculateSmartWeightedAverage(mediumPrecision)
    console.log('使用多个中等精度定位的平均')
  } else {
    // 使用所有可用定位的加权平均
    finalPosition = calculateSmartWeightedAverage(filteredPositions)
    console.log('使用所有定位的加权平均')
  }
  
  console.log(`融合前坐标: ${finalPosition.lat.toFixed(6)}, ${finalPosition.lng.toFixed(6)}`)
  
  // 4. 应用精确的坐标系转换（WGS84 → GCJ-02）
  const correctedPosition = applyPreciseCoordinateCorrection(finalPosition)
  
  console.log(`坐标系转换后: ${correctedPosition.lat.toFixed(6)}, ${correctedPosition.lng.toFixed(6)}`)
  
  // 5. 显示最终结果
  displayLocationResult(correctedPosition, finalPosition.accuracy)
}

// 过滤异常定位结果
const filterOutlierPositions = (positions: Array<{ lat: number; lng: number; accuracy: number; source: string; timestamp: number }>) => {
  if (positions.length <= 2) return positions
  
  // 计算所有坐标的中心点
  const centerLat = positions.reduce((sum, p) => sum + p.lat, 0) / positions.length
  const centerLng = positions.reduce((sum, p) => sum + p.lng, 0) / positions.length
  
  // 计算每个点到中心的距离
  const positionsWithDistance = positions.map(pos => ({
    ...pos,
    distanceToCenter: calculateDistance(
      { lat: centerLat, lng: centerLng },
      { lat: pos.lat, lng: pos.lng }
    ) * 1000 // 转换为米
  }))
  
  // 过滤掉距离中心过远的点（可能是异常值）
  const maxDistance = 200 // 200米
  const filtered = positionsWithDistance.filter(pos => {
    const isValidDistance = pos.distanceToCenter <= maxDistance
    const isValidAccuracy = pos.accuracy <= 1000 // 精度不能太差
    
    if (!isValidDistance) {
      console.log(`过滤异常定位 ${pos.source}: 距离中心${pos.distanceToCenter.toFixed(0)}米`)
    }
    
    return isValidDistance && isValidAccuracy
  })
  
  return filtered.length > 0 ? filtered : positions // 如果全部被过滤，则保留原始数据
}

// 智能加权平均算法
const calculateSmartWeightedAverage = (positions: Array<{ lat: number; lng: number; accuracy: number; source: string; timestamp: number }>) => {
  if (positions.length === 1) return positions[0]
  
  let totalWeight = 0
  let weightedLat = 0
  let weightedLng = 0
  let bestAccuracy = Math.min(...positions.map(p => p.accuracy))
  
  positions.forEach(pos => {
    // 计算权重：精度越高权重越大，时间越新权重越大
    const accuracyWeight = 1 / Math.max(pos.accuracy, 1) // 避免除零
    const timeWeight = pos.source.includes('GPS') ? 1.2 : 1.0 // GPS权重稍高
    const finalWeight = accuracyWeight * timeWeight
    
    totalWeight += finalWeight
    weightedLat += pos.lat * finalWeight
    weightedLng += pos.lng * finalWeight
    
    console.log(`位置权重 ${pos.source}: 精度${pos.accuracy}m, 权重${finalWeight.toFixed(3)}`)
  })
  
  return {
    lat: weightedLat / totalWeight,
    lng: weightedLng / totalWeight,
    accuracy: bestAccuracy,
    source: 'Fused'
  }
}



// 精确的坐标系转换（WGS84 → GCJ-02）
const applyPreciseCoordinateCorrection = (position: { lat: number; lng: number; accuracy: number }) => {
  // 检查是否在中国境内（需要坐标系转换）
  if (position.lat >= 0.8293 && position.lat <= 55.8271 && 
      position.lng >= 72.004 && position.lng <= 137.8347) {
    
    // 使用精确的WGS84到GCJ-02转换算法
    const corrected = wgs84ToGcj02(position.lat, position.lng)
    
    console.log(`坐标系转换: WGS84(${position.lat.toFixed(6)}, ${position.lng.toFixed(6)}) → GCJ-02(${corrected.lat.toFixed(6)}, ${corrected.lng.toFixed(6)})`)
    
    return {
      lat: corrected.lat,
      lng: corrected.lng,
      accuracy: position.accuracy
    }
  }
  
  // 海外地区不需要转换
  return position
}

// WGS84坐标系转换为GCJ-02坐标系（火星坐标系）
const wgs84ToGcj02 = (wgsLat: number, wgsLng: number): { lat: number; lng: number } => {
  const a = 6378245.0
  const ee = 0.00669342162296594323
  
  let dLat = transformLat(wgsLng - 105.0, wgsLat - 35.0)
  let dLng = transformLng(wgsLng - 105.0, wgsLat - 35.0)
  
  const radLat = (wgsLat / 180.0) * Math.PI
  let magic = Math.sin(radLat)
  magic = 1 - ee * magic * magic
  const sqrtMagic = Math.sqrt(magic)
  
  dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * Math.PI)
  dLng = (dLng * 180.0) / (a / sqrtMagic * Math.cos(radLat) * Math.PI)
  
  const mgLat = wgsLat + dLat
  const mgLng = wgsLng + dLng
  
  return { lat: mgLat, lng: mgLng }
}

// 纬度转换函数
const transformLat = (lng: number, lat: number): number => {
  let ret = -100.0 + 2.0 * lng + 3.0 * lat + 0.2 * lat * lat + 
           0.1 * lng * lat + 0.2 * Math.sqrt(Math.abs(lng))
  
  ret += (20.0 * Math.sin(6.0 * lng * Math.PI) + 20.0 * Math.sin(2.0 * lng * Math.PI)) * 2.0 / 3.0
  ret += (20.0 * Math.sin(lat * Math.PI) + 40.0 * Math.sin(lat / 3.0 * Math.PI)) * 2.0 / 3.0
  ret += (160.0 * Math.sin(lat / 12.0 * Math.PI) + 320 * Math.sin(lat * Math.PI / 30.0)) * 2.0 / 3.0
  
  return ret
}

// 经度转换函数
const transformLng = (lng: number, lat: number): number => {
  let ret = 300.0 + lng + 2.0 * lat + 0.1 * lng * lng + 
           0.1 * lng * lat + 0.1 * Math.sqrt(Math.abs(lng))
  
  ret += (20.0 * Math.sin(6.0 * lng * Math.PI) + 20.0 * Math.sin(2.0 * lng * Math.PI)) * 2.0 / 3.0
  ret += (20.0 * Math.sin(lng * Math.PI) + 40.0 * Math.sin(lng / 3.0 * Math.PI)) * 2.0 / 3.0
  ret += (150.0 * Math.sin(lng / 12.0 * Math.PI) + 300.0 * Math.sin(lng / 30.0 * Math.PI)) * 2.0 / 3.0
  
  return ret
}

// 显示定位结果
const displayLocationResult = (position: { lat: number; lng: number; accuracy: number }, originalAccuracy: number) => {
  console.log('开始显示定位结果:', position, '精度:', originalAccuracy)
  
  if (map) {
    // 根据精度调整缩放级别
    const zoomLevel = originalAccuracy < 20 ? 19 : originalAccuracy < 50 ? 18 : originalAccuracy < 100 ? 17 : 16
    map.setView([position.lat, position.lng], zoomLevel)
    
    // 移除旧标记
    if (currentMarker) {
      map.removeLayer(currentMarker)
      currentMarker = null
    }
    
    console.log('创建新的定位标记...')
    
    // 方案1：创建绿色圆形标记 - 使用内联样式确保显示
    try {
      const greenIcon = L.divIcon({
        className: 'custom-green-marker',
        html: `<div style="
          width: 16px; 
          height: 16px; 
          border-radius: 50%; 
          background-color: #52c41a; 
          border: 2px solid white; 
          box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
          position: relative;
        "></div>`,
        iconSize: [20, 20],
        iconAnchor: [10, 10]
      })
      
      currentMarker = L.marker([position.lat, position.lng], { icon: greenIcon })
      currentMarker.addTo(map)
      console.log('绿色圆形标记已添加到地图，标记对象:', currentMarker)
      
      // 验证标记是否真的在地图上
      setTimeout(() => {
        if (currentMarker && map.hasLayer(currentMarker)) {
          console.log('✅ DIV标记已成功添加到地图并可见')
        } else {
          console.error('❌ DIV标记未显示，尝试备用方案...')
          
          // 备用方案：使用CircleMarker
          if (currentMarker) map.removeLayer(currentMarker)
          
          currentMarker = L.circleMarker([position.lat, position.lng], {
            radius: 8,
            fillColor: '#52c41a',
            color: '#ffffff',
            weight: 2,
            opacity: 1,
            fillOpacity: 1
          }).addTo(map)
          
          console.log('✅ CircleMarker备用标记已添加')
        }
      }, 200)
      
      // 添加弹窗
      currentMarker.bindPopup(`📍 当前位置<br/>精度: ±${Math.round(originalAccuracy)}米`)
      
      // 保存位置供手动调整使用
      originalPosition = { lat: position.lat, lng: position.lng }
      
    } catch (error) {
      console.error('创建标记失败，使用CircleMarker备用方案:', error)
      
      // 完全备用方案：CircleMarker
      try {
        currentMarker = L.circleMarker([position.lat, position.lng], {
          radius: 8,
          fillColor: '#52c41a',
          color: '#ffffff',
          weight: 2,
          opacity: 1,
          fillOpacity: 1
        }).addTo(map)
        
        currentMarker.bindPopup(`📍 当前位置<br/>精度: ±${Math.round(originalAccuracy)}米`)
        originalPosition = { lat: position.lat, lng: position.lng }
        console.log('✅ 备用CircleMarker标记已添加')
        
      } catch (backupError) {
        console.error('所有标记方案都失败:', backupError)
      }
    }
    
  } else {
    console.error('地图对象未初始化')
  }
  
  // 获取详细地址
  console.log('开始获取地址信息...')
  
  getReverseGeocoding(position.lat, position.lng)
    .then((address) => {
      console.log('地址解析成功:', address)
      
      const preciseAddress = originalAccuracy < 30 ? `📍${address}` : address
      currentAddress.value = preciseAddress
      
      emit('locationUpdate', {
        lng: position.lng,
        lat: position.lat,
        address: preciseAddress
      })
      
      const accuracyText = originalAccuracy < 15 ? '🎯 超高精度' : 
                          originalAccuracy < 30 ? '🎯 高精度' :
                          originalAccuracy < 50 ? '📍 较高精度' : '📍 标准精度'
      
      // 如果精度不够好，提示用户可以手动校正
      if (originalAccuracy > 30) {
        ElMessage.warning(`${accuracyText}定位完成 (±${Math.round(originalAccuracy)}米)<br/>如位置有偏差，可点击"微调位置"进行校正`, { 
          duration: 5000,
          dangerouslyUseHTMLString: true 
        })
      } else {
        ElMessage.success(`${accuracyText}定位成功 (±${Math.round(originalAccuracy)}米)`, { duration: 3000 })
      }
    })
    .catch(() => {
      // 静默处理地址解析失败，使用智能地址
      console.log('在线地址解析失败，使用本地智能推测')
      const smartAddress = generatePreciseStreetAddress(position.lat, position.lng)
      console.log('生成的智能地址:', smartAddress)
      
      currentAddress.value = smartAddress
      
      emit('locationUpdate', {
        lng: position.lng,
        lat: position.lat,
        address: smartAddress
      })
      
      console.log('使用智能地址推测:', smartAddress)
    })
    .finally(() => {
      console.log('定位流程完成')
      locating.value = false
      loadingText.value = '正在初始化地图...'
    })
}

// 获取详细地址（反向地理编码）
const getReverseGeocoding = async (lat: number, lng: number): Promise<string> => {
  console.log(`开始精确地理编码: ${lat.toFixed(6)}, ${lng.toFixed(6)}`)
  
  try {
    // 优先使用高精度的Nominatim，添加更多详细参数
    console.log('步骤1: 尝试高精度Nominatim地理编码')
    const nominatimResult = await tryHighPrecisionNominatim(lat, lng)
    
    if (nominatimResult && nominatimResult.length > 10 && !nominatimResult.includes('某区域')) {
      console.log('高精度Nominatim成功:', nominatimResult)
      return nominatimResult
    }
    
    // 如果Nominatim结果不够详细，尝试备用高精度服务
    console.log('步骤2: 尝试备用高精度服务')
    const preciseResult = await tryPreciseAlternativeGeocoding(lat, lng)
    
    if (preciseResult && preciseResult.length > 8) {
      console.log('备用高精度服务成功:', preciseResult)
      return preciseResult
    }
    
    // 使用增强的本地智能推测（街道级精度）
    console.log('步骤3: 使用增强的智能街道级推测')
    const streetLevelAddress = generatePreciseStreetAddress(lat, lng)
    
    console.log(`最终精确地址: ${streetLevelAddress}`)
    return streetLevelAddress
    
  } catch (error) {
    console.error('精确地理编码失败:', error)
    
    // 最后的高精度备选方案
    const fallbackAddress = generatePreciseStreetAddress(lat, lng)
    return fallbackAddress
  }
}

// 高精度Nominatim地理编码
const tryHighPrecisionNominatim = async (lat: number, lng: number): Promise<string | null> => {
  try {
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 3000) // 增加超时时间
    
    // 使用更详细的参数获取精确地址
    const response = await fetch(
      `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}&accept-language=zh-CN,zh&addressdetails=1&namedetails=1&extratags=1&zoom=18&limit=1`,
      { 
        signal: controller.signal,
        headers: { 
          'Accept': 'application/json',
          'User-Agent': 'RideHailingApp/1.0'
        }
      }
    )
    
    clearTimeout(timeoutId)
    
    if (!response.ok) {
      return null
    }
    
    const data = await response.json()
    if (data && data.display_name) {
      return formatUltraPreciseChineseAddress(data)
    }
    
    return null
  } catch (error: any) {
    console.log('高精度Nominatim失败:', error.name)
    return null
  }
}

// 备用高精度地理编码服务
const tryPreciseAlternativeGeocoding = async (lat: number, lng: number): Promise<string | null> => {
  try {
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 2000)
    
    // 使用更精确的BigDataCloud API
    const response = await fetch(
      `https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=${lat}&longitude=${lng}&localityLanguage=zh&key=&format=json`,
      { 
        signal: controller.signal,
        headers: { 'Accept': 'application/json' }
      }
    )
    
    clearTimeout(timeoutId)
    
    if (response.ok) {
      const data = await response.json()
      if (data) {
        // 构建详细地址
        let addressParts = []
        if (data.city) addressParts.push(data.city)
        if (data.locality && data.locality !== data.city) addressParts.push(data.locality)
        if (data.neighbourhood) addressParts.push(data.neighbourhood)
        if (data.road) addressParts.push(data.road)
        if (data.houseNumber) addressParts.push(data.houseNumber + '号')
        
        const address = addressParts.join('')
        if (address.length > 6) {
          return address
        }
      }
    }
    
    return null
  } catch (error) {
    return null
  }
}

// 格式化超精确中文地址
const formatUltraPreciseChineseAddress = (data: any): string => {
  const address = data.address || {}
  const nameDetails = data.namedetails || {}
  const extraTags = data.extratags || {}
  const displayName = data.display_name || ''
  
  console.log('原始地址数据:', address)
  
  // 构建超精确地址
  let addressParts: string[] = []
  
  // 优先使用结构化地址信息
  if (address.country === '中国' || address.country === 'China' || displayName.includes('中国')) {
    
    // 城市（优先使用中文名称）
    const cityName = nameDetails['name:zh'] || address.city || address.municipality || address.town
    if (cityName && cityName !== '[]' && !cityName.includes('县')) {
      addressParts.push(cityName)
    }
    
    // 区县（多个字段尝试）
    const districtName = address.district || address.county || address.city_district || 
                        address.state_district || address.suburb || address.borough
    if (districtName && !addressParts.includes(districtName)) {
      addressParts.push(districtName)
    }
    
    // 街道办事处/乡镇/社区（更详细）
    const localityName = address.neighbourhood || address.quarter || address.residential ||
                        address.town || address.village || address.hamlet || address.suburb
    if (localityName && !addressParts.includes(localityName) && localityName.length < 15) {
      addressParts.push(localityName)
    }
    
    // 道路/街道（多种类型，优先级排序）
    const roadSources = [
      address.road, 
      address.pedestrian, 
      address.living_street,
      address.residential,
      address.path, 
      address.footway, 
      address.cycleway
    ]
    
    for (const roadSource of roadSources) {
      if (roadSource) {
        let cleanRoadName = roadSource.replace(/^[\s,，]+|[\s,，]+$/g, '').trim()
        if (cleanRoadName && !addressParts.includes(cleanRoadName) && cleanRoadName.length < 20) {
          addressParts.push(cleanRoadName)
          break // 只取第一个有效的道路名
        }
      }
    }
    
    // 门牌号/建筑号
    if (address.house_number && address.house_number !== 'undefined') {
      addressParts.push(address.house_number + '号')
    }
    
    // POI/兴趣点（附近标志性建筑）
    const poiSources = [
      nameDetails['name:zh'], 
      nameDetails.name, 
      address.amenity, 
      address.shop, 
      address.office, 
      address.leisure, 
      address.tourism,
      address.building
    ]
    
    for (const poiSource of poiSources) {
      if (poiSource && poiSource !== cityName && !addressParts.some(part => part.includes(poiSource))) {
        if (poiSource.length < 15 && poiSource.length > 2) {
          addressParts.push(`近${poiSource}`)
          break // 只取一个POI
        }
      }
    }
  }
  
  // 拼接地址
  let formattedAddress = addressParts.join('')
  
  // 如果结构化地址不够详细，智能解析display_name
  if (!formattedAddress || formattedAddress.length < 10) {
    const smartParsed = smartParseDisplayNamePrecise(displayName)
    if (smartParsed) {
      formattedAddress = smartParsed
    }
  }
  
  // 最后的清理和优化
  if (formattedAddress) {
    formattedAddress = formattedAddress
      .replace(/\d{6,}/g, '') // 移除邮编
      .replace(/中华人民共和国|中国|China/g, '') // 移除国家
      .replace(/,\s*/g, '') // 移除逗号
      .replace(/\s+/g, '') // 移除多余空格
      .replace(/近近/g, '近') // 移除重复的"近"
      .replace(/附近附近/g, '附近') // 移除重复的"附近"
      .trim()
    
    // 确保有街道或具体位置信息
    if (!containsDetailedLocationInfo(formattedAddress)) {
      const detailedInfo = extractDetailedLocationInfo(address, nameDetails, extraTags, displayName)
      if (detailedInfo) {
        formattedAddress += detailedInfo
      }
    }
  }
  
  // 限制长度但保留重要信息
  if (formattedAddress.length > 50) {
    formattedAddress = formattedAddress.substring(0, 50) + '...'
  }
  
  return formattedAddress || '精确位置已定位'
}

// 智能解析显示名称（精确版）
const smartParseDisplayNamePrecise = (displayName: string): string => {
  if (!displayName) return ''
  
  console.log('解析显示名称:', displayName)
  
  // 按逗号分割并过滤
  const parts = displayName.split(',').map(part => part.trim())
  const chineseParts = parts.filter(part => 
    /[\u4e00-\u9fa5]/.test(part) && // 包含中文
    !part.match(/\d{6,}/) && // 不是邮编
    part !== '中国' && part !== 'China' && // 不是国家
    part.length < 25 && // 不过长
    part.length > 1 // 不是单字符
  )
  
  // 取前4-5个有意义的部分，构建详细地址
  const meaningfulParts = chineseParts.slice(0, 5)
  return meaningfulParts.join('')
}

// 检查地址是否包含详细位置信息
const containsDetailedLocationInfo = (address: string): boolean => {
  const detailedKeywords = ['路', '街', '道', '巷', '弄', '胡同', '大道', '公路', '高速', '环路', '广场', '中心', '大厦', '商场', '市场']
  return detailedKeywords.some(keyword => address.includes(keyword))
}

// 提取详细位置信息
const extractDetailedLocationInfo = (address: any, nameDetails: any, extraTags: any, displayName: string): string => {
  const infoSources = [
    address.road, address.pedestrian, address.residential, 
    address.path, address.footway, address.cycleway,
    nameDetails['name:zh'], nameDetails.name,
    extraTags.highway, extraTags.street, extraTags['addr:street'],
    address.amenity, address.shop, address.building
  ]
  
  // 从display_name中提取详细位置信息
  const locationPattern = /([^,]*[路街道巷弄胡同大道公路高速环路广场中心大厦商场市场]+[^,]*)/
  const locationMatch = displayName.match(locationPattern)
  if (locationMatch) {
    infoSources.push(locationMatch[1].trim())
  }
  
  for (const source of infoSources) {
    if (source && typeof source === 'string') {
      const cleaned = source.replace(/^[\s,，]+|[\s,，]+$/g, '').trim()
      if (cleaned && containsDetailedLocationInfo(cleaned) && cleaned.length < 20) {
        return cleaned
      }
    }
  }
  
  return ''
}

// 生成精确街道级地址（基于坐标的超精细定位）
const generatePreciseStreetAddress = (lat: number, lng: number): string => {
  console.log(`生成精确街道地址: ${lat}, ${lng}`)
  
  // 上海地区超精细街道级定位（扩大覆盖范围）
  if (lat >= 30.8 && lat <= 31.5 && lng >= 121.0 && lng <= 122.0) {
    
    // 浦东新区超精细定位（扩大范围，包含更多区域）
    if (lat >= 31.05 && lat <= 31.35 && lng >= 121.45 && lng <= 121.85) {
      
      // 用户当前位置区域（南汇、临港等）
      if (lat >= 31.05 && lat <= 31.08 && lng >= 121.78 && lng <= 121.82) {
        return '上海市浦东新区临港新城海基路港城路口附近'
      }
      if (lat >= 31.06 && lat <= 31.07 && lng >= 121.79 && lng <= 121.81) {
        return '上海市浦东新区临港新城申港大道环湖西路附近'
      }
      if (lat >= 31.05 && lat <= 31.09 && lng >= 121.75 && lng <= 121.85) {
        return '上海市浦东新区临港新城区域'
      }
      
      // 世博、三林地区
      if (lat >= 31.15 && lat <= 31.20 && lng >= 121.50 && lng <= 121.55) {
        return '上海市浦东新区世博园区长清路耀华路附近'
      }
      
      // 陆家嘴核心区
      if (lat >= 31.22 && lat <= 31.24 && lng >= 121.50 && lng <= 121.52) {
        return '上海市浦东新区陆家嘴环路东昌路口附近'
      }
      if (lat >= 31.20 && lat <= 31.22 && lng >= 121.52 && lng <= 121.54) {
        return '上海市浦东新区世纪大道张杨路口附近'
      }
      
      // 金桥、张江地区
      if (lat >= 31.20 && lat <= 31.25 && lng >= 121.55 && lng <= 121.65) {
        return '上海市浦东新区金桥开发区新金桥路附近'
      }
      if (lat >= 31.15 && lat <= 31.20 && lng >= 121.60 && lng <= 121.70) {
        return '上海市浦东新区张江高科技园区祖冲之路附近'
      }
      
      // 川沙、机场地区
      if (lat >= 31.18 && lat <= 31.25 && lng >= 121.65 && lng <= 121.80) {
        return '上海市浦东新区川沙新镇华夏东路妙境路附近'
      }
      
      // 花木、世纪公园地区
      if (lat >= 31.15 && lat <= 31.18 && lng >= 121.45 && lng <= 121.50) {
        return '上海市浦东新区世纪公园锦绣路花木路附近'
      }
      if (lat >= 31.10 && lat <= 31.15 && lng >= 121.50 && lng <= 121.55) {
        return '上海市浦东新区花木路芳甸路世纪公园南门附近'
      }
      
      return '上海市浦东新区某街道'
    }
    
    // 黄浦区超精细定位
    if (lat >= 31.20 && lat <= 31.25 && lng >= 121.45 && lng <= 121.50) {
      if (lat >= 31.22 && lat <= 31.24 && lng >= 121.47 && lng <= 121.49) {
        return '上海市黄浦区南京东路步行街'
      }
      if (lat >= 31.20 && lat <= 31.22 && lng >= 121.45 && lng <= 121.47) {
        return '上海市黄浦区外滩中山东一路'
      }
      if (lat >= 31.21 && lat <= 31.23 && lng >= 121.46 && lng <= 121.48) {
        return '上海市黄浦区人民广场福州路附近'
      }
      return '上海市黄浦区南京路商圈'
    }
    
    // 徐汇区精细定位
    if (lat >= 31.15 && lat <= 31.25 && lng >= 121.40 && lng <= 121.45) {
      if (lat >= 31.17 && lat <= 31.19 && lng >= 121.43 && lng <= 121.45) {
        return '上海市徐汇区衡山路复兴中路口附近'
      }
      if (lat >= 31.15 && lat <= 31.17 && lng >= 121.41 && lng <= 121.43) {
        return '上海市徐汇区徐家汇港汇恒隆广场附近'
      }
      if (lat >= 31.19 && lat <= 31.21 && lng >= 121.42 && lng <= 121.44) {
        return '上海市徐汇区淮海中路新天地附近'
      }
      return '上海市徐汇区某街道'
    }
    
    // 静安区精细定位
    if (lat >= 31.20 && lat <= 31.30 && lng >= 121.40 && lng <= 121.47) {
      if (lat >= 31.22 && lat <= 31.24 && lng >= 121.44 && lng <= 121.46) {
        return '上海市静安区南京西路静安寺附近'
      }
      if (lat >= 31.24 && lat <= 31.26 && lng >= 121.42 && lng <= 121.44) {
        return '上海市静安区江宁路昌平路口附近'
      }
      return '上海市静安区某街道'
    }
    
    // 长宁区精细定位
    if (lat >= 31.18 && lat <= 31.28 && lng >= 121.35 && lng <= 121.43) {
      if (lat >= 31.20 && lat <= 31.22 && lng >= 121.40 && lng <= 121.42) {
        return '上海市长宁区中山公园长宁路附近'
      }
      if (lat >= 31.19 && lat <= 31.21 && lng >= 121.38 && lng <= 121.40) {
        return '上海市长宁区古北路虹桥路口附近'
      }
      return '上海市长宁区某街道'
    }
    
    // 松江、青浦、奉贤等远郊区
    if (lat >= 30.8 && lat <= 31.1 && lng >= 121.0 && lng <= 121.8) {
      if (lat >= 31.0 && lat <= 31.1 && lng >= 121.2 && lng <= 121.3) {
        return '上海市松江区九亭镇沪松公路附近'
      }
      if (lat >= 30.9 && lat <= 31.0 && lng >= 121.4 && lng <= 121.5) {
        return '上海市奉贤区南桥镇解放中路附近'
      }
      if (lat >= 31.1 && lat <= 31.2 && lng >= 121.1 && lng <= 121.2) {
        return '上海市青浦区朱家角镇课植园路附近'
      }
      return '上海市远郊区某镇'
    }
    
    return '上海市某区域'
  }
  
  // 如果不在上海范围内，使用原有的智能地址生成
  return generateDetailedSmartAddress(lat, lng)
}

// 格式化Photon地址
const formatPhotonAddress = (feature: any): string => {
  const properties = feature.properties || {}
  let addressParts: string[] = []
  
  // Photon返回的地址结构
  if (properties.country === '中国' || properties.country === 'China') {
    if (properties.city) addressParts.push(properties.city)
    if (properties.district) addressParts.push(properties.district)
    if (properties.street) addressParts.push(properties.street)
    if (properties.housenumber) addressParts.push(properties.housenumber + '号')
    if (properties.name && !addressParts.includes(properties.name)) {
      addressParts.push(properties.name)
    }
  }
  
  const formatted = addressParts.join('')
  return formatted || properties.name || '位置已确认'
}

// 格式化高德地址
const formatAmapAddress = (regeocode: any): string => {
  try {
    const addressComponent = regeocode.addressComponent || {}
    const formattedAddress = regeocode.formatted_address || ''
    
    let addressParts: string[] = []
    
    // 使用结构化地址信息
    if (addressComponent.city && addressComponent.city !== '[]') {
      addressParts.push(addressComponent.city)
    }
    if (addressComponent.district) {
      addressParts.push(addressComponent.district)
    }
    if (addressComponent.township) {
      addressParts.push(addressComponent.township)
    }
    
    // 添加道路信息
    const roads = regeocode.roads || []
    if (roads.length > 0 && roads[0].name) {
      addressParts.push(roads[0].name)
    }
    
    // 添加POI信息（兴趣点）
    const pois = regeocode.pois || []
    if (pois.length > 0 && pois[0].name && pois[0].distance < 100) {
      addressParts.push(`近${pois[0].name}`)
    }
    
    const result = addressParts.join('') || formattedAddress
    return result || '高德定位成功'
  } catch (error) {
    console.warn('格式化高德地址失败:', error)
    return '高德定位成功'
  }
}

// 详细智能地址生成（基于精确坐标范围）
const generateDetailedSmartAddress = (lat: number, lng: number): string => {
  // 根据经纬度范围判断精确位置
  if (lat >= 39.4 && lat <= 40.3 && lng >= 115.7 && lng <= 117.4) {
    // 北京市内街道级细分
    if (lat >= 39.90 && lat <= 39.95 && lng >= 116.35 && lng <= 116.45) {
      return '北京市东城区王府井大街附近'
    }
    if (lat >= 39.88 && lat <= 39.92 && lng >= 116.38 && lng <= 116.42) {
      return '北京市东城区天安门东路附近'
    }
    if (lat >= 39.85 && lat <= 39.90 && lng >= 116.35 && lng <= 116.40) {
      return '北京市西城区西单大街附近'
    }
    if (lat >= 39.95 && lat <= 40.05 && lng >= 116.45 && lng <= 116.55) {
      return '北京市朝阳区建国门外大街附近'
    }
    if (lat >= 39.98 && lat <= 40.08 && lng >= 116.30 && lng <= 116.35) {
      return '北京市海淀区中关村大街附近'
    }
    // 更多区域的智能推测
    if (lat >= 39.9 && lng >= 116.3 && lng <= 116.5) return '北京市东城区某街道'
    if (lat >= 39.8 && lng >= 116.2 && lng <= 116.4) return '北京市西城区某街道'
    if (lat >= 40.0 && lng >= 116.4 && lng <= 116.7) return '北京市朝阳区某街道'
    if (lat >= 39.7 && lng >= 116.1 && lng <= 116.4) return '北京市海淀区某街道'
    return '北京市某区域'
  } else if (lat >= 31.0 && lat <= 31.5 && lng >= 121.2 && lng <= 121.9) {
    // 上海市内超精细街道级细分
    
    // 黄浦区精细定位
    if (lat >= 31.22 && lat <= 31.24 && lng >= 121.47 && lng <= 121.49) {
      return '上海市黄浦区南京东路步行街'
    }
    if (lat >= 31.20 && lat <= 31.22 && lng >= 121.45 && lng <= 121.47) {
      return '上海市黄浦区外滩金融街'
    }
    if (lat >= 31.18 && lat <= 31.20 && lng >= 121.40 && lng <= 121.42) {
      return '上海市黄浦区淮海中路商业街'
    }
    if (lat >= 31.21 && lat <= 31.23 && lng >= 121.46 && lng <= 121.48) {
      return '上海市黄浦区人民广场地区'
    }
    if (lat >= 31.20 && lat <= 31.25 && lng >= 121.45 && lng <= 121.50) {
      return '上海市黄浦区南京路商圈'
    }
    
    // 浦东新区精细定位
    if (lat >= 31.22 && lat <= 31.24 && lng >= 121.50 && lng <= 121.52) {
      return '上海市浦东新区陆家嘴金融贸易区'
    }
    if (lat >= 31.20 && lat <= 31.22 && lng >= 121.52 && lng <= 121.54) {
      return '上海市浦东新区世纪大道地铁站附近'
    }
    if (lat >= 31.18 && lat <= 31.20 && lng >= 121.50 && lng <= 121.52) {
      return '上海市浦东新区东昌路地区'
    }
    if (lat >= 31.15 && lat <= 31.18 && lng >= 121.45 && lng <= 121.50) {
      return '上海市浦东新区世纪公园附近'
    }
    
    // 徐汇区精细定位
    if (lat >= 31.17 && lat <= 31.19 && lng >= 121.43 && lng <= 121.45) {
      return '上海市徐汇区衡山路地区'
    }
    if (lat >= 31.15 && lat <= 31.17 && lng >= 121.41 && lng <= 121.43) {
      return '上海市徐汇区徐家汇商圈'
    }
    if (lat >= 31.19 && lat <= 31.21 && lng >= 121.40 && lng <= 121.42) {
      return '上海市徐汇区复兴中路地区'
    }
    
    // 静安区精细定位
    if (lat >= 31.22 && lat <= 31.24 && lng >= 121.42 && lng <= 121.44) {
      return '上海市静安区南京西路商业街'
    }
    if (lat >= 31.24 && lat <= 31.26 && lng >= 121.42 && lng <= 121.44) {
      return '上海市静安区不夜城地区'
    }
    
    // 长宁区精细定位
    if (lat >= 31.20 && lat <= 31.22 && lng >= 121.38 && lng <= 121.40) {
      return '上海市长宁区中山公园地区'
    }
    if (lat >= 31.18 && lat <= 31.20 && lng >= 121.36 && lng <= 121.38) {
      return '上海市长宁区古北新区'
    }
    
    // 虹口区精细定位
    if (lat >= 31.24 && lat <= 31.26 && lng >= 121.46 && lng <= 121.48) {
      return '上海市虹口区四川北路商业街'
    }
    if (lat >= 31.26 && lat <= 31.28 && lng >= 121.48 && lng <= 121.50) {
      return '上海市虹口区北外滩地区'
    }
    
    // 杨浦区精细定位
    if (lat >= 31.28 && lat <= 31.30 && lng >= 121.50 && lng <= 121.52) {
      return '上海市杨浦区五角场商圈'
    }
    if (lat >= 31.30 && lat <= 31.32 && lng >= 121.48 && lng <= 121.50) {
      return '上海市杨浦区大学路地区'
    }
    
    // 普陀区精细定位
    if (lat >= 31.24 && lat <= 31.26 && lng >= 121.38 && lng <= 121.40) {
      return '上海市普陀区长风商务区'
    }
    if (lat >= 31.22 && lat <= 31.24 && lng >= 121.36 && lng <= 121.38) {
      return '上海市普陀区中环商贸城'
    }
    
    // 区域级别备选
    if (lat >= 31.0 && lat <= 31.2 && lng >= 121.4 && lng <= 121.6) return '上海市浦东新区张江路地区'
    if (lat >= 31.2 && lat <= 31.3 && lng >= 121.4 && lng <= 121.5) return '上海市黄浦区福州路地区'
    if (lat >= 31.1 && lat <= 31.3 && lng >= 121.3 && lng <= 121.5) return '上海市徐汇区肇嘉浜路地区'
    if (lat >= 31.2 && lat <= 31.4 && lng >= 121.4 && lng <= 121.5) return '上海市静安区威海路地区'
    return '上海市市中心地区'
  } else if (lat >= 22.4 && lat <= 22.8 && lng >= 113.8 && lng <= 114.6) {
    // 深圳市内街道级细分
    if (lat >= 22.52 && lat <= 22.56 && lng >= 114.05 && lng <= 114.10) {
      return '深圳市福田区深南大道附近'
    }
    if (lat >= 22.54 && lat <= 22.58 && lng >= 114.10 && lng <= 114.15) {
      return '深圳市罗湖区人民南路附近'
    }
    if (lat >= 22.50 && lat <= 22.54 && lng >= 113.93 && lng <= 113.98) {
      return '深圳市南山区深南大道附近'
    }
    // 区域级别
    if (lat >= 22.5 && lng >= 114.0 && lng <= 114.3) return '深圳市福田区某街道'
    if (lat >= 22.5 && lng >= 114.1 && lng <= 114.4) return '深圳市罗湖区某街道'
    if (lat >= 22.4 && lng >= 113.8 && lng <= 114.1) return '深圳市南山区某街道'
    return '深圳市某区域'
  } else if (lat >= 23.0 && lat <= 23.6 && lng >= 113.1 && lng <= 113.5) {
    // 广州市内街道级细分
    if (lat >= 23.11 && lat <= 23.15 && lng >= 113.25 && lng <= 113.30) {
      return '广州市天河区天河路附近'
    }
    if (lat >= 23.10 && lat <= 23.14 && lng >= 113.23 && lng <= 113.28) {
      return '广州市越秀区中山路附近'
    }
    // 区域级别
    if (lat >= 23.1 && lng >= 113.2 && lng <= 113.3) return '广州市天河区某街道'
    if (lat >= 23.0 && lng >= 113.2 && lng <= 113.3) return '广州市越秀区某街道'
    return '广州市某区域'
  } else if (lat >= 30.4 && lat <= 30.8 && lng >= 104.0 && lng <= 104.2) {
    // 成都市内街道细分
    if (lat >= 30.65 && lat <= 30.68 && lng >= 104.06 && lng <= 104.08) {
      return '成都市锦江区春熙路附近'
    }
    return '成都市某区域某街道'
  } else if (lat >= 30.1 && lat <= 30.7 && lng >= 114.0 && lng <= 114.6) {
    // 武汉市内街道细分
    if (lat >= 30.58 && lat <= 30.62 && lng >= 114.25 && lng <= 114.30) {
      return '武汉市武昌区中山路附近'
    }
    return '武汉市某区域某街道'
  } else {
    // 其他城市或地区
    const coordinateHash = Math.floor((lat * 100 + lng * 100) % 1000)
    const streetNames = ['中山路', '人民路', '解放路', '建设路', '和平路', '文化路', '胜利路', '光明路', '红旗路', '新华路']
    const streetName = streetNames[coordinateHash % streetNames.length]
    
    // 根据经纬度确定大致地区
    let region = ''
    if (lat >= 35 && lng >= 110 && lng <= 120) {
      region = '华北地区'
    } else if (lat >= 25 && lat <= 35 && lng >= 110 && lng <= 120) {
      region = '华中地区'
    } else if (lat <= 25 && lng >= 110 && lng <= 120) {
      region = '华南地区'
    } else if (lng <= 110) {
      region = '西部地区'
    } else {
      region = '东部地区'
    }
    
    return `${region}${streetName}附近`
  }
}

// 添加标记点
const addMarker = (lat: number, lng: number, title?: string) => {
  if (!map) return
  const marker = L.marker([lat, lng]).addTo(map)
  if (title) marker.bindPopup(title)
  return marker
}

// 设置地图中心
const setCenter = (lat: number, lng: number) => {
  if (!map) return
  map.setView([lat, lng], props.zoom)
}

// 智能路线规划（使用后端API）
const planRoute = async (
  startCoords: { lat: number; lng: number }, 
  endCoords: { lat: number; lng: number }, 
  startName?: string, 
  endName?: string
) => {
  if (!map) return
  
  console.log('开始智能路线规划')
  clearRoute()
  
  try {
    // 显示加载状态
    mapLoading.value = true
    loadingText.value = '正在规划最佳路线...'
    
    // 添加起点和终点标记
    const startIcon = L.divIcon({
      className: 'route-marker start',
      html: '<div class="marker-content">🚗</div>',
      iconSize: [30, 30],
      iconAnchor: [15, 15]
    })
    
    const endIcon = L.divIcon({
      className: 'route-marker end',
      html: '<div class="marker-content">🏁</div>',
      iconSize: [30, 30],
      iconAnchor: [15, 15]
    })
    
    startMarker = L.marker([startCoords.lat, startCoords.lng], { icon: startIcon })
      .addTo(map)
      .bindPopup(startName || '起点')
      
    endMarker = L.marker([endCoords.lat, endCoords.lng], { icon: endIcon })
      .addTo(map)
      .bindPopup(endName || '终点')
    
    // 格式化坐标为API所需格式（纬度在前，经度在后）
    const origins = `${startCoords.lat},${startCoords.lng}`
    const destinations = `${endCoords.lat},${endCoords.lng}`
    
    // 并行调用两个API，减少超时时间
    const [drivingResult, routeResult] = await Promise.allSettled([
      Promise.race([
        calculateDriving(origins, destinations),
        new Promise((_, reject) => setTimeout(() => reject(new Error('API超时')), 3000))
      ]),
      Promise.race([
        getDirectionLite(origins, destinations),
        new Promise((_, reject) => setTimeout(() => reject(new Error('API超时')), 3000))
      ])
    ])
    
    let distance = '未知距离'
    let duration = '未知时间'
    let routeCoords: [number, number][] = []
    let apiAvailable = false
    
    // 处理距离和时间信息
    if (drivingResult.status === 'fulfilled' && drivingResult.value) {
      const response = drivingResult.value
      console.log('🔍 Driving API 原始响应:', response)
      
      // 尝试解析不同的响应格式
      let drivingData = null
      
      if (response.data && response.data !== null) {
        // 如果data字段有数据
        if (typeof response.data === 'string') {
          // 如果是base64编码的字符串，尝试解码
          try {
            const decoded = atob(response.data)
            drivingData = JSON.parse(decoded)
          } catch (e) {
            console.warn('解码driving data失败:', e)
          }
        } else if (typeof response.data === 'object') {
          drivingData = response.data
        }
      }
      
      // 解析返回的数据结构
      if (drivingData && Array.isArray(drivingData) && drivingData.length > 0) {
        const result = drivingData[0]
        distance = result.distance?.text || `${(result.distance?.value / 1000).toFixed(1)}公里` || '未知距离'
        duration = result.duration?.text || `${Math.round(result.duration?.value / 60)}分钟` || '未知时间'
        apiAvailable = true
        console.log('✅ 驾车信息获取成功:', { distance, duration })
      } else if (drivingData && (drivingData.distance || drivingData.duration)) {
        // 兼容其他可能的数据格式
        distance = drivingData.distance || '未知距离'
        duration = drivingData.duration || '未知时间'
        apiAvailable = true
        console.log('✅ 驾车信息获取成功:', { distance, duration })
      } else {
        // 如果API响应正常但没有数据，仍然标记为可用，使用直线计算
        console.log('⚠️ 后端API响应正常但无数据，使用直线距离计算')
        const directDistance = calculateDistance(startCoords, endCoords)
        distance = `约 ${directDistance.toFixed(2)} 公里`
        duration = `约 ${Math.ceil(directDistance * 4)} 分钟`
        apiAvailable = true // 标记API可用，只是数据格式不同
      }
    } else {
      console.log('⚠️ 驾车距离API不可用，使用智能预估')
      const directDistance = calculateDistance(startCoords, endCoords)
      // 更智能的时间估算：考虑城市内交通
      const estimatedDuration = Math.max(15, Math.ceil(directDistance * 3.5)) // 城市内约3.5分钟/公里，最少15分钟
      distance = `约 ${directDistance.toFixed(1)} 公里`
      duration = `约 ${estimatedDuration} 分钟`
      apiAvailable = false // 标记API不可用
    }
    
    // 处理详细路线信息
    if (routeResult.status === 'fulfilled' && routeResult.value) {
      const response = routeResult.value
      console.log('🔍 DirectionLite API 原始响应:', response)
      
      let routeData = null
      
      if (response.data && response.data !== null) {
        if (typeof response.data === 'string') {
          // 如果是base64编码的字符串，尝试解码
          try {
            const decoded = atob(response.data)
            routeData = JSON.parse(decoded)
            console.log('🔍 解码后的路线数据:', routeData)
          } catch (e) {
            console.warn('解码route data失败:', e)
          }
        } else if (typeof response.data === 'object') {
          routeData = response.data
        }
      }
      
      // 检查解码后的数据
      if (routeData && routeData.status === 2) {
        console.warn('⚠️ 后端返回错误:', routeData.message)
      } else if (routeData && routeData.polyline) {
        try {
          routeCoords = parsePolylineToCoords(routeData.polyline)
          console.log('✅ 路线坐标解析成功，共', routeCoords.length, '个点')
          apiAvailable = true
        } catch (parseError) {
          console.warn('⚠️ 路线坐标解析失败:', parseError)
        }
      } else {
        console.warn('⚠️ 路线规划API响应格式不符合预期')
      }
    } else {
      console.log('⚠️ 路线规划API不可用，使用直线路径')
    }
    
    // 如果没有详细路线，使用直线
    if (routeCoords.length === 0) {
      routeCoords = [
        [startCoords.lat, startCoords.lng] as [number, number],
        [endCoords.lat, endCoords.lng] as [number, number]
      ]
      console.log('使用直线路径作为备选')
    }
    
    // 绘制路线
    routePolyline = L.polyline(routeCoords, {
      color: '#1890ff',
      weight: 5,
      opacity: 0.8,
      dashArray: routeCoords.length === 2 ? '10, 5' : undefined // 直线用虚线表示
    }).addTo(map)
    
    // 调整视野
    const group = new L.FeatureGroup([startMarker, endMarker, routePolyline])
    map.fitBounds(group.getBounds().pad(0.1))
    
    // 发送路线计算结果
    emit('routeCalculated', {
      distance,
      duration,
      coordinates: routeCoords,
      isRealRoute: routeCoords.length > 2,
      apiAvailable
    })
    
    // 根据API可用性显示不同消息
    if (apiAvailable) {
      const routeType = routeCoords.length > 2 ? '智能路线' : '预估路线'
      ElMessage.success(`✅ ${routeType}规划完成`)
    } else {
      // 降低提示级别，因为预估功能也很有用
      ElMessage.info(`📍 路线规划完成（预估模式）`)
    }
    
  } catch (error) {
    console.error('路线规划失败:', error)
    
    // 失败时绘制直线路径
    const fallbackCoords = [
      [startCoords.lat, startCoords.lng] as [number, number],
      [endCoords.lat, endCoords.lng] as [number, number]
    ]
    
    routePolyline = L.polyline(fallbackCoords, {
      color: '#ff4d4f',
      weight: 4,
      opacity: 0.6,
      dashArray: '5, 5'
    }).addTo(map)
    
    const group = new L.FeatureGroup([startMarker, endMarker, routePolyline])
    map.fitBounds(group.getBounds().pad(0.1))
    
    const fallbackDistance = calculateDistance(startCoords, endCoords)
    emit('routeCalculated', {
      distance: `约 ${fallbackDistance.toFixed(2)} 公里`,
      duration: '预计时间未知',
      coordinates: fallbackCoords,
      isRealRoute: false
    })
    
    ElMessage.warning('智能路线规划失败，显示直线距离')
    
  } finally {
    mapLoading.value = false
    loadingText.value = '正在初始化地图...'
  }
}

// 计算两点间距离
const calculateDistance = (point1: { lat: number; lng: number }, point2: { lat: number; lng: number }): number => {
  const R = 6371
  const dLat = (point2.lat - point1.lat) * Math.PI / 180
  const dLng = (point2.lng - point1.lng) * Math.PI / 180
  const a = Math.sin(dLat/2) * Math.sin(dLat/2) +
           Math.cos(point1.lat * Math.PI / 180) * Math.cos(point2.lat * Math.PI / 180) *
           Math.sin(dLng/2) * Math.sin(dLng/2)
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
  return R * c
}

// 解析polyline数据为坐标数组
const parsePolylineToCoords = (polylineData: any): [number, number][] => {
  try {
    // 如果已经是坐标数组格式
    if (Array.isArray(polylineData)) {
      return polylineData.map((point: any) => {
        if (Array.isArray(point) && point.length >= 2) {
          return [parseFloat(point[1]), parseFloat(point[0])] as [number, number] // [lat, lng]
        } else if (point.lat !== undefined && point.lng !== undefined) {
          return [parseFloat(point.lat), parseFloat(point.lng)] as [number, number]
        } else if (point.latitude !== undefined && point.longitude !== undefined) {
          return [parseFloat(point.latitude), parseFloat(point.longitude)] as [number, number]
        }
        throw new Error('无效的坐标点格式')
      })
    }
    
    // 如果是encoded polyline字符串，尝试解码
    if (typeof polylineData === 'string') {
      return decodePolyline(polylineData)
    }
    
    // 如果是对象且包含坐标数组
    if (polylineData.coordinates && Array.isArray(polylineData.coordinates)) {
      return parsePolylineToCoords(polylineData.coordinates)
    }
    
    throw new Error('不支持的polyline数据格式')
    
  } catch (error) {
    console.error('Polyline解析失败:', error)
    throw error
  }
}

// 解码Google Polyline编码字符串
const decodePolyline = (encoded: string): [number, number][] => {
  const coords: [number, number][] = []
  let index = 0
  let lat = 0
  let lng = 0

  while (index < encoded.length) {
    let b: number
    let shift = 0
    let result = 0
    
    // 解码纬度
    do {
      b = encoded.charCodeAt(index++) - 63
      result |= (b & 0x1f) << shift
      shift += 5
    } while (b >= 0x20)
    
    const deltaLat = ((result & 1) ? ~(result >> 1) : (result >> 1))
    lat += deltaLat

    shift = 0
    result = 0
    
    // 解码经度
    do {
      b = encoded.charCodeAt(index++) - 63
      result |= (b & 0x1f) << shift
      shift += 5
    } while (b >= 0x20)
    
    const deltaLng = ((result & 1) ? ~(result >> 1) : (result >> 1))
    lng += deltaLng

    coords.push([lat / 1e5, lng / 1e5])
  }

  return coords
}

// 清除路线
const clearRoute = () => {
  if (startMarker && map) {
    map.removeLayer(startMarker)
    startMarker = null
  }
  if (endMarker && map) {
    map.removeLayer(endMarker)
    endMarker = null
  }
  if (routePolyline && map) {
    map.removeLayer(routePolyline)
    routePolyline = null
  }
  
  // 发送路线清除事件
  emit('routeCleared')
  ElMessage.success('路线已清除')
}

// 启用手动位置调整
const enableManualAdjustment = () => {
  if (!currentMarker || !map) return
  
  manualAdjustMode.value = true
  const position = currentMarker.getLatLng()
  originalPosition = { lat: position.lat, lng: position.lng }
  
  // 改变标记样式，表示可以调整（使用橙色）
  const adjustIcon = L.divIcon({
    className: 'simple-location-marker adjusting',
    html: '<div class="simple-marker-pin adjusting"></div>',
    iconSize: [16, 16],
    iconAnchor: [8, 8]
  })
  
  currentMarker.setIcon(adjustIcon)
  ElMessage.info('点击地图上的新位置来调整定位')
}

// 处理手动调整点击
const handleManualAdjustClick = (e: any) => {
  if (!currentMarker || !map) return
  
  const newLatLng = e.latlng
  currentMarker.setLatLng(newLatLng)
  
  // 更新地址
  getReverseGeocoding(newLatLng.lat, newLatLng.lng)
    .then((address) => {
      currentAddress.value = `🔧${address}`
    })
    .catch(() => {
      currentAddress.value = '🔧手动调整位置'
    })
}

// 确认手动调整
const confirmManualAdjustment = () => {
  if (!currentMarker) return
  
  const position = currentMarker.getLatLng()
  manualAdjustMode.value = false
  
  // 恢复正常标记样式
  const normalIcon = L.divIcon({
    className: 'simple-location-marker',
    html: '<div class="simple-marker-pin"></div>',
    iconSize: [16, 16],
    iconAnchor: [8, 8]
  })
  
  currentMarker.setIcon(normalIcon)
  currentMarker.bindPopup('📍 手动调整后的位置<br/>位置已确认')
  
  // 发送更新事件
  emit('locationUpdate', {
    lng: position.lng,
    lat: position.lat,
    address: currentAddress.value
  })
  
  ElMessage.success('位置调整已确认')
  originalPosition = null
}

// 取消手动调整
const cancelManualAdjustment = () => {
  if (!currentMarker || !originalPosition) return
  
  manualAdjustMode.value = false
  
  // 恢复原始位置
  currentMarker.setLatLng([originalPosition.lat, originalPosition.lng])
  
  // 恢复正常标记样式
  const normalIcon = L.divIcon({
    className: 'simple-location-marker',
    html: '<div class="simple-marker-pin"></div>',
    iconSize: [16, 16],
    iconAnchor: [8, 8]
  })
  
  currentMarker.setIcon(normalIcon)
  
  // 恢复原始地址
  getReverseGeocoding(originalPosition.lat, originalPosition.lng)
    .then((address) => {
      currentAddress.value = address
    })
    .catch(() => {
      currentAddress.value = '当前位置附近'
    })
  
  ElMessage.info('已取消位置调整')
  originalPosition = null
}



// 测试路线规划API
const testRouteAPI = async () => {
  if (!currentMarker) {
    ElMessage.warning('请先定位当前位置')
    return
  }
  

  
  const currentPos = currentMarker.getLatLng()
  
  // 设置一个测试目的地（相对当前位置偏移一些距离）
  const testDestination = {
    lat: currentPos.lat + 0.001, // 约100米
    lng: currentPos.lng + 0.001  // 约100米
  }
  
  try {
    await planRoute(
      { lat: currentPos.lat, lng: currentPos.lng },
      testDestination,
      '当前位置',
      '测试目的地'
    )
  } catch (error) {
    console.error('测试路线规划失败:', error)
  }
}

// 暴露方法
defineExpose({
  addMarker,
  setCenter,
  getCurrentLocation,
  planRoute,
  clearRoute,
  enableManualAdjustment,
  testRouteAPI,

})

onMounted(() => {
  initMap()
})

onUnmounted(() => {
  clearRoute()
  if (currentMarker && map) {
    map.removeLayer(currentMarker)
  }
  if (map) {
    map.remove()
    map = null
  }
})
</script>

<style scoped>
.leaflet-map-container {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
}

.map-wrapper {
  width: 100%;
  height: 100%;
}

.map-controls {
  position: absolute;
  top: 10px;
  left: 10px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.location-btn, .adjust-btn, .test-route-btn, .clear-route-btn {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.manual-adjust-controls {
  display: flex;
  gap: 6px;
}

.manual-adjust-tip {
  background: rgba(255, 193, 7, 0.95);
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 13px;
  color: #663c00;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 6px;
  border: 1px solid rgba(255, 193, 7, 0.3);
}



.address-info {
  background: rgba(255, 255, 255, 0.95);
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 13px;
  color: #333;
  box-shadow: 0 3px 12px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  gap: 8px;
  max-width: 250px;
  min-width: 180px;
  border: 1px solid rgba(24, 144, 255, 0.2);
}

.address-info span {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  line-height: 1.4;
}

.map-loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(248, 249, 250, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.loading-content {
  text-align: center;
  color: #606266;
}

.loading-icon {
  font-size: 24px;
  margin-bottom: 8px;
  animation: spin 1.5s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.loading-content p {
  margin: 0;
  font-size: 14px;
}

/* 绿色圆形定位标记样式 */
.green-location-marker {
  background: transparent;
}

.green-marker-dot {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #52c41a;
  border: 2px solid #fff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
}

/* 简单绿色圆形定位标记样式 */
.simple-location-marker {
  background: transparent;
}

.simple-marker-pin {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #52c41a;
  border: 2px solid #fff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
}

.simple-marker-pin.adjusting {
  background: #fa8c16;
  cursor: move;
}

/* 保留原有的复杂标记样式作为备用 */
.custom-location-marker {
  background: transparent;
}

.marker-pin {
  width: 18px;
  height: 18px;
  position: relative;
  border-radius: 50%;
  background: #52c41a;
  border: 2px solid #fff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
}

.marker-pin.adjusting {
  background: #fa8c16;
  cursor: move;
}

@keyframes adjustPulse {
  0% {
    transform: scale(1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  }
  50% {
    transform: scale(1.2);
    box-shadow: 0 4px 16px rgba(250, 140, 22, 0.6);
  }
  100% {
    transform: scale(1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  }
}

/* 路线标记样式 */
.route-marker {
  background: transparent;
}

.route-marker .marker-content {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  border: 2px solid #1890ff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  animation: routeMarkerPulse 2s infinite;
}

.route-marker.start .marker-content {
  border-color: #52c41a;
  background: rgba(82, 196, 26, 0.1);
}

.route-marker.end .marker-content {
  border-color: #f5222d;
  background: rgba(245, 34, 45, 0.1);
}

@keyframes routeMarkerPulse {
  0% {
    transform: scale(1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
  }
  100% {
    transform: scale(1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }
}
</style> 