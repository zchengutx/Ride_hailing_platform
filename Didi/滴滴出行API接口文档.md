# 滴滴出行 API 接口文档

## 概述

本文档描述了滴滴出行系统的所有API接口，包括用户相关服务和司机相关服务。所有接口均采用RESTful API设计规范，使用JSON格式进行数据交换。

### 基础信息
- **协议**: HTTPS
- **基础URL**: `https://api.didi.com`
- **数据格式**: JSON
- **字符编码**: UTF-8

### 通用响应格式
```json
{
    "code": 200,
    "msg": "success",
    "data": {}
}
```

### 状态码说明
| HTTP状态码 | 业务码 | 说明 |
|-----------|-------|------|
| 200 | 200 | 请求成功 |
| 400 | 400 | 请求参数错误 |
| 401 | 401 | 未授权访问 |
| 403 | 403 | 禁止访问 |
| 404 | 404 | 资源不存在 |
| 500 | 500+ | 服务器错误 |

---

## 用户服务接口

### 1. 发送短信验证码

#### 接口描述
发送短信验证码到用户手机，用于用户注册和登录验证。

#### 接口信息
- **URL**: `/user/sendSms`
- **Method**: `POST`
- **认证**: 无需认证

#### 请求参数
```json
{
    "mobile": "13800138000",
    "send_sms_code": "register"
}
```

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| mobile | string | 是 | 手机号码，11位数字 |
| send_sms_code | string | 是 | 验证码类型标识 |

#### 响应示例
**成功响应**:
```json
{
    "code": 200,
    "msg": "短信发送成功"
}
```

**失败响应**:
```json
{
    "code": 400,
    "msg": "手机号码格式不正确"
}
```

#### 错误码说明
| 错误码 | 说明 |
|-------|------|
| 400 | 手机号码格式不正确 |
| 500 | 服务器异常 |

---

### 2. 用户登录注册

#### 接口描述
用户通过手机号和验证码进行登录，如果用户不存在则自动注册。

#### 接口信息
- **URL**: `/user/login`
- **Method**: `POST`
- **认证**: 无需认证

#### 请求参数
```json
{
    "mobile": "13800138000",
    "send_sms_code": "123456"
}
```

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| mobile | string | 是 | 手机号码 |
| send_sms_code | string | 是 | 短信验证码 |

#### 响应示例
**成功响应**:
```json
{
    "code": 200,
    "msg": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user_id": 12345,
        "expires_at": 1702728000
    }
}
```

**失败响应**:
```json
{
    "code": 301,
    "msg": "验证码错误，剩余尝试次数：3"
}
```

#### 错误码说明
| 错误码 | 说明 |
|-------|------|
| 301 | 验证码错误或过期 |
| 400 | 参数格式错误 |
| 500 | 服务器异常 |

---

### 3. 实名认证

#### 接口描述
用户提交真实身份信息进行实名认证。

#### 接口信息
- **URL**: `/user/realName`
- **Method**: `POST`
- **认证**: 需要JWT Token

#### 请求头
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

#### 请求参数
```json
{
    "user_name": "张三",
    "sex": "男",
    "age": 28,
    "id_card": "110101199001011234"
}
```

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| user_name | string | 是 | 真实姓名 |
| sex | string | 是 | 性别（男/女） |
| age | int64 | 是 | 年龄 |
| id_card | string | 是 | 身份证号码 |

#### 响应示例
**成功响应**:
```json
{
    "code": 200,
    "msg": "实名信息提交成功"
}
```

**失败响应**:
```json
{
    "code": 400,
    "msg": "身份证格式不正确"
}
```

#### 错误码说明
| 错误码 | 说明 |
|-------|------|
| 400 | 参数格式错误 |
| 404 | 用户未登录 |
| 500 | 服务器异常 |

---

### 4. 用户叫车

#### 接口描述
用户发起叫车请求，系统会根据起始位置匹配周边3公里内的在线司机。

#### 接口信息
- **URL**: `/user/takeACar`
- **Method**: `POST`
- **认证**: 需要JWT Token

#### 请求头
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

#### 请求参数
```json
{
    "start_location": "北京市朝阳区三里屯,39.918058,116.397026",
    "end_location": "北京市海淀区中关村",
    "cart_type": "快车"
}
```

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| start_location | string | 是 | 起始位置，格式："地址,纬度,经度" |
| end_location | string | 是 | 目的地地址 |
| cart_type | string | 是 | 车型（快车/专车/豪华车等） |

#### 响应示例
**成功响应**:
```json
{
    "code": 200,
    "msg": "订单创建成功！订单号：20231215143052001，已通知5位司机，请等待司机接单"
}
```

**部分成功响应**:
```json
{
    "code": 200,
    "msg": "订单创建成功（订单号：20231215143052001），但附近3公里内暂无可用的快车司机，系统将持续为您寻找"
}
```

**失败响应**:
```json
{
    "code": 400,
    "msg": "起始位置坐标无效，请确保定位功能已开启并重新获取位置"
}
```

