package s3_test

// import (
// 	"context"
// 	"net/url"
// 	"testing"
// 	"time"

// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/storage/s3"
// 	"github.com/joho/godotenv"
// 	"github.com/minio/minio-go/v7"
// )

// // Пример использования curl для загрузки файла:
// //
// //	curl -X PUT \
// //	  -H "Content-Type: image/jpeg" \
// //	  --upload-file /path/to/your/file.jpg \
// //	  "PRESIGNED_URL"
// func TestUploadFileToS3(t *testing.T) {
// 	godotenv.Load("../../../.env")
// 	// 1. Создаем S3 хранилище с настройками по умолчанию из переменных окружения
// 	s3Storage := s3.NewS3Storage(nil).WithDefaults()

// 	// 2. Создаем контекст с таймаутом
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	// 3. Загружаем фотку
// 	bucketName := "odnoi-krovi-photos"                                                                      // Имя бакета
// 	objectName := "sam.jpg"                                                                                 // Имя файла в S3
// 	filePath := "/Users/ruslan/Documents/zed-projects/go-projects/odnoi-krovi-app/backend/tests/s3/Sam.jpg" // Полный путь к файлу на компьютере
// 	contentType := "image/jpeg"                                                                             // MIME-тип фотки

// 	// Загружаем файл
// 	info, err := s3Storage.Client().FPutObject(
// 		ctx,
// 		bucketName,
// 		objectName,
// 		filePath,
// 		minio.PutObjectOptions{
// 			ContentType: contentType,
// 		},
// 	)

// 	if err != nil {
// 		t.Errorf("Ошибка загрузки: %v", err)
// 		return
// 	}

// 	t.Logf("Фотка успешно загружена!")
// 	t.Logf("Имя файла: %s", objectName)
// 	t.Logf("Размер: %d байт", info.Size)
// }

// // Пример использования curl для скачивания файла:
// // curl -o sam-download.jpg "PRESIGNED_URL"
// //
// // Пример curl команды для проверки заголовков:
// // curl -I "PRESIGNED_URL"
// func TestPresignedGetObject(t *testing.T) {
// 	godotenv.Load("../../../.env")
// 	s3Storage := s3.NewS3Storage(nil).WithDefaults()
// 	ctx := context.Background()

// 	bucketName := "odnoi-krovi-photos"
// 	objectName := "sam.jpg"

// 	// Проверяем, существует ли файл
// 	_, err := s3Storage.Client().StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
// 	if err != nil {
// 		t.Skipf("Файл %s не существует в бакете %s, пропускаем тест", objectName, bucketName)
// 		return
// 	}

// 	// Создаем presigned URL для скачивания файла на 15 минут
// 	expiry := 15 * time.Minute
// 	reqParams := make(url.Values)

// 	// Можно добавить дополнительные параметры, например, для скачивания с определенным именем файла
// 	reqParams.Set("response-content-disposition", "attachment; filename=\"sam-download.jpg\"")

// 	presignedURL, err := s3Storage.Client().PresignedGetObject(
// 		ctx,
// 		bucketName,
// 		objectName,
// 		expiry,
// 		reqParams,
// 	)

// 	if err != nil {
// 		t.Errorf("Ошибка создания presigned URL: %v", err)
// 		return
// 	}

// 	t.Logf("Presigned URL для скачивания создан!")
// 	t.Logf("URL: %s", presignedURL.String())
// 	t.Logf("Срок действия: %v", expiry)

// 	// Пример curl команды для скачивания файла
// 	t.Logf("\nПример curl команды для скачивания файла:")
// 	t.Logf("curl -o sam-download.jpg \"%s\"", presignedURL.String())

// 	// Пример curl команды для проверки заголовков
// 	t.Logf("\nПример curl команды для проверки заголовков:")
// 	t.Logf("curl -I \"%s\"", presignedURL.String())
// }

// // Пример использования curl для загрузки файла:
// //
// //	curl -X PUT \
// //	  -H "Content-Type: image/jpeg" \
// //	  --upload-file /path/to/your/file.jpg \
// //	  "PRESIGNED_URL"
// func TestPresignedPutObject(t *testing.T) {
// 	godotenv.Load("../../../.env")                   // Загружаем переменные окружения из файла .env
// 	s3Storage := s3.NewS3Storage(nil).WithDefaults() // Создаем клиент S3 с настройками по умолчанию
// 	ctx := context.Background()                      // Создаем базовый контекст

