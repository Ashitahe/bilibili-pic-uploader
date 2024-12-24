package uploader

import (
	"bilibili-uploader/internal/compression"
	"bilibili-uploader/internal/models"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type Uploader struct {
	CSRF       string
	Cookie     string
	Compressor *compression.Compressor
}

func New(csrf, cookie string, enableCompression bool, quality float32) *Uploader {
	var compressor *compression.Compressor
	if enableCompression {
		compressor = compression.New(quality)
	}
	
	return &Uploader{
		CSRF:       csrf,
		Cookie:     cookie,
		Compressor: compressor,
	}
}

func (u *Uploader) UploadImage(imagePath string) (*models.Response, *models.CompressionInfo, error) {
	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, nil, fmt.Errorf("读取文件失败: %v", err)
	}

	var compressionInfo *models.CompressionInfo
	
	// 如果启用压缩，进行压缩处理
	if u.Compressor != nil {
		compressedData, info, err := u.Compressor.CompressImage(fileData)
		if err != nil {
			return nil, nil, fmt.Errorf("压缩图片失败: %v", err)
		}
		fileData = compressedData
		compressionInfo = info
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("csrf", u.CSRF); err != nil {
		return nil, nil, fmt.Errorf("写入csrf失败: %v", err)
	}

	part, err := writer.CreateFormFile("binary", "read-editor-"+fmt.Sprintf("%d", time.Now().UnixNano()))
	if err != nil {
		return nil, nil, fmt.Errorf("创建文件字段失败: %v", err)
	}

	if _, err := part.Write(fileData); err != nil {
		return nil, nil, fmt.Errorf("写入图片数据失败: %v", err)
	}

	if err := writer.Close(); err != nil {
		return nil, nil, fmt.Errorf("关闭writer失败: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.bilibili.com/x/article/creative/article/upcover", body)
	if err != nil {
		return nil, nil, fmt.Errorf("创建请求失败: %v", err)
	}

	u.setHeaders(req, writer.FormDataContentType())

	req.Header.Set("X-Csrf-Token", u.CSRF)

	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	fmt.Printf("Http响应状态: %v\n", resp.Status)

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("创建gzip解压器失败: %v", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	fmt.Printf("响应体: %s\n", string(bodyBytes))

	var result models.Response
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Code == 0 {
		result.Data.URL = strings.Replace(result.Data.URL, "http", "https", -1)
		result.Data.URL = fmt.Sprintf("https://images.weserv.nl/?url=%s", result.Data.URL)
	}

	return &result, compressionInfo, nil
}

func (u *Uploader) setHeaders(req *http.Request, contentType string) {
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cookie", u.Cookie)
	req.Header.Set("Origin", "https://member.bilibili.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://member.bilibili.com/")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "macOS")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	// 只接受 gzip 压缩
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("DNT", "1")
	req.Header.Set("X-Csrf-Token", u.CSRF)
}