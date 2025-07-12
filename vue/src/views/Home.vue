<template>
  <div class="home-container">
    <!-- 顶部栏 -->
    <div class="top-bar">
      <div class="city-select">
        <el-icon><Location /></el-icon>
        <span>北京</span>
        <el-icon class="arrow"><ArrowDown /></el-icon>
      </div>
      <div class="user-info" @click="showLoginDialog = true">
        <el-icon><User /></el-icon>
        <span>登录</span>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="main-content">
      <!-- 地图区域 -->
      <div class="map-section">
        <LeafletMap 
          ref="mapRef"
          height="200px"
          :zoom="16"
          :autoLocate="shouldAutoLocate"
          @location-update="handleLocationUpdate"
          @map-click="handleMapClick"
          @route-calculated="handleRouteCalculated"
        />
      </div>

      <!-- 打车卡片 -->
      <div class="ride-card">
        <div class="ride-header">
          <h2>现在叫车</h2>
          <div class="weather">
            <el-icon><Sunny /></el-icon>
            <span>晴 28°</span>
          </div>
        </div>

        <div class="location-inputs">
          <div class="location-item">
            <div class="location-dot start"></div>
            <div 
              class="location-search-input-wrapper"
              @click="openAddressInput('start')"
            >
              <el-input
                v-model="startLocation"
                placeholder="您在哪儿？"
                readonly
                class="location-search-input"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
            </div>
          </div>
          <div class="swap-btn" @click="swapLocations">
            <el-icon><Switch /></el-icon>
          </div>
          <div class="location-item">
            <div class="location-dot end"></div>
            <div 
              class="location-search-input-wrapper"
              @click="openAddressInput('end')"
            >
              <el-input
                v-model="endLocation"
                placeholder="您要去哪儿？"
                readonly
                class="location-search-input"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
            </div>
          </div>
        </div>

        <!-- 路线信息显示 -->
        <div v-if="routeInfo" class="route-info" :class="{ 'route-error': isRouteError }">
          <div class="route-info-item">
            <el-icon class="route-icon"><Position /></el-icon>
            <span class="route-text">距离：{{ routeInfo.distance }}</span>
          </div>
          <div class="route-info-item">
            <el-icon class="route-icon"><Clock /></el-icon>
            <span class="route-text">预计：{{ routeInfo.duration }}</span>
          </div>
          <div class="route-clear-btn" @click="clearCurrentRoute">
            <el-icon><Close /></el-icon>
          </div>
          
          <!-- 错误状态提示 -->
          <div v-if="isRouteError" class="route-error-tip">
            <p>💡 提示：距离计算接口正在开发中</p>
          </div>
        </div>

        <div class="ride-options">
          <div 
            v-for="option in rideOptions" 
            :key="option.type"
            class="ride-option"
            :class="{ active: selectedRide === option.type }"
            @click="selectedRide = option.type"
          >
            <div class="option-icon">
              <el-icon><component :is="option.icon" /></el-icon>
            </div>
            <div class="option-info">
              <div class="option-name">{{ option.name }}</div>
              <div class="option-price">{{ option.price }}</div>
            </div>
          </div>
        </div>

        <el-button 
          type="primary" 
          size="large" 
          class="call-car-btn"
          @click="callCar"
        >
          立即叫车
        </el-button>
      </div>

      <!-- 功能菜单 -->
      <div class="feature-menu">
        <div 
          v-for="feature in features" 
          :key="feature.name"
          class="feature-item"
          @click="handleFeatureClick(feature)"
        >
          <div class="feature-icon">
            <el-icon><component :is="feature.icon" /></el-icon>
          </div>
          <span>{{ feature.name }}</span>
        </div>
      </div>
    </div>

    <!-- 登录对话框 -->
    <el-dialog 
      v-model="showLoginDialog" 
      title="登录" 
      width="300px"
      :show-close="false"
    >
      <LoginForm @login-success="handleLoginSuccess" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Location,
  ArrowDown,
  User,
  Sunny,
  Switch,
  Car,
  Van,
  Truck,
  Bell,
  Gift,
  CreditCard,
  More,
  Search,
  Position,
  Clock,
  Close
} from '@element-plus/icons-vue'
import LoginForm from '@/components/LoginForm.vue'
import LeafletMap from '@/components/LeafletMap.vue'
import { createHistoricalSearch } from '@/api/index'

