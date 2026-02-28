package s3

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain"
	"github.com/aws/smithy-go"
)

// S3Storage представляет собой клиент для работы с S3-совместимым хранилищем
type S3Storage struct {
	cfg struct {
		bucketName string
		expire     time.Duration
		endpoint   string
		region     string
	}
	client *s3.Client
	fs     *domain.MediaService
}

// S3Config содержит конфигурацию для подключения к S3
type S3Config struct {
	endpoint        string
	bucketName      string
	accessKeyID     string
	secretAccessKey string
	region          string
	secure          bool
	expire          time.Duration
}

// S3Builder представляет собой билдер для создания S3Storage с различными настройками
type S3Builder struct {
	config S3Config
	fs     *domain.MediaService
}

// NewS3Storage создает новый билдер для S3Storage
func NewS3Storage(fservice *domain.MediaService) *S3Builder {
	return &S3Builder{
		fs: fservice,
		config: S3Config{
			secure: true,          // По умолчанию используем безопасное соединение
			region: "ru-central1", // По умолчанию регион VK Cloud
		},
	}
}

// WithDefaults настраивает билдер с настройками по умолчанию из переменных окружения
func (b *S3Builder) WithDefaults() *S3Storage {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKeyID := os.Getenv("S3_ACCESS_KEY")
	secretAccessKey := os.Getenv("S3_SECRET_KEY")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	region := os.Getenv("S3_REGION")
	if region == "" {
		region = "ru-central1"
	}

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
	b.config.region = region
	b.config.expire = time.Minute * 10

	return b.build(context.Background())
}

// WithCustoms настраивает билдер с пользовательскими настройками, включая имя бакета и время жизни ссылки
func (b *S3Builder) WithCustoms(endpoint, accessKeyID, secretAccessKey, bucketName, region string, expire time.Duration, secure ...bool) *S3Storage {
	b.config.endpoint = endpoint
	b.config.bucketName = bucketName
	b.config.accessKeyID = accessKeyID
	b.config.secretAccessKey = secretAccessKey
	b.config.region = region
	b.config.expire = expire

	if len(secure) > 0 {
		b.config.secure = secure[0]
	}

	return b.build(context.Background())
}

// build создает и возвращает экземпляр S3Storage на основе текущей конфигурации
func (b *S3Builder) build(ctx context.Context) *S3Storage {
	protocol := "https"
	if !b.config.secure {
		protocol = "http"
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     b.config.accessKeyID,
				SecretAccessKey: b.config.secretAccessKey,
			}, nil
		})),
		config.WithRegion(b.config.region),
	)
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("%s://%s", protocol, b.config.endpoint))
		o.UsePathStyle = true
	})

	storageCfg := struct {
		bucketName string
		expire     time.Duration
		endpoint   string
		region     string
	}{
		bucketName: b.config.bucketName,
		expire:     b.config.expire,
		endpoint:   b.config.endpoint,
		region:     b.config.region,
	}

	return &S3Storage{
		cfg:    storageCfg,
		client: s3Client,
		fs:     b.fs,
	}
}

// Client возвращает s3.Client для прямого доступа к API
func (s *S3Storage) Client() *s3.Client {
	return s.client
}

// FileService возвращает сервис для работы с файлами
func (s *S3Storage) FileService() *domain.MediaService {
	return s.fs
}