#### 错误码说明
| 错误码 | 说明 |
|-------|------|
| 400 | 参数错误或坐标无效 |
| 403 | 用户未实名认证 |
| 404 | 用户未登录 |
| 500 | 服务器异常 |

---

## 司机服务接口

> **注意**: 所有司机接口都需要JWT Token认证

### 5. 司机申请

#### 接口描述
用户提交司机申请，包含车辆信息和驾驶证信息。

#### 接口信息
- **URL**: `/driver/callACar`
- **Method**: `POST`
- **认证**: 需要JWT Token

#### 请求头
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

#### 请求参数
```json
{
    "address": "北京市朝阳区",
    "driving_license_number": "110101199001011234",
    "quasi_driving_type": "C1",
    "driving_age": 5,
    "car_num": "京A12345",
    "car_type": "快车",
    "vehicle_mileage": "50000",
    "serving_the_city": "北京"
}
```

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| address | string | 是 | 家庭地址 |
| driving_license_number | string | 是 | 驾驶证号码 |
| quasi_driving_type | string | 是 | 准驾车型 |
| driving_age | uint64 | 是 | 驾龄（年） |
| car_num | string | 是 | 车牌号码 |
| car_type | string | 是 | 车辆类型 |
| vehicle_mileage | string | 是 | 车辆行驶里程 |
| serving_the_city | string | 是 | 服务城市 |

#### 响应示例
**成功响应**:
```json
{
    "code": 200,
    "msg": "申请成功"
}
```

**失败响应**:
```json
{
    "code": 601,
    "msg": "用户未实名,不允许申请"
}
```

#### 错误码说明
| 错误码 | 说明 |
|-------|------|
| 400 | 必要参数不能为空 |
| 404 | 用户未登录 |
| 500 | 服务器异常 |
| 601 | 用户未实名认证 |

---

### 6. 司机审核

#### 接口描述
管理员对司机申请进行审核，通过或拒绝申请。

#### 接口信息
- **URL**: `/driver/driverAudit`
- **Method**: `POST`
- **认证**: 需要JWT Token（管理员权限）

#### 请求头
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

#### 请求参数
```json
{
    "driver_id": 12345,
    "audit_status": "已审核"
}
```

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| driver_id | uint64 | 是 | 司机申请ID |
| audit_status | string | 是 | 审核状态（已审核/审核拒绝） |

#### 响应示例
**成功响应**:
```json
{
    "code": 200,
    "msg": "审核成功"
}
```

**失败响应**:
```json
{
    "code": 505,
    "msg": "该信息已审核"
}
```

#### 错误码说明
| 错误码 | 说明 |
|-------|------|
| 400 | 审核状态不能为空 |
| 404 | 该记录不存在 |
| 500 | 服务器异常 |
| 505 | 该信息已审核 |

---

### 7. 司机添加

#### 接口描述
将审核通过的司机申请转为正式司机账户。

#### 接口信息
- **URL**: `/driver/addDriver`
- **Method**: `POST`
- **认证**: 需要JWT Token（管理员权限）

#### 请求头
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

#### 请求参数
```json
{
    "driver_id": 12345
}
```

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| driver_id | uint64 | 是 | 司机申请ID |

#### 响应示例
**成功响应**:
```json
{
    "code": 200,
    "msg": "司机添加成功"
}
```

**失败响应**:
```json
{
    "code": 505,
    "msg": "用户审核为通过,不允许添加"
}
```

#### 错误码说明
| 错误码 | 说明 |
|-------|------|
| 404 | 该信息不存在 |
| 500 | 服务器异常 |
| 505 | 用户审核未通过 |

---

### 8. 司机状态管理

#### 接口描述
司机更新自己的在线状态（在线/离线）。

#### 接口信息
- **URL**: `/driver/driverOnline`
- **Method**: `POST`
- **认证**: 需要JWT Token

#### 请求头
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

#### 请求参数
```json
{
    "driver_id": 12345,
    "driver_status": "在线"
}
```

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| driver_id | uint64 | 是 | 司机ID |
| driver_status | string | 是 | 司机状态（在线/离线/忙碌） |

#### 响应示例
**成功响应**:
```json
{
    "code": 200,
    "msg": "状态更改成功"
}
```

**失败响应**:
```json
{
    "code": 505,
    "msg": "操作异常"
}
```

#### 错误码说明
| 错误码 | 说明 |
|-------|------|
| 400 | 司机状态不能为空 |
| 404 | 司机不存在 |
| 500 | 服务器异常 |
| 505 | 状态重复或操作异常 |

---

### 9. 司机接单

#### 接口描述
司机接取用户发布的叫车订单，需要在3公里范围内且状态为在线。

#### 接口信息
- **URL**: `/driver/receivingOrder`
- **Method**: `POST`
- **认证**: 需要JWT Token

#### 请求头
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

#### 请求参数
```json
{
    "order_id": "20231215143052001",
    "current_latitude": 39.918058,
    "current_longitude": 116.397026
}
```

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| order_id | string | 是 | 订单ID |
| current_latitude | float64 | 是 | 司机当前纬度 |
| current_longitude | float64 | 是 | 司机当前经度 |

