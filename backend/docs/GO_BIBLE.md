# 📖 Go Bible — Библия разработчика на Golang

Собрание ключевых идиом, паттернов и лучших практик языка Go. Только то, что касается самого языка — без внешних библиотек.

---

## Содержание

1. [Стиль кода и поток управления](#1-стиль-кода-и-поток-управления)
2. [Соглашения об именовании](#2-соглашения-об-именовании)
3. [Обработка ошибок](#3-обработка-ошибок)
4. [Структуры и интерфейсы](#4-структуры-и-интерфейсы)
5. [Context](#5-context)
6. [Конкурентность](#6-конкурентность)
7. [Структуры данных](#7-структуры-данных)
8. [Безопасность и защитное программирование](#8-безопасность-и-защитное-программирование)
9. [Тестирование](#9-тестирование)
10. [Паттерны проектирования](#10-паттерны-проектирования)

---

## 1. Стиль кода и поток управления

### 1.1 Ранний возврат — держи код плоским

> **Золотое правило:** ошибки и крайние случаи обрабатывай первыми. Счастливый путь должен быть на минимальном уровне вложенности.

```go
// ❌ Плохо — лесенка из вложенных if
func process(data []byte) (*Result, error) {
    if data != nil {
        parsed, err := parse(data)
        if err == nil {
            result, err := transform(parsed)
            if err == nil {
                return result, nil
            } else {
                return nil, err
            }
        } else {
            return nil, err
        }
    } else {
        return nil, errors.New("empty data")
    }
}

// ✅ Хорошо — ранний возврат, код плоский
func process(data []byte) (*Result, error) {
    if len(data) == 0 {
        return nil, errors.New("empty data")
    }

    parsed, err := parse(data)
    if err != nil {
        return nil, fmt.Errorf("parsing: %w", err)
    }

    return transform(parsed)
}
```

### 1.2 Убирай ненужный `else`

Если тело `if` заканчивается на `return`/`break`/`continue` — `else` не нужен.

```go
// ❌ Плохо
if user.IsAdmin {
    return fullAccess
} else {
    return limitedAccess
}

// ✅ Хорошо
if user.IsAdmin {
    return fullAccess
}
return limitedAccess
```

### 1.3 Default-then-override вместо цепочек else-if

```go
// ❌ Плохо — цепочка else-if скрывает дефолтное значение
var level slog.Level
if debug {
    level = slog.LevelDebug
} else if verbose {
    level = slog.LevelWarn
} else {
    level = slog.LevelInfo
}

// ✅ Хорошо — дефолт задан явно, switch для переопределения
level := slog.LevelInfo
switch {
case debug:
    level = slog.LevelDebug
case verbose:
    level = slog.LevelWarn
}
```

### 1.4 Сложные условия — выноси в именованные булевы переменные

```go
// ❌ Плохо — стена из || нечитаема
if user.Role == RoleAdmin || resource.OwnerID == user.ID ||
   (resource.IsPublic && user.IsVerified) || permissions.Contains(PermOverride) {
    allow()
}

// ✅ Хорошо — имена раскрывают смысл
isAdmin := user.Role == RoleAdmin
isOwner := resource.OwnerID == user.ID
isPublicVerified := resource.IsPublic && user.IsVerified

if isAdmin || isOwner || isPublicVerified || permissions.Contains(PermOverride) {
    allow()
}
```

### 1.5 Switch вместо цепочек if-else

```go
// ❌ Плохо
if status == StatusActive {
    activate()
} else if status == StatusInactive {
    deactivate()
} else if status == StatusPending {
    wait()
}

// ✅ Хорошо
switch status {
case StatusActive:
    activate()
case StatusInactive:
    deactivate()
case StatusPending:
    wait()
default:
    panic(fmt.Sprintf("unexpected status: %d", status))
}
```

### 1.6 Область видимости переменных через `if`

```go
// Переменная живёт только внутри блока if
if err := validate(input); err != nil {
    return err
}
// err здесь уже не существует
```

### 1.7 Функции: короткие и сфокусированные

- **Одна функция — одна задача.**
- **Не больше 4 параметров.** Если больше — используй структуру-конфиг.
- **Порядок параметров:** `context.Context` первым, потом входные данные, потом выходные.

```go
func FetchUser(ctx context.Context, id string) (*User, error)
func SendEmail(ctx context.Context, msg EmailMessage) error
```

### 1.8 Предпочитай `range` вместо индексов

```go
// ✅ Хорошо
for _, user := range users {
    process(user)
}

// Go 1.22+ — range по числу
for i := range 10 {
    fmt.Println(i)
}
```

---

## 2. Соглашения об именовании

### 2.1 MixedCaps — и только так

Go использует **MixedCaps** (или **mixedCaps**). Никаких `snake_case` и `ALL_CAPS`.

```go
// ✅ Правильно
MaxPacketSize
userCount
parseHTTPResponse

// ❌ Неправильно
MAX_PACKET_SIZE   // C/Python стиль
max_packet_size   // snake_case
```

### 2.2 Капитализация управляет видимостью

- `UpperCamelCase` — экспортируется (публичное)
- `lowerCamelCase` — не экспортируется (приватное)

### 2.3 Избегай заикания (stuttering)

Имя пакета всегда присутствует в месте вызова — не повторяй его в имени типа.

```go
// ✅ Хорошо
http.Client       // не http.HTTPClient
json.Decoder      // не json.JSONDecoder
user.New()        // не user.NewUser()

// В пакете dbpool:
type Pool struct{}   // не DBPool
```

### 2.4 Шпаргалка по именованию

| Элемент | Соглашение | Пример |
|---|---|---|
| Пакет | lowercase, одно слово | `json`, `http` |
| Файл | lowercase, можно `_` | `user_handler.go` |
| Интерфейс | Метод + `-er` | `Reader`, `Closer`, `Stringer` |
| Структура | MixedCaps, существительное | `Request`, `FileHeader` |
| Константа | MixedCaps | `MaxRetries`, `defaultTimeout` |
| Ошибка (переменная) | `Err` префикс | `ErrNotFound`, `ErrTimeout` |
| Ошибка (тип) | `Error` суффикс | `PathError`, `SyntaxError` |
| Конструктор | `New` или `NewTypeName` | `ring.New`, `http.NewRequest` |
| Булево поле | `is`/`has`/`can` префикс | `isReady`, `hasPermission` |
| Акроним | Все заглавные или все строчные | `URL`, `HTTPServer`, `xmlParser` |
| Функция-опция | `With` + имя поля | `WithPort()`, `WithLogger()` |
| Enum (iota) | Префикс типа, 0 = unknown | `StatusUnknown`, `StatusReady` |

### 2.5 Частые ошибки именования

```go
// ❌ ALL_CAPS константы
const MAX_RETRIES = 3

// ✅ MixedCaps
const MaxRetries = 3

// ❌ Get-префикс у геттера
func (u *User) GetName() string

// ✅ Без Get
func (u *User) Name() string

// ❌ Булево поле без префикса
type Config struct {
    connected bool
}

// ✅ С префиксом
type Config struct {
    isConnected bool
}

// ❌ Acronym с混合ным регистром
type HttpClient struct{}

// ✅ Все заглавные
type HTTPClient struct{}
```

---

## 3. Обработка ошибок

### 3.1 Ошибки — это значения

В Go ошибки возвращают, а не бросают. **Никогда не игнорируй возвращённую ошибку.**

```go
// ❌ Плохо — ошибка проигнорирована
data, _ := os.ReadFile("config.json")

// ✅ Хорошо
data, err := os.ReadFile("config.json")
if err != nil {
    return fmt.Errorf("reading config: %w", err)
}
```

### 3.2 Оборачивай ошибки с контекстом

Используй `%w` для создания цепочки ошибок, `%v` — только на границах системы.

```go
func GetUser(ctx context.Context, id string) (*User, error) {
    user, err := db.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("getting user %s: %w", id, err)
    }
    return user, nil
}
```

### 3.3 Строки ошибок: lowercase, без пунктуации

```go
// ✅ Правильно
errors.New("user not found")
fmt.Errorf("parsing config: %w", err)

// ❌ Неправильно
errors.New("User not found.")
fmt.Errorf("Parsing config: %w.", err)
```

### 3.4 Проверка ошибок: `errors.Is` и `errors.As`

```go
// Проверка на конкретную sentinel-ошибку
if errors.Is(err, ErrNotFound) {
    // обработать "не найдено"
}

// Извлечение типизированной ошибки
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Проблемный путь:", pathErr.Path)
}
```

### 3.5 Правило одного обработчика

**Ошибка либо логируется, либо возвращается — но никогда и то и другое одновременно.**

```go
// ❌ Плохо — двойное логирование
func process(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        log.Printf("failed to read %s: %v", path, err) // логируем
        return err                                      // и возвращаем
    }
    // ...
}

// ✅ Хорошо — возвращаем с контекстом, логирует верхний уровень
func process(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("reading %s: %w", path, err)
    }
    // ...
}
```

### 3.6 Sentinel-ошибки и кастомные типы

```go
// Sentinel-ошибка (для проверки через errors.Is)
var ErrNotFound = errors.New("user not found")

// Кастомный тип (для переноса данных)
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s: %s", e.Field, e.Message)
}
```

### 3.7 `errors.Join` для нескольких ошибок

```go
var errs []error
for _, item := range items {
    if err := validate(item); err != nil {
        errs = append(errs, err)
    }
}
if len(errs) > 0 {
    return errors.Join(errs...)
}
```

### 3.8 Panic — только для невосстановимых багов

```go
// ✅ Паника уместна — нарушен инвариант программы
func mustParse(input string) *Config {
    cfg, err := parse(input)
    if err != nil {
        panic("invalid config: " + err.Error())
    }
    return cfg
}

// ❌ Паника неуместна — сетевой сбой
func fetch(url string) []byte {
    resp, err := http.Get(url)
    if err != nil {
        panic(err) // так нельзя! Верни ошибку
    }
    // ...
}
```

---

## 4. Структуры и интерфейсы

### 4.1 Маленькие интерфейсы

> «Чем больше интерфейс, тем слабее абстракция.»

Интерфейсы должны иметь 1-3 метода. Большие контракты — компоновать из маленьких.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Композиция
type ReadWriter interface {
    Reader
    Writer
}
```

### 4.2 Интерфейс определяет потребитель, а не реализация

```go
// ✅ Интерфейс в пакете-потребителе (notification)
type Sender interface {
    Send(to, body string) error
}

type Service struct {
    sender Sender
}

// Пакет email экспортирует конкретную структуру Client
// Ему не нужно знать про интерфейс Sender
```

### 4.3 Принимай интерфейсы, возвращай структуры

```go
// ✅ Хорошо
func NewService(store UserStore) *Service { ... }

// ❌ Плохо — никогда не возвращай интерфейс из конструктора
func NewService(store UserStore) ServiceInterface { ... }
```

### 4.4 Не создавай интерфейсы преждевременно

Начинай с конкретного типа. Выделяй интерфейс когда появится **вторая реализация** или **потребность в моке для тестов**.

```go
// ❌ Преждевременный интерфейс
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
}
type userRepository struct { db *sql.DB }

// ✅ Начни с конкретного типа
type UserRepository struct { db *sql.DB }
```

### 4.5 Нулевое значение должно быть полезным

```go
// ✅ Нулевое значение готово к использованию
var buf bytes.Buffer
buf.WriteString("hello")

var mu sync.Mutex
mu.Lock()

// ❌ Нулевое значение сломано — nil map паникует при записи
type Registry struct {
    items map[string]Item
}

// ✅ Исправление — ленивая инициализация
func (r *Registry) Register(name string, item Item) {
    if r.items == nil {
        r.items = make(map[string]Item)
    }
    r.items[name] = item
}
```

### 4.6 Проверка интерфейса на этапе компиляции

```go
var _ io.ReadWriter = (*MyBuffer)(nil)
```

Если `MyBuffer` перестанет удовлетворять `io.ReadWriter` — код не скомпилируется.

### 4.7 Безопасное приведение типов

```go
// ✅ Безопасно — comma-ok
s, ok := val.(string)
if !ok {
    // обработать несовпадение типа
}

// ❌ Опасно — паника при несовпадении
s := val.(string)
```

### 4.8 Type Switch

```go
switch v := val.(type) {
case string:
    fmt.Println("строка:", v)
case int:
    fmt.Println("число:", v)
case io.Reader:
    io.Copy(os.Stdout, v)
default:
    fmt.Printf("неизвестный тип %T\n", v)
}
```

### 4.9 Опциональное поведение через type assertion

```go
// Проверяем, поддерживает ли Writer интерфейс Flusher
type Flusher interface {
    Flush() error
}

func writeData(w io.Writer, data []byte) error {
    if _, err := w.Write(data); err != nil {
        return err
    }
    if f, ok := w.(Flusher); ok {
        return f.Flush()
    }
    return nil
}
```

### 4.10 Встраивание (embedding) — композиция, не наследование

```go
type Logger struct {
    *slog.Logger
}

type Server struct {
    Logger          // Продвигает методы Logger
    addr    string
}

s := Server{Logger: Logger{slog.Default()}, addr: ":8080"}
s.Info("starting", "addr", s.addr) // метод продвинут из Logger
```

### 4.11 Pointer vs Value receiver

| Pointer `(s *Server)` | Value `(s Server)` |
|---|---|
| Метод меняет получателя | Получатель маленький и неизменяемый |
| Получатель содержит `sync.Mutex` | Получатель — базовый тип (int, string) |
| Получатель — большая структура | Метод только читает данные |

**Правило:** если хоть один метод использует pointer receiver — **все** методы должны использовать pointer receiver.

### 4.12 Структурные теги

```go
type Order struct {
    ID        string    `json:"id"         db:"id"`
    UserID    string    `json:"user_id"    db:"user_id"`
    Total     float64   `json:"total"      db:"total"`
    Items     []Item    `json:"items"      db:"-"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    Internal  string    `json:"-"          db:"-"`
}
```

| Директива | Значение |
|---|---|
| `json:"name"` | Имя поля в JSON |
| `json:"name,omitempty"` | Опустить если zero-value |
| `json:"-"` | Исключить из JSON |
| `db:"column"` | Колонка в БД |

---

## 5. Context

### 5.1 Главное правило: пробрасывай контекст

Контекст должен проходить через **весь** цикл запроса: HTTP-обработчик → сервис → БД → внешние API.

```go
// ❌ Плохо — создаёт новый контекст, разрывая цепочку
func (s *OrderService) Create(ctx context.Context, order Order) error {
    return s.db.ExecContext(context.Background(), "INSERT ...", order.ID)
}

// ✅ Хорошо — пробрасывает контекст вызывающего кода
func (s *OrderService) Create(ctx context.Context, order Order) error {
    return s.db.ExecContext(ctx, "INSERT ...", order.ID)
}
```

### 5.2 `ctx` всегда первый параметр

```go
func FetchUser(ctx context.Context, id string) (*User, error)
```

### 5.3 Когда что использовать

| Ситуация | Используй |
|---|---|
| Точка входа (main, init, тест) | `context.Background()` |
| Нужен контекст, но вызывающий код его не даёт | `context.TODO()` |
| Внутри HTTP-обработчика | `r.Context()` |
| Нужен таймаут | `context.WithTimeout(parentCtx, duration)` |
| Нужна ручная отмена | `context.WithCancel(parentCtx)` |

### 5.4 Всегда вызывай cancel

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // освобождает ресурсы
```

### 5.5 Context values — только для request-scoped метаданных

```go
// Ключ должен быть неэкспортируемым типом (чтобы избежать коллизий)
type contextKey string

const requestIDKey contextKey = "requestID"

// Запись
ctx = context.WithValue(ctx, requestIDKey, "abc-123")

// Чтение
requestID, ok := ctx.Value(requestIDKey).(string)
```

### 5.6 `context.WithoutCancel` для фоновой работы

```go
// Операция, которая должна пережить родительский запрос
// (например, аудит-лог)
go func() {
    ctx := context.WithoutCancel(r.Context())
    auditLog(ctx, event)
}()
```

---

## 6. Конкурентность

### 6.1 Каждая горутина должна иметь чёткий выход

```go
// ❌ Плохо — горутина-сирота, нет механизма остановки
go func() {
    for {
        doWork()
        time.Sleep(time.Second)
    }
}()

// ✅ Хорошо — контекст для остановки
go func() {
    for {
        select {
        case <-ctx.Done():
            return
        case <-time.After(time.Second):
            doWork()
        }
    }
}()
```

### 6.2 Каналы: отправляй копии, не указатели

```go
// ❌ Плохо — указатель создаёт неявную общую память
ch <- &result

// ✅ Хорошо — копия значения
ch <- result
```

### 6.3 Канал закрывает только отправитель

```go
func producer(out chan<- int) {
    defer close(out) // отправитель закрывает
    for i := 0; i < 10; i++ {
        out <- i
    }
}
```

### 6.4 Указывай направление канала

```go
func producer(out chan<- int) { ... }  // только отправка
func consumer(in <-chan int) { ... }   // только получение
```

### 6.5 Всегда включай `ctx.Done()` в select

```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-workCh:
    process(result)
}
```

### 6.6 Что выбрать: канал, мьютекс или atomic?

| Сценарий | Инструмент |
|---|---|
| Передача данных между горутинами | Канал |
| Защита полей структуры | `sync.Mutex` / `sync.RWMutex` |
| Простые счётчики, флаги | `sync/atomic` |
| Много читателей, мало писателей (map) | `sync.Map` |
| Кэширование дорогих вычислений | `sync.Once` / `singleflight` |

### 6.7 WaitGroup vs errgroup

```go
// Просто ждём завершения — sync.WaitGroup
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func(item Item) {
        defer wg.Done()
        process(item)
    }(item)
}
wg.Wait()

// Нужны ошибки и отмена — errgroup
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(5) // максимум 5 одновременных горутин

for _, item := range items {
    item := item
    g.Go(func() error {
        return process(ctx, item)
    })
}

if err := g.Wait(); err != nil {
    // первая ошибка
}
```

### 6.8 Чеклист перед запуском горутины

- [ ] Как она завершится? (контекст, канал, сигнал)
- [ ] Могу ли я её остановить?
- [ ] Могу ли я дождаться её завершения? (WaitGroup/errgroup)
- [ ] Кто владеет каналами?
- [ ] Может, синхронный код здесь достаточен?

### 6.9 Типичные ошибки

```go
// ❌ time.After в цикле — утечка памяти
for {
    select {
    case <-time.After(time.Second): // каждый вызов создаёт новый таймер
        doWork()
    }
}

// ✅ Переиспользуй таймер
timer := time.NewTimer(time.Second)
for {
    select {
    case <-timer.C:
        doWork()
        timer.Reset(time.Second)
    }
}

// ❌ wg.Add внутри горутины — Wait может проскочить
go func() {
    wg.Add(1) // поздно!
    defer wg.Done()
}()

// ✅ wg.Add до go
wg.Add(1)
go func() {
    defer wg.Done()
}()
```

---

## 7. Структуры данных

### 7.1 Слайс: всегда преаллоцируй если знаешь размер

```go
// ❌ Плохо — многократный рост и копирование
var users []User
for _, id := range ids {
    users = append(users, fetchUser(id))
}

// ✅ Хорошо — один аллокейт
users := make([]User, 0, len(ids))
for _, id := range ids {
    users = append(users, fetchUser(id))
}
```

### 7.2 Слайс: инициализируй явно, не nil

```go
// ❌ nil-слайс сериализуется в JSON как null
var users []User

// ✅ Пустой слайс сериализуется как []
users := []User{}
// или
users := make([]User, 0)
```

### 7.3 Map: всегда инициализируй перед записью

```go
// ❌ Паника: запись в nil map
var m map[string]int
m["key"] = 1

// ✅ Инициализация
m := make(map[string]int)
// или с преаллокацией
m := make(map[string]int, 100)
```

### 7.4 Осторожно с append — слайсы делят память

```go
a := make([]int, 3, 5) // len=3, cap=5
b := append(a, 4)       // b делит backing array с a
b[0] = 99               // a[0] тоже стал 99!

// ✅ Защита: full slice expression
b := append(a[:len(a):len(a)], 4) // принудительное копирование
```

### 7.5 strings.Builder для конкатенации в циклах

```go
// ❌ Плохо — каждая конкатенация аллоцирует новую строку
var result string
for _, s := range parts {
    result += s
}

// ✅ Хорошо
var builder strings.Builder
builder.Grow(totalSize) // опциональная преаллокация
for _, s := range parts {
    builder.WriteString(s)
}
result := builder.String()
```

### 7.6 bytes.Buffer vs strings.Builder

| | `strings.Builder` | `bytes.Buffer` |
|---|---|---|
| Для чего | Только сборка строк | Двунаправленный I/O |
| `String()` | Без копирования | Копирует байты |
| Интерфейсы | `io.Writer` | `io.Reader` + `io.Writer` |

### 7.7 Семантика копирования

| Тип | Поведение при копировании |
|---|---|
| `int`, `float`, `bool`, `string` | Полная копия (независима) |
| `array`, `struct` | Полная копия |
| `slice` | Копируется заголовок, backing array общий |
| `map` | Копируется ссылка (те же данные) |
| `channel` | Копируется ссылка (тот же канал) |
| `*T` (указатель) | Копируется адрес (те же данные) |

---

## 8. Безопасность и защитное программирование

### 8.1 Ловушка nil interface

Интерфейс хранит пару `(тип, значение)`. Он равен `nil` только когда **оба** nil.

```go
// ❌ Опасно — возвращается не-nil интерфейс
func getHandler() http.Handler {
    var h *MyHandler // nil-указатель
    if !enabled {
        return h // интерфейс = {type: *MyHandler, value: nil} — НЕ nil!
    }
    return h
}

// ✅ Правильно
func getHandler() http.Handler {
    if !enabled {
        return nil // интерфейс = {type: nil, value: nil} — nil
    }
    return &MyHandler{}
}
```

### 8.2 Поведение nil-значений

| Тип | Индексация | Запись | Len/Cap | Range |
|---|---|---|---|---|
| Map | Zero-value | **panic** | 0 | 0 итераций |
| Slice | **panic** | **panic** | 0 | 0 итераций |
| Channel | Блокирует навсегда | Блокирует навсегда | 0 | Блокирует навсегда |

### 8.3 Защитное копирование

```go
// ❌ Плохо — вызывающий код может изменить внутренности
type Config struct {
    Hosts []string
}

// ✅ Хорошо — неэкспортируемое поле + копия в геттере
type Config struct {
    hosts []string
}

func (c *Config) Hosts() []string {
    return slices.Clone(c.hosts)
}
```

### 8.4 defer в цикле — утечка ресурсов

```go
// ❌ Плохо — defer выполняется при выходе из функции, не из итерации
for _, path := range paths {
    f, _ := os.Open(path)
    defer f.Close() // все файлы остаются открытыми до конца функции
    process(f)
}

// ✅ Хорошо — выносим тело цикла в отдельную функцию
for _, path := range paths {
    if err := processOne(path); err != nil {
        return err
    }
}

func processOne(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()
    return process(f)
}
```

### 8.5 Целочисленные конверсии — проверяй границы

```go
// ❌ Опасно — молчаливое переполнение
var val int64 = 3_000_000_000
i32 := int32(val) // -1294967296 — тихий wraparound

// ✅ Проверка перед конвертацией
if val > math.MaxInt32 || val < math.MinInt32 {
    return fmt.Errorf("value %d overflows int32", val)
}
i32 := int32(val)
```

### 8.6 Сравнение float через epsilon

```go
// ❌ Плохо — неточное сравнение
var a, b, c float64 = 0.1, 0.2, 0.3
a+b == c // false!

// ✅ Хорошо
const epsilon = 1e-9
math.Abs((a+b)-c) < epsilon // true
```

### 8.7 Деление на ноль

```go
// Целочисленное деление на ноль — panic
func avg(total, count int) (int, error) {
    if count == 0 {
        return 0, errors.New("division by zero")
    }
    return total / count, nil
}
```

### 8.8 sync.Once для ленивой инициализации

```go
type DB struct {
    once sync.Once
    conn *sql.DB
}

func (db *DB) connection() *sql.DB {
    db.once.Do(func() {
        db.conn, _ = sql.Open("postgres", connStr)
    })
    return db.conn
}
```

---

## 9. Тестирование

### 9.1 Table-driven tests — идиоматичный подход

```go
func TestCalculatePrice(t *testing.T) {
    tests := []struct {
        name      string
        quantity  int
        unitPrice float64
        expected  float64
    }{
        {
            name:      "single item",
            quantity:  1,
            unitPrice: 10.0,
            expected:  10.0,
        },
        {
            name:      "bulk discount",
            quantity:  100,
            unitPrice: 10.0,
            expected:  900.0,
        },
        {
            name:      "zero quantity",
            quantity:  0,
            unitPrice: 10.0,
            expected:  0.0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CalculatePrice(tt.quantity, tt.unitPrice)
            if got != tt.expected {
                t.Errorf("CalculatePrice(%d, %.2f) = %.2f, want %.2f",
                    tt.quantity, tt.unitPrice, got, tt.expected)
            }
        })
    }
}
```

### 9.2 Именование тестов

```go
func TestAdd(t *testing.T) { ... }               // тест функции
func TestMyStruct_MyMethod(t *testing.T) { ... } // тест метода
func BenchmarkAdd(b *testing.B) { ... }          // бенчмарк
func ExampleAdd() { ... }                        // пример (проверяется go test)
func FuzzAdd(f *testing.F) { ... }               // фазз-тест
```

### 9.3 Параллельные тесты

```go
func TestParallel(t *testing.T) {
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // запускаем параллельно
            // ...
        })
    }
}
```

### 9.4 Интеграционные тесты через build tags

```go
//go:build integration

package mypackage

func TestDatabaseIntegration(t *testing.T) {
    // тест с реальной БД
}
```

```bash
go test -tags=integration ./...
```

### 9.5 Детектор утечек горутин

```go
func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}
```

### 9.6 Примеры как документация

```go
func ExampleCalculatePrice() {
    price := CalculatePrice(100, 10.0)
    fmt.Printf("Price: %.2f\n", price)
    // Output: Price: 900.00
}
```

### 9.7 Покрытие кода

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 9.8 Всегда запускай с race detector

```bash
go test -race ./...
```

---

## 10. Паттерны проектирования

### 10.1 Functional Options (предпочтительный паттерн конструирования)

```go
type Server struct {
    addr         string
    readTimeout  time.Duration
    writeTimeout time.Duration
    maxConns     int
}

type Option func(*Server)

func WithReadTimeout(d time.Duration) Option {
    return func(s *Server) { s.readTimeout = d }
}

func WithMaxConns(n int) Option {
    return func(s *Server) { s.maxConns = n }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{
        addr:         addr,
        readTimeout:  5 * time.Second,  // дефолты
        writeTimeout: 10 * time.Second,
        maxConns:     100,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Использование
srv := NewServer(":8080",
    WithReadTimeout(30*time.Second),
    WithMaxConns(500),
)
```

### 10.2 Избегай `init()`

```go
// ❌ Плохо — скрытое глобальное состояние
var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
}

// ✅ Хорошо — явный конструктор, внедрение зависимости
func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}
```

### 10.3 Enum через iota — начинай с Unknown

```go
type Status int

const (
    StatusUnknown Status = iota // 0 = не установлено
    StatusActive                // 1
    StatusInactive              // 2
    StatusSuspended             // 3
)
```

### 10.4 Регулярки компилируй один раз

```go
// ✅ На уровне пакета — компилируется один раз
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) bool {
    return emailRegex.MatchString(email)
}
```

### 10.5 Таймаут на каждый внешний вызов

```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

resp, err := httpClient.Do(req.WithContext(ctx))
```

### 10.6 defer Close сразу после открытия

```go
f, err := os.Open(path)
if err != nil {
    return err
}
defer f.Close() // сразу, не через 50 строк!
```

### 10.7 Dependency Injection через интерфейсы

```go
type UserStore interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

type UserService struct {
    store UserStore
}

func NewUserService(store UserStore) *UserService {
    return &UserService{store: store}
}

// В тестах — мок, в проде — реальная БД
```

### 10.8 //go:embed для статических ресурсов

```go
import "embed"

//go:embed templates/*
var templateFS embed.FS

//go:embed version.txt
var version string
```

---

## Философия Go

> **«Clear is better than clever.»** — Ясность важнее хитроумности.

> **«A little copying is better than a little dependency.»** — Лучше скопировать немного кода, чем тащить зависимость.

> **«Don't design with interfaces, discover them.»** — Не проектируй интерфейсы заранее, открывай их по мере необходимости.

> **«Reflection is never clear.»** — Рефлексия никогда не бывает ясной. Избегай `reflect` без крайней необходимости.

- **Не абстрагируйся преждевременно** — выделяй абстракцию когда паттерн устоялся.
- **Минимизируй публичную поверхность** — каждое экспортируемое имя — это обязательство.
- **Экспортируй с умом** — сделать публичным позже легко, сделать приватным — ломающее изменение.

---

## Быстрые команды

```bash
go test ./...                          # все тесты
go test -race ./...                    # с детектором гонок
go test -coverprofile=coverage.out     # покрытие
go test -run TestName ./...            # конкретный тест
go test -tags=integration ./...        # интеграционные тесты
go test -bench=. -benchmem ./...       # бенчмарки
go vet ./...                           # статический анализ
```