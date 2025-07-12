package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func GetHostIp() string {
	resp, err := http.Get("http://api.ipify.org")
	if err != nil {
		fmt.Printf("请求失败: %s\n", err)
		return ""
	}
	defer resp.Body.Close() // 确保关闭响应体

	// 读取响应体内容
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %s\n", err)
		return ""
	}

	return string(ip)
}

// 定义几个 IP 查询服务的 URL
var ipServices = []string{
	"https://api.ipify.org",
	"https://ipinfo.io/ip",
	"https://icanhazip.com",
	"https://api.myip.com",
}

// IPServiceResponse 定义了从某些 IP 查询服务获取的响应结构
type IPServiceResponse struct {
	IP string `json:"ip"`
}

// getExternalIP 函数尝试从多个服务获取外网 IP，提高可靠性
func GetExternalIP() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var lastErr error

	// 尝试从多个服务获取 IP，只要有一个成功就返回
	for _, serviceURL := range ipServices {
		ip, err := fetchIPFromService(ctx, serviceURL)
		if err == nil {
			// 验证 IP 格式是否有效
			if net.ParseIP(strings.TrimSpace(ip)) != nil {
				return strings.TrimSpace(ip), nil
			}
			lastErr = fmt.Errorf("从 %s 获取的 IP 格式无效: %s", serviceURL, ip)
		} else {
			lastErr = err
		}
		fmt.Printf("尝试 %s 失败: %v\n", serviceURL, err)
	}

	if lastErr != nil {
		return "", fmt.Errorf("所有 IP 查询服务都失败: %w", lastErr)
	}

	return "", fmt.Errorf("无法获取外网 IP")
}

// fetchIPFromService 函数从指定的服务获取 IP
func fetchIPFromService(ctx context.Context, url string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP 请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 处理不同服务的响应格式
	ip := string(body)

	// 如果是 JSON 格式的响应，尝试解析
	if strings.Contains(url, "api.myip.com") {
		var response IPServiceResponse
		if err := json.Unmarshal(body, &response); err == nil {
			ip = response.IP
		}
	}

	return strings.TrimSpace(ip), nil
}
