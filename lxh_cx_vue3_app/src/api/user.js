import apiClient from './config'

// 用户相关接口
export const userApi = {
  // 发送短信验证码
  sendSms: (data) => {
    return apiClient.post('/user/sendSms', data)
  },

  // 用户注册
  register: (data) => {
    return apiClient.post('/user/register', data)
  },

  // 用户登录
  login: (data) => {
    return apiClient.post('/user/login', data)
  },

  // 获取用户信息列表
  getUserInfoList: () => {
    return apiClient.get('/user/userInfoList')
  },

  // 绑定手机号
  bindMobile: (data) => {
    return apiClient.post('/user/bindMobile', data)
  },

  // 刷新token
  refreshToken: () => {
    return apiClient.get('/user/refresh_token')
  },

  // 微信登录相关
  wechat: {
    // 获取微信登录二维码
    getQRCode: () => {
      return apiClient.get('/user/weChat')
    },
    // 检查签名
    checkSignature: (params) => {
      return apiClient.get('/user/checkSignature', { params })
    },
    // 回调
    callback: (params) => {
      return apiClient.get('/user/callback', { params })
    }
  }
}

export default userApi 