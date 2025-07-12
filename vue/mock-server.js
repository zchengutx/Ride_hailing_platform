import express from 'express'
import cors from 'cors'
const app = express()

// 中间件
app.use(cors())
app.use(express.json())

// 模拟 /api/driving 接口
app.post('/api/driving', (req, res) => {
  console.log('收到driving请求:', req.body)
  
  const { origins, destinations } = req.body
  
  // 模拟响应数据
  res.json({
    code: 200,
    message: 'success',
    data: {
      distance: '约2.5公里',
      duration: '约8分钟'
    }
  })
})

// 模拟 /api/directionlite 接口
app.post('/api/directionlite', (req, res) => {
  console.log('收到directionlite请求:', req.body)
  
  const { origins, destinations } = req.body
  
  // 解析起点和终点坐标
  const [startLng, startLat] = origins.split(',').map(Number)
  const [endLng, endLat] = destinations.split(',').map(Number)
  
  // 生成模拟路线坐标（简单的直线加一些弯曲）
  const steps = 10
  const polyline = []
  
  for (let i = 0; i <= steps; i++) {
    const ratio = i / steps
    const lat = startLat + (endLat - startLat) * ratio + (Math.sin(ratio * Math.PI) * 0.0001)
    const lng = startLng + (endLng - startLng) * ratio + (Math.cos(ratio * Math.PI) * 0.0001)
    polyline.push([lng, lat])
  }
  
  res.json({
    code: 200,
    message: 'success',
    data: {
      polyline: polyline
    }
  })
})

// 健康检查
app.get('/health', (req, res) => {
  res.json({ status: 'ok', message: 'Mock server is running' })
})

const PORT = 8888
app.listen(PORT, () => {
  console.log(`🚀 Mock服务器启动成功!`)
  console.log(`📍 地址: http://localhost:${PORT}`)
  console.log(`🔗 健康检查: http://localhost:${PORT}/health`)
  console.log(`📋 支持的接口:`)
  console.log(`   POST /api/driving - 距离时间计算`)
  console.log(`   POST /api/directionlite - 详细路线规划`)
}) 