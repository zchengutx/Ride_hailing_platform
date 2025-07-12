<template>
  <div class="address-input-page">
    <!-- 顶部导航栏 -->
    <div class="header">
      <div class="header-left" @click="goBack">
        <el-icon class="back-icon"><ArrowLeft /></el-icon>
      </div>
      <div class="header-title">选择地点</div>
      <div class="header-right">
        <!-- 登录状态显示 -->
        <div v-if="isLoggedIn" class="user-status">
          <el-icon class="user-icon"><User /></el-icon>
          <span class="login-status">已登录</span>
          <span class="logout-btn" @click="logout">退出登录</span>
        </div>
        <div v-else class="login-hint" @click="showLoginDialog = true">
          <span class="login-status">未登录</span>
          <span>登录</span>
        </div>
      </div>
    </div>

    <!-- 路线输入区域 -->
    <div class="route-input-section">
      <div class="route-line">
        <div class="route-dots">
          <div class="route-dot start"></div>
          <div class="route-line-connector"></div>
          <div class="route-dot end"></div>
        </div>
        <div class="route-inputs">
          <!-- 起点输入 -->
          <div class="input-group start-input">
            <el-input
              v-model="startLocation"
              placeholder="您在哪儿？"
              class="location-input"
              :class="{ focused: currentFocus === 'start' }"
              readonly
              @click="focusInput('start')"
            >
              <template #suffix>
                <el-icon v-if="startLocation" @click="clearStart" class="clear-icon"><Close /></el-icon>
              </template>
            </el-input>
            <el-button 
              v-if="!startLocation || startLocation === '定位中...'"
              @click="getCurrentLocation" 
              :loading="locating"
              type="primary" 
              size="small" 
              class="locate-btn"
            >
              <el-icon><Location /></el-icon>
              {{ locating ? '定位中' : '定位' }}
            </el-button>
          </div>

          <!-- 终点输入 -->
          <div class="input-group end-input">
            <el-input
              ref="endInputRef"
              v-model="endLocation"
              placeholder="您要去哪儿？"
              class="location-input"
              :class="{ focused: currentFocus === 'end' }"
              @focus="focusInput('end')"
              @input="handleSearch"
              clearable
              @clear="clearSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </div>

      <!-- 交换按钮 -->
      <div class="swap-button" @click="swapLocations">
        <el-icon><Switch /></el-icon>
      </div>
    </div>

    <!-- 登录提示横幅 -->
    <div v-if="!isLoggedIn" class="login-banner" @click="showLoginDialog = true">
      <div class="banner-content">
        <el-icon class="banner-icon"><User /></el-icon>
        <span>登录后可保存常用地址，享受更便捷的服务</span>
        <el-icon class="banner-arrow"><ArrowRight /></el-icon>
      </div>
    </div>

    <!-- 搜索结果 -->
    <div v-if="showSearchResults" class="search-results">
      <div class="results-header">
        <span>搜索结果</span>
      </div>
      <div class="results-list">
        <div
          v-for="(item, index) in searchResults"
          :key="index"
          class="result-item"
          @click="selectLocation(item)"
        >
          <div class="result-icon">
            <el-icon><LocationInformation /></el-icon>
          </div>
          <div class="result-content">
            <div class="result-name">{{ item.name }}</div>
            <div class="result-address">{{ item.address }}</div>
          </div>
          <!-- 需要登录的标识 -->
          <div v-if="!isLoggedIn" class="login-required">
            <el-icon class="lock-icon"><Lock /></el-icon>
          </div>
        </div>
      </div>
    </div>

    <!-- 历史记录和推荐 -->
    <div v-else class="history-section">
      <!-- 选择提示 -->
      <div class="selection-tip">
        <span v-if="currentFocus === 'start'">请选择出发地点</span>
        <span v-else>请选择目的地点</span>
        <span v-if="!isLoggedIn" class="login-tip">（选择地址需要登录）</span>
      </div>
      
      <!-- 最近搜索 -->
      <div v-if="recentSearches.length > 0 && isLoggedIn" class="section">
        <div class="section-header">
          <span>最近搜索</span>
          <el-button text @click="clearHistory">清除</el-button>
        </div>
        <div class="item-list">
          <div
            v-for="(item, index) in recentSearches"
            :key="index"
            class="history-item"
            @click="selectLocation(item)"
          >
            <div class="item-icon">
              <el-icon><Clock /></el-icon>
            </div>
            <div class="item-content">
              <div class="item-name">{{ item.name }}</div>
              <div class="item-address">{{ item.address }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 热门地点 -->
      <div class="section">
        <div class="section-header">
          <span>热门地点</span>
        </div>
        <div class="item-list">
          <div
            v-for="(item, index) in hotLocations"
            :key="index"
            class="hot-item"
            @click="selectLocation(item)"
          >
            <div class="item-icon">
              <el-icon><Star /></el-icon>
            </div>
            <div class="item-content">
              <div class="item-name">{{ item.name }}</div>
              <div class="item-address">{{ item.address }}</div>
            </div>
            <div class="item-tag">{{ item.tag }}</div>
            <!-- 需要登录的标识 -->
            <div v-if="!isLoggedIn" class="login-required">
              <el-icon class="lock-icon"><Lock /></el-icon>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 确认按钮 -->
    <div v-if="startLocation && endLocation && startLocation !== '定位中...'" class="confirm-section">
      <el-button 
        type="primary" 
        size="large" 
        class="confirm-btn"
        @click="confirmRoute"
      >
        确认路线
      </el-button>
    </div>

    <!-- 登录对话框 -->
    <el-dialog 
      v-model="showLoginDialog" 
      title="登录使用网约车服务" 
      width="320px"
      :show-close="false"
      class="login-dialog"
    >
      <div class="login-content">
        <div class="login-tips">
          <p>登录后可以：</p>
          <ul>
            <li>保存常用地址</li>
            <li>查看出行记录</li>
            <li>享受会员优惠</li>
            <li>获得更好的服务体验</li>
          </ul>
        </div>
        <LoginForm @login-success="handleLoginSuccess" />
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showLoginDialog = false" size="large">
            稍后再说
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft,
  Location,
  Search,
  Switch,
  Close,
  LocationInformation,
  Clock,
  Star,
  User,
  Lock,
  ArrowRight
} from '@element-plus/icons-vue'
import LoginForm from '@/components/LoginForm.vue'
import { addAddrName, createHistoricalSearch, getDirectionLite } from '@/api/index'

