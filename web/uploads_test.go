package web

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestHEICImageMimeType(t *testing.T) {
	fpath := "/Users/elie/Downloads/IMG_0779.HEIC"
	f, err := os.Open(fpath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Load the image file as bytes
	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	// Detect the content type
	contentType := http.DetectContentType(content)
	t.Logf("Content type: %s", contentType)
}

func TestHEICToJPEGConversion(t *testing.T) {
	fpath := "/Users/elie/Downloads/IMG_0779.HEIC"
	f, err := os.Open(fpath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Load the image file as bytes
	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	// Convert the image to JPEG
	jpegdata, err := convertHeicToJpeg(content, 80)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("JPEG data size: %d", len(jpegdata))

	// Save the JPEG data to a file
	f, err = os.Create("/Users/elie/Downloads/IMG_0779.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(jpegdata)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHEICConversionMemoryLeaks(t *testing.T) {
	fpath := "/Users/elie/Downloads/IMG_0779.HEIC"
	f, err := os.Open(fpath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Load the image file as bytes
	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	// Continously convert the image and check for memory leaks
	for i := 0; i < 1000; i++ {
		jpegdata, err := convertHeicToJpeg(content, 80)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("JPEG data size: %d", len(jpegdata))

		// Check memory usage
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		t.Logf("Allocated memory: %d", m.Alloc)

		if i == 0 {
			t.Logf("First iteration, waiting for 1 second")
			time.Sleep(1 * time.Second)
		} else {
			// Check if the memory usage is not increasing too much
			if m.Alloc > 100000000 {
				t.Fatalf("Memory usage is too high: %d", m.Alloc)
			}
		}
	}
}
