# 后端路线规划API对接指南

## 概述

已成功对接后端路线规划API，包括距离时间计算和详细路线获取功能。

## 后端接口说明

### 1. 距离时间计算接口
- **接口**: `POST /api/driving`
- **功能**: 获取两点之间的驾车距离和预计时间
- **请求参数**:
  ```json
  {
    "origins": "116.404,39.915",      // 起点经纬度 "lng,lat"
    "destinations": "116.407,39.918"  // 终点经纬度 "lng,lat"
  }
  ```
- **响应格式**:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": {
      "distance": "约2.5公里",
      "duration": "约8分钟"
    }
  }
  ```

### 2. 详细路线规划接口
- **接口**: `POST /api/directionlite`
- **功能**: 获取详细的行驶路线坐标
- **请求参数**:
  ```json
  {
    "origins": "116.404,39.915",      // 起点经纬度 "lng,lat"
    "destinations": "116.407,39.918"  // 终点经纬度 "lng,lat"
  }
  ```
- **响应格式**:
  ```json
  {
    "code": 200,
    "message": "success",
    "data": {
      "polyline": "路线坐标数据"  // 可以是坐标数组或编码字符串
    }
  }
  ```

## 前端集成功能

### 1. 智能路线规划
- 自动调用两个后端API获取完整路线信息
- 支持真实路线和直线路线显示
- 优雅的错误处理和降级方案

### 2. 地图交互
- 🚗 起点标记（绿色边框）
- 🏁 终点标记（红色边框）
- 蓝色实线：真实驾车路线
- 蓝色虚线：直线距离（API失败时的备选）
- 红色虚线：错误时的应急显示

### 3. 用户界面
- **重新定位按钮**: 获取当前GPS位置
- **微调位置按钮**: 手动调整定位精度
- **测试路线按钮**: 快速测试路线规划API
- **清除路线按钮**: 清除当前显示的路线

## 使用方法

### 1. 基本路线规划
```typescript
// 在组件中调用
const mapRef = ref()

const planRouteExample = () => {
  const start = { lat: 39.915, lng: 116.404 }  // 天安门
  const end = { lat: 39.918, lng: 116.407 }    // 故宫
  
  mapRef.value?.planRoute(start, end, '天安门', '故宫')
}
```

### 2. 监听路线计算结果
```vue
<template>
  <LeafletMap 
    ref="mapRef"
    @routeCalculated="handleRouteResult"
  />
</template>

<script setup>
const handleRouteResult = (routeInfo) => {
  console.log('距离:', routeInfo.distance)
  console.log('时间:', routeInfo.duration)
  console.log('是否真实路线:', routeInfo.isRealRoute)
  console.log('路线坐标:', routeInfo.coordinates)
}
</script>
```

### 3. 测试功能
1. 点击"重新定位"获取当前位置
2. 点击"测试路线"自动规划到附近目的地
3. 查看控制台输出的API调用结果
4. 点击"清除路线"清除显示

## 支持的数据格式

### Polyline坐标格式
系统支持多种坐标数据格式：

1. **坐标数组** (推荐):
   ```json
   [
     [116.404, 39.915],  // [lng, lat]
     [116.405, 39.916],
     [116.406, 39.917]
   ]
   ```

2. **对象数组**:
   ```json
   [
     {"lat": 39.915, "lng": 116.404},
     {"latitude": 39.916, "longitude": 116.405}
   ]
   ```

3. **Google Polyline编码字符串**:
   ```json
   "u{~vFvyys@fS]"
   ```

## 错误处理

系统具有完善的错误处理机制：

1. **API超时**: 5秒超时，自动降级到直线路线
2. **网络错误**: 显示直线距离和预估时间
3. **数据格式错误**: 自动解析多种坐标格式
4. **坐标无效**: 验证并提示用户

## 性能优化

1. **并行API调用**: 同时请求距离和路线数据
2. **智能缓存**: 避免重复计算相同路线
3. **按需加载**: 只在需要时调用API
4. **降级方案**: API失败时提供备选方案

## 注意事项

1. 确保后端API正常运行
2. 检查网络连接状态
3. 坐标格式必须为 "经度,纬度"
4. 建议添加适当的API调用频率限制
5. 大批量路线规划时考虑添加队列机制

## 调试技巧

1. 打开浏览器控制台查看API调用日志
2. 使用"测试路线"功能验证API连通性
3. 检查Network面板中的API请求响应
4. 确认后端接口返回的数据格式正确 