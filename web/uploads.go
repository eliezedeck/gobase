package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"

	"github.com/eliezedeck/gobase/azure"
	"github.com/eliezedeck/gobase/random"
	"github.com/strukturag/libheif/go/heif"

	"github.com/labstack/echo/v4"
)

const (
	ErrorCodeBadUploadFileSize = 1300000
	ErrorCodeUnknownImageType  = 1300001
)

var (
	ErrorNotAnHEICImage = errors.New("not an HEIC image")
)

func UploadImageToAzureBlob(c echo.Context, bs *azure.BlobService, azureContainer string, formName string, maxFileSize int64) (string, error) {
	fh, err := c.FormFile(formName)
	if err != nil {
		return "", err // HTTP 500
	}
	if fh.Size > maxFileSize {
		return "", ErrorWithCode(c, ErrorCodeBadUploadFileSize, "Image size too big")
	}
	if fh.Size <= 512 {
		return "", ErrorWithCode(c, ErrorCodeBadUploadFileSize, "Image size too small")
	}

	fc, err := fh.Open()
	if err != nil {
		return "", err // HTTP 500
	}
	defer fc.Close()
	content, err := io.ReadAll(fc)
	if err != nil {
		return "", err // HTTP 500
	}
	contentType := http.DetectContentType(content)
	fileext := "jpg"
	switch contentType {
	case "image/png":
		fileext = "png"
	case "image/jpg", "image/jpeg":
		fileext = "jpg"
	case "image/tiff":
		fileext = "tiff"
	case "image/svg+xml":
		fileext = "svg"
	default:
		// We are going to detect if this is a HEIC file
		if len(content) > 12 && content[4] == 'f' && content[5] == 't' && content[6] == 'y' && content[7] == 'p' && content[8] == 'h' && content[9] == 'e' && content[10] == 'i' && content[11] == 'c' {
			jpegdata, err := convertHeicToJpeg(content, 85)
			if err != nil {
				if err == ErrorNotAnHEICImage {
					return "", ErrorWithCode(c, ErrorCodeUnknownImageType, "Only Image (JPEG, PNG, TIFF, HEIC) files are supported")
				} else {
					return "", err
				}
			}
			content = jpegdata
			contentType = "image/jpeg"
			fileext = "jpg"
		} else {
			return "", ErrorWithCode(c, ErrorCodeUnknownImageType, "Only Image (JPEG, PNG, TIFF, HEIC) files are supported")
		}
	}

	// Upload the file to Azure Blob
	filename := fmt.Sprintf("%s.%s", random.String(128), fileext)
	fu, err := bs.UploadInContainer(context.Background(), azureContainer, filename, contentType, bytes.NewReader(content))
	if err != nil {
		return "", err // HTTP 500
	}

	return fu.String(), nil
}

func convertHeicToJpeg(data []byte, quality int) ([]byte, error) {
	// Check if the data is an HEIC image
	ctx, err := heif.NewContext()
	if err != nil {
		return nil, err
	}
	err = ctx.ReadFromMemory(data)
	if err != nil {
		return nil, ErrorNotAnHEICImage
	}

	// Decode the HEIC image
	handle, err := ctx.GetPrimaryImageHandle()
	if err != nil {
		return nil, err
	}

	img, err := handle.DecodeImage(heif.ColorspaceUndefined, heif.ChromaUndefined, nil)
	if err != nil {
		return nil, err
	}

	// Convert HEIC image to Go's native image.Image
	goImg, err := img.GetImage()
	if err != nil {
		return nil, err
	}

	// Encode the image as JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, goImg, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