// GetPresignedURLs возвращает информацию для загрузки нескольких фотографий
// Возвращает: слайс UploadInfo, error
func (s *S3Storage) GetPresignedURLs(ctx context.Context, count int64, id string) ([]domain.UploadInfo, error) {
	var format string
	var contentType string
	year := time.Now().Year()

	switch {
	case strings.HasPrefix(id, "USR"):
		format = fmt.Sprintf("users/%d/%%s/photos/%%d.jpg", year)
		contentType = "image/jpeg"
	case strings.HasPrefix(id, "PET"):
		format = fmt.Sprintf("pets/%d/%%s/photos/%%d.jpg", year)
		contentType = "image/jpeg"
	case strings.HasPrefix(id, "BLS"):
		format = fmt.Sprintf("blood_requests/%d/%%s/photos/%%d.jpg", year)
		contentType = "image/jpeg"
	default:
		return nil, fmt.Errorf("неподдерживаемый тип файла для id: %s", id)
	}

	presigner := s3.NewPresignClient(s.client)
	uploadInfos := make([]domain.UploadInfo, count)

	for p := range count {
		path := fmt.Sprintf(format, id, p+1)

		req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:      &s.cfg.bucketName,
			Key:         &path,
			ContentType: aws.String(contentType),
		}, s3.WithPresignExpires(s.cfg.expire))
		if err != nil {
			return nil, fmt.Errorf("ошибка создания presigned URL для загрузки %d: %v", p, err)
		}
		uploadInfos[p] = domain.UploadInfo{
			UploadURL:  req.URL,
			ObjectPath: path,
		}
	}

	return uploadInfos, nil
}

// CheckObjectExists проверяет существование объекта в S3
func (s *S3Storage) CheckObjectExists(ctx context.Context, objectPath string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &s.cfg.bucketName,
		Key:    &objectPath,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NotFound" {
			return false, nil
		}
		return false, fmt.Errorf("ошибка проверки существования объекта: %v", err)
	}
	return true, nil
}

// GetAvatarPublicURL возвращает публичный URL для просмотра аватарки
func (s *S3Storage) GetAvatarPublicURL(id string) string {
	var format string
	year := time.Now().Year()
	switch {
	case strings.HasPrefix(id, "USR"):
		format = fmt.Sprintf("users/%d/%%s/avatar.jpg", year)
	case strings.HasPrefix(id, "PET"):
		format = fmt.Sprintf("pets/%d/%%s/avatar.jpg", year)
	case strings.HasPrefix(id, "BLS"):
		format = fmt.Sprintf("blood_requests/%d/%%s/avatar.jpg", year)
	default:
		return ""
	}

	path := fmt.Sprintf(format, id)

	// Встроенная логика GetURL
	protocol := "https"
	if strings.Contains(s.cfg.endpoint, "http://") {
		protocol = "http"
	}

	return fmt.Sprintf("%s://%s/%s/%s",
		protocol,
		s.cfg.endpoint,
		s.cfg.bucketName,
		path)
}

// SetObjectPublicACL устанавливает публичный ACL для объекта
func (s *S3Storage) SetObjectPublicACL(ctx context.Context, objectPath string) error {
	_, err := s.client.PutObjectAcl(ctx, &s3.PutObjectAclInput{
		Bucket: &s.cfg.bucketName,
		Key:    &objectPath,
		ACL:    types.ObjectCannedACLPublicRead,
	})
	return err
}

// ConfirmUploads подтверждает загрузку массива файлов: проверяет существование и устанавливает публичный ACL
func (s *S3Storage) ConfirmUploads(ctx context.Context, paths []string) error {
	for _, path := range paths {
		// Проверить существование файла
		exists, err := s.CheckObjectExists(ctx, path)
		if err != nil {
			return fmt.Errorf("failed to check existence of %s: %w", path, err)
		}
		if !exists {
			return fmt.Errorf("file %s does not exist", path)
		}

		// Установить публичный ACL
		err = s.SetObjectPublicACL(ctx, path)
		if err != nil {
			return fmt.Errorf("failed to set public ACL for %s: %w", path, err)
		}
	}
	return nil
}

// GetPublicURLFromPath возвращает публичный URL для объекта по пути
func (s *S3Storage) GetPublicURLFromPath(path string) string {
	protocol := "https"
	if strings.Contains(s.cfg.endpoint, "http://") {
		protocol = "http"
	}

	return fmt.Sprintf("%s://%s/%s/%s",
		protocol,
		s.cfg.endpoint,
		s.cfg.bucketName,
		path)
}