#### 响应示例
**成功响应**:
```json
{
    "code": 200,
    "msg": "接单成功！订单号：20231215143052001，乘客距离您0.8公里，请尽快前往接客",
    "order_details": "{\"order_id\":\"20231215143052001\",\"user_id\":\"123\",\"user_name\":\"张三\",\"user_mobile\":\"13800138000\",\"start_location\":\"北京市朝阳区\",\"end_location\":\"北京市海淀区\",\"cart_type\":\"快车\",\"latitude\":\"39.918058\",\"longitude\":\"116.397026\",\"distance\":\"0.8\",\"assigned_time\":1702728000,\"driver_id\":\"456\",\"driver_latitude\":39.920000,\"driver_longitude\":116.400000,\"status\":\"assigned\"}"
}
```

**失败响应**:
```json
{
    "code": 409,
    "msg": "订单已被其他司机接取"
}
```

```json
{
    "code": 403,
    "msg": "您距离乘客太远（5.2公里），无法接单"
}
```

#### 订单详情字段说明
| 字段名 | 类型 | 说明 |
|-------|------|------|
| order_id | string | 订单ID |
| user_id | string | 用户ID |
| user_name | string | 用户姓名 |
| user_mobile | string | 用户手机号 |
| start_location | string | 起始位置 |
| end_location | string | 目的地 |
| cart_type | string | 车型 |
| latitude | string | 起始位置纬度 |
| longitude | string | 起始位置经度 |
| distance | string | 司机到乘客距离（公里） |
| assigned_time | int64 | 接单时间戳 |
| driver_id | string | 司机ID |
| driver_latitude | float64 | 司机接单时纬度 |
| driver_longitude | float64 | 司机接单时经度 |
| status | string | 订单状态 |

#### 错误码说明
| 错误码 | 说明 |
|-------|------|
| 400 | 参数错误或坐标无效 |
| 403 | 司机状态异常或距离超出范围 |
| 404 | 司机不存在或订单不存在 |
| 409 | 订单已被接取或状态冲突 |
| 500 | 服务器异常 |

---

## 错误处理

### 通用错误响应格式
```json
{
    "code": 错误码,
    "msg": "错误描述",
    "data": null
}
```

### 认证错误
当JWT Token无效或过期时：
```json
{
    "code": 401,
    "msg": "未授权访问",
    "data": null
}
```

### 参数验证错误
当请求参数不符合要求时：
```json
{
    "code": 400,
    "msg": "参数绑定失败",
    "data": "具体错误信息"
}
```

---

## 接口测试示例

### 完整业务流程测试

#### 1. 用户注册登录流程
```bash
# 1. 发送验证码
curl -X POST https://api.didi.com/user/sendSms \
  -H "Content-Type: application/json" \
  -d '{
    "mobile": "13800138000",
    "send_sms_code": "register"
  }'

# 2. 用户登录
curl -X POST https://api.didi.com/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "mobile": "13800138000",
    "send_sms_code": "123456"
  }'
```

#### 2. 用户实名认证和叫车
```bash
# 3. 实名认证
curl -X POST https://api.didi.com/user/realName \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_name": "张三",
    "sex": "男",
    "age": 28,
    "id_card": "110101199001011234"
  }'

# 4. 用户叫车
curl -X POST https://api.didi.com/user/takeACar \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_location": "北京市朝阳区三里屯,39.918058,116.397026",
    "end_location": "北京市海淀区中关村",
    "cart_type": "快车"
  }'
```

#### 3. 司机接单流程
```bash
# 5. 司机接单
curl -X POST https://api.didi.com/driver/receivingOrder \
  -H "Authorization: Bearer DRIVER_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "20231215143052001",
    "current_latitude": 39.920000,
    "current_longitude": 116.400000
  }'
```

---

## 开发注意事项

### 1. 认证要求
- 除了`/user/sendSms`和`/user/login`外，所有接口都需要JWT Token认证
- JWT Token通过`Authorization: Bearer <token>`头部传递
- Token有效期为24小时，过期需要重新登录

### 2. 地理位置格式
- 起始位置格式：`"地址名称,纬度,经度"`
- 纬度范围：-90° 到 +90°
- 经度范围：-180° 到 +180°
- 精度建议保留6位小数

### 3. 接单距离限制
- 司机接单距离限制为3公里
- 系统会自动计算司机与乘客的距离
- 超过3公里的司机无法接单

### 4. 状态管理
- 司机状态：在线/离线/忙碌
- 订单状态：waiting/pushed/assigned/completed
- 只有在线状态的司机才能接收订单推送

### 5. 错误重试
- 网络错误建议重试，最多3次
- 业务错误（4xx）不建议重试
- 服务器错误（5xx）可以短暂重试

---

## 变更日志

| 版本 | 日期 | 变更内容 |
|------|------|---------|
| v1.0 | 2023-12-15 | 初始版本，包含所有核心接口 |

---

## 联系方式

如有接口相关问题，请联系：
- **技术支持**: tech@didi.com
- **API文档**: https://docs.didi.com
- **问题反馈**: https://github.com/didi/api-issues 