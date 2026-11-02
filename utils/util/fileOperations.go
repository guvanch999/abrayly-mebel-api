package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"
)

// allowedImageExtensions lists the upload types that are stored and later
// served from /uploads. Anything else (e.g. .html or .svg, which a browser
// could execute) is rejected.
var allowedImageExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".webp": {},
	".gif":  {},
	".avif": {},
}

// ImageExtension returns the lower-cased extension of an uploaded file name
// or a 400 ResponseError if it is not an allowed image type.
func ImageExtension(filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if _, ok := allowedImageExtensions[ext]; !ok {
		return "", exceptions.NewResponseError(
			exceptions.ErrBadRequest,
			fmt.Errorf("unsupported image type %q: allowed are .jpg, .jpeg, .png, .webp, .gif, .avif", ext),
		)
	}
	return ext, nil
}

// SaveUploadedFile stores an uploaded image as <dir>/<name><ext> and returns
// that relative path (which is also its URL path under /uploads).
func SaveUploadedFile(file *multipart.FileHeader, name, dir string) (string, error) {
	ext, err := ImageExtension(file.Filename)
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	fileUrl := path.Join(dir, name+ext)

	dst, err := os.Create(fileUrl)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return fileUrl, nil
}

// DeleteFile removes a previously stored upload.
func DeleteFile(path string) error {
	return os.Remove(path)
}
