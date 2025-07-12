package utils

// Description:
//
// 使用凭据初始化账号Client
//
// @return Client
////
//// @throws Exception
//func CreateClient() (_result *dysmsapi20170525.Client, _err error) {
//	// 工程代码建议使用更安全的无AK方式，凭据配置方式请参见：https://help.aliyun.com/document_detail/378661.html。
//	_, _err = credential.NewCredential(nil)
//	if _err != nil {
//		return _result, _err
//	}
//	config := &openapi.Config{
//		// 您的AccessKey ID
//		AccessKeyId: tea.String(config2.DataConfig.ALiYun.AccessKeyID),
//		// 您的AccessKey Secret
//		AccessKeySecret: tea.String(config2.DataConfig.ALiYun.AccessKeySecret),
//	}
//	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
//	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
//	_result = &dysmsapi20170525.Client{}
//	_result, _err = dysmsapi20170525.NewClient(config)
//	return _result, _err
//}
//
//func SendSms(mobile, code string) (sms *dysmsapi20170525.SendSmsResponse, _err error) {
//	client, _err := CreateClient()
//	if _err != nil {
//		return sms, _err
//	}
//
//	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
//		PhoneNumbers:  &mobile,
//		TemplateParam: &code,
//	}
//	sms, _err = client.SendSmsWithOptions(sendSmsRequest, &util.RuntimeOptions{})
//	return sms, _err
//}