interface LocationItem {
  name: string
  address: string
  lat: number
  lng: number
  tag?: string
  distance?: string
}

const router = useRouter()
const route = useRoute()

// 响应式数据
const startLocation = ref('')
const endLocation = ref('')
const currentFocus = ref<'start' | 'end' | null>('end') // 默认焦点在终点
const locating = ref(false)
const searchResults = ref<LocationItem[]>([])
const recentSearches = ref<LocationItem[]>([])
const searchTimeout = ref<NodeJS.Timeout>()
const endInputRef = ref()

// 登录相关状态
const isLoggedIn = computed(() => {
  return !!localStorage.getItem('token')
})
const showLoginDialog = ref(false)
const pendingLocationSelection = ref<LocationItem | null>(null)
const loginAndConfirming = ref(false)

// 保存选中的地址坐标信息
const startLocationData = ref<LocationItem | null>(null)
const endLocationData = ref<LocationItem | null>(null)

// 计算属性
const showSearchResults = computed(() => {
  return (currentFocus.value === 'end' && endLocation.value && searchResults.value.length > 0) ||
         (currentFocus.value === 'start' && searchResults.value.length > 0)
})

// 热门地点数据
const hotLocations = ref<LocationItem[]>([
  { name: '上海浦东国际机场', address: '上海市浦东新区机场大道900号', lat: 31.1434, lng: 121.8052, tag: '机场' },
  { name: '上海虹桥机场', address: '上海市闵行区虹桥路2550号', lat: 31.1979, lng: 121.3364, tag: '机场' },
  { name: '上海火车站', address: '上海市静安区秣陵路100号', lat: 31.2496, lng: 121.4553, tag: '火车站' },
  { name: '陆家嘴金融中心', address: '上海市浦东新区陆家嘴环路1000号', lat: 31.2415, lng: 121.5056, tag: '商圈' },
  { name: '外滩', address: '上海市黄浦区中山东一路', lat: 31.2401, lng: 121.4900, tag: '景点' },
  { name: '南京路步行街', address: '上海市黄浦区南京东路', lat: 31.2385, lng: 121.4737, tag: '商圈' }
])

