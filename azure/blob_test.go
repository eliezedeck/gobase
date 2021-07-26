package azure

import (
	"context"
	"strings"
	"testing"

	"github.com/eliezedeck/gobase/config"
)

func TestFullCycle(t *testing.T) {
	bs, err := NewBlobService(config.MustGetEnvValue("AZURE_ACCOUNT_NAME"), config.MustGetEnvValue("AZURE_ACCOUNT_KEY"))
	if err != nil {
		t.Error(err)
		return
	}

	// Upload
	content := strings.NewReader("hello world!")
	u, err := bs.UploadInContainer(context.Background(), "test-container", "hello-world.txt", "text/plain", content)
	if err != nil {
		t.Error(err)
		return
	}
	surl := u.String()

	// Test that the file exists
	exists, err := bs.IsFileExists(surl)
	if err != nil {
		t.Error(err)
		return
	}
	if !exists {
		t.Fatal("Upload file didn't exist!")
		return
	}

	// Test that a non-existing file doesn't exist
	exists, err = bs.IsFileExists(surl + "-force-not-exist")
	if err != nil {
		t.Error(err)
		return
	}
	if exists {
		t.Fatal("Wrong file URL return status 'exists'")
		return
	}

	// Delete
	if err := bs.DeleteByUrl(context.Background(), surl); err != nil {
		t.Error(err)
		return
	}
}
