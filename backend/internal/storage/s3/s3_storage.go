package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Storage представляет собой клиент для работы с S3-совместимым хранилищем
type S3Storage struct {
	cfg struct {
		bucketName string
		expire     time.Duration
	}
	client *minio.Client
	fs     *services.FileService
}

// S3Config содержит конфигурацию для подключения к S3
type S3Config struct {
	endpoint        string
	bucketName      string
	accessKeyID     string
	secretAccessKey string
	secure          bool
	expire          time.Duration
}

// S3Builder представляет собой билдер для создания S3Storage с различными настройками
type S3Builder struct {
	config S3Config
	fs     *services.FileService
}

// NewS3Storage создает новый билдер для S3Storage
func NewS3Storage(fservice *services.FileService) *S3Builder {
	return &S3Builder{
		fs: fservice,
		config: S3Config{
			secure: true, // По умолчанию используем безопасное соединение
		},
	}
}

// WithDefaults настраивает билдер с настройками по умолчанию из переменных окружения
func (b *S3Builder) WithDefaults() *S3Storage {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKeyID := os.Getenv("S3_ACCESS_KEY")
	secretAccessKey := os.Getenv("S3_SECRET_KEY")
	bucketName := os.Getenv("S3_BUCKET_NAME")

	missingVars := []string{}
	if endpoint == "" {
		missingVars = append(missingVars, "S3_ENDPOINT")
	}
	if accessKeyID == "" {
		missingVars = append(missingVars, "S3_ACCESS_KEY")
	}
	if secretAccessKey == "" {
		missingVars = append(missingVars, "S3_SECRET_KEY")
	}
	if bucketName == "" {
		missingVars = append(missingVars, "S3_BUCKET_NAME")
	}

	if len(missingVars) > 0 {
		panic("Необходимо установить следующие переменные окружения: " + strings.Join(missingVars, ", "))
	}

	b.config.endpoint = endpoint
	b.config.bucketName = bucketName
	b.config.accessKeyID = accessKeyID
	b.config.secretAccessKey = secretAccessKey
	b.config.expire = time.Minute * 10

	return b.build()
}

// WithCustoms настраивает билдер с пользовательскими настройками, включая имя бакета и время жизни ссылки
func (b *S3Builder) WithCustoms(endpoint, accessKeyID, secretAccessKey, bucketName string, expire time.Duration, secure ...bool) *S3Storage {
	b.config.endpoint = endpoint
	b.config.accessKeyID = accessKeyID
	b.config.secretAccessKey = secretAccessKey
	b.config.bucketName = bucketName
	b.config.expire = expire

	if len(secure) > 0 {
		b.config.secure = secure[0]
	}

	return b.build()
}

// build создает и возвращает экземпляр S3Storage на основе текущей конфигурации
func (b *S3Builder) build() *S3Storage {
	creds := credentials.NewStaticV4(b.config.accessKeyID, b.config.secretAccessKey, "")

	minioClient, err := minio.New(b.config.endpoint, &minio.Options{
		Creds:  creds,
		Secure: b.config.secure,
	})
	if err != nil {
		panic(err)
	}

	cfg := struct {
		bucketName string
		expire     time.Duration
	}{
		bucketName: b.config.bucketName,
		expire:     b.config.expire,
	}

	return &S3Storage{
		cfg:    cfg,
		client: minioClient,
		fs:     b.fs,
	}
}

// Client возвращает minio.Client для прямого доступа к API minio
func (s *S3Storage) Client() *minio.Client {
	return s.client
}

// FileService возвращает сервис для работы с файлами
func (s *S3Storage) FileService() *services.FileService {
	return s.fs
}

// Upload загружает файл в S3
func (s *S3Storage) Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.cfg.bucketName, filename, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки файла: %v", err)
	}
	return filename, nil
}

// UploadPublic загружает файл в S3 с публичным доступом
func (s *S3Storage) UploadPublic(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.cfg.bucketName, filename, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	})
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки файла: %v", err)
	}

	return filename, nil
}

// Download скачивает файл из S3
func (s *S3Storage) Download(ctx context.Context, filepath string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.cfg.bucketName, filepath, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("ошибка скачивания файла: %v", err)
	}
	return obj, nil
}

// Delete удаляет файл из S3
func (s *S3Storage) Delete(ctx context.Context, filepath string) error {
	err := s.client.RemoveObject(ctx, s.cfg.bucketName, filepath, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("ошибка удаления файла: %v", err)
	}
	return nil
}

// GetURL возвращает публичный URL к файлу
func (s *S3Storage) GetURL(filepath string) string {
	if s.cfg.bucketName == "" || filepath == "" {
		return ""
	}

	// Формируем публичный URL
	protocol := "https"
	if !strings.Contains(s.client.EndpointURL().String(), "https://") {
		protocol = "http"
	}

	return fmt.Sprintf("%s://%s/%s/%s",
		protocol,
		s.client.EndpointURL().Host,
		s.cfg.bucketName,
		filepath)
}

// GetPresignedURL возвращает временную ссылку с ограниченным сроком действия
func (s *S3Storage) GetPresignedURL(ctx context.Context, filepath string, expire time.Duration) (string, error) {
	if expire == 0 {
		expire = s.cfg.expire
	}

	url, err := s.client.PresignedGetObject(ctx, s.cfg.bucketName, filepath, expire, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания временной ссылки: %v", err)
	}

	return url.String(), nil
}

