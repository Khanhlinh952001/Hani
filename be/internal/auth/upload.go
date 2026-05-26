package auth

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"be/internal/modules/users"

	"github.com/gin-gonic/gin"
)

const (
	maxAvatarBytes  = 5 << 20 // 5MB
	avatarFormField = "avatar"
)

var avatarExtByMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

func UploadAvatarHandler(c *gin.Context) {
	userID, ok := UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := c.FormFile(avatarFormField)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing avatar file"})
		return
	}
	if file.Size > maxAvatarBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 5MB)"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read upload"})
		return
	}
	defer src.Close()

	head := make([]byte, 512)
	n, _ := io.ReadFull(src, head)
	if n == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty file"})
		return
	}

	contentType := http.DetectContentType(head[:n])
	ext, ok := avatarExtByMIME[contentType]
	if !ok {
		ext = extFromFilename(file.Filename)
	}
	if ext == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only JPEG, PNG, WebP or GIF allowed"})
		return
	}

	uploadDir := filepath.Join("uploads", "avatars")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create upload dir"})
		return
	}

	for _, oldExt := range avatarExtByMIME {
		_ = os.Remove(filepath.Join(uploadDir, fmt.Sprintf("%d%s", userID, oldExt)))
	}

	destName := fmt.Sprintf("%d%s", userID, ext)
	destPath := filepath.Join(uploadDir, destName)

	dst, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save file"})
		return
	}
	defer dst.Close()

	if _, err := dst.Write(head[:n]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save file"})
		return
	}
	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save file"})
		return
	}

	avatarURL := "/uploads/avatars/" + destName
	patch := &users.User{Avatar: avatarURL}
	if err := users.UpdateUserService(strconv.Itoa(userID), patch, ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := users.GetUserByIDService(strconv.Itoa(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// extFromFilename is a fallback when MIME sniffing is inconclusive.
func extFromFilename(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".jpg", ".jpeg":
		return ".jpg"
	case ".png":
		return ".png"
	case ".webp":
		return ".webp"
	case ".gif":
		return ".gif"
	default:
		return ""
	}
}
