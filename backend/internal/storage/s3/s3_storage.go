package s3

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Storage представляет собой клиент для работы с S3-совместимым хранилищем
type S3Storage struct {
	client *minio.Client
	fs     *services.FileService
}

// S3Config содержит конфигурацию для подключения к S3
type S3Config struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	secure          bool
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

	if len(missingVars) > 0 {
		panic("Необходимо установить следующие переменные окружения: " + strings.Join(missingVars, ", "))
	}

	b.config.endpoint = endpoint
	b.config.accessKeyID = accessKeyID
	b.config.secretAccessKey = secretAccessKey

	return b.build()
}

// WithCustoms настраивает билдер с пользовательскими настройками
func (b *S3Builder) WithCustoms(endpoint, accessKeyID, secretAccessKey string, secure ...bool) *S3Storage {
	b.config.endpoint = endpoint
	b.config.accessKeyID = accessKeyID
	b.config.secretAccessKey = secretAccessKey

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

	return &S3Storage{
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

func (s *S3Storage) Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// Заглушка: возвращаем пустую строку и nil-ошибку
	return "", nil
}

func (s *S3Storage) Download(ctx context.Context, filepath string) (io.ReadCloser, error) {
	// Заглушка: возвращаем nil и nil-ошибку
	return nil, nil
}

func (s *S3Storage) Delete(ctx context.Context, filepath string) error {
	// Заглушка: возвращаем nil-ошибку
	return nil
}

func (s *S3Storage) GetURL(filepath string) string {
	// Заглушка: возвращаем пустую строку
	return ""
}