// 初始化
onMounted(() => {
  if (!isLoggedIn.value) {
    router.replace({ name: 'Login', query: { redirect: router.currentRoute.value.fullPath } })
    return
  }
  
  loadRecentSearches()
  
  // 从路由参数获取初始值
  if (route.query.startLocation) {
    startLocation.value = route.query.startLocation as string
  }
  if (route.query.endLocation) {
    endLocation.value = route.query.endLocation as string
  }
  
  // 如果有焦点参数，设置焦点
  if (route.query.focus) {
    currentFocus.value = route.query.focus as 'start' | 'end'
  } else {
    // 默认焦点设置为终点输入框
    currentFocus.value = 'end'
  }
  
  // 初始化已有地址的坐标数据
  if (startLocation.value && !startLocationData.value) {
    startLocationData.value = findLocationData(startLocation.value)
  }
  if (endLocation.value && !endLocationData.value) {
    endLocationData.value = findLocationData(endLocation.value)
  }
  
  if (!startLocation.value) {
    getCurrentLocation()
  }
  
  // 如果是编辑终点且终点输入框存在，自动聚焦
  if (currentFocus.value === 'end' && endInputRef.value) {
    setTimeout(() => {
      endInputRef.value?.focus()
    }, 300)
  }
})

// 检查登录状态
const checkLoginStatus = () => {
  // 现在 isLoggedIn 是 computed，无需手动赋值
  const token = localStorage.getItem('token')
  if (token) {
    console.log('用户已登录')
  }
}

// 监听搜索输入
watch(endLocation, (newValue) => {
  if (newValue && currentFocus.value === 'end') {
    if (searchTimeout.value) {
      clearTimeout(searchTimeout.value)
    }
    
    searchTimeout.value = setTimeout(() => {
      performSearch(newValue)
    }, 300)
  } else {
    searchResults.value = []
  }
})

// 返回上一页
const goBack = () => {
  router.back()
}

// 设置输入焦点
const focusInput = (type: 'start' | 'end') => {
  if (type === 'end' && !isLoggedIn.value) {
    router.push({ name: 'Login', query: { redirect: router.currentRoute.value.fullPath } })
    return
  }
  currentFocus.value = type
  if (type === 'end' && endLocation.value) {
    performSearch(endLocation.value)
  } else if (type === 'start') {
    // 为起点显示热门地点
    performSearch('')
  }
}

// 获取当前位置
const getCurrentLocation = () => {
  if (!navigator.geolocation) {
    ElMessage.error('您的浏览器不支持定位功能')
    return
  }
  
  locating.value = true
  startLocation.value = '定位中...'
  
  navigator.geolocation.getCurrentPosition(
    (position) => {
      const { latitude, longitude } = position.coords
      
      // 反向地理编码获取地址
      reverseGeocode(latitude, longitude)
        .then((address) => {
          startLocation.value = address
          // 保存定位的坐标信息
          startLocationData.value = {
            name: address,
            address: address,
            lat: latitude,
            lng: longitude
          }
          locating.value = false
        })
        .catch(() => {
          const fallbackAddress = '当前位置'
          startLocation.value = fallbackAddress
          // 保存定位的坐标信息
          startLocationData.value = {
            name: fallbackAddress,
            address: fallbackAddress,
            lat: latitude,
            lng: longitude
          }
          locating.value = false
        })
    },
    () => {
      locating.value = false
      startLocation.value = ''
      ElMessage.error('定位失败，请手动输入起点')
    },
    { enableHighAccuracy: true, timeout: 10000 }
  )
}

// 反向地理编码
const reverseGeocode = async (lat: number, lng: number): Promise<string> => {
  try {
    // 生成智能地址（基于坐标范围）
    const smartAddress = generateSmartAddress(lat, lng)
    return smartAddress
  } catch (error) {
    console.error('地理编码失败:', error)
    return '当前位置'
  }
}

