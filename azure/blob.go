package azure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/eliezedeck/gobase/logging"
	"go.uber.org/zap"
)

var (
	ErrInvalidContainerName = errors.New("invalid container name")
)

type BlobService struct {
	accountName    string
	accountKey     string
	pipeline       pipeline.Pipeline
	baseUrl        *url.URL
	defaultService azblob.ServiceURL

	containerURLs      map[string]*azblob.ContainerURL
	containerURLsMutex *sync.Mutex
}

func NewBlobService(accountName, accountKey string) (*BlobService, error) {
	result := &BlobService{
		accountName:        accountName,
		accountKey:         accountKey,
		containerURLs:      make(map[string]*azblob.ContainerURL),
		containerURLsMutex: &sync.Mutex{},
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		logging.L.Error("Azure Blob invalid credential", zap.Error(err))
		return nil, err
	}
	result.pipeline = azblob.NewPipeline(credential, azblob.PipelineOptions{
		Retry: azblob.RetryOptions{
			Policy:        azblob.RetryPolicyExponential,
			MaxTries:      3,
			TryTimeout:    60 * time.Second,
			RetryDelay:    2 * time.Second,
			MaxRetryDelay: 15 * time.Second,
		},
	})
	result.baseUrl, _ = url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	result.defaultService = result.GetNewServiceURL()
	return result, nil
}

func (b *BlobService) GetNewServiceURL() azblob.ServiceURL {
	return azblob.NewServiceURL(*b.baseUrl, b.pipeline)
}

// EnsureContainer always attempts to create the container and will not return any error if it already exists. Otherwise
// it will create the container and cache it. Either way, it will be cached because it is a time-wise expensive op.
//
// If you only want to upload a file to a container, see UploadInContainer() below.
func (b *BlobService) EnsureContainer(ctx context.Context, container string) (*azblob.ContainerURL, error) {
	if strings.ToLower(container) != container {
		return nil, ErrInvalidContainerName
	}

	b.containerURLsMutex.Lock()
	defer b.containerURLsMutex.Unlock()

	if courl, found := b.containerURLs[container]; found {
		return courl, nil
	}

	// This is why we need to cache the result because there is no API to check the existence of a container
	containerURL := b.defaultService.NewContainerURL(container)
	if _, err := containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessBlob); err != nil {
		if serr, ok := err.(azblob.StorageError); ok && serr.ServiceCode() == azblob.ServiceCodeContainerAlreadyExists {
			logging.L.Debug("Azure Storage: Container already exists", zap.String("container", container))
			// Already exists, no need to worry
		} else {
			return nil, err
		}
	}

	// Cache this result
	logging.L.Debug("Azure Storage: Cached container URL", zap.String("container", container))
	b.containerURLs[container] = &containerURL
	return &containerURL, nil
}

func (b *BlobService) UploadInContainer(ctx context.Context, container, filename, contentType string, data io.ReadSeeker) (*url.URL, error) {
	courl, err := b.EnsureContainer(ctx, container)
	if err != nil {
		return nil, err
	}

	blurl := courl.NewBlockBlobURL(filename)
	if _, err = blurl.Upload(
		ctx,
		data,
		azblob.BlobHTTPHeaders{ContentType: contentType},
		azblob.Metadata{},
		azblob.BlobAccessConditions{},
		azblob.DefaultAccessTier,
		nil,
		azblob.ClientProvidedKeyOptions{},
	); err != nil {
		return nil, err
	}

	finalurl := blurl.URL()
	logging.L.Debug("Azure Storage: file uploaded", zap.String("url", finalurl.String()))

	return &finalurl, err
}

func (b *BlobService) DeleteByUrl(ctx context.Context, fullUrl string) error {
	u, err := url.Parse(fullUrl)
	if err != nil {
		return err
	}

	blurl := azblob.NewBlockBlobURL(*u, b.pipeline)
	if _, err := blurl.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{}); err != nil {
		return err
	}
	logging.L.Debug("Azure Storage: file deleted", zap.String("url", fullUrl))
	return nil
}

func (b *BlobService) IsFileExists(fullUrl string) (bool, error) {
	// For efficiency, we'll simply do a HEAD request with the assumption that all files are public
	resp, err := http.DefaultClient.Head(fullUrl)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
