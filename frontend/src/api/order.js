import request from './request'

// 路程规划接口
export function getPathPlanning(origin, destination) {
  return request({
    url: '/order/pathPlanning',
    method: 'get',
    params: {
      origin,
      destination
    }
  })
} 