// 生成智能地址
const generateSmartAddress = (lat: number, lng: number): string => {
  // 根据经纬度范围判断大致位置
  if (lat >= 31.0 && lat <= 31.5 && lng >= 121.2 && lng <= 121.9) {
    // 上海市内
    if (lat >= 31.22 && lat <= 31.24 && lng >= 121.47 && lng <= 121.49) {
      return '上海市黄浦区南京东路步行街附近'
    }
    if (lat >= 31.20 && lat <= 31.22 && lng >= 121.45 && lng <= 121.47) {
      return '上海市黄浦区外滩附近'
    }
    if (lat >= 31.22 && lat <= 31.24 && lng >= 121.50 && lng <= 121.52) {
      return '上海市浦东新区陆家嘴附近'
    }
    if (lat >= 31.13 && lat <= 31.16 && lng >= 121.79 && lng <= 121.82) {
      return '上海市浦东新区浦东国际机场附近'
    }
    if (lat >= 31.19 && lat <= 31.21 && lng >= 121.33 && lng <= 121.35) {
      return '上海市闵行区虹桥机场附近'
    }
    if (lat >= 31.24 && lat <= 31.26 && lng >= 121.44 && lng <= 121.46) {
      return '上海市静安区上海火车站附近'
    }
    
    // 根据区域返回大致位置
    if (lng >= 121.5) return '上海市浦东新区'
    if (lng <= 121.4) return '上海市徐汇区'
    return '上海市黄浦区'
  } else if (lat >= 39.4 && lat <= 40.3 && lng >= 115.7 && lng <= 117.4) {
    // 北京市内
    if (lat >= 39.90 && lat <= 39.92 && lng >= 116.38 && lng <= 116.42) {
      return '北京市东城区天安门附近'
    }
    if (lat >= 39.95 && lat <= 40.10 && lng >= 116.58 && lng <= 116.65) {
      return '北京市朝阳区首都国际机场附近'
    }
    return '北京市朝阳区'
  } else if (lat >= 22.4 && lat <= 22.8 && lng >= 113.8 && lng <= 114.6) {
    // 深圳市内
    return '深圳市福田区'
  } else if (lat >= 23.0 && lat <= 23.6 && lng >= 113.1 && lng <= 113.5) {
    // 广州市内
    return '广州市天河区'
  } else {
    // 其他地区
    return '当前位置附近'
  }
}

// 执行搜索
const performSearch = async (query: string) => {
  // 如果没有搜索词，显示热门地点（用于起点选择）
  if (!query || query.length < 1) {
    if (currentFocus.value === 'start') {
      searchResults.value = hotLocations.value.slice(0, 10)
    } else {
      searchResults.value = []
    }
    return
  }

  const queryLower = query.toLowerCase()
  
  // 在热门地点中搜索
  const hotResults = hotLocations.value.filter(location => 
    location.name.toLowerCase().includes(queryLower) ||
    location.address.toLowerCase().includes(queryLower)
  )

  // 在最近搜索中搜索
  const recentResults = recentSearches.value.filter(location =>
    location.name.toLowerCase().includes(queryLower) ||
    location.address.toLowerCase().includes(queryLower)
  )

  // 合并搜索结果并去重
  const combined = [...hotResults, ...recentResults]
  const unique = combined.filter((item, index, self) =>
    index === self.findIndex(t => t.name === item.name)
  )

  // 如果没有匹配的结果，尝试通过地理编码搜索
  if (unique.length === 0 && query.length >= 2) {
    try {
      const geocodeResult = await geocodeAddress(query)
      if (geocodeResult) {
        unique.push(geocodeResult)
      }
    } catch (error) {
      console.log('地理编码搜索失败:', error)
    }
  }

  searchResults.value = unique.slice(0, 10)
}

// 地理编码：通过地址获取坐标
const geocodeAddress = async (address: string): Promise<LocationItem | null> => {
  try {
    const response = await fetch(
      `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(address)}&limit=1&accept-language=zh-CN`
    )
    const data = await response.json()
    
    if (data && data.length > 0) {
      const result = data[0]
      return {
        name: address,
        address: result.display_name || address,
        lat: parseFloat(result.lat),
        lng: parseFloat(result.lon)
      }
    }
    return null
  } catch (error) {
    console.error('地理编码失败:', error)
    return null
  }
}

// 选择位置
const selectLocation = (location: LocationItem) => {
  // 检查是否已登录
  if (!isLoggedIn.value) {
    // 如果未登录，保存用户想要选择的位置并显示登录对话框
    pendingLocationSelection.value = location
    showLoginDialog.value = true
    ElMessage.warning('请先登录后再选择地址')
    return
  }
  
  // 已登录，正常执行地址选择
  doSelectLocation(location)
}