const router = useRouter()
const route = useRoute()

const startLocation = ref('定位中...')
const endLocation = ref('')
const selectedRide = ref('express')
const showLoginDialog = ref(false)
const mapRef = ref()

// 当前位置信息
const currentLocation = ref<{ lng: number; lat: number; address: string } | null>(null)

// 地址搜索相关
const searchLocations = ref<Array<{ value: string; name: string; address: string; lat: number; lng: number }>>([])
const isSearching = ref(false)
const endLocationData = ref<{ lng: number; lat: number; address: string } | null>(null)

// 路线信息
const routeInfo = ref<{ distance: string; duration: string } | null>(null)

// 判断路线是否处于错误状态
const isRouteError = computed(() => {
  if (!routeInfo.value) return false
  const errorStates = ['服务暂不可用', '网络错误', '计算失败', '计算中...']
  return errorStates.some(state => 
    routeInfo.value?.distance.includes(state) || routeInfo.value?.duration.includes(state)
  )
})

// 检查用户是否已登录
const isLoggedIn = computed(() => {
  return !!localStorage.getItem('token')
})

// 保存位置到历史记录
const saveLocationToHistory = async (location: any) => {
  console.log('🔍 saveLocationToHistory 被调用，地址:', location.name)
  console.log('🔍 当前登录状态:', isLoggedIn.value)
  console.log('🔍 localStorage token:', localStorage.getItem('token') ? '存在' : '不存在')
  
  // 只有登录用户才保存历史记录
  if (!isLoggedIn.value) {
    console.log('❌ 用户未登录，跳过保存历史记录')
    return
  }

  try {
    // 提取区域信息
    const extractDistrict = (address: string): string => {
      const match = address.match(/(.*?(区|县))/)
      return match ? match[0] : address
    }

    let idAddr = '未知区域'
    if (currentLocation.value?.address) {
      idAddr = extractDistrict(currentLocation.value.address)
    }

    // 构建历史记录数据
    const historyData = {
      addrName: location.name || location.address || '未知目的地',
      idAddr: idAddr,
      lat: location.lat || 0,
      lng: location.lng || 0,
      searchTime: new Date().toISOString()
    }

    console.log('💾 首页保存历史搜索记录:', historyData)

    // 调用后端API保存历史记录
    const response = await createHistoricalSearch(historyData)
    
    if (response && response.code === 200) {
      console.log('✅ 首页历史记录保存成功')
    } else {
      console.warn('⚠️ 首页历史记录保存失败:', response?.message)
    }
  } catch (error: any) {
    console.error('❌ 首页保存历史记录时发生错误:', error)
    
    // 静默处理错误，不干扰用户正常使用
    if (error.response?.status === 401 || error.response?.status === 403) {
      console.log('Token已过期，需要重新登录')
      localStorage.removeItem('token')
    }
  }
}

// 添加响应式变量控制是否自动定位
const shouldAutoLocate = ref(true)

const rideOptions = [
  {
    type: 'express',
    name: '快车',
    price: '预估 ¥15',
    icon: 'Car'
  },
  {
    type: 'comfort',
    name: '舒适型',
    price: '预估 ¥22',
    icon: 'Van'
  },
  {
    type: 'business',
    name: '商务车',
    price: '预估 ¥45',
    icon: 'Truck'
  }
]

const features = [
  { name: '代叫车', icon: 'Bell' },
  { name: '礼品车', icon: 'Gift' },
  { name: '企业用车', icon: 'CreditCard' },
  { name: '更多', icon: 'More' }
]

const swapLocations = () => {
  // 交换显示的地址
  const tempLocation = startLocation.value
  startLocation.value = endLocation.value
  endLocation.value = tempLocation
  
  // 交换坐标数据
  const tempData = currentLocation.value
  currentLocation.value = endLocationData.value
  endLocationData.value = tempData
  
  // 如果交换后起点有坐标，更新地图中心
  if (currentLocation.value && mapRef.value) {
    mapRef.value.setCenter(currentLocation.value.lat, currentLocation.value.lng)
  }
  
  ElMessage.success('起点和终点已交换')
  
  // 重新规划路线
  setTimeout(() => {
    checkAndPlanRoute()
  }, 300)
}

