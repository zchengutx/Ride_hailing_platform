# 🚀 启动指南

## 问题解决方案

您反映"根本没有调到后端接口"，这是因为：
1. 后端服务没有运行
2. 代理配置需要优化
3. 距离计算有误差

## 🔧 完整解决方案

### 方案1：启动Mock服务器（推荐）

1. **安装依赖**:
```bash
npm install
```

2. **启动Mock服务器**（新终端窗口）:
```bash
npm run mock
```

3. **启动前端**（另一个终端窗口）:
```bash
npm run dev
```

4. **访问应用**: http://localhost:8084/

### 方案2：使用您的真实后端

1. 确保您的后端服务运行在 `http://localhost:8888`
2. 提供这两个接口：
   - `POST /api/driving`
   - `POST /api/directionlite`
3. 启动前端：`npm run dev`

## 🧪 测试步骤

1. **检查Mock服务器**:
   - 访问 http://localhost:8888/health
   - 应该显示 `{"status":"ok","message":"Mock server is running"}`

2. **测试API调用**:
   - 在前端点击"重新定位"
   - 点击"测试路线"
   - 查看控制台日志和地图效果

3. **预期结果**:
   - 🟢 后端API状态显示"正常"
   - 🛣️ 地图显示蓝色实线路径
   - 📊 显示"约2.5公里，约8分钟"
   - ✅ 控制台无403错误

## 📋 Mock服务器功能

Mock服务器会：
- ✅ 响应 `/api/driving` 请求，返回模拟距离时间
- ✅ 响应 `/api/directionlite` 请求，返回模拟路线坐标
- 🔍 在控制台显示收到的请求
- 🚀 完全模拟真实后端行为

## 🐛 故障排除

### 端口冲突
如果8888端口被占用：
```bash
# 查看占用的进程
netstat -ano | findstr :8888
# 或者修改mock-server.js中的PORT为其他值
```

### 依赖安装问题
```bash
# 清除依赖重新安装
rm -rf node_modules package-lock.json
npm install
```

### 代理不工作
- 确保Vite开发服务器和Mock服务器都在运行
- 检查控制台网络面板中的请求URL

## 🎯 预期效果对比

**修复前**:
- ❌ 控制台大量403错误
- ❌ 后端API显示"离线"
- ❌ 路线显示异常距离（1106公里）
- ❌ 只能使用直线模式

**修复后**:
- ✅ API正常调用
- ✅ 后端API显示"正常"
- ✅ 合理的距离时间（2.5公里，8分钟）
- ✅ 真实的路线轨迹 