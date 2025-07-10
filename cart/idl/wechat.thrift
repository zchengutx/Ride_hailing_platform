namespace go cart.wechat

// 微信签名验证请求
struct SignReq {
    1: string signature    // 微信加密签名
    2: string timestamp    // 时间戳
    3: string nonce        // 随机数
    4: string echostr      // 随机字符串
}

// 微信签名验证响应
struct SignResp {
    1: i16 code
    2: string message
    3: optional string echostr  // 验证成功时返回
}

// 获取授权二维码请求
struct GetQRCodeReq {
    // 此接口无特殊参数
}

// 获取授权二维码响应
struct GetQRCodeResp {
    1: i16 code
    2: string message
    3: optional binary qr_code_data  // 二维码图片数据
    4: optional string auth_url      // 授权URL
}

// 微信授权回调请求
struct CallbackReq {
    1: string code         // 微信授权码
    2: optional string state // 状态参数
}

// 用户信息结构
struct WeChatUserInfo {
    1: string openid       // 用户唯一标识
    2: string nickname     // 用户昵称
    3: string headimgurl   // 用户头像URL
    4: i16 sex             // 用户性别，1为男性，2为女性，0为未知
    5: string country      // 用户所在国家
    6: string province     // 用户所在省份
    7: string city         // 用户所在城市
    8: string language     // 用户语言
    9: list<string> privilege // 用户特权信息
}

// 微信授权回调响应
struct CallbackResp {
    1: i16 code
    2: string message
    3: optional WeChatUserInfo user_info  // 用户信息
    4: optional string access_token       // 访问令牌
}

// 微信服务
service WeChatService {
    // 微信签名验证
    SignResp Sign(1: SignReq req)
    
    // 获取授权二维码
    GetQRCodeResp GetQRCode(1: GetQRCodeReq req)
    
    // 处理微信授权回调
    CallbackResp Callback(1: CallbackReq req)
} 