const callCar = () => {
  if (!endLocation.value) {
    ElMessage.warning('请输入或选择目的地')
    return
  }
  
  if (!currentLocation.value) {
    ElMessage.warning('请等待获取当前位置或手动选择出发地')
    return
  }
  
  // 显示叫车信息
  const startAddr = currentLocation.value.address
  const endAddr = endLocationData.value?.address || endLocation.value
  
  console.log('叫车信息:', `从 ${startAddr} 到 ${endAddr}`)
  
  showLoginDialog.value = true
}

const handleFeatureClick = (feature: any) => {
  ElMessage.info(`点击了${feature.name}`)
}

const handleLoginSuccess = () => {
  showLoginDialog.value = false
  ElMessage.success('登录成功！')
  // 登录成功后刷新页面状态，但不重新加载整个页面
  setTimeout(() => {
    location.reload()
  }, 500)
}

// 打开地址输入页面
const openAddressInput = (focus?: 'start' | 'end') => {
  router.push({
    name: 'AddressInput',
    query: {
      startLocation: startLocation.value,
      endLocation: endLocation.value,
      focus: focus || 'end'
    }
  })
}

// 处理位置更新
const handleLocationUpdate = (location: { lng: number; lat: number; address: string }) => {
  currentLocation.value = location
  startLocation.value = location.address
  console.log('位置更新:', location.address)
}

// 处理地图点击
const handleMapClick = (event: any) => {
  console.log('地图点击:', event)
  // 可以在这里添加地图点击后的逻辑，比如设置目的地
}

// 处理路线计算完成
const handleRouteCalculated = (info: { distance: string; duration: string; coordinates: any[] }) => {
  routeInfo.value = {
    distance: info.distance,
    duration: info.duration
  }
  console.log('路线信息更新:', `距离: ${info.distance}, 时间: ${info.duration}`)
}

// 检查并触发路线规划
const checkAndPlanRoute = () => {
  if (currentLocation.value && endLocationData.value && mapRef.value) {
    console.log('开始规划路线:', `${currentLocation.value?.address} -> ${endLocationData.value?.address}`)
    mapRef.value.planRoute(
      { lat: currentLocation.value.lat, lng: currentLocation.value.lng },
      { lat: endLocationData.value.lat, lng: endLocationData.value.lng },
      currentLocation.value.address,
      endLocationData.value.address
    )
  }
}

// 清除路线
const clearCurrentRoute = () => {
  if (mapRef.value) {
    mapRef.value.clearRoute()
  }
  routeInfo.value = null
}

// 搜索起点地址（优化为主要使用本地数据库）
const searchStartLocation = async (queryString: string, callback: (suggestions: any[]) => void) => {
  if (!queryString || queryString.length < 2) {
    callback([])
    return
  }

  isSearching.value = true
  
  // 直接使用本地热门地址库
  const popularLocations = getPopularLocations(queryString)
  callback(popularLocations)
  
  isSearching.value = false
}

