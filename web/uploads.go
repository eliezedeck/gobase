package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/eliezedeck/gobase/azure"
	"github.com/eliezedeck/gobase/random"

	"github.com/labstack/echo/v4"
)

func UploadImageToAzureBlob(c echo.Context, bs *azure.BlobService, azureContainer string, formName string, maxFileSize int64) (string, error) {
	fh, err := c.FormFile(formName)
	if err != nil {
		return "", err // HTTP 500
	}
	if fh.Size > maxFileSize {
		return "", Error(c, "Image size too big")
	}
	if fh.Size <= 512 {
		return "", Error(c, "Image size too small")
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
	case "image/jpg":
		fileext = "jpg"
	case "image/tiff":
		fileext = "tiff"
	case "image/svg+xml":
		fileext = "svg"
	default:
		return "", Error(c, "Only Image (JPEG, PNG, TIFF) files are supported")
	}

	// Upload the file to Azure Blob
	filename := fmt.Sprintf("%s.%s", random.String(128), fileext)
	fu, err := bs.UploadInContainer(context.Background(), azureContainer, filename, contentType, bytes.NewReader(content))
	if err != nil {
		return "", err // HTTP 500
	}

	return fu.String(), nil
}
