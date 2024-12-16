package storage

import (
	"bilibili-uploader/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SaveResults 保存结果到JSON文件
func SaveResults(results []models.UploadResult, outputPath string) error {
    jsonData, err := json.MarshalIndent(results, "", "    ")
    if err != nil {
        return fmt.Errorf("转换JSON失败: %v", err)
    }
    
    return os.WriteFile(outputPath, jsonData, 0644)
}

// IsImageFile 判断文件是否为图片
func IsImageFile(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    switch ext {
    case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
        return true
    default:
        return false
    }
}