// 获取热门地址（降级方案）
const getPopularLocations = (query: string) => {
  const popularPlaces = [
    // 测试路线对 - 方便演示
    { name: '上海浦东机场T1航站楼', address: '上海市浦东新区机场大道900号T1', lat: 31.1434, lng: 121.8052 },
    { name: '陆家嘴金融中心', address: '上海市浦东新区陆家嘴环路1000号', lat: 31.2415, lng: 121.5056 },
    { name: '上海火车站南广场', address: '上海市静安区秣陵路100号', lat: 31.2496, lng: 121.4553 },
    { name: '外滩观景台', address: '上海市黄浦区中山东一路500号', lat: 31.2401, lng: 121.4900 },
    { name: '虹桥机场T2航站楼', address: '上海市闵行区虹桥路2550号T2', lat: 31.1979, lng: 121.3364 },
    { name: '南京路步行街入口', address: '上海市黄浦区南京东路20号', lat: 31.2385, lng: 121.4737 },
    
    // 上海地区 - 机场
    { name: '上海浦东国际机场', address: '上海市浦东新区机场大道900号', lat: 31.1434, lng: 121.8052 },
    { name: '上海浦东机场', address: '上海市浦东新区机场大道900号', lat: 31.1434, lng: 121.8052 },
    { name: '浦东机场', address: '上海市浦东新区机场大道900号', lat: 31.1434, lng: 121.8052 },
    { name: '浦东国际机场', address: '上海市浦东新区机场大道900号', lat: 31.1434, lng: 121.8052 },
    { name: '上海虹桥国际机场', address: '上海市闵行区虹桥路2550号', lat: 31.1979, lng: 121.3364 },
    { name: '上海虹桥机场', address: '上海市闵行区虹桥路2550号', lat: 31.1979, lng: 121.3364 },
    { name: '虹桥机场', address: '上海市闵行区虹桥路2550号', lat: 31.1979, lng: 121.3364 },
    { name: '虹桥国际机场', address: '上海市闵行区虹桥路2550号', lat: 31.1979, lng: 121.3364 },
    
    // 上海地区 - 火车站
    { name: '上海火车站', address: '上海市静安区秣陵路100号', lat: 31.2496, lng: 121.4553 },
    { name: '上海站', address: '上海市静安区秣陵路100号', lat: 31.2496, lng: 121.4553 },
    { name: '上海南站', address: '上海市徐汇区石龙路100号', lat: 31.1416, lng: 121.4346 },
    { name: '南站', address: '上海市徐汇区石龙路100号', lat: 31.1416, lng: 121.4346 },
    { name: '虹桥火车站', address: '上海市闵行区申虹路', lat: 31.1936, lng: 121.3200 },
    { name: '虹桥站', address: '上海市闵行区申虹路', lat: 31.1936, lng: 121.3200 },
    
    // 上海地区 - 浦东新区热门地点
    { name: '陆家嘴', address: '上海市浦东新区陆家嘴金融贸易区', lat: 31.2415, lng: 121.5056 },
    { name: '陆家嘴金融中心', address: '上海市浦东新区陆家嘴金融贸易区', lat: 31.2415, lng: 121.5056 },
    { name: '东方明珠', address: '上海市浦东新区世纪大道1号', lat: 31.2400, lng: 121.4997 },
    { name: '东方明珠塔', address: '上海市浦东新区世纪大道1号', lat: 31.2400, lng: 121.4997 },
    { name: '上海中心大厦', address: '上海市浦东新区银城中路501号', lat: 31.2333, lng: 121.5052 },
    { name: '金茂大厦', address: '上海市浦东新区世纪大道88号', lat: 31.2356, lng: 121.5054 },
    { name: '环球金融中心', address: '上海市浦东新区世纪大道100号', lat: 31.2342, lng: 121.5058 },
    { name: '世纪公园', address: '上海市浦东新区锦绣路1001号', lat: 31.2194, lng: 121.5500 },
    { name: '浦东嘉里城', address: '上海市浦东新区花木路1378号', lat: 31.2078, lng: 121.5502 },
    { name: '第一八佰伴', address: '上海市浦东新区张杨路501号', lat: 31.2276, lng: 121.5207 },
    
    // 上海地区 - 市中心景点
    { name: '外滩', address: '上海市黄浦区中山东一路', lat: 31.2401, lng: 121.4900 },
    { name: '南京路步行街', address: '上海市黄浦区南京东路', lat: 31.2385, lng: 121.4737 },
    { name: '南京路', address: '上海市黄浦区南京东路', lat: 31.2385, lng: 121.4737 },
    { name: '人民广场', address: '上海市黄浦区人民大道', lat: 31.2283, lng: 121.4759 },
    { name: '人民公园', address: '上海市黄浦区人民大道', lat: 31.2283, lng: 121.4759 },
    { name: '田子坊', address: '上海市黄浦区泰康路210弄', lat: 31.2103, lng: 121.4737 },
    { name: '新天地', address: '上海市黄浦区太仓路181弄', lat: 31.2176, lng: 121.4763 },
    { name: '城隍庙', address: '上海市黄浦区豫园老街', lat: 31.2247, lng: 121.4920 },
    { name: '豫园', address: '上海市黄浦区安仁街132号', lat: 31.2247, lng: 121.4920 },
    
    // 上海地区 - 商圈
    { name: '淮海路', address: '上海市黄浦区淮海中路', lat: 31.2210, lng: 121.4737 },
    { name: '徐家汇', address: '上海市徐汇区徐家汇商圈', lat: 31.1886, lng: 121.4363 },
    { name: '五角场', address: '上海市杨浦区五角场商圈', lat: 31.2990, lng: 121.5130 },
    { name: '中山公园', address: '上海市长宁区中山公园', lat: 31.2223, lng: 121.4209 },
    
    // 北京地区
    { name: '北京首都国际机场', address: '北京市朝阳区首都机场', lat: 40.0799, lng: 116.6031 },
    { name: '北京大兴国际机场', address: '北京市大兴区榆垡镇', lat: 39.5098, lng: 116.4100 },
    { name: '天安门广场', address: '北京市东城区东长安街', lat: 39.9042, lng: 116.3974 },
    { name: '北京西站', address: '北京市丰台区莲花池东路', lat: 39.8951, lng: 116.3225 },
    { name: '北京南站', address: '北京市丰台区永外大街', lat: 39.8656, lng: 116.3789 },
    
    // 深圳地区
    { name: '深圳宝安国际机场', address: '深圳市宝安区航城大道', lat: 22.6390, lng: 113.8111 },
    { name: '深圳北站', address: '深圳市龙华区致远中路', lat: 22.6097, lng: 114.0298 },
    
    // 广州地区
    { name: '广州白云国际机场', address: '广州市白云区机场路', lat: 23.3924, lng: 113.2988 },
    { name: '广州南站', address: '广州市番禺区南站北路', lat: 22.9888, lng: 113.2723 }
  ]
  
  const queryLower = query.toLowerCase()
  
  // 多种匹配策略
  const matches = popularPlaces.filter(place => {
    const nameLower = place.name.toLowerCase()
    const addressLower = place.address.toLowerCase()
    
    // 1. 精确匹配
    if (nameLower === queryLower) return true
    
    // 2. 名称包含查询
    if (nameLower.includes(queryLower)) return true
    
    // 3. 地址包含查询
    if (addressLower.includes(queryLower)) return true
    
    // 4. 去掉常见词汇的匹配（如：机场、火车站、地铁站等）
    const cleanQuery = queryLower
      .replace(/机场|火车站|高铁站|地铁站|汽车站|客运站/g, '')
      .trim()
    
    if (cleanQuery && nameLower.includes(cleanQuery)) return true
    
    // 5. 拼音首字母匹配（简化版）
    const queryChars = queryLower.split('')
    let nameIndex = 0
    for (const char of queryChars) {
      const found = nameLower.indexOf(char, nameIndex)
      if (found === -1) return false
      nameIndex = found + 1
    }
    
    return true
  })
  
  return matches
    .map(place => ({
      value: place.address,
      name: place.name,
      address: place.address,
      lat: place.lat,
      lng: place.lng
    }))
    .slice(0, 8) // 增加返回结果数量
}