// GenerateAvatarURL генерирует URL для загрузки аватарки
func (s *S3Storage) GenerateAvatarURL(ctx context.Context, id string) (string, error) {
	var format string
	switch {
	case strings.HasPrefix(id, "USR"):
		format = "users/%s/avatar.jpg"
	case strings.HasPrefix(id, "PET"):
		format = "pets/%s/avatar.jpg"
	default:
		return "", fmt.Errorf("неподдерживаемый тип файла")
	}

	path := fmt.Sprintf(format, id)

	presignedURL, err := s.Client().PresignedPutObject(
		ctx,
		s.cfg.bucketName,
		path,
		s.cfg.expire,
	)

	if err != nil {
		return "", fmt.Errorf("ошибка создания presigned URL для загрузки: %v", err)
	}

	return presignedURL.String(), nil
}

// GetAvatarUploadInfo возвращает информацию для загрузки аватарки
// Возвращает: uploadURL (подписанная ссылка), objectPath (путь в S3), error
func (s *S3Storage) GetAvatarUploadInfo(ctx context.Context, id string) (string, string, error) {
	var format string
	switch {
	case strings.HasPrefix(id, "USR"):
		format = "users/%s/avatar.jpg"
	case strings.HasPrefix(id, "PET"):
		format = "pets/%s/avatar.jpg"
	default:
		return "", "", fmt.Errorf("неподдерживаемый тип файла")
	}

	path := fmt.Sprintf(format, id)

	presignedURL, err := s.Client().PresignedPutObject(
		ctx,
		s.cfg.bucketName,
		path,
		s.cfg.expire,
	)

	if err != nil {
		return "", "", fmt.Errorf("ошибка создания presigned URL для загрузки: %v", err)
	}

	return presignedURL.String(), path, nil
}

// MakeAvatarPublic делает аватарку публичной
func (s *S3Storage) MakeAvatarPublic(ctx context.Context, id string) error {
	var format string
	switch {
	case strings.HasPrefix(id, "USR"):
		format = "users/%s/avatar.jpg"
	case strings.HasPrefix(id, "PET"):
		format = "pets/%s/avatar.jpg"
	default:
		return fmt.Errorf("неподдерживаемый тип файла")
	}

	path := fmt.Sprintf(format, id)
	return s.SetObjectPublic(ctx, path)
}

// CheckObjectExists проверяет существование объекта в S3
func (s *S3Storage) CheckObjectExists(ctx context.Context, objectPath string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.cfg.bucketName, objectPath, minio.StatObjectOptions{})
	if err != nil {
		// Проверяем, является ли ошибка "объект не найден"
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("ошибка проверки существования объекта: %v", err)
	}
	return true, nil
}

// CheckAvatarExists проверяет существование аватарки по ID
func (s *S3Storage) CheckAvatarExists(ctx context.Context, id string) (bool, error) {
	var format string
	switch {
	case strings.HasPrefix(id, "USR"):
		format = "users/%s/avatar.jpg"
	case strings.HasPrefix(id, "PET"):
		format = "pets/%s/avatar.jpg"
	default:
		return false, fmt.Errorf("неподдерживаемый тип файла")
	}

	path := fmt.Sprintf(format, id)
	return s.CheckObjectExists(ctx, path)
}

// GetAvatarPublicURL возвращает публичный URL для просмотра аватарки
func (s *S3Storage) GetAvatarPublicURL(id string) string {
	var format string
	switch {
	case strings.HasPrefix(id, "USR"):
		format = "users/%s/avatar.jpg"
	case strings.HasPrefix(id, "PET"):
		format = "pets/%s/avatar.jpg"
	default:
		return ""
	}

	path := fmt.Sprintf(format, id)
	return s.GetURL(path)
}

// UploadAvatar загружает аватарку с автоматическим определением типа
func (s *S3Storage) UploadAvatar(ctx context.Context, id string, file io.Reader, makePublic bool) (string, error) {
	var format string
	switch {
	case strings.HasPrefix(id, "USR"):
		format = "users/%s/avatar.jpg"
	case strings.HasPrefix(id, "PET"):
		format = "pets/%s/avatar.jpg"
	default:
		return "", fmt.Errorf("неподдерживаемый тип файла")
	}

	path := fmt.Sprintf(format, id)

	if makePublic {
		return s.UploadPublic(ctx, file, path, "image/jpeg")
	}

	return s.Upload(ctx, file, path, "image/jpeg")
}

// SetObjectPublic устанавливает публичный доступ к объекту через ACL
func (s *S3Storage) SetObjectPublic(ctx context.Context, filepath string) error {
	// Копируем объект с новыми ACL настройками
	srcOpts := minio.CopySrcOptions{
		Bucket: s.cfg.bucketName,
		Object: filepath,
	}

	dstOpts := minio.CopyDestOptions{
		Bucket: s.cfg.bucketName,
		Object: filepath,
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	}

	_, err := s.client.CopyObject(ctx, dstOpts, srcOpts)
	if err != nil {
		return fmt.Errorf("ошибка установки публичного доступа: %v", err)
	}

	return nil
}

// SetObjectPrivate устанавливает приватный доступ к объекту через ACL
func (s *S3Storage) SetObjectPrivate(ctx context.Context, filepath string) error {
	// Копируем объект с приватными ACL настройками
	srcOpts := minio.CopySrcOptions{
		Bucket: s.cfg.bucketName,
		Object: filepath,
	}

	dstOpts := minio.CopyDestOptions{
		Bucket: s.cfg.bucketName,
		Object: filepath,
		UserMetadata: map[string]string{
			"x-amz-acl": "private",
		},
	}

	_, err := s.client.CopyObject(ctx, dstOpts, srcOpts)
	if err != nil {
		return fmt.Errorf("ошибка установки приватного доступа: %v", err)
	}

	return nil
}
