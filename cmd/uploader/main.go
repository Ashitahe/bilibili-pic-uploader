package main

import (
	"bilibili-uploader/internal/config"
	"bilibili-uploader/internal/models"
	"bilibili-uploader/internal/storage"
	"bilibili-uploader/internal/uploader"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("读取配置失败: %v\n", err)
		os.Exit(1)
	}

	// 验证配置
	if cfg.InputDir == "" {
		fmt.Println("错误：配置文件中未指定输入目录(input_dir)")
		os.Exit(1)
	}

	csrf := config.ExtractCSRF(cfg.Cookie)
	if csrf == "" {
		fmt.Println("错误：无法从cookie中提取bili_jct值")
		os.Exit(1)
	}

	// 使用默认输出文件名（如果未在配置中指定）
	if cfg.OutputFile == "" {
		cfg.OutputFile = "upload_results.json"
	}

	// 创建上传器
	uploader := uploader.New(
		csrf, 
		cfg.Cookie,
		cfg.Compression.Enabled,
		cfg.Compression.Quality,
	)
	
	var results []models.UploadResult
	
	// 处理上传
	err = processUploads(cfg.InputDir, uploader, &results)
	if err != nil {
		fmt.Printf("处理上传失败: %v\n", err)
		os.Exit(1)
	}

	// 保存结果
	if err := storage.SaveResults(results, cfg.OutputFile); err != nil {
		fmt.Printf("保存结果失败: %v\n", err)
		os.Exit(1)
	}

	printSummary(results, cfg.OutputFile)
}

func processUploads(inputDir string, uploader *uploader.Uploader, results *[]models.UploadResult) error {
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && storage.IsImageFile(path) {
			fmt.Printf("正在处理: %s\n", path)
			
			result, compressionInfo, err := uploader.UploadImage(path)
			uploadResult := models.UploadResult{
				LocalPath:       path,
				RemoteURL:       result.Data.URL,
				Success:         true,
				CompressionInfo: compressionInfo,
			}
			
			if err != nil {
				uploadResult.Success = false
				uploadResult.Error = err.Error()
			} else if result.Code == 0 {
				uploadResult.Success = true
				uploadResult.RemoteURL = result.Data.URL
			} else {
				uploadResult.Success = false
				uploadResult.Error = result.Message
			}
			
			*results = append(*results, uploadResult)
			
			fmt.Printf("处理完成: %s, 状态: %v\n", path, uploadResult.Success)
			if !uploadResult.Success {
				fmt.Printf("错误信息: %s\n", uploadResult.Error)
			}

			time.Sleep(1 * time.Second)
		}
		return nil
	})
}

func printUsage() {
	fmt.Println("使用方法:")
	fmt.Println("  程序 [-config <配置文件路径>]")
	fmt.Println("\n配置文件格式(config.json):")
	fmt.Println("  {")
	fmt.Println("    \"cookie\": \"你的B站cookie\",")
	fmt.Println("    \"input_dir\": \"要上传的图片目录\",")
	fmt.Println("    \"output_file\": \"结果输出文件路径\"")
	fmt.Println("  }")
}

func printSummary(results []models.UploadResult, outputJSON string) {
	successCount := 0
	var totalOriginalSize, totalCompressedSize int64
	
	for _, result := range results {
		if result.Success {
			successCount++
			if result.CompressionInfo != nil {
				totalOriginalSize += result.CompressionInfo.OriginalSize
				totalCompressedSize += result.CompressionInfo.CompressedSize
			}
		}
	}
	
	fmt.Printf("\n处理完成，共处理 %d 个文件\n", len(results))
	fmt.Printf("成功：%d 个\n", successCount)
	
	if totalOriginalSize > 0 {
		fmt.Printf("总压缩前大小：%s\n", models.FormatFileSize(totalOriginalSize))
		fmt.Printf("总压缩后大小：%s\n", models.FormatFileSize(totalCompressedSize))
		fmt.Printf("平均压缩率：%.2f%%\n", float64(totalOriginalSize-totalCompressedSize)/float64(totalOriginalSize)*100)
	}
	
	fmt.Printf("结果已保存到：%s\n", outputJSON)
}