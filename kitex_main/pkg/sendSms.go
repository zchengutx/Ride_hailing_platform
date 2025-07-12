package pkg

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
	"github.com/spf13/viper"
)

// Description:
//
// 使用凭据初始化账号Client
//
// @return Client
//
// @throws Exception
func CreateClient() (_result *dysmsapi20170525.Client, _err error) {
	// 工程代码建议使用更安全的无AK方式，凭据配置方式请参见：https://help.aliyun.com/document_detail/378661.html。
	credential, _err := credential.NewCredential(nil)
	if _err != nil {
		return _result, _err
	}

	config := &openapi.Config{
		AccessKeyId:     tea.String(viper.GetString("LTAI5tLNaojcRjuW2u5pFf2o")),
		AccessKeySecret: tea.String(viper.GetString("VkeNMo2lXCUZefIdoWcZWceLghzmJ4")),
		Credential:      credential,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func SendSms(mobile, code string) (options *dysmsapi20170525.SendSmsResponse, _err error) {
	client, _err := CreateClient()
	if _err != nil {
		return options, _err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers: tea.String(mobile),
		SignName:     tea.String(code),
	}

	options, err := client.SendSmsWithOptions(sendSmsRequest, &util.RuntimeOptions{})

	return options, err
}