// 处理起点地址选择
const handleStartLocationSelect = async (item: any) => {
  startLocation.value = item.name || item.address
  currentLocation.value = {
    lng: item.lng,
    lat: item.lat,
    address: item.address
  }
  
  // 更新地图中心位置
  if (mapRef.value) {
    mapRef.value.setCenter(item.lat, item.lng)
    // 清除之前的单独标记，因为路线规划会添加起点终点标记
    mapRef.value.clearRoute()
  }
  
  ElMessage.success(`已选择：${item.name}`)
  
  // 保存历史记录（仅登录用户）
  await saveLocationToHistory(item)
  
  // 检查是否可以规划路线
  setTimeout(() => {
    checkAndPlanRoute()
  }, 300)
}

// 处理起点地址输入
const handleStartLocationInput = (value: string) => {
  if (!value) {
    // 如果输入为空，可以重新获取当前位置
    if (mapRef.value) {
      mapRef.value.getCurrentLocation()
    }
  }
}

// 处理起点地址回车键
const handleStartLocationEnter = async () => {
  const query = startLocation.value?.trim()
  if (!query || query.length < 2) {
    ElMessage.warning('请输入至少2个字符的地址')
    return
  }

  ElMessage.info('正在搜索并定位...')
  
  // 搜索地址（优先使用本地数据库）
  const searchResults = getPopularLocations(query)
  
  if (searchResults.length > 0) {
    // 自动选择第一个结果
    const firstResult = searchResults[0]
    handleStartLocationSelect(firstResult)
  } else {
    // 如果本地没有匹配结果，尝试网络搜索
    try {
      const networkResults = await searchLocationByQuery(query)
      if (networkResults.length > 0) {
        const firstResult = networkResults[0]
        handleStartLocationSelect(firstResult)
      } else {
        ElMessage.error('未找到匹配的地址，请尝试输入常见地点名称（如：浦东机场、上海火车站）')
      }
    } catch (error) {
      ElMessage.error('搜索失败，请尝试输入常见地点名称（如：浦东机场、上海火车站）')
    }
  }
}

