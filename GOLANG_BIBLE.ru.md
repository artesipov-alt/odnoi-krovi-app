# 📖 Голанг Библия — Азбука и напоминалка для разработчика

Практический минимум правил, чтобы писать идиоматичный, современный, безопасный и
производительный Go. Собрано агентами из официальных навыков (`golang-*`) и
переведено на русский. Не исчерпывающий справочник, а чек-лист «что я забыл» перед
ревью и коммитом.

> Код, идентификаторы и команды — на английском (Go). Объяснения — на русском.

---

## Содержание

1. [Основы: стиль, именование, современный Go и безопасность](#1-основы-стиль-именование-современный-go-и-безопасность)
2. [Типы, интерфейсы и архитектура](#2-типы-интерфейсы-и-архитектура)
3. [Конкурентность, context и структуры данных](#3-конкурентность-context-и-структуры-данных)
4. [Ошибки, тестирование, линтинг и отладка](#4-ошибки-тестирование-линтинг-и-отладка)
5. [Базы данных, безопасность, производительность, наблюдаемость и документация](#5-базы-данных-безопасность-производительность-наблюдаемость-и-документация)
6. [Инструменты: CLI, gRPC, зависимости и рефакторинг](#6-инструменты-cli-grpc-зависимости-и-рефакторинг)

---

## 1. Основы: стиль, именование, современный Go и безопасность

Практический минимум правил, чтобы писать идиоматичный, современный и безопасный Go. Собрано из навыков `golang-code-style`, `golang-naming`, `golang-modernize` и `golang-safety`.

### 1.1 Стиль кода (golang-code-style)

- **Форматируй через `gofmt`/`goimports`** — никаких ручных правок форматирования.
- **Отступы — только табы**, ширина таба 8 (стандарт Go).
- **Строки ≤ 80–100 символов**, но не ломай читаемость ради длины.
- **Имена пакетов** — короткие, строчные, без подчёркиваний и `util`/`common`: `http`, `user`, а не `user_service`.
- **Группируй импорты**: стандартная библиотека / сторонние / локальные (через пустую строку).
- **Обрабатывай ошибки сразу**, не откладывай:

  ```go
  // Делай
  f, err := os.Open(name)
  if err != nil {
      return err
  }
  defer f.Close()
  ```

- **`defer`** — для освобождения ресурсов (закрытие файлов, анлоки), ставь сразу после успешного получения ресурса.
- **Контекст передавай первым аргументом**: `func Do(ctx context.Context, req Request) error`.
- **Не используй `panic`** для обычных ошибок — только для невозможных ситуаций при инициализации.
- **Комментарии к экспортируемым сущностям** обязательны и начинаются с имени: `// User представляет зарегистрированного пользователя.`

### 1.2 Именование (golang-naming)

- **Краткость важна**: `i` лучше чем `index`, `buf` лучше чем `buffer`, если область видимости мала.
- **Имя должно описывать содержимое, а не тип**: `users` а не `userSlice`, `userMap`.
- **Геттеры/сеттеры** — без `Get`/`Set`: `user.Name()`, `user.SetName()`.
- **Булевы переменные** — с префиксом утверждения: `isReady`, `hasToken`, `canEdit` (не `ready`, `token`).
- **Не используй венгерскую нотацию** (`strName`, `iCount`) — тип виден из объявления.
- **Избегай зарезервированных/теневых имён**: не называй переменную `len`, `nil`, `error`.
- **Константы** — `CamelCase` с заглавной для экспорта: `MaxRetries`; для локальных — строчные.
- **Интерфейсы из одного метода** называй с суффиксом `-er`: `Reader`, `Closer`, `Notifier`.
- **Пакеты** — одно слово, не множественное: `time`, а не `times`.

### 1.3 Современный Go (golang-modernize)

- **Используй дженерики** (Go 1.18+) для повторяющихся типобезопасных функций:

  ```go
  func Map[T, U any](s []T, f func(T) U) []U {
      r := make([]U, len(s))
      for i, v := range s {
          r[i] = f(v)
      }
      return r
  }
  ```

- **`any` вместо `interface{}`** (с Go 1.18).
- **`errors.Join`** (Go 1.20) для объединения нескольких ошибок:

  ```go
  // Делай
  return errors.Join(err1, err2)
  ```

- **`fmt.Errorf` с `%w`** для оборачивания ошибок (можно потом `errors.Is`/`As`):

  ```go
  return fmt.Errorf("загрузка конфигурации: %w", err)
  ```

- **`errors.Is` / `errors.As`** вместо сравнения строк или приведения типов.
- **`slices` и `maps`** (Go 1.21+) вместо ручных циклов: `slices.Contains`, `slices.Index`, `maps.Keys`.
- **`min`/`max`/`clear`** встроенные (Go 1.21+) вместо самописных хелперов.
- **`for range` по карте/срезу** — не создавай лишних индексных переменных.
- **`context.WithTimeout`/`WithCancel`** + `defer cancel()` — всегда освобождай контекст.

### 1.4 Безопасность (golang-safety)

- **Никогда не игнорируй ошибку** присваиванием `_`, если она значима:

  ```go
  // Не делай
  json.Unmarshal(data, &v) // ошибка проигнорирована

  // Делай
  if err := json.Unmarshal(data, &v); err != nil {
      return fmt.Errorf("разбор JSON: %w", err)
  }
  ```

- **Проверяй `nil` перед разыменованием** указателей и интерфейсов.
- **Закрывай `Body` ответов HTTP** (`defer resp.Body.Close()`), иначе утечка соединений.
- **Не используй `unsafe`** без крайней необходимости — ломает типобезопасность и GC-гарантии.
- **`crypto/rand`** для секретов/токенов, НЕ `math/rand`:

  ```go
  // Делай (криптографически стойко)
  b := make([]byte, 32)
  rand.Read(b)
  ```

- **Не логируй секреты** (пароли, токены, ключи) — используй редактирование/маскировку.
- **Ограничивай размер ввода** (чтение из сети/файла), чтобы избежать OOM/DoS.
- **Избегай `time.Now()` в тестах** — инъектируй время или используй `clock` для детерминизма.
- **Гонки данных**: защищай общий доступ через `sync.Mutex`/`RWMutex` или каналы, запускай `go test -race`.

---

## 2. Типы, интерфейсы и архитектура

Практический раздел по проектированию типов, интерфейсов и структуры Go-проектов. Охватывает композицию, внедрение зависимостей и раскладку пакетов — то, что чаще всего забывают.

### 2.1 Структуры и интерфейсы (`golang-structs-interfaces`)

- **Держи интерфейсы маленькими** — 1–3 метода. Большие интерфейсы = слабая абстракция. Собирай большие из маленьких:
  ```go
  type Reader interface { Read(p []byte) (n int, err error) }
  type Writer interface { Write(p []byte) (n int, err error) }
  type ReadWriter interface { Reader; Writer } // композиция
  ```
- **Определяй интерфейс там, где он потребляется**, а не где реализуется. Реализатор не должен знать о контракте.
  ```go
  // пакет notification — определяет только то, что нужно
  type Sender interface { Send(to, body string) error }
  type Service struct { sender Sender }
  ```
- **Принимай интерфейсы, возвращай структуры** (accept interfaces, return structs).
  - Делай: `func NewService(store UserStore) *Service { ... }`
  - Не делай: `func NewService(store UserStore) ServiceInterface { ... }`
- **Не создавай интерфейсы заранее.** Жди 2+ реализации или нужду в моке. Начинай с конкретного типа:
  - Не делай: интерфейс `UserRepository` при единственной реализации.
  - Делай: `type UserRepository struct { db *sql.DB }`, выдели интерфейс позже.
- **Делай нулевое значение полезным** — структура должна работать без явной инициализации (ленивая инициализация map в методах, а не `nil` map).
- **Проверка интерфейса на этапе компиляции** — бесплатно, ловит поломку сборки:
  ```go
  var _ io.ReadWriter = (*MyBuffer)(nil)
  ```
- **Безопасные type assertion** — всегда comma-ok:
  - Делай: `s, ok := val.(string)`
  - Не делай: `s := val.(string)` // паникует, если тип не совпал
- **Встраивание (embedding) vs именованное поле:**
  - Встраивай (`Logger`), когда внешний тип «является» улучшенной версией и должен экспонировать весь API.
  - Именованное поле (`store *DataStore`), когда тип лишь «имеет» зависимость и не должен светить её методы.
- **Избегай `any`/`interface{}`**, где хватит дженериков (Go 1.18+):
  - Не делай: `func Contains(slice []any, target any) bool`
  - Делай: `func Contains[T comparable](slice []T, target T) bool`
- **Поля сериализуемых структур** — теги обязательны на экспортируемых полях:
  ```go
  type Order struct {
      ID        string    `json:"id" db:"id"`
      Total     float64   `json:"total" db:"total"`
      DeletedAt time.Time `json:"-" db:"deleted_at"` // "-" — исключить из JSON
  }
  ```
- **Получатели (receivers):** метод меняет receiver → указатель; маленький/неизменяемый → значение. **Будь последователен** — если один метод указательный, все должны быть.
- **`noCopy`** для структур с mutex/каналами — `go vet` поймает случайное копирование. Передавай такие структуры только по указателю.

### 2.2 Паттерны проектирования (`golang-design-patterns`)

- **Функциональные опции** для конструкторов — масштабируются без breaking changes:
  ```go
  type Option func(*Server)
  func WithReadTimeout(d time.Duration) Option {
      return func(s *Server) { s.readTimeout = d }
  }
  func NewServer(addr string, opts ...Option) *Server {
      s := &Server{addr: addr, readTimeout: 5 * time.Second} // значения по умолчанию
      for _, opt := range opts { opt(s) }
      return s
  }
  // srv := NewServer(":8080", WithReadTimeout(30*time.Second))
  ```
  - Опция, чья валидация может упасть, **должна возвращать error** — лови плохой конфиг при создании, а не в рантайме.
- **Избегай `init()`** — неявный запуск, нет возврата ошибки, ломает тесты. Инициализируй явно в конструкторе.
  - Не делай: `func init() { db, _ = sql.Open(...) }`
  - Делай: `func NewUserRepository(db *sql.DB) *UserRepository { ... }`
- **Enum — начинай с 1** (или `Unknown` на 0). Нулевое значение = невалидное/несет состояние:
  ```go
  type Status int
  const (
      StatusUnknown Status = iota // 0 — невалидно
      StatusActive                // 1
  )
  ```
- **Сначала обрабатывай ошибки**, ранний return — счастливый путь остаётся плоским.
- **Panic — для багов, не для ожидаемых ошибок.** Сетевые сбои/неверный ввод → `error`.
- **`defer Close()` сразу после открытия** — не через 50 строк.
- **Таймаут на каждый внешний вызов:**
  ```go
  ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
  defer cancel()
  resp, err := httpClient.Do(req.WithContext(ctx))
  ```
- **Регулярку компилируй один раз** на уровне пакета (`regexp.MustCompile`), а не внутри функции.
- **`//go:embed`** для статических ассетов — встраивается при компиляции, нет файлового I/O в рантайме.
- **Архитектура:** не навязывай слои маленькому проекту. Держи домен чистым (без зависимостей от фреймворков), валидируй на границах, делай невалидные состояния непредставимыми.

### 2.3 Раскладка проекта (`golang-project-layout`)

- **Не переусложняй** — 100-строчному CLI не нужны слои абстракции и DI. Спроси архитектуру у разработчика заранее.
- **`cmd/` — только точки входа.** Один подкаталог на `main`, минимум логики (флаги, сборка зависимостей, `Run()`). Бизнес-логика — в `internal/` или `pkg/`.
  ```
  cmd/
  ├── server/main.go   // API-сервер
  ├── worker/main.go   // фоновый воркер
  └── migrate/main.go  // миграции БД
  ```
- **`internal/`** — приватный код приложения (неэкспортируемый). **`pkg/`** — только если код полезен внешним потребителям.
- **Не используй** `src/`, `utils/`, `common/`, `helpers/` и `main.go` в корне — это не Go-стиль.
- **Имя модуля** (`go.mod`) = URL репозитория, lowercase, через дефис: `github.com/jdoe/payment-processor`. Пакеты — lowercase, единственное число, совпадают с каталогом.
- **12-факторное приложение** для сервисов: конфиг через env, логи в stdout, stateless, graceful shutdown, admin-задачи как отдельные бинари (напр. `cmd/migrate/`).
- **Тесты** — рядом с кодом (`_test.go`), фикстуры в `testdata/`.

### 2.4 Внедрение зависимостей (`golang-dependency-injection`)

- **Внедряй через конструкторы** — никогда не используй глобальные переменные и `init()` для настройки сервисов.
  - Делай:
    ```go
    func NewUserService(db UserStore, mailer Mailer, logger *slog.Logger) *UserService {
        return &UserService{db: db, mailer: mailer, logger: logger}
    }
    ```
  - Не делай: `func NewUserService() *UserService { db, _ := sql.Open(...) }` // скрытая зависимость
- **Интерфейсы у потребителя, конкретика у реализатора** — принимай интерфейсы на границах потребления.
- **Маленькие проекты (< 10 сервисов)** — ручное конструкторное DI, без библиотек.
- **Никаких глобальных реестров / service locator.** Контейнер живёт только в composition root (`main()`), **никогда не передавай его как зависимость**.
- **Держи граф зависимостей плоским** — глубокие цепочки A→B→C→D→E = проблема дизайна.
- **Мокай на границе интерфейса** — в тестах подставляй stub:
  ```go
  type MockUserStore struct{ users map[string]*User }
  func (m *MockUserStore) FindByID(ctx context.Context, id string) (*User, error) {
      u, ok := m.users[id]
      if !ok { return nil, ErrNotFound }
      return u, nil
  }
  svc := NewUserService(mock, nil, slog.Default()) // без реальной БД
  ```
- **Ленивая инициализация** — создавай сервис при первом запросе; синглтоны для stateful (БД, кэш), транзиенты для stateless.
- **Выбор библиотеки DI** (по размеру): <10 сервисов — ручное; 10–20 — рассмотри библиотеку; 20+ с lifecycle/health checks — рекомендуется (google/wire — codegen, uber-go/fx — reflection, samber/do — generics).

---

## 3. Конкурентность, context и структуры данных

Практический раздел по самому забываемому в Go: утечки горутин, владение каналами, корректное распространение `context` и внутреннее устройство слайсов/мап. Читай как чек-лист перед ревью.

### 3.1 Конкурентность

- **У каждой горутины должен быть явный выход.** Без `context`, done-канала или `WaitGroup` она течёт и накапливается до падения процесса.
- **Всегда включай `<-ctx.Done()` в `select`.** Без этого горутина не заметит отмену и утечёт.
  ```go
  for {
      select {
      case <-ctx.Done():
          return // обязательно
      case task, ok := <-in:
          if !ok {
              return
          }
          handle(ctx, task)
      }
  }
  ```
- **Канал закрывает только отправитель.** Закрытие со стороны получателя → panic, если отправитель пишет после закрытия.
- **Указывай направление канала** (`chan<-`, `<-chan`) — компилятор запретит неверное использование.
  ```go
  func produce(ch chan<- int) { ... } // только отправка
  func consume(ch <-chan int)  { ... } // только приём
  ```
- **Отправляй копии, а не указатели** по каналам — указатель создаёт невидимую общую память.
- **По умолчанию — небуферизованные каналы.** Большие буферы маскируют отсутствие backpressure; обосновывай `N` в комментариях. Буфер 1 — для сигнальных каналов: `done := make(chan struct{}, 1)`.
- **Не создавай `time.After` в горячих циклах** — каждый вызов аллоцирует таймер. Переиспользуй `time.NewTimer` + `Reset`.
  ```go
  timer := time.NewTimer(5 * time.Second)
  defer timer.Stop()
  for {
      select {
      case msg := <-ch:
          timer.Stop()
          timer.Reset(5 * time.Second)
          handle(msg)
      case <-timer.C:
          handleTimeout()
          timer.Reset(5 * time.Second)
      }
  }
  ```
- **`wg.Add` вызывай ДО `go`.** Иначе `Wait` может вернуться раньше времени.
- **Ограничивай распараллеливание.** Бесконечный `go func()` → используй `errgroup.SetLimit(n)` или семафор.
- Восстанавливай panic на границах горутин (`defer recover()` в проде).

**errgroup vs WaitGroup**

| Нужно | Использовать |
| --- | --- |
| Дождаться, ошибки не важны | `sync.WaitGroup` |
| Дождаться + первый error | `errgroup.Group` |
| Отменить остальных при ошибке | `errgroup.WithContext` |
| Ограничить конкурентность | `errgroup.SetLimit(n)` |

**Чем защищать данные**

- Канал — передача данных между горутинами (явная передача владения).
- `sync.Mutex`/`RWMutex` — защита полей структуры; критические секции короткие, не держи мьютекс через I/O.
- `sync/atomic` — счётчики/флаги; предпочитай типизированные `atomic.Int64`, `atomic.Bool` (Go 1.19+).
- `sync.Map` — много читателей, мало писателей (конкурентное чтение/запись в обычную map — hard crash).
- `x/sync/singleflight` — дедупликация одновременных вызовов (защита от cache stampede).

**Делай / Не делай**
- ✅ Делай: лови утечки тестами через `go.uber.org/goleak`, гонки через `go test -race ./...`.
- ❌ Не делай: fire-and-forget горутины без механизма остановки.

### 3.2 Context

- **`ctx` — первый параметр, имя `ctx context.Context`.** Распространяй один и тот же контекст по всей цепочке: HTTP handler → service → DB → внешние API. Отмена родителя автоматически отменяет всё внизу.
  ```go
  // ✗ плохо — рвёт цепочку
  return s.db.ExecContext(context.Background(), "INSERT ...")
  // ✓ хорошо — пробрасываем ctx вызывающего
  return s.db.ExecContext(ctx, "INSERT ...")
  ```
- **Никогда не храни `context` в структуре** — передавай явно через параметры.
- **Не передавай `nil`.** Если контекста ещё нет — `context.TODO()`.
- **`context.Background()`** только на верхнем уровне (main, init, тесты). Никогда не создавай его посреди пути запроса.
- **Всегда вызывай `cancel()`** на всех путях для `WithCancel`/`WithTimeout`/`WithDeadline` (обычно `defer cancel()` сразу).
  ```go
  func fetch(ctx context.Context) error {
      ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
      defer cancel() // иначе утечка таймера/ресурсов
      return doWork(ctx)
  }
  ```
- **Вложенные таймауты: побеждает короткий.** Ребёнок с 10s над родителем с 2s всё равно истечёт через 2s.
- **Слушай отмену в `select`** (`<-ctx.Done()`) и проверяй `ctx.Err()` в CPU-ботных циклах.
- **`context.AfterFunc` (Go 1.21+)** — колбэк в отдельной горутине при отмене (для cleanup без блокировки).
- **`context.WithoutCancel` (Go 1.21+)** — фоновая работа, которая должна пережить запрос (аудит-логи). Сохраняет значения (trace_id), но отвязывает отмену.
  ```go
  auditCtx := context.WithoutCancel(ctx)
  go h.auditService.LogOrderCreated(auditCtx, order)
  ```
- **Значения в context:** ключи — только неэкспортируемые типы (защита от коллизий); храни только метаданные запроса (request_id, user_id), НЕ параметры функций.

### 3.3 Структуры данных

- **Преаллоцируй слайсы и мапы**, когда размер известен/оценим — `make([]T, 0, n)` / `make(map[K]V, n)`. Иначе каждый рост копирует backing array (O(n)).
  ```go
  users := make([]User, 0, len(ids))
  m := make(map[string]*User, len(users))
  s = slices.Grow(s, additionalNeeded) // дозарастить (Go 1.21+)
  ```
- **Слайс — заголовок из 3 слов** (ptr, len, cap). Несколько слайсов могут делить backing array. Рост: <256 — удвоение, ≥256 — ~+25%. **Не полагайся на тайминг роста** — алгоритм менялся между версиями.
- **Полные копии vs ссылки** (семантика копирования):

| Тип | Копия | Независимость |
| --- | --- | --- |
| `int`, `float`, `bool`, `string`, `array`, `struct` | значение (глубокая) | полная |
| `slice` | копируется заголовок, array общий | нужен `slices.Clone` |
| `map` | ссылка | нужен `maps.Clone` |
| `channel` | ссылка | тот же канал |
| `*T` | адрес | тот же объект |

- **Мапа — хеш-таблица с бакетами по 8 элементов.** Это reference type: присваивание копирует указатель, не данные. Мапы никогда не уменьшаются. Для больших значений в мапе используй `map[K]*V`, чтобы избежать копирования при доступе.
- **Массивы** — только для фиксированных размеров (digest, IPv4, матрицы). Копируются целиком; сравнимы → можно как ключ мапы `[2]int`.
- **`strings.Builder`** для сборки строк (без копии в `String()`); **`bytes.Buffer`** для двунаправленного I/O (реализует `io.Reader`+`io.Writer`). Оба поддерживают `Grow`.
- **`container/`**: `heap` — приоритетная очередь, `list` — только при частых вставках в середину, `ring` — циклический буфер. Линкед-листы плохи по cache locality — бенчмаркай прежде чем брать вместо слайса.
- **Дженерики:** используй самое узкое ограничение — `comparable` для ключей, `cmp.Ordered` для сортировки.
- **`unsafe.Pointer`** — только 6 валидных паттернов спека; никогда не храни в `uintptr` между операторами (GC может переместить объект → висячая ссылка).
- **`weak.Pointer[T]` (Go 1.24+)** — для кешей/каноникизации, позволяет GC забирать записи.

**Делай / Не делай**
- ✅ Делай: `make([]T, 0, n)` перед циклом `append`; `map[K]*V` для крупных значений.
- ❌ Не делай: `bytes.Buffer` для чистой конкатенации строк; растить слайс в цикле без преаллокации.

---

## 4. Ошибки, тестирование, линтинг и отладка

Практический раздел-напоминалка по работе с ошибками, тестам, линтингу и отладке в Go. Охватывает то, что разработчики чаще всего забывают: обёртку ошибок, sentinel-ошибки, структурный логгинг, table-driven тесты и инструменты отладки.

### 4.1 Обработка ошибок (golang-error-handling)

- **Всегда проверяйте ошибку** сразу после возврата функции. Не игнорируйте через `_`.

```go
// Делай
f, err := os.Open("file.txt")
if err != nil {
    return err
}

// Не делай
f, _ := os.Open("file.txt")
```

- **Оборачивайте ошибки через `%w`** для сохранения цепочки (stack).

```go
// Делай
if err != nil {
    return fmt.Errorf("failed to read config: %w", err)
}
```

- **Проверяйте тип ошибки через `errors.Is` / `errors.As`**, а не сравнением строк или `==`.

```go
// Делай
if errors.Is(err, os.ErrNotExist) {
    // файл не найден
}

var pathErr *os.PathError
if errors.As(err, &pathErr) {
    // извлекаем вложенную ошибку
}
```

- **Sentinel-ошибки** (предопределённые, напр. `var ErrNotFound = errors.New("not found")`) — экспортируйте как `var`, а не `const`, чтобы их можно было сравнивать через `errors.Is`.

```go
var ErrNotFound = errors.New("not found")
```

- **Правило единственной обработки**: ошибку либо обрабатывают локально, либо прокидывают наверх — не делайте и то, и другое (не логируйте и сразу же `return err` на каждом слое).

- **Контекст добавляйте при обёртке**, но не дублируйте сообщения одного и того же уровня.

- **Не используйте `panic`** для обычного потока управления; только для невосстановимых ситуаций (init, программа не может стартовать).

### 4.2 Тестирование (golang-testing)

- **Table-driven тесты** — стандарт для Go: один тест, много case через срез.

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"zero", 0, 0, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

- **Параллельные тесты**: вызывайте `t.Parallel()` внутри `t.Run`, если кейсы независимы (ускоряет прогон).

```go
t.Run("case", func(t *testing.T) {
    t.Parallel()
    // ...
})
```

- **Fuzzing** (`go test -fuzz=F`) — для поиска краевых входных данных.

```go
func FuzzParse(f *testing.F) {
    f.Add("input")
    f.Fuzz(func(t *testing.T, s string) {
        Parse(s)
    })
}
```

- **Проверка утечек горутин — `goleak`**: добавьте в TestMain.

```go
func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}
```

- **Детектор гонок**: запускайте тесты/бинарь с `-race`.

```bash
go test -race ./...
go build -race .
```

- Не используйте `time.Sleep` для синхронизации в тестах — используйте каналы/WaitGroup/`httptest`.

### 4.3 Линтинг (golang-lint)

- **golangci-lint** — агрегатор линтеров; конфиг в `.golangci.yml`.

```yaml
run:
  timeout: 5m
linters:
  enable:
    - govet
    - errcheck
    - staticcheck
    - ineffassign
    - unused
issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
```

- **Запуск**: `golangci-lint run ./...`

- **`//nolint`** — подавление конкретного линтера с обязательным пояснением.

```go
//nolint:errcheck // ответ игнорируем намеренно
resp.Body.Close()
```

- **`//nolint:lint1,lint2`** — перечисляйте только нужные линтеры, не используйте `//nolint:all` без причины.

- Не отключайте линтеры глобально в CI ради скорости — лучше точечные `exclude-rules`.

- Обязательно гоняйте линтер в CI, а не только локально.

### 4.4 Отладка и troubleshooting (golang-troubleshooting)

- **Структурное логгирование через `slog`** (стандартная библиотека), а не `fmt.Println`.

```go
logger.Info("request handled",
    slog.String("path", r.URL.Path),
    slog.Int("status", status),
)
```

- **Детектор гонок**: `go run -race` / `go test -race` — ловит data race в рантайме.

- **`GODEBUG`** — включает внутреннюю диагностику рантайма.

```bash
GODEBUG=gctrace=1 ./app      # лог сборок мусора
GODEBUG=http2debug=2 ./app   # дебаг HTTP/2
```

- **Delve (`dlv`)** — отладчик Go: `dlv debug`, `dlv test`, точки останова `b main.main`, `c`, `n`, `p var`.

```bash
dlv debug ./cmd/app
(dlv) break main.main
(dlv) continue
(dlv) print user
```

- **Профилирование**: `go tool pprof` по `net/http/pprof` — CPU, heap, goroutine leaks.

- **При зависании/утечке горутин**: `curl localhost:6060/debug/pprof/goroutine?debug=2` — стек всех горутин.

- **`panic` в логах** — читайте полный stack trace сверху вниз до своего кода; используйте `recover()` только в goroutine-обработчиках (например, HTTP middleware).

---

## 5. Базы данных, безопасность, производительность, наблюдаемость и документация

Практический «напоминалка» по тому, что разработчики чаще всего забывают при работе с БД, безопасностью, производительностью, метриками/логами и документацией.

### 5.1 Базы данных (database/sql, sqlx, pgx)

- **Параметризованные запросы — всегда.** Никогда не конкатенируй пользовательский ввод в SQL.
  ```go
  // ✗ ПЛОХО — SQL-инъекция
  query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)

  // ✓ ХОРОШО (PostgreSQL)
  err := db.GetContext(ctx, &user, "SELECT id, name, email FROM users WHERE email = $1", email)
  // ✓ ХОРОШО (MySQL)
  err := db.GetContext(ctx, &user, "SELECT id, name, email FROM users WHERE email = ?", email)
  ```
- **Динамический IN:** используй `sqlx.In`, а не ручную сборку строки.
  ```go
  query, args, err := sqlx.In("SELECT * FROM users WHERE id IN (?)", ids)
  query = db.Rebind(query)
  err = db.SelectContext(ctx, &users, query, args...)
  ```
- **Имена колонок/сортировки** — только через allowlist, не интерполируй из ввода.
- **Передавай `ctx`** во все операции (`QueryContext`, `ExecContext`, `GetContext`). Без контекста запрос не прервётся при отключении клиента.
- **Закрывай rows:** `defer rows.Close()` сразу после `QueryContext`, и проверяй `rows.Err()` после цикла.
- **`sql.ErrNoRows`** лови явно через `errors.Is` и переводи в доменный error (not found ≠ ошибка).
- **Не используй `db.Query` для не-SELECT** (нет строк → утечка соединения). Для записи — `db.Exec`.
- **NULLable-колонки:** указатели (`*string`, `*time.Time`) или `sql.NullXxx`. В sqlx — тег `db:"column_name"`, в pgx — `pgx.RowToStructByName`.
- **Транзакции** для многошаговых операций: `BeginTxx`/`Commit`. Для данных, которые собираешься менять — `SELECT ... FOR UPDATE` (защита от гонок). Для финансов — изоляция `serializable`.
- **Пул соединений** — обязательно настрой:
  ```go
  db.SetMaxOpenConns(25)
  db.SetMaxIdleConns(10)
  db.SetConnMaxLifetime(5 * time.Minute)
  db.SetConnMaxIdleTime(1 * time.Minute)
  ```
- **Батчи** — разумными порциями: не по одной строке (много раунд-трипов) и не миллионами сразу (локи/память).
- **ORM — нет.** Используй sqlx/pgx. Миграции — golang-migrate/Flyway/Atlas (не ручные/AI-сгенерированные).

### 5.2 Безопасность

- **SQL-инъекция:** только параметризованные запросы (см. 5.1).
- **Command injection:** `exec.Command("sh", "-c", script)` — плохо; передавай аргументы отдельно: `exec.Command("prog", arg1, arg2)`.
- **Секреты:** никогда не хардкодь в исходниках (попадают в git-историю/CI/бэкапы). Используй env-переменные или secret manager.
- **Криптография:** не изобретай свою. Токены/рандом — `crypto/rand`, а не `math/rand` (предсказуемо). Пароли — Argon2id или bcrypt, не MD5/SHA1. AES — только с GCM (аутентификация), не ECB/CBC.
- **Сравнение секретов:** не через `==` (утечка по таймингу). Используй `crypto/subtle.ConstantTimeCompare`.
- **Ошибки наружу:** возвращай общие сообщения, детали логируй сервер-сайд. Не возвращай стектрейсы/ошибки БД клиенту.
- **Клиентские заголовки** (`X-Forwarded-For`, `X-Is-Admin`) и client-side авторизация — не доверяй; проверяй сервер-сайд.
- **Path traversal:** Go 1.24+ — `os.Root`; до — `filepath.IsLocal` + `filepath.Rel`, не полагайся на `filepath.Clean`+`HasPrefix`.
- **Race conditions:** прогоняй `go test -race ./...`, фиксируй все находки.
- **Инструменты:** `go tool gosec ./...`, `go tool govulncheck ./...`.

### 5.3 Производительность

- **Сначала профилируй (pprof), потом оптимизируй.** Интуиция о бутылочном горлышке ошибается ~80% времени.
- **Снижение аллокаций — самый большой ROI** (GC не бесплатен). Преаллоцируй слайсы/мапы, используй `strings.Builder`.
- **`http.Client` без настроенного Transport:** `MaxIdleConnsPerHost` по умолчанию = 2 — подними под свой уровень конкурентности.
- **Логи в горячих циклах** — мешают инлайнингу и аллоцируют. Используй `slog.LogAttrs`.
- **`panic`/`recover` как управление потоком** — не используй (аллокация стектрейса). Возвращай `error`.
- **`reflect.DeepEqual` в проде** — в 50–200 раз медленнее типизированного; используй `slices.Equal`, `maps.Equal`, `bytes.Equal`.
- **GC в контейнерах:** ставь `GOMEMLIMIT` = 80–90% памяти контainerа, чтобы не словить OOM-kill.
- **Кэширование / singleflight** — для повторяющейся дорогой работы (один и тот же fetch много раз).
- **Пулы/батчи** — см. 5.1 (пул соединений, batch операции).

### 5.4 Наблюдаемость (логи, метрики, трейсы)

- **Структурированные логи** через `log/slog` (JSON в проде). Уровни: Debug/Info/Warn/Error по назначению. Коррелируй с трейсами: `slog.InfoContext(ctx, ...)`.
  ```go
  // ✗ ПЛОХО — логируешь И возвращаешь (залогируется несколько раз)
  // ✓ ХОРОШО — верни с контекстом, залогируй один раз на верхнем уровне
  return fmt.Errorf("querying users: %w", err)
  ```
- **Метрики (Prometheus):** latency/error-rate на каждом HTTP-эндпоинте. **Histogram, а не Summary** (агрегируется на сервере, `histogram_quantile`).
  ```go
  prometheus.NewHistogram(prometheus.HistogramOpts{
      Name:    "http_request_duration_seconds",
      Buckets: prometheus.DefBuckets,
  })
  ```
- **Кардинальность меток:** НЕ используй неограниченные значения (user ID, полные URL) как лейблы — только ограниченные (`routePattern`).
  ```go
  // ✗ ПЛОХО
  httpRequests.WithLabelValues(r.Method, r.URL.Path, userID).Inc()
  // ✓ ХОРОШО
  httpRequests.WithLabelValues(r.Method, routePattern).Inc()
  ```
- **OpenTelemetry трейсы:** настраивай TracerProvider сразу, добавляй span на каждую операцию (метод, БД-запрос, внешний вызов), `span.RecordError(err)`. Прокидывай `ctx` везде (иначе трейс обрывается).
- **Корреляция:** `otelslog` бридж внедряет `trace_id`/`span_id` в логи; exemplars связывают метрику с трейсом.
- **Профилирование (pprof/Pyroscope):** включай через env-переменные без редеплоя.
- **«Фича не готова, пока не наблюдаема»** — метрики + логи + спаны + дашборды/алерты.

### 5.5 Документация (godoc)

- **Каждый экспортируемый символ** должен иметь doc-комментарий, начинающийся с имени и глагола. Пиши **почему/когда**, а не пересказ сигнатуры.
  ```go
  // CalculateDiscount вычисляет итоговую цену после применения ступенчатых скидок.
  // Возвращает ErrInvalidPrice, если basePrice отрицательно.
  func CalculateDiscount(basePrice float64, quantity int, tiers []DiscountTier) (float64, error)
  ```
- **Комментарий к пакету** (`// Package foo ...`) — обязателен.
- **Что нужно всем проектам:** doc-комментарии, package comment, `README.md`, `LICENSE`, getting started, рабочие примеры.
- **Для библиотек:** `ExampleXxx` тест-функции (исполняемая документация), Go Playground демо (`// Play: ...`), `go doc` для превью на pkg.go.dev.
- **Для приложений/CLI:** исчерпывающий `--help`, документация env-переменных/флагов, способы установки.
- **AI-дружественность:** `llms.txt` в корне репозитория + структурированные doc-комментарии.
- **CHANGELOG** — формат Keep a Changelog; пиши «что изменилось для пользователя», без раздувания.

---

## 6. Инструменты: CLI, gRPC, зависимости и рефакторинг

Практический раздел по инструментам: как правильно строить CLI, работать с gRPC, управлять зависимостями и безопасно рефакторить Go-код.

### 6.1 CLI (cobra + viper)

- **Флаги через cobra, конфиг через viper.** Не смешивайте ручной `flag` и viper.
- **Делай:** регистрируйте флаги в `PersistentPreRun`/`init`, а не внутри `RunE`.
- **Не делай:** не читайте `os.Args` вручную — используйте `cmd.Flags().StringVar(&cfg, "host", "localhost", "хост")`.
- **Обязательно вызывайте `Execute()` и возвращайте код выхода:**

```go
func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1) // корректный код ошибки
    }
}
```

- **Делай:** пишите ошибки в `os.Stderr`, а результат — в `os.Stdout`.
- **Коды выхода:** `0` — успех, `1` — общая ошибка, `2` — ошибка использования (usage). Используйте `cobra.Command{ SilenceUsage: true, SilenceErrors: true }`.
- **Тестирование CLI:** выносите логику из `RunE` в чистую функцию и тестируйте её; перехватывайте вывод через `bytes.Buffer`:

```go
var out bytes.Buffer
cmd.SetOut(&out)
cmd.SetArgs([]string{"--host", "example.com"})
cmd.Execute()
```

- **Не делай:** не пишите бизнес-логику прямо в `RunE` — её невозможно протестировать без запуска бинарника.

### 6.2 gRPC

- **Используйте правильные статус-коды** через `status`/`codes`, а не голые `error`:

```go
import "google.golang.org/grpc/codes"
import "google.golang.org/grpc/status"

if user == nil {
    return nil, status.Errorf(codes.NotFound, "пользователь не найден: %s", id)
}
```

- **Делай:** оборачивывайте ошибки `status.Errorf`, чтобы клиент получал структурированный код.
- **Перехватчики (interceptors)** для сквозной логики — логи, метрики, auth, recovery:

```go
grpc.NewServer(
    grpc.UnaryInterceptor(otel.UnaryServerInterceptor()),
    grpc.StreamInterceptor(otel.StreamServerInterceptor()),
)
```

- **Не делай:** не дублируйте проверки токена в каждом методе — выносите в unary/stream interceptor.
- **Стриминг:** для больших данных используйте `stream.Send()` порциями, а не один огромный ответ.
- **Локальное тестирование без сети — `bufconn`:**

```go
lis := bufconn.Listen(1024 * 1024)
srv := grpc.NewServer()
go srv.Serve(lis)
conn, _ := grpc.Dial("bufnet",
    grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
        return lis.Dial()
    }),
    grpc.WithTransportCredentials(insecure.NewCredentials()))
```

- **Не делай:** не поднимаете реальный порт в unit-тестах — это делает тесты медленными и гоночными.
- **Делай:** версионируйте API (например, пакет `v1`) и не ломайте контракты несовместимо.

### 6.3 Управление зависимостями

- **`go.mod` — единственный источник правды.** Не правьте `go.sum` руками.
- **Minimal Version Selection (MVS):** Go берёт минимально необходимую версию, совместимую со всеми требованиями. Не бойтесь старых версий — это норма.
- **Делай:** `go mod tidy` перед коммитом — убирает лишнее, добавляет нужное.
- **Не делай:** не фиксируйте лишние/неиспользуемые зависимости; не пишите `replace` для публичных модулей без крайней нужды.
- **Обновления:** `go get -u ./...` (всё) или `go get pkg@v1.2.3` (точечно). Проверяйте `go.sum`.
- **Сканирование уязвимостей:**

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

- **Монорепо / несколько модулей — `go.work`:**

```bash
go work init ./service-a ./service-b
go work use ./shared
```

- **Не делай:** не держите `vendor/` без `go mod vendor` и `-mod=vendor` в CI; не игнорируйте предупреждения `go mod verify`.

### 6.4 Рефакторинг

- **Безопасный процесс:** сначала покрытие тестами → маленькие шаги → прогон `go test` после каждого шага.
- **Сначала — сеть безопасности (coverage):**

```bash
go test -coverprofile=cover.out ./...
go tool cover -func=cover.out
```

- **Переименование через `gopls`, а не руками:** `gofmt`/`gopls rename` учитывает весь модуль и переименовывает использования корректно.
- **Инлайн (inline) и экстракт** через `gopls` (code action) — меньше опечаток, чем копипаста.
- **Делай:** рефакторь небольшими коммитами, чтобы `git bisect` и откат работали.
- **Не делай:** не меняйте публичный API и поведение одновременно — разделяйте.
- **После правок — `go vet` и `go build ./...`** перед пушем.
- **Делай:** удаляйте мёртвый код (`deadcode`, `unused`) и проверяйте, что coverage не упал.

---

> Собрано из навыков `golang-code-style`, `golang-naming`, `golang-modernize`, `golang-safety`,
> `golang-structs-interfaces`, `golang-design-patterns`, `golang-project-layout`,
> `golang-dependency-injection`, `golang-concurrency`, `golang-context`, `golang-data-structures`,
> `golang-error-handling`, `golang-testing`, `golang-lint`, `golang-troubleshooting`,
> `golang-database`, `golang-security`, `golang-performance`, `golang-observability`,
> `golang-documentation`, `golang-cli`, `golang-grpc`, `golang-dependency-management`,
> `golang-refactoring`.
