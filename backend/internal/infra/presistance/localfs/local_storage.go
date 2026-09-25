// Package localfs — файловое хранилище на локальном диске для standalone-режима
// (передача экземпляра ПО на экспертизу, локальная разработка без S3).
// Реализует filestorage.Repository: загрузка через собственный PUT-роут
// с HMAC-токеном (аналог presigned URL), раздача — через GET-роут.
// Включается переменной окружения FILE_STORAGE_MODE=local.
package localfs

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
)

const (
	// uploadExpire — срок действия ссылки на загрузку (аналог presigned URL).
	uploadExpire = 10 * time.Minute

	// pathPattern допускает только пути, которые генерирует само хранилище
	// (см. objectPathFormat / avatarPathFormat), и отсекает path traversal.
	pathPattern = `^(users|pets|blood_requests)/[0-9]{4}/[A-Z]+-[0-9A-Za-z]+/(photos/[0-9]+\.jpg|avatar\.jpg)$`
)

var pathRe = regexp.MustCompile(pathPattern)

// LocalStorage — реализация filestorage.Repository на локальной файловой системе.
type LocalStorage struct {
	dir        string
	publicBase string
	secret     []byte
}

// NewFromEnv создает LocalStorage из переменных окружения:
//   - LOCAL_STORAGE_DIR — корневая директория (по умолчанию ./media);
//   - LOCAL_STORAGE_PUBLIC_URL — публичный префикс URL (по умолчанию /api/v1/files);
//   - LOCAL_STORAGE_SECRET, при отсутствии — JWT_SECRET_KEY.
func NewFromEnv() *LocalStorage {
	dir := os.Getenv("LOCAL_STORAGE_DIR")
	if dir == "" {
		dir = "./media"
	}

	publicBase := strings.TrimRight(os.Getenv("LOCAL_STORAGE_PUBLIC_URL"), "/")
	if publicBase == "" {
		publicBase = "/api/v1/files"
	}

	secret := os.Getenv("LOCAL_STORAGE_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET_KEY")
	}

	return &LocalStorage{
		dir:        dir,
		publicBase: publicBase,
		secret:     []byte(secret),
	}
}

// RegisterRoutes подключает роуты загрузки и раздачи файлов в mux.
// Вызывать на apiMux (путь без префикса /api, который срезает StripPrefix).
func (s *LocalStorage) RegisterRoutes(mux *http.ServeMux) {
	routeBase := strings.TrimPrefix(s.publicBase, "/api")

	mux.HandleFunc(routeBase+"/upload", s.handleUpload)
	mux.Handle(routeBase+"/", http.StripPrefix(routeBase, http.FileServer(http.Dir(s.dir))))
}

