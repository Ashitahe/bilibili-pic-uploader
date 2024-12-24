package models

import "fmt"

// Response B站API响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
			Size int    `json:"size"`
			URL  string `json:"url"`
	} `json:"data"`
}

// CompressionInfo 压缩信息
type CompressionInfo struct {
	OriginalSize      int64   `json:"original_size"`
	CompressedSize    int64   `json:"compressed_size"`
	OriginalSizeStr   string  `json:"original_size_str"`
	CompressedSizeStr string  `json:"compressed_size_str"`
	CompressionRate   float64 `json:"compression_rate"`
}

// UploadResult 上传结果
type UploadResult struct {
	LocalPath       string          `json:"local_path"`
	RemoteURL       string          `json:"remote_url"`
	Success         bool            `json:"success"`
	Error           string          `json:"error,omitempty"`
	CompressionInfo *CompressionInfo `json:"compression_info,omitempty"`
}

// Config 配置文件结构
type Config struct {
	Cookie     string `json:"cookie"`
	InputDir   string `json:"input_dir"`
	OutputFile string `json:"output_file"`
	Compression struct {
		Enabled bool    `json:"enabled"`
		Quality float32 `json:"quality"`
	} `json:"compression"`
}

// FormatFileSize 格式化文件大小
func FormatFileSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
	)

	switch {
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}