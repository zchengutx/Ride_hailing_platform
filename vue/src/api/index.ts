import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const request = axios.create({
  baseURL: '',  // 使用代理，不需要指定完整URL
  timeout: 5000, // 减少超时时间到5秒
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    console.log('📤 发送请求:', {
      url: config.url,
      method: config.method,
      data: config.data,
      headers: config.headers
    })
    
    // 添加token到请求头 - 优先从localStorage获取，备选从cookie获取
    let token = localStorage.getItem('token')
    
    // 如果localStorage没有token，尝试从cookie获取
    if (!token || token === 'undefined' || token === 'null') {
      const cookieToken = document.cookie
        .split('; ')
        .find(row => row.startsWith('x-token='))
        ?.split('=')[1]
      
      if (cookieToken) {
        token = cookieToken
        // 同步保存到localStorage
        localStorage.setItem('token', cookieToken)
        console.log('🍪 从cookie获取token并同步到localStorage')
      }
    }
    
    if (token && token !== 'undefined' && token !== 'null') {
      config.headers.Authorization = token // 只传 token，无 Bearer 前缀
      console.log('🔑 添加认证token:', token.substring(0, 20) + '...')
    } else {
      console.warn('⚠️ 未找到有效的token')
      // 删除无效的Authorization头
      delete config.headers.Authorization
    }
    return config
  },
  (error) => {
    console.error('📤 请求拦截器错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    console.log('📥 收到响应:', {
      url: response.config.url,
      status: response.status,
      statusText: response.statusText,
      data: response.data
    })
    const { data } = response
    
    // 只有当响应中确实包含新token时才更新
    if (data && (data.token || data.data) && data.token !== 'undefined') {
      const newToken = data.token || data.data
      if (newToken && typeof newToken === 'string' && newToken.length > 10) {
        localStorage.setItem('token', newToken)
        // 同步更新cookie
        document.cookie = `x-token=${newToken}; expires=${new Date(Date.now() + 7*24*60*60*1000).toUTCString()}; path=/`
        console.log('🔄 API响应更新token:', newToken.substring(0, 20) + '...')
      }
    }
    
    return data
  },
  (error) => {
    console.error('📥 响应错误:', {
      url: error.config?.url,
      status: error.response?.status,
      statusText: error.response?.statusText,
      data: error.response?.data,
      message: error.message
    })
    
    // 静默处理的错误接口（不显示错误提示）
    const silentUrls = [
      '/api/wechat/status/',
      '/api/geocoding'
    ]
    
    const shouldSilent = silentUrls.some(url => error.config?.url?.includes(url))
    
    if (shouldSilent) {
      // 完全静默处理，不在控制台记录普通的连接错误
      if (error.response?.status !== 403 && error.response?.status !== 404) {
        console.warn(`API请求失败 [${error.config?.url}]:`, error.response?.status, error.response?.statusText)
      }
      return Promise.reject(error)
    }
    
    // 对于路线规划API，静默处理错误，不显示用户提示（已有fallback机制）
    if (error.config?.url?.includes('/api/driving') || error.config?.url?.includes('/api/directionlite')) {
      // 只在控制台记录，不显示用户错误提示，因为地图组件会处理fallback
      console.warn(`路线API调用失败 [${error.config?.url}]:`, {
        status: error.response?.status,
        message: error.message,
        code: error.code
      })
      return Promise.reject(error)
    }
    
    ElMessage.error(error.response?.data?.message || '请求失败')
    return Promise.reject(error)
  }
)

// API接口定义
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 发送短信验证码
export const sendSms = (phone: string) => {
  return request.post('/api/sendSms', { 
    mobile: phone,
    source: 'login'  // 添加source参数，表示来源是登录
  })
}

// 登录
export const login = (data: { phone: string; code: string }) => {
  return request.post('/api/login', { mobile: data.phone, send_sms_code: data.code })
}

// 微信登录相关
export const getWechatAuth = () => {
  return request.get('/api/wechat')
}

export const wechatLogin = () => {
  return request.get('/api/wechat/Login')
}

export const wechatCallback = (code: string) => {
  return request.get('/api/callback', { params: { code } })
}

// 检查微信扫码登录状态
export const checkWechatLoginStatus = (qrKey: string) => {
  return request.get(`/api/wechat/status/${qrKey}`)
}

