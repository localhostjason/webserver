package uv

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
)

func SaveFile(src io.Reader, dstFile string) error {

	t, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upload package :%v", err))
	}
	defer func() { _ = t.Close() }()

	_, err = io.Copy(t, src)
	return err
}
func DownloadFile(c *gin.Context, filename string) {
	c.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", filepath.Base(filename)))
	c.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filename)
}
