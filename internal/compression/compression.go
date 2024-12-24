package compression

import (
	"bilibili-uploader/internal/models"
	"bytes"
	"fmt"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

type Compressor struct {
	Quality int
}

func New(quality float32) *Compressor {
	// 将 float32 转换为 int (0-100)
	intQuality := int(quality)
	if intQuality <= 0 || intQuality > 100 {
		intQuality = 75 // 默认质量
	}
	
	return &Compressor{
		Quality: intQuality,
	}
}

func (c *Compressor) CompressImage(data []byte) ([]byte, *models.CompressionInfo, error) {
	originalSize := int64(len(data))
	
	// 解码图片
	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, nil, fmt.Errorf("解码图片失败: %v", err)
	}

	// 创建输出buffer
	buf := new(bytes.Buffer)

	// 以指定质量编码为JPEG
	err = jpeg.Encode(buf, img, &jpeg.Options{
		Quality: c.Quality,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("压缩图片失败: %v", err)
	}

	compressedData := buf.Bytes()
	compressedSize := int64(len(compressedData))
	
	compressionInfo := &models.CompressionInfo{
		OriginalSize:      originalSize,
		CompressedSize:    compressedSize,
		OriginalSizeStr:   models.FormatFileSize(originalSize),
		CompressedSizeStr: models.FormatFileSize(compressedSize),
		CompressionRate:   float64(originalSize-compressedSize) / float64(originalSize) * 100,
	}

	return compressedData, compressionInfo, nil
} 