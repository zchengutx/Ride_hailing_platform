<template>
  <div class="home">
    <div class="map-container" ref="mapContainer"></div>
    <div class="search-box">
      <!-- 搜索框和结果列表 -->
    </div>
    <div class="order-box">
      <div class="location-info">
        <div class="start-point">
          <span>起点：</span>
          <span>{{ startAddress }}</span>
        </div>
        <div class="end-point">
          <span>终点：</span>
          <span>{{ endAddress }}</span>
        </div>
      </div>
      <button class="order-btn" @click="handleOrder">立即打车</button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { getPathPlanning } from '@/api/order'

export default {
  name: 'Home',
  setup() {
    const mapContainer = ref(null)
    const map = ref(null)
    const startPoint = ref(null)
    const endPoint = ref(null)
    const startAddress = ref('')
    const endAddress = ref('')
    const geolocation = ref(null)

    // 初始化地图
    const initMap = () => {
      // 设置浦东新区为默认中心点
      const defaultCenter = new BMap.Point(121.7589, 31.0456) // 惠南镇坐标
      map.value = new BMap.Map(mapContainer.value, {
        enableHighAccuracy: true,
        coordsType: BMAP_COORDS_BD09  // 使用百度坐标系
      })
      
      // 先定位到惠南，确保用户体验
      map.value.centerAndZoom(defaultCenter, 15)
      
      // 创建定位控件
      const locationControl = new BMap.GeolocationControl({
        // 允许浏览器定位
        enableAutoLocation: true,
        // 定位按钮的停靠位置
        anchor: BMAP_ANCHOR_TOP_RIGHT,
        // 定位按钮的水平偏移值
        offset: new BMap.Size(10, 10),
        // 定位成功后将定位结果显示在地图上
        showAddressInfo: true,
        // 用户手动点击定位按钮后进行定位
        locationIcon: new BMap.Icon("/icons/location.png", new BMap.Size(32, 32))
      })
      map.value.addControl(locationControl)

      // 创建新的定位对象
      const newGeolocation = new BMap.Geolocation()
      // 启用SDK辅助定位
      newGeolocation.enableSDKLocation()
      
      // 获取当前位置
      newGeolocation.getCurrentPosition(
        function(r) {
          if (this.getStatus() == BMAP_STATUS_SUCCESS) {
            // 定位成功
            const mk = new BMap.Marker(r.point)
            map.value.addOverlay(mk)
            map.value.panTo(r.point)
            
            // 获取详细地址
            const geoc = new BMap.Geocoder()
            geoc.getLocation(r.point, (rs) => {
              const addComp = rs.addressComponents
              console.log('当前位置：', {
                district: addComp.district,
                street: addComp.street,
                streetNumber: addComp.streetNumber,
                point: r.point
              })
            })
          } else {
            console.error('定位失败，错误码：' + this.getStatus())
            // 定位失败时使用默认位置（惠南）
            map.value.centerAndZoom(defaultCenter, 15)
          }
        },
        {
          enableHighAccuracy: true,    // 启用高精度
          timeout: 10000,              // 超时时间，单位毫秒
          maximumAge: 0,               // 禁用缓存
          SDKLocation: true            // 启用SDK辅助定位
        }
      )

      // 启用滚轮缩放和自动适应容器尺寸
      map.value.enableScrollWheelZoom()
      map.value.enableAutoResize()
      
      // 添加地图点击事件
      map.value.addEventListener('click', handleMapClick)
    }

    // 处理地图点击
    const handleMapClick = (e) => {
      if (!startPoint.value) {
        // 设置起点
        startPoint.value = e.point
        const marker = new BMap.Marker(e.point, {
          icon: new BMap.Icon('/icons/start-point.png', new BMap.Size(32, 32)),
          title: '起点'
        })
        map.value.addOverlay(marker)
        
        // 获取地址信息
        const geoc = new BMap.Geocoder()
        geoc.getLocation(e.point, (rs) => {
          const addComp = rs.addressComponents
          startAddress.value = addComp.district + addComp.street + addComp.streetNumber
        })
      } else if (!endPoint.value) {
        // 设置终点
        endPoint.value = e.point
        const marker = new BMap.Marker(e.point, {
          icon: new BMap.Icon('/icons/end-point.png', new BMap.Size(32, 32)),
          title: '终点'
        })
        map.value.addOverlay(marker)
        
        // 获取地址信息
        const geoc = new BMap.Geocoder()
        geoc.getLocation(e.point, (rs) => {
          const addComp = rs.addressComponents
          endAddress.value = addComp.district + addComp.street + addComp.streetNumber
        })
      }
    }

    // 处理打车请求
    const handleOrder = async () => {
      if (!startPoint.value || !endPoint.value) {
        alert('请先选择起点和终点')
        return
      }

      try {
        // 格式化经纬度为 "latitude,longitude" 格式
        const origin = `${startPoint.value.lat.toFixed(6)},${startPoint.value.lng.toFixed(6)}`
        const destination = `${endPoint.value.lat.toFixed(6)},${endPoint.value.lng.toFixed(6)}`
        
        const response = await getPathPlanning(origin, destination)
        if (response.code === 200) {
          // 成功提示
          alert('行程规划成功')
        } else {
          alert('请求失败：' + response.msg)
        }
      } catch (error) {
        console.error('行程规划错误：', error)
        alert('行程规划失败，请稍后重试')
      }
    }

    onMounted(() => {
      initMap()
    })

    return {
      mapContainer,
      startAddress,
      endAddress,
      handleOrder
    }
  }
}
</script>

<style scoped>
.home {
  position: relative;
  width: 100%;
  height: 100vh;
}

.map-container {
  width: 100%;
  height: 100%;
}

.order-box {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  padding: 20px;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
}

.location-info {
  margin-bottom: 15px;
}

.start-point, .end-point {
  margin: 5px 0;
}

.order-btn {
  width: 100%;
  padding: 12px;
  background: #18a45b;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
}

.order-btn:hover {
  background: #148a4c;
}
</style> 