// 搜索终点地址（复用起点搜索逻辑）
const searchEndLocation = async (queryString: string, callback: (suggestions: any[]) => void) => {
  searchStartLocation(queryString, callback)
}

// 处理终点地址选择
const handleEndLocationSelect = async (item: any) => {
  endLocation.value = item.name || item.address
  endLocationData.value = {
    lng: item.lng,
    lat: item.lat,
    address: item.address
  }
  
  ElMessage.success(`已选择目的地：${item.name}`)
  
  // 保存历史记录（仅登录用户）
  await saveLocationToHistory(item)
  
  // 检查是否可以规划路线
  setTimeout(() => {
    checkAndPlanRoute()
  }, 300)
}

// 处理终点地址输入
const handleEndLocationInput = (value: string) => {
  if (!value) {
    endLocationData.value = null
  }
}

// 处理终点地址回车键
const handleEndLocationEnter = async () => {
  const query = endLocation.value?.trim()
  if (!query || query.length < 2) {
    ElMessage.warning('请输入至少2个字符的地址')
    return
  }

  ElMessage.info('正在搜索并定位...')
  
  // 搜索地址（优先使用本地数据库）
  const searchResults = getPopularLocations(query)
  
  if (searchResults.length > 0) {
    // 自动选择第一个结果
    const firstResult = searchResults[0]
    handleEndLocationSelect(firstResult)
  } else {
    // 如果本地没有匹配结果，尝试网络搜索
    try {
      const networkResults = await searchLocationByQuery(query)
      if (networkResults.length > 0) {
        const firstResult = networkResults[0]
        handleEndLocationSelect(firstResult)
      } else {
        ElMessage.error('未找到匹配的地址，请尝试输入常见地点名称（如：浦东机场、上海火车站）')
      }
    } catch (error) {
      ElMessage.error('搜索失败，请尝试输入常见地点名称（如：浦东机场、上海火车站）')
    }
  }
}

// 通用的地址搜索函数
const searchLocationByQuery = async (query: string): Promise<any[]> => {
  // 优先使用本地热门地址库，避免网络请求问题
  const localResults = getPopularLocations(query)
  
  if (localResults.length > 0) {
    return localResults
  }
  
  // 如果本地没有匹配结果，尝试网络搜索
  try {
    const response = await fetch(
      `/api/search?q=${encodeURIComponent(query)}&format=json&limit=5&accept-language=zh-CN&countrycodes=CN`
    )
    
    if (response.ok) {
      const data = await response.json()
      return data.map((item: any) => ({
        value: item.display_name,
        name: item.name || item.display_name.split(',')[0],
        address: item.display_name,
        lat: parseFloat(item.lat),
        lng: parseFloat(item.lon)
      }))
    }
  } catch (error) {
    console.warn('网络地址搜索失败，使用本地数据库:', error)
  }
  
  // 如果网络搜索也失败，返回空数组
  return []
}