// 	bucketName := "odnoi-krovi-photos"     // Имя бакета в S3
// 	objectName := "uploaded-test-file.jpg" // Имя файла, который будет загружен

// 	// Создаем presigned URL для загрузки файла на 30 минут
// 	expiry := 30 * time.Minute

// 	presignedURL, err := s3Storage.Client().PresignedPutObject(
// 		ctx,
// 		bucketName,
// 		objectName,
// 		expiry,
// 	)

// 	if err != nil {
// 		t.Errorf("Ошибка создания presigned URL для загрузки: %v", err) // Логируем ошибку при создании URL
// 		return                                                          // Прекращаем выполнение теста при ошибке
// 	}

// 	// Логируем успешное создание presigned URL
// 	t.Logf("Presigned URL для загрузки создан!")
// 	t.Logf("URL: %s", presignedURL.String()) // Выводим сам URL
// 	t.Logf("Срок действия: %v", expiry)      // Выводим срок действия ссылки
// 	t.Logf("Бакет: %s", bucketName)          // Выводим имя бакета
// 	t.Logf("Имя файла в S3: %s", objectName) // Выводим имя файла внутри бакета
// }

// // Пример использования curl для получения метаданных файла:
// // curl -I "PRESIGNED_URL"
// //
// // Пример с подробным выводом:
// // curl -v -I "PRESIGNED_URL"
// func TestPresignedHeadObject(t *testing.T) {
// 	godotenv.Load("../../../.env")
// 	s3Storage := s3.NewS3Storage(nil).WithDefaults()
// 	ctx := context.Background()

// 	bucketName := "odnoi-krovi-photos"
// 	objectName := "sam.jpg"

// 	// Проверяем, существует ли файл
// 	_, err := s3Storage.Client().StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
// 	if err != nil {
// 		t.Skipf("Файл %s не существует в бакете %s, пропускаем тест", objectName, bucketName)
// 		return
// 	}

// 	// Создаем presigned URL для получения метаданных на 10 минут
// 	expiry := 10 * time.Minute
// 	reqParams := make(url.Values)

// 	presignedURL, err := s3Storage.Client().PresignedHeadObject(
// 		ctx,
// 		bucketName,
// 		objectName,
// 		expiry,
// 		reqParams,
// 	)

// 	if err != nil {
// 		t.Errorf("Ошибка создания presigned URL для метаданных: %v", err)
// 		return
// 	}

// 	t.Logf("Presigned URL для метаданных создан!")
// 	t.Logf("URL: %s", presignedURL.String())
// 	t.Logf("Срок действия: %v", expiry)

// 	// Пример curl команды для получения метаданных
// 	t.Logf("\nПример curl команды для получения метаданных:")
// 	t.Logf("curl -I \"%s\"", presignedURL.String())

// 	// Пример с дополнительными заголовками
// 	t.Logf("\nПример с подробным выводом:")
// 	t.Logf("curl -v -I \"%s\"", presignedURL.String())
// }

// // Пример использования curl для загрузки через POST форму:
// //
// //	curl -X POST \
// //	  -F "key=..." \
// //	  -F "policy=..." \
// //	  -F "x-amz-signature=..." \
// //	  -F "file=@/path/to/your/file.jpg" \
// //	  "POST_URL"
// func TestPresignedPostPolicy(t *testing.T) {
// 	godotenv.Load("../../../.env")
// 	s3Storage := s3.NewS3Storage(nil).WithDefaults()
// 	ctx := context.Background()

// 	bucketName := "odnoi-krovi-photos"
// 	objectName := "form-uploaded-file.jpg"

// 	// Создаем политику POST
// 	policy := minio.NewPostPolicy()

// 	// Устанавливаем параметры политики
// 	policy.SetBucket(bucketName)
// 	policy.SetKey(objectName)

// 	// Срок действия - 1 час
// 	policy.SetExpires(time.Now().UTC().Add(1 * time.Hour))

// 	// Можно установить дополнительные условия
// 	// policy.SetContentType("image/jpeg")
// 	// policy.SetContentLengthRange(1024, 10485760) // от 1KB до 10MB

