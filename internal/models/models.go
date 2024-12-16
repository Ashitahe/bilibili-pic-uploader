package models

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

// UploadResult 上传结果
type UploadResult struct {
    LocalPath string `json:"local_path"`
    RemoteURL string `json:"remote_url"`
    Success   bool   `json:"success"`
    Error     string `json:"error,omitempty"`
}

// Config 配置文件结构
type Config struct {
	Cookie     string `json:"cookie"`
	InputDir   string `json:"input_dir"`
	OutputFile string `json:"output_file"`
}