func (s *LocalStorage) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	pathParam := q.Get("path")
	expParam := q.Get("exp")
	token := q.Get("token")

	if !pathRe.MatchString(pathParam) {
		http.Error(w, "недопустимый путь", http.StatusBadRequest)
		return
	}

	exp, err := strconv.ParseInt(expParam, 10, 64)
	if err != nil || time.Now().Unix() > exp {
		http.Error(w, "срок действия ссылки истек", http.StatusForbidden)
		return
	}

	if !hmac.Equal([]byte(token), []byte(s.signPath(pathParam, exp))) {
		http.Error(w, "недействительный токен", http.StatusForbidden)
		return
	}

	fullPath := filepath.Join(s.dir, filepath.FromSlash(pathParam))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		slog.ErrorContext(r.Context(), "localfs: не удалось создать директорию", "error", err, "path", pathParam)
		http.Error(w, "ошибка записи файла", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	if err := os.WriteFile(fullPath, body, 0o644); err != nil {
		slog.ErrorContext(r.Context(), "localfs: не удалось записать файл", "error", err, "path", pathParam)
		http.Error(w, "ошибка записи файла", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// signPath подписывает путь и время истечения HMAC-SHA256 с секретом хранилища.
func (s *LocalStorage) signPath(path string, exp int64) string {
	mac := hmac.New(sha256.New, s.secret)
	fmt.Fprintf(mac, "%s|%d", path, exp)
	return hex.EncodeToString(mac.Sum(nil))
}

// objectPathFormat возвращает формат пути к фото для сущности с заданным ID.
// Форматы идентичны S3-хранилищу.
func objectPathFormat(id string) (string, error) {
	year := time.Now().Year()
	switch {
	case strings.HasPrefix(id, "USR"):
		return fmt.Sprintf("users/%d/%%s/photos/%%d.jpg", year), nil
	case strings.HasPrefix(id, "PET"):
		return fmt.Sprintf("pets/%d/%%s/photos/%%d.jpg", year), nil
	case strings.HasPrefix(id, "BLS"):
		return fmt.Sprintf("blood_requests/%d/%%s/photos/%%d.jpg", year), nil
	default:
		return "", fmt.Errorf("неподдерживаемый тип файла для id: %s", id)
	}
}

// avatarPathFormat возвращает формат пути к аватару. Идентичен S3-хранилищу.
func avatarPathFormat(id string) (string, error) {
	year := time.Now().Year()
	switch {
	case strings.HasPrefix(id, "USR"):
		return fmt.Sprintf("users/%d/%%s/avatar.jpg", year), nil
	case strings.HasPrefix(id, "PET"):
		return fmt.Sprintf("pets/%d/%%s/avatar.jpg", year), nil
	case strings.HasPrefix(id, "BLS"):
		return fmt.Sprintf("blood_requests/%d/%%s/avatar.jpg", year), nil
	default:
		return "", fmt.Errorf("неподдерживаемый тип файла для id: %s", id)
	}
}

// GetPresignedURLs возвращает ссылки на загрузку и пути файлов.
// Ссылка подписана токеном с ограниченным сроком действия (аналог presigned URL).
func (s *LocalStorage) GetPresignedURLs(ctx context.Context, count int64, id string) ([]filestorage.UploadInfo, error) {
	format, err := objectPathFormat(id)
	if err != nil {
		return nil, err
	}

	exp := time.Now().Add(uploadExpire).Unix()
	uploadInfos := make([]filestorage.UploadInfo, count)

	for p := range count {
		path := fmt.Sprintf(format, id, p+1)
		uploadURL := fmt.Sprintf("%s/upload?path=%s&exp=%d&token=%s",
			s.publicBase, url.QueryEscape(path), exp, s.signPath(path, exp))

		uploadInfos[p] = filestorage.UploadInfo{
			UploadURL:  uploadURL,
			ObjectPath: path,
		}
	}

	return uploadInfos, nil
}

// CheckObjectExists проверяет существование файла на диске.
func (s *LocalStorage) CheckObjectExists(ctx context.Context, objectPath string) (bool, error) {
	_, err := os.Stat(filepath.Join(s.dir, filepath.FromSlash(objectPath)))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("ошибка проверки существования объекта: %v", err)
	}
	return true, nil
}

// GetAvatarPublicURL возвращает публичный URL для просмотра аватарки.
func (s *LocalStorage) GetAvatarPublicURL(id string) string {
	format, err := avatarPathFormat(id)
	if err != nil {
		return ""
	}
	return s.GetPublicURLFromPath(fmt.Sprintf(format, id))
}

// GetPublicURLFromPath возвращает публичный URL для объекта по пути.
func (s *LocalStorage) GetPublicURLFromPath(path string) string {
	return fmt.Sprintf("%s/%s", s.publicBase, path)
}

// SetObjectPublicACL — no-op: публичный доступ обеспечивается GET-роутом раздачи.
func (s *LocalStorage) SetObjectPublicACL(ctx context.Context, objectPath string) error {
	return nil
}

// ConfirmUploads подтверждает загрузку массива файлов проверкой их существования.
func (s *LocalStorage) ConfirmUploads(ctx context.Context, paths []string) error {
	for _, path := range paths {
		exists, err := s.CheckObjectExists(ctx, path)
		if err != nil {
			return fmt.Errorf("failed to check existence of %s: %w", path, err)
		}
		if !exists {
			return fmt.Errorf("file %s does not exist", path)
		}
	}
	return nil
}

// BuildPhotoURLs преобразует пути к фото в полные публичные URL.
func (s *LocalStorage) BuildPhotoURLs(paths []string, updatedAt time.Time) []string {
	if len(paths) == 0 {
		return []string{}
	}
	result := make([]string, len(paths))
	for i, path := range paths {
		if path == "" {
			result[i] = ""
		} else {
			result[i] = s.GetPublicURLFromPath(path) + "?t=" + strconv.FormatInt(updatedAt.Unix(), 10)
		}
	}
	return result
}