// 	// Генерируем URL и данные формы
// 	postURL, formData, err := s3Storage.Client().PresignedPostPolicy(ctx, policy)
// 	if err != nil {
// 		t.Errorf("Ошибка создания POST политики: %v", err)
// 		return
// 	}

// 	t.Logf("POST политика создана!")
// 	t.Logf("URL для загрузки: %s", postURL.String())

// 	// Выводим данные формы
// 	t.Logf("\nДанные формы:")
// 	for key, value := range formData {
// 		t.Logf("  %s: %s", key, value)
// 	}
// }

// // Пример использования curl для загрузки файла:
// // echo 'This is a test file for presigned URL workflow' > test-file.txt
// //
// //	curl -X PUT \
// //	  -H "Content-Type: text/plain" \
// //	  --data-binary "This is a test file for presigned URL workflow" \
// //	  "UPLOAD_URL"
// //
// // Пример использования curl для скачивания файла:
// // curl "DOWNLOAD_URL" -o downloaded-file.txt
// //
// // Пример использования curl для получения метаданных:
// // curl -I "HEAD_URL"
// func TestCompletePresignedWorkflow(t *testing.T) {
// 	godotenv.Load("../../../.env")
// 	s3Storage := s3.NewS3Storage(nil).WithDefaults()
// 	ctx := context.Background()

// 	bucketName := "odnoi-krovi-photos"
// 	objectName := "workflow-test-file.txt"
// 	testContent := "This is a test file for presigned URL workflow\n"

// 	// Шаг 1: Создаем presigned URL для загрузки
// 	t.Logf("=== Шаг 1: Создание presigned URL для загрузки ===")
// 	uploadURL, err := s3Storage.Client().PresignedPutObject(
// 		ctx,
// 		bucketName,
// 		objectName,
// 		15*time.Minute,
// 	)
// 	if err != nil {
// 		t.Errorf("Ошибка создания upload URL: %v", err)
// 		return
// 	}
// 	t.Logf("Upload URL: %s", uploadURL.String())

// 	// Пример curl для загрузки (закомментирован, так как требует реального файла)
// 	t.Logf("\nПример curl для загрузки файла:")
// 	t.Logf("echo '%s' > test-file.txt", testContent)
// 	t.Logf("curl -X PUT \\")
// 	t.Logf("  -H \"Content-Type: text/plain\" \\")
// 	t.Logf("  --data-binary \"%s\" \\", testContent)
// 	t.Logf("  \"%s\"", uploadURL.String())

// 	// Шаг 2: Создаем presigned URL для скачивания
// 	t.Logf("\n=== Шаг 2: Создание presigned URL для скачивания ===")
// 	downloadURL, err := s3Storage.Client().PresignedGetObject(
// 		ctx,
// 		bucketName,
// 		objectName,
// 		15*time.Minute,
// 		nil,
// 	)
// 	if err != nil {
// 		t.Errorf("Ошибка создания download URL: %v", err)
// 		return
// 	}
// 	t.Logf("Download URL: %s", downloadURL.String())

// 	t.Logf("\nПример curl для скачивания файла:")
// 	t.Logf("curl \"%s\" -o downloaded-file.txt", downloadURL.String())

// 	// Шаг 3: Создаем presigned URL для метаданных
// 	t.Logf("\n=== Шаг 3: Создание presigned URL для метаданных ===")
// 	headURL, err := s3Storage.Client().PresignedHeadObject(
// 		ctx,
// 		bucketName,
// 		objectName,
// 		15*time.Minute,
// 		nil,
// 	)
// 	if err != nil {
// 		t.Errorf("Ошибка создания head URL: %v", err)
// 		return
// 	}
// 	t.Logf("Head URL: %s", headURL.String())

// 	t.Logf("\nПример curl для получения метаданных:")
// 	t.Logf("curl -I \"%s\"", headURL.String())

// 	// Шаг 4: Очистка (удаление тестового файла)
// 	t.Logf("\n=== Шаг 4: Очистка тестового файла ===")
// 	err = s3Storage.Client().RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
// 	if err != nil {
// 		t.Logf("Примечание: не удалось удалить тестовый файл: %v", err)
// 	} else {
// 		t.Logf("Тестовый файл удален")
// 	}
// }