// 实际执行地址选择的函数
const doSelectLocation = (location: LocationItem) => {
  // 如果当前正在编辑起点输入框，则设置起点
  if (currentFocus.value === 'start') {
    startLocation.value = location.name
    startLocationData.value = location
  } else {
    // 默认情况下设置为终点（目的地）
    endLocation.value = location.name
    endLocationData.value = location
    currentFocus.value = null // 清除焦点，避免继续编辑
  }
  
  saveToRecentSearches(location)
  searchResults.value = []
  
  ElMessage.success(`已选择：${location.name}`)
}

// 交换起点终点
const swapLocations = () => {
  // 交换地址文本
  const tempLocation = startLocation.value
  startLocation.value = endLocation.value
  endLocation.value = tempLocation
  
  // 交换坐标数据
  const tempData = startLocationData.value
  startLocationData.value = endLocationData.value
  endLocationData.value = tempData
  
  ElMessage.success('起点和终点已交换')
}

// 清除起点
const clearStart = () => {
  startLocation.value = ''
  startLocationData.value = null
}

// 清除搜索
const clearSearch = () => {
  endLocation.value = ''
  endLocationData.value = null
  searchResults.value = []
  currentFocus.value = null
}

// 处理搜索输入
const handleSearch = (value: string) => {
  endLocation.value = value
  // 如果手动输入，清空之前保存的坐标数据
  if (value !== endLocationData.value?.name) {
    endLocationData.value = null
  }
  performSearch(value)
}

// 确认路线
const confirmRoute = async () => {
  if (!startLocationData.value || !endLocationData.value) {
    ElMessage.warning('请选择完整的起点和终点')
    return
  }

  try {
    // 调用路线规划API
    const response = await getDirectionLite(
      { lat: startLocationData.value.lat, lng: startLocationData.value.lng },
      { lat: endLocationData.value.lat, lng: endLocationData.value.lng }
    )

    if (response && response.data) {
      // 返回到首页并传递路线信息
      router.push({
        name: 'Home',
        query: {
          planRoute: 'true',
          startLocation: startLocation.value,
          endLocation: endLocation.value,
          startLat: startLocationData.value.lat.toString(),
          startLng: startLocationData.value.lng.toString(),
          endLat: endLocationData.value.lat.toString(),
          endLng: endLocationData.value.lng.toString(),
          routeData: JSON.stringify(response.data)
        }
      })
    } else {
      ElMessage.error('路线规划失败，请重试')
    }
  } catch (error) {
    console.error('路线规划失败:', error)
    ElMessage.error('路线规划失败，请重试')
  }
}

// 查找地点的坐标数据
const findLocationData = (locationName: string): LocationItem | null => {
  // 在热门地点中查找
  const found = hotLocations.value.find(item => 
    item.name === locationName || item.address === locationName
  )
  
  if (found) return found
  
  // 在最近搜索中查找
  const recentFound = recentSearches.value.find(item => 
    item.name === locationName || item.address === locationName
  )
  
  if (recentFound) return recentFound
  
  // 如果没找到，返回默认坐标（上海市中心）
  return {
    name: locationName,
    address: locationName,
    lat: 31.2304,
    lng: 121.4737
  }
}

// 加载最近搜索
const loadRecentSearches = () => {
  // 只有登录用户才加载历史记录
  if (!isLoggedIn.value) {
    recentSearches.value = []
    return
  }
  
  const saved = localStorage.getItem('recentSearches')
  if (saved) {
    try {
      recentSearches.value = JSON.parse(saved)
    } catch (e) {
      recentSearches.value = []
    }
  }
}

// 提取区/县级地址（简单正则，上海市浦东新区/闵行区/静安区等）
function extractDistrict(address: string): string {
  const match = address.match(/(.*?(区|县))/)
  return match ? match[0] : address
}