// 计算驾车距离和时间
export const calculateDriving = async (origins: string, destinations: string) => {
  console.log('🚗 调用 calculateDriving API:', { origins, destinations })
  try {
    const response = await request.post('/api/driving', {
      origins,
      destinations
    })
    console.log('✅ calculateDriving API 响应:', response)
    return response
  } catch (error: any) {
    console.error('❌ calculateDriving API 失败:', error)
    console.warn('> 使用 mock 距离/时间 数据以继续流程')
    // 根据经纬度简单返回固定距离/时间，用于开发环境
    return {
      data: {
        distance: '约 3.2 km',
        duration: '约 8 分钟'
      }
    }
  }
}

// 获取详细路线规划（包含路径坐标）
export const getDirectionLite = async (startPoint: { lat: number, lng: number }, endPoint: { lat: number, lng: number }) => {
  console.log('🗺️ 调用 getDirectionLite API:', { startPoint, endPoint })
  try {
    // 将坐标转换为字符串格式："纬度,经度"
    const origins = `${startPoint.lat},${startPoint.lng}`
    const destinations = `${endPoint.lat},${endPoint.lng}`
    
    const response = await request.post('/api/directionlite', {
      origins,
      destinations
    })
    console.log('✅ getDirectionLite API 响应:', response)
    return response
  } catch (error: any) {
    console.error('❌ getDirectionLite API 失败:', error)
    console.warn('> 使用 mock 路线坐标数据以继续流程')
    // 返回简单的直线路径作为fallback
    return {
      data: {
        coordinates: [
          [startPoint.lat, startPoint.lng],
          [endPoint.lat, endPoint.lng]
        ],
        isRealRoute: false
      }
    }
  }
}

// 新增：添加地址名称接口
export const addAddrName = (addrName: string, idAddr: string) => {
  return request.post('/api/AddAddrName', {
    addrName,
    idAddr
  })
}

// 根据经纬度获取城市信息
 export const getCityByLocation = async (latitude: number, longitude: number) => {
   console.log('📍 调用 getCityByLocation API (使用/driving接口):', { latitude, longitude })
   try {
     // 调用后端提供的/driving接口获取位置信息
     const response = await request.post('/api/driving', {
       origins: `${latitude},${longitude}`,
          destinations: `${latitude},${longitude}` // 修正经纬度顺序为纬度,经度
     })
     console.log('✅ getCityByLocation API 响应:', response)
     return response
   } catch (error: any) {
     console.error('❌ getCityByLocation API 失败:', error)
     console.error('错误详情:', {
       status: error.response?.status,
       statusText: error.response?.statusText,
       data: error.response?.data,
       message: error.message
     })
     // mock 返回，避免前端卡住
     return { data: { city: '上海' } }
   }
 }

 // 计算车型价格接口
export const computePrice = (Distance: string, VehicleType: number) => {
  return request.post('/api/ComputePrice', {
    Distance,
    VehicleType
  })
}

// 新增：创建历史搜索记录接口
export const createHistoricalSearch = (locationData: {
  addrName: string;
  idAddr: string;
  lat?: number;
  lng?: number;
  searchTime?: string;
}) => {
  console.log('🔍 保存历史搜索记录:', locationData)
  
  // 确保传递完整的参数，并验证数据类型
  const params = {
    addrName: String(locationData.addrName || '未知地址'),    // 改回驼峰命名
    idAddr: String(locationData.idAddr || '未知区域'),        // 改回驼峰命名
    lat: Number(locationData.lat) || 0,                      // 使用简单的lat/lng
    lng: Number(locationData.lng) || 0,                      // 使用简单的lat/lng
    searchTime: locationData.searchTime || new Date().toISOString()  // 改回驼峰命名
  }
  
  console.log('📤 最终发送参数:', params)
  
  return request.post('/api/CreateHistoricalSearch', params)
}

export const getUserInfo = async () => {
  try {
    const response = await axios.post('/api/infoUser');
    console.log('用户信息:', response.data);
    return response.data;
  } catch (error) {
    console.error('获取用户信息失败:', error);
    throw error;
  }
};

// 获取当前位置（通过IP）
export const getCurrentLocationByIp = async () => {
  console.log('🌍 通过IP获取当前位置')
  try {
    const response = await request.post('/api/driving')
    console.log('✅ 获取位置成功:', response)
    return response
  } catch (error: any) {
    console.error('❌ 获取位置失败:', error)
    // 返回默认坐标（上海市中心）
    return {
      data: {
        lat: 31.2304,
        lng: 121.4737,
        address: '上海市'
      }
    }
  }
}

export default request