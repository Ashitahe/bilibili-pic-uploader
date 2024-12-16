package config

import (
	"bilibili-uploader/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// LoadConfig 读取配置文件
func LoadConfig(configPath string) (*models.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// ExtractCSRF 从cookie中提取bili_jct值
func ExtractCSRF(cookie string) string {
	parts := strings.Split(cookie, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "bili_jct=") {
			return strings.TrimPrefix(part, "bili_jct=")
		}
	}
	return ""
}