// 保存到最近搜索
const saveToRecentSearches = async (location: LocationItem) => {
  console.log('🔍 saveToRecentSearches 被调用，地址:', location.name)
  console.log('🔍 当前登录状态:', isLoggedIn.value)
  console.log('🔍 localStorage token:', localStorage.getItem('token') ? '存在' : '不存在')
  
  // 检查登录状态
  if (!isLoggedIn.value) {
    console.log('❌ 用户未登录，跳过保存历史记录')
    return
  }

  try {
    // 提取区域信息，优先使用起点地址的区域信息
    let idAddr = extractDistrict(startLocation.value)
    if (!idAddr || idAddr.trim() === '') {
      idAddr = extractDistrict(location.address) || '未知区域'
    }

    // 构建完整的历史记录参数
    const historyData = {
      addrName: location.name || location.address || '未知目的地',
      idAddr: idAddr,
      lat: location.lat || 0,
      lng: location.lng || 0,
      searchTime: new Date().toISOString()
    }

    console.log('💾 准备保存历史搜索记录:', historyData)

    // 调用后端API保存历史记录
    const response = await createHistoricalSearch(historyData)
    
    if (response && response.code === 200) {
      console.log('✅ 历史记录保存成功:', response.message || '保存成功')
      
      // 更新本地历史记录显示
      recentSearches.value = recentSearches.value.filter(item => item.name !== location.name)
      recentSearches.value.unshift(location)
      
      // 限制历史记录数量
      if (recentSearches.value.length > 10) {
        recentSearches.value = recentSearches.value.slice(0, 10)
      }
      
      // 同步保存到localStorage作为备份
      localStorage.setItem('recentSearches', JSON.stringify(recentSearches.value))
      
      console.log('📱 本地历史记录已更新')
    } else {
      console.warn('⚠️ 历史记录保存失败:', response?.message || '未知错误')
      ElMessage.warning('历史记录保存失败，但不影响正常使用')
    }
  } catch (error: any) {
    console.error('❌ 保存历史记录时发生错误:', error)
    
    // 根据错误状态码提供不同提示
    if (error.response?.status === 401 || error.response?.status === 403) {
      ElMessage.error('登录已过期，请重新登录')
      // 清除过期token并跳转登录
      localStorage.removeItem('token')
      setTimeout(() => {
        router.push('/login')
      }, 1500)
    } else if (error.response?.status === 500) {
      ElMessage.warning('服务器繁忙，历史记录保存失败')
    } else {
      ElMessage.warning('网络异常，历史记录保存失败')
    }
    
    // 即使后端保存失败，也更新本地记录以不影响用户体验
    recentSearches.value = recentSearches.value.filter(item => item.name !== location.name)
    recentSearches.value.unshift(location)
    if (recentSearches.value.length > 10) {
      recentSearches.value = recentSearches.value.slice(0, 10)
    }
    localStorage.setItem('recentSearches', JSON.stringify(recentSearches.value))
  }
}

// 清除历史记录
const clearHistory = () => {
  recentSearches.value = []
  localStorage.removeItem('recentSearches')
  ElMessage.success('历史记录已清除')
}

// 登录成功后的处理，自动继续路线规划
const handleLoginSuccess = () => {
  showLoginDialog.value = false
  ElMessage.success('登录成功')
  // 登录后自动执行确认路线
  confirmRoute()
  // 如果有待选择的位置，则自动选择
  if (pendingLocationSelection.value) {
    setTimeout(() => {
      doSelectLocation(pendingLocationSelection.value!)
      pendingLocationSelection.value = null
    }, 500)
  }
  // 加载用户的历史记录
  loadRecentSearches()
}

// 退出登录
const logout = () => {
  localStorage.removeItem('token')
  location.reload()
}

watch(isLoggedIn, (val) => {
  if (!val) {
    router.replace({ name: 'Login', query: { redirect: router.currentRoute.value.fullPath } })
  }
})
</script>

<style scoped>
.address-input-page {
  min-height: 100vh;
  background-color: #f8f9fa;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: white;
  border-bottom: 1px solid #eee;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left {
  width: 40px;
  display: flex;
  justify-content: center;
  cursor: pointer;
}

.back-icon {
  font-size: 20px;
  color: #333;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.header-right {
  width: 40px;
}

.route-input-section {
  position: relative;
  background: white;
  padding: 20px 16px;
  margin-bottom: 8px;
}

