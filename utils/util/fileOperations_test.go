package util

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
)

func TestImageExtension(t *testing.T) {
	tests := []struct {
		filename string
		want     string
		wantErr  bool
	}{
		{filename: "sofa.jpg", want: ".jpg"},
		{filename: "SOFA.JPEG", want: ".jpeg"},
		{filename: "banner.Png", want: ".png"},
		{filename: "photo.webp", want: ".webp"},
		{filename: "anim.gif", want: ".gif"},
		{filename: "next-gen.avif", want: ".avif"},
		{filename: "page.html", wantErr: true},
		{filename: "vector.svg", wantErr: true},
		{filename: "script.jpg.js", wantErr: true},
		{filename: "no-extension", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got, err := ImageExtension(tt.filename)
			if tt.wantErr {
				var respErr *exceptions.ResponseError
				if !errors.As(err, &respErr) || respErr.Code != http.StatusBadRequest {
					t.Fatalf("error = %v, want a 400 ResponseError", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("extension = %q, want %q", got, tt.want)
			}
		})
	}
}

// newFileHeader builds a *multipart.FileHeader the same way Echo's
// c.FormFile does when parsing a multipart request.
func newFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	form, err := multipart.NewReader(&body, writer.Boundary()).ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("ReadForm: %v", err)
	}
	t.Cleanup(func() { _ = form.RemoveAll() })

	return form.File["image"][0]
}

func TestSaveUploadedFileStoresImage(t *testing.T) {
	dir := filepath.ToSlash(filepath.Join(t.TempDir(), "uploads", "products"))
	content := []byte("fake image bytes")

	path, err := SaveUploadedFile(newFileHeader(t, "Sofa.JPG", content), "product-id", dir)
	if err != nil {
		t.Fatalf("SaveUploadedFile: %v", err)
	}

	if want := dir + "/product-id.jpg"; path != want {
		t.Errorf("path = %q, want %q", path, want)
	}
	stored, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read stored file: %v", err)
	}
	if !bytes.Equal(stored, content) {
		t.Errorf("stored content = %q, want %q", stored, content)
	}

	if err := DeleteFile(path); err != nil {
		t.Fatalf("DeleteFile: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("file still exists after DeleteFile (stat err = %v)", err)
	}
}

func TestSaveUploadedFileRejectsNonImages(t *testing.T) {
	dir := t.TempDir()

	_, err := SaveUploadedFile(newFileHeader(t, "evil.html", []byte("<script>alert(1)</script>")), "id", dir)
	if err == nil {
		t.Fatal("expected an error for an .html upload")
	}

	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		t.Fatalf("ReadDir: %v", readErr)
	}
	if len(entries) != 0 {
		t.Errorf("rejected upload left %d file(s) on disk", len(entries))
	}
}
