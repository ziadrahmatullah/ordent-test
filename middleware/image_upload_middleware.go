package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/ordent-test/apperror"
)

func ImageUploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const (
			MB = 1 << 20
		)

		type Sizer interface {
			Size() int64
		}
		file, _, err := c.Request.FormFile("image")
		if err != nil {
			if err == http.ErrMissingFile && c.Request.Method == http.MethodPut {
				c.Next()
				return
			}
			c.Error(apperror.NewClientError(err))
			c.Abort()
			return
		}
		defer file.Close()
		if err := c.Request.ParseMultipartForm(0.5 * MB); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 0.5*MB)

		fileHeader := make([]byte, 512)

		if _, err := file.Read(fileHeader); err != nil {

			c.Error(err)
			c.Abort()
			return
		}

		if _, err := file.Seek(0, 0); err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		if file.(Sizer).Size() > 500000 {
			c.Error(apperror.NewClientError(errors.New("image must below 500 kb")))
			c.Abort()
			return
		}
		imageType := http.DetectContentType(fileHeader)
		if imageType == "image/png" {
			ctx := context.WithValue(c.Request.Context(), "image", file)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			return
		}
		c.Error(apperror.NewClientError(errors.New("image type must be png")))
		c.Abort()
	}
}