.route-line {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.route-dots {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 12px;
}

.route-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.route-dot.start {
  background: #52c41a;
}

.route-dot.end {
  background: #ff4d4f;
}

.route-line-connector {
  width: 2px;
  height: 40px;
  background: #ddd;
  margin: 4px 0;
}

.route-inputs {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.input-group {
  position: relative;
}

.location-input {
  font-size: 16px;
}

.locate-btn {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 12px;
}

.swap-button {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  width: 36px;
  height: 36px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 16px;
  color: #666;
  transition: all 0.2s;
}

.swap-button:hover {
  background: #f5f5f5;
  border-color: #1890ff;
  color: #1890ff;
}

.search-results {
  background: white;
  flex: 1;
}

.results-header {
  padding: 12px 16px;
  border-bottom: 1px solid #eee;
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.result-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.result-item:hover {
  background: #f5f5f5;
}

.result-icon {
  width: 40px;
  height: 40px;
  background: #f0f8ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  color: #1890ff;
  font-size: 18px;
}

.result-content {
  flex: 1;
}

.result-name {
  font-size: 16px;
  color: #333;
  margin-bottom: 4px;
  font-weight: 500;
}

.result-address {
  font-size: 14px;
  color: #666;
}

.history-section {
  flex: 1;
  background: white;
}

.section {
  margin-bottom: 8px;
  background: white;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #eee;
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.history-item,
.hot-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.history-item:hover,
.hot-item:hover {
  background: #f5f5f5;
}

.item-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  font-size: 18px;
}

.history-item .item-icon {
  background: #f0f0f0;
  color: #999;
}

.hot-item .item-icon {
  background: #fff7e6;
  color: #faad14;
}

.item-content {
  flex: 1;
}

.item-name {
  font-size: 16px;
  color: #333;
  margin-bottom: 4px;
  font-weight: 500;
}

.item-address {
  font-size: 14px;
  color: #666;
}

.item-tag {
  background: #e6f7ff;
  color: #1890ff;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
}

.confirm-section {
  padding: 16px;
  background: white;
  border-top: 1px solid #eee;
}

.confirm-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
}

.selection-tip {
  padding: 12px 16px;
  background: #f0f8ff;
  border-left: 4px solid #1890ff;
  margin-bottom: 8px;
  font-size: 14px;
  color: #1890ff;
  font-weight: 500;
}

/* 登录相关样式 */
.user-status {
  width: 32px;
  height: 32px;
  background: #f0f8ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #1890ff;
}

.login-hint {
  cursor: pointer;
  color: #1890ff;
  font-size: 14px;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.login-hint:hover {
  background: #f0f8ff;
}

.login-banner {
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f7ff 100%);
  padding: 16px;
  margin-bottom: 8px;
  cursor: pointer;
  border-bottom: 1px solid #e6f7ff;
  transition: all 0.3s ease;
}

.login-banner:hover {
  background: linear-gradient(135deg, #e6f7ff 0%, #d6f4ff 100%);
}

.banner-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #1890ff;
  font-size: 14px;
  font-weight: 500;
}

.banner-icon {
  font-size: 20px;
  color: #1890ff;
}

.banner-arrow {
  font-size: 12px;
  color: #1890ff;
}

.selection-tip {
  padding: 12px 16px;
  font-size: 14px;
  color: #666;
  border-bottom: 1px solid #eee;
}

.login-tip {
  color: #ff4d4f;
  font-size: 12px;
  margin-left: 4px;
}

.login-required {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: #fff2f0;
  border-radius: 50%;
  margin-left: 8px;
}

.lock-icon {
  font-size: 12px;
  color: #ff4d4f;
}

/* 登录对话框样式 */
:deep(.login-dialog) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.login-dialog .el-dialog__header) {
  padding: 20px 20px 0 20px;
  background: linear-gradient(135deg, #1890ff 0%, #40a9ff 100%);
  color: white;
}

:deep(.login-dialog .el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 18px;
}

.login-content {
  padding: 20px 0;
}

.login-tips {
  margin-bottom: 20px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border-left: 4px solid #1890ff;
}

.login-tips p {
  margin: 0 0 12px 0;
  font-weight: 600;
  color: #333;
  font-size: 16px;
}

.login-tips ul {
  margin: 0;
  padding-left: 20px;
}

.login-tips li {
  margin-bottom: 8px;
  color: #666;
  font-size: 14px;
  line-height: 1.5;
}

.dialog-footer {
  display: flex;
  justify-content: center;
  padding: 0 20px 20px 20px;
}

.dialog-footer .el-button {
  width: 100%;
  border-radius: 8px;
  background: #f5f5f5;
  border: none;
  color: #999;
  font-weight: 500;
}

.dialog-footer .el-button:hover {
  background: #e6f7ff;
  color: #1890ff;
}

.logout-btn {
  margin-left: 8px;
  color: #ff4d4f;
  cursor: pointer;
  font-size: 14px;
}

.login-status {
  margin-left: 6px;
  color: #52c41a;
  font-size: 14px;
}

.login-hint .login-status {
  color: #ff4d4f;
}
</style> 