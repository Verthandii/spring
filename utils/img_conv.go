package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
)

// ConvertFormatFromImageByte 转换图片格式图片的Byte
func ConvertFormatFromImageByte(picture []byte, toFormat string) ([]byte, error) {
	if !InSlice(toFormat, []string{"png", "jpeg"}) {
		return nil, errors.New("toFormat Error")
	}
	img, format, err := image.Decode(bytes.NewReader(picture))
	if err != nil {
		return nil, err
	}
	if format == toFormat {
		// return picture, nil
	}
	switch toFormat {
	case "png":
		var buffer bytes.Buffer
		err = png.Encode(&buffer, img)
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	case "jpeg":
		var buffer bytes.Buffer
		err = jpeg.Encode(&buffer, img, &jpeg.Options{Quality: 100})
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}
	return nil, errors.New("Unreachable")
}

// ConvertFormatFromImageBase64 转换图片格式图片的base64
func ConvertFormatFromImageBase64(pictureBase64 string, toFormat string) (string, error) {
	if strings.Contains(pictureBase64, ",") {
		pictureBase64 = strings.Split(pictureBase64, ",")[1]
	}
	decoded, err := base64.StdEncoding.DecodeString(pictureBase64)
	if err != nil {
		return "", err
	}
	imageBytes, err := ConvertFormatFromImageByte(decoded, toFormat)
	if err != nil {
		return "", err
	}
	pngBase64 := base64.StdEncoding.EncodeToString(imageBytes)
	return pngBase64, nil
}

// CompressJPEG 压缩JPEG图像到最大尺寸
func CompressJPEG(input []byte, maxSize int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	quality := 100 // 初始质量
	for quality > 0 {
		buf.Reset()
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, err
		}
		if buf.Len() <= maxSize {
			break
		}
		quality -= 5 // 每次降低5的质量
	}
	if quality <= 0 {
		return nil, fmt.Errorf("无法在保持适当质量的情况下压缩到所需大小")
	}
	return buf.Bytes(), nil
}
