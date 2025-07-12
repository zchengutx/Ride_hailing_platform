import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(process.cwd(), 'src'),
    },
  },
  server: {
    port: 8080, // 恢复默认端口
    host: '0.0.0.0',
    open: true,
    proxy: {
      // 通用API代理 - 处理所有 /api 请求
      '/api': {
        target: 'http://localhost:8888',
        changeOrigin: true,
        secure: false,
        ws: false,
        configure: (proxy, options) => {
          // 处理预检请求
          proxy.on('proxyReq', (proxyReq, req, res) => {
            // 详细日志，便于调试
            console.log('🔄 代理请求:', req.method, req.url, 'to', 'http://localhost:8888' + req.url)
            console.log('📤 请求头:', req.headers)
            
            // 确保正确设置请求头
            proxyReq.setHeader('Host', 'localhost:8888')
            // 动态设置Origin，支持不同端口
            const actualPort = req.headers.host?.split(':')[1] || '8080'
            proxyReq.setHeader('Origin', `http://localhost:${actualPort}`)
            proxyReq.setHeader('Referer', `http://localhost:${actualPort}/`)
            
            // 处理 OPTIONS 预检请求
            if (req.method === 'OPTIONS') {
              res.setHeader('Access-Control-Allow-Origin', '*')
              res.setHeader('Access-Control-Allow-Methods', 'GET,PUT,POST,DELETE,OPTIONS')
              res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization, Content-Length, X-Requested-With')
              res.setHeader('Access-Control-Max-Age', '86400')
              res.writeHead(200)
              res.end()
              return
            }
          })
          
          proxy.on('proxyRes', (proxyRes, req, res) => {
            // 详细日志，便于调试
            console.log('📥 代理响应:', req.url, '状态码:', proxyRes.statusCode)
            console.log('📥 响应头:', proxyRes.headers)
            
            // 为所有响应添加CORS头
            res.setHeader('Access-Control-Allow-Origin', '*')
            res.setHeader('Access-Control-Allow-Methods', 'GET,PUT,POST,DELETE,OPTIONS')
            res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization, Content-Length, X-Requested-With')
            res.setHeader('Access-Control-Allow-Credentials', 'true')
          })
          
          proxy.on('error', (err, req, res) => {
            // 只记录严重错误
            console.warn('❌ 代理错误:', req.url, err.message)
            if (!res.headersSent) {
              res.writeHead(500, { 'content-type': 'application/json' })
              res.end(JSON.stringify({ error: 'Proxy error: ' + err.message }))
            }
          })
        }
      },
      // 地理编码相关代理
      '/geocoding': {
        target: 'https://nominatim.openstreetmap.org',
        changeOrigin: true,
        secure: true,
        rewrite: (path) => {
          const url = new URL(path, 'http://localhost')
          const lat = url.searchParams.get('lat')
          const lon = url.searchParams.get('lon')
          return `/reverse?format=json&lat=${lat}&lon=${lon}&accept-language=zh-CN`
        }
      },
      '/search': {
        target: 'https://nominatim.openstreetmap.org',
        changeOrigin: true,
        secure: true,
        rewrite: (path) => {
          const url = new URL(path, 'http://localhost')
          const q = url.searchParams.get('q')
          const format = url.searchParams.get('format')
          const limit = url.searchParams.get('limit')
          const countrycodes = url.searchParams.get('countrycodes')
          return `/search?q=${q}&format=${format}&limit=${limit}&countrycodes=${countrycodes}&accept-language=zh-CN`
        }
      }
    }
  },
}) 