// 初始化
onMounted(() => {
  console.log('Home组件已加载，Leaflet地图组件将自动初始化')
  
  // 检查是否从地址输入页面返回
  if (route.query.planRoute === 'true') {
    console.log('✅ 从地址选择页面返回，准备规划路线')
    
    // 禁用自动定位，避免重新定位
    shouldAutoLocate.value = false
    console.log('🔴 已禁用自动定位，shouldAutoLocate =', shouldAutoLocate.value)
    
    // 从地址输入页面返回，需要规划路线
    if (route.query.startLocation) {
      startLocation.value = route.query.startLocation as string
    }
    if (route.query.endLocation) {
      endLocation.value = route.query.endLocation as string
    }
    
    // 获取坐标信息
    const startLat = route.query.startLat ? parseFloat(route.query.startLat as string) : null
    const startLng = route.query.startLng ? parseFloat(route.query.startLng as string) : null
    const endLat = route.query.endLat ? parseFloat(route.query.endLat as string) : null
    const endLng = route.query.endLng ? parseFloat(route.query.endLng as string) : null
    
    console.log('🚗 路线信息:', { startLat, startLng, endLat, endLng })
    
    // 更新当前位置信息
    if (startLat && startLng) {
      currentLocation.value = {
        lat: startLat,
        lng: startLng,
        address: startLocation.value
      }
    }
    
    if (endLat && endLng) {
      endLocationData.value = {
        lat: endLat,
        lng: endLng,
        address: endLocation.value
      }
    }

    // 获取路线数据并规划路线
    if (startLat && startLng && endLat && endLng && mapRef.value) {
      // 规划路线
      mapRef.value.planRoute(
        { lat: startLat, lng: startLng },
        { lat: endLat, lng: endLng },
        startLocation.value,
        endLocation.value
      )
    }
    
    // 清除URL参数
    router.replace({ name: 'Home' })
  } else {
    // 正常加载时启用自动定位
    shouldAutoLocate.value = true
    console.log('🟢 正常加载，启用自动定位，shouldAutoLocate =', shouldAutoLocate.value)
  }
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: #f8f8f8;
}

.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: white;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.city-select {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #333;
  font-weight: 500;
  cursor: pointer;
}

.arrow {
  font-size: 12px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #ff7e00;
  cursor: pointer;
}

.main-content {
  padding: 0 20px 20px;
}

.map-section {
  height: 200px;
  margin: 20px 0;
  border-radius: 12px;
  overflow: hidden;
}

.ride-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}

.ride-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.ride-header h2 {
  margin: 0;
  color: #333;
  font-size: 20px;
}

.weather {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #666;
  font-size: 14px;
}

.location-inputs {
  position: relative;
  margin-bottom: 16px;
}

.location-item {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.location-search-input {
  flex: 1;
}

.location-search-input-wrapper {
  flex: 1;
  cursor: pointer;
}

.search-suggestion {
  padding: 8px 0;
}

.suggestion-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.suggestion-address {
  font-size: 12px;
  color: #666;
  line-height: 1.4;
}

.location-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 12px;
  flex-shrink: 0;
}

.location-dot.start {
  background: #67c23a;
}

.location-dot.end {
  background: #ff7e00;
}

.swap-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 32px;
  height: 32px;
  background: #f5f5f5;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
}

.swap-btn:hover {
  background: #e5e5e5;
}

.route-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f7fa 100%);
  border: 1px solid #81d4fa;
  border-radius: 8px;
  padding: 8px 12px;
  margin-bottom: 16px;
  position: relative;
}

.route-info-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #0277bd;
  font-size: 12px;
}

.route-icon {
  font-size: 14px;
  color: #0288d1;
}

.route-text {
  font-weight: 500;
}

.route-clear-btn {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
  color: #666;
  font-size: 12px;
}

.route-clear-btn:hover {
  background: #ff5722;
  color: white;
  transform: scale(1.1);
}

.route-info.route-error {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  border-color: #f59e0b;
}

.route-info.route-error .route-icon,
.route-info.route-error .route-text {
  color: #d97706;
}

.route-error-tip {
  margin-top: 8px;
  padding: 4px 8px;
  background: rgba(217, 119, 6, 0.1);
  border-radius: 4px;
  width: 100%;
}

.route-error-tip p {
  margin: 0;
  font-size: 10px;
  color: #92400e;
  text-align: center;
}

.ride-options {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  overflow-x: auto;
}

.ride-option {
  flex: 1;
  min-width: 80px;
  padding: 12px 8px;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.ride-option.active {
  border-color: #ff7e00;
  background: #fff7f0;
}

.option-icon {
  margin-bottom: 6px;
  color: #ff7e00;
}

.option-name {
  font-size: 14px;
  color: #333;
  margin-bottom: 2px;
}

.option-price {
  font-size: 12px;
  color: #666;
}

.call-car-btn {
  width: 100%;
  height: 48px;
  background: linear-gradient(135deg, #ff7e00 0%, #ff6600 100%);
  border: none;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
}

.feature-menu {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  background: white;
  padding: 20px;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.feature-item:hover {
  background: #f5f5f5;
  border-radius: 8px;
}

.feature-icon {
  width: 40px;
  height: 40px;
  background: #f0f0f0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
  color: #ff7e00;
}

.feature-item span {
  font-size: 12px;
  color: #666;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-content {
    padding: 0 16px 16px;
  }
  
  .ride-card {
    padding: 16px;
  }
  
  .feature-menu {
    padding: 16px;
  }
}
</style> 