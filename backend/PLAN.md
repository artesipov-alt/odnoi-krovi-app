## 1. Контекст и цель

**Бизнес-задача.** Владельцы питомцев с подходящей группой крови и регионом теперь могут включить настройку «Open for Contact» в `DonorPreference`. Когда такая настройка включена, реципиент, открывая свою заявку на странице поиска крови, должен видеть список таких питомцев как **потенциальных доноров** и иметь возможность **самому выбрать** конкретного донора. Выбор реципиента = финальное принятие (донору уходит уведомление «вас выбрал реципиент»).

**Текущее поведение (что НЕ меняем).**
- Старый flow «донор сам откликается → реципиент принимает» остаётся работать как есть.
- `apply_request.go` (отклик донора) и `accept_response.go` (принятие реципиентом) не трогаем.
- `SuitableDonors` (число подходящих) в DTO оставляем для обратной совместимости с фронтом.

**Что НЕ входит в задачу.**
- Изменение/удаление старого flow (донор сам откликается).
- Гео-фильтрация по расстоянию (только регионы из `preferred_location_ids`).
- Rate-limit и пагинация с курсором (на старте простой offset/limit).
- Изменение семантики `SuitableDonors`.
- Двухстороннее подтверждение (приглашение с ожиданием ответа донора) — НЕ делаем.

---

## 2. Зафиксированные решения

| # | Вопрос | Решение |
|---|---|---|
| 1 | Имя нового поля в `DonorPreference` | `OpenForContact` (Go) / `open_for_contact` (БД) / «Open for Contact» (UI) |
| 2 | Когда реципиент выбирает донора | Финальное принятие. `DonorResponse` создаётся сразу со статусом `accepted`. Никаких новых статусов. |
| 3 | Старый flow (донор сам откликается) | Оставляем работающим. В списке потенциальных исключаем уже откликнувшихся. |
| 4 | Регионы для фильтрации | Используем существующее `preferred_location_ids` в `DonorPreference`. Семантика для потенциальных — **позитивная**: пустой массив = донор НЕ открыт (в отличие от `GetPetsByBloodGroupAndRegion`, где пустой = «любая локация»). |

---

## 3. Архитектурные договорённости

### 3.1 Слои и направление зависимостей

Код добавляется строго по слоям DDD, сверху вниз:

```
Ent schema → миграция БД
  → domain model (DonorPreference в user-контексте)
    → repository (новый метод FindPotentialDonors)
      → application (query + cmd handlers)
        → transport (HTTP endpoints + DTO + mapper)
          → DI wiring (cmd/api/main.go)
```

### 3.2 Что переиспользуем без изменений

- `internal/domain/bloodsearch/bloodcounter_svc.go::RecalculateBloodAmount` — пересчёт объёмов крови.
- `internal/domain/bloodsearch/model/bloodsearch_model.go::RecalculateStatus` — пересчёт статуса заявки.
- `internal/domain/donor/model/donor_model.go::NewDonorResponse` — конструктор (но вызываем напрямую со статусом `accepted`, см. п. 6.2).
- `internal/domain/bonus/bonus_service.go::AssignBonuses` — начисление бонусов донору.
- `internal/domain/pet/enrich/enricher.go::PetEnricher.RecalculateAll` / `RecalculateOne` — пересчёт актуального `PetStatus` (донор/не донор) с учётом последних `DonorResponse` и активных `BloodRequest`. **Использовать обязательно**: в БД поле `pet.status` может быть устаревшим, поэтому мы **не доверяем** SQL-фильтру по `pet.status` как единственному критерию «донор может сдавать сейчас».
- `internal/domain/bloodsearch/events/ApplyDonor` — событие «донор принят реципиентом» (переиспользуем для нового flow без изменений).
- `internal/transport/http/dto/donor_dto.go::DonorApplication` — DTO карточки донора. Используем **тот же тип** для потенциальных, с пустым `ID`. Создавать новый DTO не нужно.

### 3.3 Аутентификация и авторизация

Новый endpoint `POST /v1/blood-request/{req_id}/donor/select` должен работать только от лица владельца питомца-реципиента:
- Извлекаем `userID` из контекста (middleware `Auth` уже кладёт).
- Внутри handler-а получаем `BloodRequest` по `req_id`, проверяем `BloodRequest.PetID → Pet.OwnerID == userID`. Если нет — `apperrors.Forbidden(...)`.

GET `GET /v1/blood-requests/by-pet/{pet_id}` (уже существующий) — потенциальные доноры возвращаются только владельцу питомца-реципиента. HTTP-layer вытягивает `userID` из контекста и передаёт в query-handler; query-handler требует `callerUserID` (пустой — `Unauthorized`).

---

## 4. Слой 1 — Схема БД ✅

### 4.1 Файл `backend/internal/infra/ent/schema/donor_preference.go`

В методе `Fields()` добавлено после `notification_frequency`:

```go
field.Bool("open_for_contact").
    Default(true).  // ← изменено с false на true по решению пользователя
    Comment("Разрешает реципиентам находить донора как потенциального и приглашать"),
```

### 4.2 Перегенерация Ent

Выполнено:
```bash
cd backend
go generate ./internal/infra/ent/...
```

Обновлены `internal/infra/ent/donorpreference/...` и `internal/infra/ent/migrate/schema.go`.

### 4.3 Миграция

⚠️ SQL-файл миграции для `open_for_contact` НЕ создавался — `go generate` только обновил схему. Применяется через atlas или вручную.

**Открытая задача:** приложить `atlas migrate new open_for_contact` или SQL-файл. Подалажётся в §15.

---

## 5. Слой 2 — Domain model ✅

### 5.1 Файл `backend/internal/domain/user/model/user_model.go`

В структуру `DonorPreference` добавлено поле `OpenForContact`.
В структуру `DonorPreferenceParams` добавлено поле `OpenForContact`.
В `DefaultDonorPreference()` проставлено `OpenForContact: true` (в плане было `false`, поправлено под default в БД).

### 5.2 Файл `backend/internal/infra/presistance/domainmapper/user_mapper.go`

Поле `dp.OpenForContact` проброшено в маппинг `Ent → Domain`.

### 5.3 Дополнительно: Transport DTO + DTO mapper

- `internal/transport/http/dto/user_dto.go` — поле `OpenForContact` добавлено в `DonorPreference` (response) и `DonorPreferenceParams` (input).
- `internal/transport/http/dtomapper/user_mapper_dto.go` — поле проброшено в маппинг `Domain → DTO`.

---

## 6. Слой 3 — Application handlers

### 6.1 Расширение `GetByPetIDHandler`

**Файл:** `backend/internal/application/bloodsearch/query/get_by_pet_id.go`

**✅ Статус: реализовано.**

Реализовано:
- Добавлена зависимость `petEnricher enrich.PetEnricher`.
- `Handle(ctx, callerUserID, petID)` — query-handler теперь принимает `callerUserID` и требует, чтобы `recipientPet.OwnerID == callerUserID` (иначе `Forbidden`).
- Батчевый `Fetch` + поштучный `Recalculate` для каждого потенциального донора; индивидуальный `RecoveryPeriodMonths` берётся из `PotentialDonor.RecoveryPeriodMonths`, который доменный маппер извлекает из `Owner.DonorPreference` (без доп. SQL-запроса).
- После enrich фильтрует оставшихся до `PetStatusDonor`.

Новая сигнатура:

```go
type GetByPetIDResult struct {
    BloodRequest    *model.BloodRequestWithApplications
    SuitableDonors  int
    PotentialDonors []*model.Pet
}

type GetByPetIDHandler struct {
    bloodRepo    bloodsearch.Repository
    petRepo      pet.Repository
    petEnricher  enrich.PetEnricher   // <-- НОВОЕ
    bloodCounter *bloodsearch.BloodCounterService
}

func NewGetByPetIDHandler(
    bloodRepo bloodsearch.Repository,
    petRepo pet.Repository,
    petEnricher enrich.PetEnricher,    // <-- НОВОЕ
) *GetByPetIDHandler {
    return &GetByPetIDHandler{
        bloodRepo:    bloodRepo,
        petRepo:      petRepo,
        petEnricher:  petEnricher,
        bloodCounter: bloodsearch.NewBloodCounterService(),
    }
}

func (h *GetByPetIDHandler) Handle(ctx context.Context, petID string) (*GetByPetIDResult, error) {
    bloodReq, err := h.bloodRepo.GetByPetID(ctx, petID)
    if err != nil {
        return nil, err
    }

    donated, reserved := h.bloodCounter.RecalculateBloodAmount(bloodReq.BloodRequest, bloodReq.DonorApplications)
    bloodReq.BloodRequest.SetBloodVolume(donated, reserved)

    suitableDonors, err := h.petRepo.CountSuitableDonors(ctx, bloodReq.BloodRequest.BloodGroupNames)
    if err != nil {
        return nil, err
    }

    // ====== НОВАЯ ЛОГИКА: потенциальные доноры ======
    potentialDonors, err := h.petRepo.FindPotentialDonors(ctx, pet.PotentialDonorsCriteria{
        PetType:          bloodReq.BloodRequest.Type(),       // метод домена BloodRequest
        BloodGroups:      bloodReq.BloodRequest.BloodGroupNames,
        Regions:          bloodReq.BloodRequest.Regions,
        ExcludeRequestID: bloodReq.BloodRequest.ID,
        ExcludePetID:     bloodReq.BloodRequest.PetID,
        Limit:            50,
        Offset:           0,
    })
    if err != nil {
        return nil, err
    }

    // Пересчитываем актуальный PetStatus (БД может быть устаревшей)
    fc, err := h.petEnricher.RecalculateAll(ctx, potentialDonors, enrich.Options{
        RecoveryPeriodMonths: 0, // 0 = берём дефолт питомца, см. enricher.go
    })
    if err != nil {
        return nil, err
    }
    _ = fc

    // Фильтруем только реальных доноров
    actual := make([]*model.Pet, 0, len(potentialDonors))
    for _, p := range potentialDonors {
        if p.PetStatus == model.PetStatusDonor {
            actual = append(actual, p)
        }
    }

    return &GetByPetIDResult{
        BloodRequest:    bloodReq,
        SuitableDonors:  suitableDonors,
        PotentialDonors: actual,
    }, nil
}
```

> **Замечание про `BloodRequest.Type()`.** В реализации **отказались** от геттера: `recipientPet.Type` приходит из `petRepo.GetByID` без отдельного запроса. Геттер можно вернуть, если в будущем потребуется избежать этого.

> **Замечание про `RecoveryPeriodMonths`.** В реализации передаём **доновский** `RecoveryPeriodMonths` из `DonorPreference` владельца донора — не 0. Актуальный период восстановления нужен, чтобы пересчёт `RecoveryDays` был корректным.

### 6.2 Новый cmd `SelectDonorHandler`

**Новый файл:** `backend/internal/application/bloodsearch/cmd/select_donor.go`

```go
package cmd

import (
    "context"
    "log/slog"
    "time"

    "github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
    "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
    "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/events"
    bloodmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
    "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
    "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
    donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
    "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
    petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
    "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
    "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
    "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type SelectDonorHandler struct {
    bloodRepo    bloodsearch.Repository
    donorRepo    donor.Repository
    petRepo      pet.Repository
    petEnricher  pet.Enricher      // интерфейс из internal/application/pet/enrich/enricher.go
    userRepo     user.Repository
    bonusSvc     *bonus.BonusService
    publisher    ports.EventPublisher
    txManager    *presistance.TxManager
    bloodCounter *bloodsearch.BloodCounterService
}

func NewSelectDonorHandler(
    bloodRepo bloodsearch.Repository,
    donorRepo donor.Repository,
    petRepo pet.Repository,
    petEnricher pet.Enricher,
    userRepo user.Repository,
    bonusSvc *bonus.BonusService,
    publisher ports.EventPublisher,
    txManager *presistance.TxManager,
) *SelectDonorHandler {
    return &SelectDonorHandler{
        bloodRepo:    bloodRepo,
        donorRepo:    donorRepo,
        petRepo:      petRepo,
        petEnricher:  petEnricher,
        userRepo:     userRepo,
        bonusSvc:     bonusSvc,
        publisher:    publisher,
        txManager:    txManager,
        bloodCounter: bloodsearch.NewBloodCounterService(),
    }
}

// Handle — реципиент выбирает конкретного донора из списка потенциальных.
// Создаёт DonorResponse сразу со статусом accepted (минуя pending),
// пересчитывает объёмы крови и статус заявки.
func (h *SelectDonorHandler) Handle(
    ctx context.Context,
    callerUserID string,
    requestID string,
    donorPetID string,
    compensationType string,
    taxiCompensation bool,
) (*donormodel.DonorResponse, error) {
    // 1. Получаем заявку
    req, err := h.bloodRepo.GetByID(ctx, requestID)
    if err != nil {
        return nil, err
    }
    if !req.IsActive() {
        return nil, apperrors.ErrInvalidBloodRequestStatus.WithMessage("blood request is not active")
    }

    // 2. Проверяем, что caller — владелец реципиента
    recipientPet, err := h.petRepo.GetByID(ctx, req.BloodRequest.PetID, pet.PetPreloadOptions{})
    if err != nil {
        return nil, apperrors.Internal(err, "failed to get recipient pet")
    }
    if recipientPet.OwnerID != callerUserID {
        return nil, apperrors.Forbidden("only the recipient owner can select donors")
    }

    // 3. Получаем питомца-донора и проверяем актуальный статус
    donorPet, err := h.petRepo.GetByID(ctx, donorPetID, pet.PetPreloadOptions{})
    if err != nil {
        return nil, apperrors.NotFound("donor pet not found")
    }
    // Пересчёт актуального статуса (БД может быть устаревшей)
    h.petEnricher.RecalculateOne(donorPet, nil, req, pet.EnrichOptions{RecoveryPeriodMonths: 0})
    if donorPet.PetStatus != petmodel.PetStatusDonor {
        return nil, apperrors.ErrInvalidBloodRequestStatus.WithMessage("donor pet is not currently available for donation")
    }

    // 4. Создаём DonorResponse напрямую со статусом accepted
    donorResponse, err := donormodel.NewDonorResponse(
        req.BloodRequest.ID,
        donorPet.ID,
        compensationType,
        donorPet.CalculateDonationAmount(),
        taxiCompensation,
    )
    if err != nil {
        return nil, apperrors.Validation(err.Error(), map[string]any{"field": "donor_response"})
    }
    donorResponse.Status = donormodel.DonorResponseStatusAccepted  // <-- сразу accepted, не pending

    // 5. В транзакции: создаём response + пересчитываем объёмы + статус заявки
    var petID string
    err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
        saved, err := h.donorRepo.CreateDonorResponse(txCtx, donorResponse)
        if err != nil {
            return err
        }
        donorResponse = saved

        if err := h.bonusSvc.AssignBonuses(txCtx, donorPet.OwnerID, donorPet.Type, donorResponse.ID); err != nil {
            return err
        }

        // Перечитываем заявку со всеми приложениями, чтобы пересчитать объёмы
        fresh, err := h.bloodRepo.GetByApplicationID(txCtx, donorResponse.ID, false)
        if err != nil {
            return err
        }
        // ^ GetByApplicationID находит BloodRequest по donor_response.ID.
        // Альтернатива: bloodRepo.GetByID + отдельный запрос списка applications.
        // Проверить, какая сигнатура реально есть; если GetByApplicationID не ищет
        // "все applications для данного request, включая только что созданный" — заменить
        // на h.bloodRepo.GetByID(txCtx, requestID) + h.donorRepo.GetByRequestID(txCtx, requestID).

        donated, reserved := h.bloodCounter.RecalculateBloodAmount(fresh.BloodRequest, fresh.DonorApplications)
        fresh.BloodRequest.SetBloodVolume(donated, reserved)
        fresh.RecalculateStatus()

        if err := h.bloodRepo.UpdateStatus(txCtx, fresh.BloodRequest.ID, fresh.BloodRequest.Status); err != nil {
            return err
        }
        petID = fresh.BloodRequest.PetID
        return nil
    })
    if err != nil {
        return nil, err
    }

    // 6. Уведомление донору (вне транзакции — ошибка некритична)
    recipientUser, err := h.userRepo.GetByID(ctx, recipientPet.OwnerID, user.UserPreloadOptions{WithIdentities: true})
    if err != nil {
        slog.Error("failed to get recipient user for notification", "err", err)
        return donorResponse, nil
    }
    donorUser, err := h.userRepo.GetByID(ctx, donorPet.OwnerID, user.UserPreloadOptions{WithIdentities: true})
    if err != nil {
        slog.Error("failed to get donor user for notification", "err", err)
        return donorResponse, nil
    }

    donorMaxID, donorTgID := donorUser.MessengerContacts()
    recipientMaxID, recipientTgID := recipientUser.MessengerContacts()

    donorData := events.DonorData{
        UserName:         donorUser.FullName,
        PetName:          donorPet.Name,
        Phone:            donorUser.Phone,
        BloodGroup:       donorPet.BloodGroupName,
        ProviderMaxID:    donorMaxID,
        ProviderTelegram: donorTgID,
    }
    recipientData := events.RecipientData{
        UserName:         recipientUser.FullName,
        PetName:          recipientPet.Name,
        Phone:            recipientUser.Phone,
        BloodGroup:       recipientPet.BloodGroupName,
        Volume:           req.BloodRequest.BloodVolumeNeeded,
        ProviderMaxID:    recipientMaxID,
        ProviderTelegram: recipientTgID,
    }
    if err := h.publisher.PublishEvent(ctx, ports.EventDonorApply, events.ApplyDonor{
        DonorData:     donorData,
        RecipientData: recipientData,
        CreatedAt:     time.Now(),
    }); err != nil {
        slog.Error("failed to publish donor selected notification", "err", err, "donorResponseID", donorResponse.ID)
    }

    _ = petID // если нужно для последующего enrich — пригодится
    return donorResponse, nil
}
```

> **Адаптации под реальный код:**
> - `pet.Enricher` / `pet.EnrichOptions` — проверить фактические имена интерфейса и опций в `internal/application/pet/enrich/enricher.go`. В этом пакете интерфейс называется `PetEnricher`, опции — `enrich.Options`. ✅ Поправлено.
> - `apperrors.Forbidden` / `apperrors.NotFound` — ✅ присутствуют в `internal/apperrors/errors.go`, используются как есть.
> - `donorPet.CalculateDonationAmount()` — ✅ метод уже есть, используется.
> - В транзакции перечитка заявки: выбран `h.bloodRepo.GetByID(txCtx, requestID)` + пересчёт объёмов на `fresh.DonorApplications`. Отдельный `donorRepo.GetDonorResponsesByRequestID` **не используется** — `BloodRequestWithApplications` уже приходит со всеми приложениями.

### 6.3 Возможные вспомогательные доменные методы

Если `BloodRequest.Type()` не существует — добавить в `internal/domain/bloodsearch/model/bloodsearch_model.go`:

```go
// Type возвращает тип питомца-реципиента.
// Требует предварительной загрузки через PetPreloadOptions.WithAll или отдельный запрос.
func (b *BloodRequest) Type() commonmodel.PetType {
    return b.PetType // поле, если уже есть; иначе хранить в BloodRequest
}
```

Если поля `PetType` в `BloodRequest` нет — добавить его в модель и в Ent-схему `blood_request.go` (миграция + перегенерация). Альтернатива — передавать тип в `FindPotentialDonorsCriteria` напрямую из `recipientPet.Type` в вызывающем коде, но это лишний запрос.

---

## 7. Слой 4 — Repository ✅

### 7.1 Расширение интерфейса `PetReadRepository`

**Файл:** `backend/internal/domain/pet/pet_repo.go`

В `PetReadRepository` добавлен метод `FindPotentialDonors`.
Добавлена структура `PotentialDonorsCriteria`.

### 7.2 Реализация в Ent

**Файл:** `backend/internal/infra/presistance/pg/pet_repo_ent.go`

Реализован `FindPotentialDonors`.

**Фактические отличия от плана:**
- **Ребро** в Ent называется `donations`, а не `donor_responses`. Использован `entpet.HasDonationsWith(...)` вместо `entpet.HasDonorResponsesWith(...)`.
- Добавлен импорт `entdonorresponse`.
- Убран комментарий про `else` для пустого `Regions` (упрощён).

---

## 8. Слой 5 — Transport (HTTP)

### 8.1 Новые DTO

**Файл:** `backend/internal/transport/http/dto/blood_request_dto.go`

В структуру `BloodRequestDetail` (строка ~140) добавить поле:

```go
PotentialDonors []DonorApplication `json:"potentialDonors,omitempty" doc:"Потенциальные доноры (ещё не откликнулись на заявку)"`
```

> `DonorApplication` уже определён в `backend/internal/transport/http/dto/donor_dto.go`. У него поле `ID` используется как `response_id`. Для потенциальных оно будет пустым (`""`). Это и есть фиксированное решение — не создавать новый DTO.

Новые input/output структуры в `blood_request_dto.go`:

```go
// ============================================
// Recipient selects donor (potential donor flow)
// ============================================

body:
```json
{ "donorId": "PET-ABCDEABCDE" }
```

Семантика: действие адресовано конкретной заявке (`req_id` в path), донор передаётся в body. Условия донации (тип компенсации, такси) не в body — берутся из `DonorPreference` владельца донора.

### 8.2 Новый endpoint в handler-е

```go
type SelectDonorInput struct {
    commondto.BloodRequestIDPath  // {req_id} из URL
    Body SelectDonorBody
}

type SelectDonorBody struct {
    DonorID string `json:"donorId" doc:"ID питомца-донора" minLength:"1" example:"PET-ABCDEABCDE"`
}

type SelectDonorOutput struct {
    Body DonorApplication
}
```

### 8.2 Новый endpoint в handler-е

**Файл:** `backend/internal/transport/http/bloodsearch_handler.go`

В структуре `BloodRequestHandler` добавить поле:

```go
selectDonorHandler *cmd.SelectDonorHandler
```

В конструктор добавить параметр.

В методе `Register(api huma.API)` (строка ~75) добавить после существующих маршрутов:

```go
huma.Register(api, huma.Operation{
    OperationID: "select-donor",
    Method:      http.MethodPost,
    Path:        "/v1/blood-request/{req_id}/donor/select",
    Summary:     "Реципиент выбирает донора из списка потенциальных",
    Tags:        []string{"BloodRequests"},
}, h.SelectDonor)
```

Реализация:

```go
func (h *BloodRequestHandler) SelectDonor(ctx context.Context, input *dto.SelectDonorInput) (*dto.SelectDonorOutput, error) {
    userID, ok := middleware.UserIDFromContext(ctx)
    if !ok {
        return nil, apperrors.Unauthorized("missing user context")
    }

    resp, err := h.selectDonorHandler.Handle(
        ctx,
        userID,
        input.Body.RequestID,
        input.Body.DonorID,
        input.Body.CompensationType,
        input.Body.TaxiCompensation,
    )
    if err != nil {
        return nil, err
    }

    // Маппим DonorResponse → DonorApplication. Контакты донора НЕ включаем —
    // маппер должен использовать облегчённый набор полей.
    return &dto.SelectDonorOutput{Body: h.bloodRequestMapper.DonorResponseToApplication(resp)}, nil
}
```

> **Имя метода `UserIDFromContext`.** Проверить фактический API middleware в `backend/internal/transport/http/middleware/`. Обычно это `middleware.UserIDFromContext(ctx)` или прямой доступ к значению.

### 8.3 Расширение `GetBloodRequestByPetID`

```go
func (h *BloodRequestHandler) GetBloodRequestByPetID(ctx context.Context, input *commondto.PetIDPath) (*dto.GetBloodRequestByPetIDOutput, error) {
    userID, ok := middleware.UserIDFromContext(ctx)
    if !ok {
        return nil, apperrors.Unauthorized("missing user context")
    }
    // Проверить, что caller — владелец питомца-реципиента
    // (если этой проверки ещё нет — добавить через petRepo.GetByID)
    // ... existing check ...

    result, err := h.getByPetIDHandler.Handle(ctx, input.ID)
    if err != nil {
        return nil, err
    }

    return &dto.GetBloodRequestByPetIDOutput{
        Body: h.bloodRequestMapper.ToResponseWithPotential(
            result.BloodRequest,
            &result.SuitableDonors,
            result.PotentialDonors,
        ),
    }, nil
}
```

### 8.4 Расширение mapper

**Файл:** `backend/internal/transport/http/dtomapper/bloodreq_mapper_dto.go`

Добавить новый метод `ToResponseWithPotential` (или расширить существующий `ToResponse`):

```go
func (m *BloodRequestMapper) ToResponseWithPotential(
    req *model.BloodRequestWithApplications,
    suitableDonors *int,
    potentialDonors []*petmodel.Pet,
) dto.BloodRequestDetail {
    detail := m.ToResponse(req, suitableDonors)  // существующий метод

    if len(potentialDonors) == 0 {
        return detail
    }

    potential := make([]dto.DonorApplication, 0, len(potentialDonors))
    for _, p := range potentialDonors {
        photos := m.storage.BuildPhotoURLs(p.Photos, time.Now())
        potential = append(potential, dto.DonorApplication{
            // ID оставляем пустым — это и есть маркер "потенциальный"
            ID:               "",
            RequestID:        req.BloodRequest.ID,
            DonorID:          p.ID,
            DonorName:        p.Name,
            DonorPhotos:      photos,
            DonorBloodGroup:  p.BloodGroupName,
            Amount:           p.CalculateDonationAmount(),
            WarnFactors:      buildWarnFactors(p.WarnFactors),
            CompensationType: "",  // потенциальный донор ещё не выбрал условия
            TaxiCompensation: false,
            Status:           "",
            IsConfirmed:      false,
        })
    }
    detail.PotentialDonors = potential
    return detail
}
```

> `buildWarnFactors` — вспомогательная функция, конвертирующая `[]string` в `[]dto.RestrictionFactor` через `petmodel.GetFactorDescription`. Логика уже есть в существующем `ToResponse` (см. как маппер работает для `WarnFactors` у `DonorApplication`). Вынести в общую приватную функцию или скопировать.

Дополнительный маппер для ответа на `SelectDonor`:

```go
func (m *BloodRequestMapper) DonorResponseToApplication(resp *donormodel.DonorResponse) dto.DonorApplication {
    return dto.DonorApplication{
        ID:               resp.ID,
        RequestID:        resp.RequestID,
        DonorID:          resp.DonorID,
        DonorName:        resp.DonorName,
        DonorPhotos:      resp.DonorPhotos,
        DonorBloodGroup:  resp.DonorBloodGroup,
        Amount:           resp.Amount,
        WarnFactors:      nil, // не загружены в этом сценарии
        CompensationType: resp.CompensationType,
        TaxiCompensation: resp.TaxiCompensation,
        Status:           string(resp.Status),
        IsConfirmed:      resp.IsConfirmed,
        CreatedAt:        resp.CreatedAt,
        UpdatedAt:        resp.UpdatedAt,
    }
}
```

---

## 9. Слой 6 — DI wiring

**Файл:** `backend/cmd/api/main.go` (или эквивалентный файл сборки зависимостей)

Найти место, где создаются handlers, и добавить:

```go
// PetEnricher уже создаётся где-то выше
// petEnricher := enrich.New(donorRepo, bloodRepo)

// Обновлённый GetByPetIDHandler
getByPetIDHandler := query.NewGetByPetIDHandler(
    bloodRepo,
    petRepo,
    petEnricher,
)

// Новый SelectDonorHandler
selectDonorHandler := cmd.NewSelectDonorHandler(
    bloodRepo,
    donorRepo,
    petRepo,
    petEnricher,
    userRepo,
    bonusSvc,
    publisher,
    txManager,
)

// Передать в BloodRequestHandler
bloodRequestHandler := http.NewBloodRequestHandler(
    // ... существующие параметры ...
    getByPetIDHandler,
    selectDonorHandler,
    // ...
)
```

> Точные имена конструкторов и порядок параметров — взять из существующего кода. Не добавлять без проверки — структура wiring в проекте может отличаться.

---

## 10. Тесты

### 10.1 Repository: `FindPotentialDonors`

**Файл:** `backend/internal/infra/presistance/pg/pet_repo_ent_test.go` (создать, если нет).

Интеграционные тесты с тестовой БД (enttest). Покрыть:

| Сценарий | Ожидание |
|---|---|
| Все фильтры совпадают, донор открыт | Возвращается |
| `open_for_contact = false` | Не возвращается |
| `pet.status != donor` (в БД) | Возвращается, но после enrich отфильтровывается (проверяется на уровне handler-а) |
| Не пересекается `blood_group` | Не возвращается |
| Не пересекается регион | Не возвращается |
| `preferred_location_ids = []` | Не возвращается (позитивная семантика) |
| У донора уже есть `DonorResponse` на эту заявку | Не возвращается |
| `ExcludePetID` совпадает с донором | Не возвращается |
| Limit/Offset работают | Возвращается нужная страница |

### 10.2 Query handler: `GetByPetID`

**Файл:** `backend/internal/application/bloodsearch/query/get_by_pet_id_test.go`.

С моками репозиториев. Покрыть:

| Сценарий | Ожидание |
|---|---|
| Список потенциальных не пуст, после enrich все `PetStatusDonor` | Возвращаются все |
| После enrich один из питомцев стал `Recovering` | Отфильтрован |
| `FindPotentialDonors` возвращает ошибку | Проброс ошибки |
| `bloodRepo.GetByPetID` возвращает ошибку | Проброс ошибки |

### 10.3 Cmd handler: `SelectDonor`

**Файл:** `backend/internal/application/bloodsearch/cmd/select_donor_test.go`.

С моками. Покрыть:

| Сценарий | Ожидание |
|---|---|
| Успешный flow | `DonorResponse` создан со статусом `accepted`, бонусы начислены, объёмы пересчитаны, событие опубликовано |
| Заявка не активна | `ErrInvalidBloodRequestStatus` |
| Caller не владелец реципиента | Forbidden error |
| Донор не имеет статуса `donor` (по актуальному enrich) | Validation error |
| Внутри `WithTx` упала вставка | Откат транзакции, ошибка проброшена |
| После `WithTx` не удалось получить контакты (Redis упал) | Ответ всё равно возвращается, событие логируется как failed |
| Переполнение объёма (overcommit) | Транзакция откатывается, ошибка Conflict |

### 10.4 HTTP handler: `SelectDonor`

Сквозной тест через `httptest`. Проверить:
- 401 без авторизации.
- 403 если caller не владелец.
- 200 с правильным телом при успехе.

### 10.5 DTO mapper

Юнит-тест на `ToResponseWithPotential`:
- Пустой `potentialDonors` → поле `PotentialDonors` = nil или пустой массив.
- Не пустой → корректные `DonorApplication` с пустым `ID` и `Status`.

---

## 11. Чек-лист приёмки

### Schema / DB
- [ ] Поле `open_for_contact` добавлено в `donor_preferences` (миграция применена)
- [ ] Ent-код перегенерирован
- [ ] Существующие записи имеют `open_for_contact = false` (default)

### Domain
- [ ] `usermodel.DonorPreference.OpenForContact` существует
- [ ] `user_mapper.go` пробрасывает поле
- [ ] `DefaultDonorPreference()` возвращает `OpenForContact: false`

### Repository
- [ ] `PetReadRepository.FindPotentialDonors` объявлен в интерфейсе
- [ ] Реализация в `EntPetRepository` соответствует SQL-фильтрам из п. 7.2
- [ ] Покрыт тестами по п. 10.1

### Application
- [ ] `GetByPetIDHandler.Handle` возвращает `*GetByPetIDResult` со списком потенциальных
- [ ] Список отфильтрован по актуальному `PetStatus == donor` после `RecalculateAll`
- [ ] `SelectDonorHandler.Handle` существует и работает по транзакционной последовательности п. 6.2
- [ ] Авторизация (caller == владелец реципиента) реализована
- [ ] Бонусы начисляются донору
- [ ] Событие `EventDonorApply` публикуется донору

### Transport
- [ ] `BloodRequestDetail.PotentialDonors` присутствует в DTO и OpenAPI
- [ ] `POST /v1/blood-request/{req_id}/donor/select` зарегистрирован
- [ ] DTO `SelectDonorInput` / `SelectDonorOutput` валидны (Huma-аннотации)
- [ ] Mapper `ToResponseWithPotential` возвращает потенциальных с пустым `ID`/`Status`

### DI
- [ ] `getByPetIDHandler` собирается с `petEnricher`
- [ ] `selectDonorHandler` собирается со всеми зависимостями
- [ ] Приложение стартует без ошибок

### Тесты
- [ ] `go test ./...` проходит локально
- [ ] Интеграционные тесты с тестовой БД — все сценарии из п. 10.1
- [ ] Линтер (`golangci-lint`) — 0 ошибок

### Документация
- [ ] OpenAPI (`docs/openapi.json`) перегенерирована
- [ ] Новые endpoint-ы видны в Scalar UI

---

## 12. Что НЕ делаем (жёсткий список)

1. Не трогаем `apply_request.go`, `accept_response.go`, `cancel_donation.go`, `confirm_donation.go`, `reject_donation.go`.
2. Не добавляем новый статус в `DonorResponseStatus`.
3. Не создаём новый DTO-тип для потенциальных — переиспользуем `DonorApplication` с пустым `ID`.
4. Не удаляем поле `SuitableDonors` из DTO.
5. Не вводим гео-фильтрацию (только регионы).
6. Не вводим rate-limit / курсорную пагинацию (offset/limit достаточно).
7. Не делаем двухстороннее подтверждение (приглашение с ожиданием донора).

---

## 13. Риски (напоминание)

1. **Дублирование в UI**: один и тот же питомец может попасть и в `Responses` (уже откликнулся), и в `PotentialDonors`. Решено через `NOT EXISTS` в SQL-фильтре — но только для текущей заявки. Если донор откликнулся на ДРУГУЮ заявку реципиента, он всё равно может появиться здесь. Это допустимо (другой `request_id`).

2. **Гонка при одновременном выборе**: два реципиента (или два запроса) одновременно выбирают донора. `txManager.WithTx` сериализует на уровне БД; `RecalculateBloodAmount` после создания `DonorResponse` поймает overcommit и откатит транзакцию. Дополнительная защита через оптимистичную блокировку по `BloodRequest.UpdatedAt` — вне scope, отметить как TODO.

3. **Приватность**: потенциальные доноры видны только владельцу питомца-реципиента с активной заявкой. Это обеспечивается проверкой авторизации в `GetBloodRequestByPetID`. Контакты донора (телефон, мессенджер) в списке потенциальных НЕ отдаём — только базовые поля карточки. Полные контакты — только после успешного `SelectDonor`.

---

## 14. Порядок реализации (рекомендуемый)

1. ✅ Схема БД + Ent schema → `go generate` (п. 4).
2. ✅ Domain model `DonorPreference.OpenForContact` + mapper + DTO + DTO mapper (п. 5 + расширение).
3. ✅ `FindPotentialDonors` в репозитории (п. 7).
   - **Отличие от плана**: возвращает `[]*model.PotentialDonor`, а не `[]*model.Pet`.
   - Добавлен `PotentialDonor` — обёртка над `Pet` + `CompensationType`, `TaxiCompensation`, `RecoveryPeriodMonths` из `Owner.DonorPreference`.
   - Добавлен `PetToPotentialDonor`/`PetToPotentialDonorSlice` в domainmapper.
4. ✅ Расширение `GetByPetIDHandler` (п. 6.1).
   - **Отличие от плана**: возвращает `*GetByPetIDResult` с `PotentialDonors []*petmodel.PotentialDonor`.
   - Использует `petEnricher.Fetch` + per-pet `Recalculate` с индивидуальным `RecoveryPeriodMonths`.
   - Тип питомца получает через `petRepo.GetByID` (в `BloodRequest` нет поля `PetType`).
5. ✅ DTO `BloodRequestDetail.PotentialDonors` + mapper (п. 8.1, 8.4).
   - `ToResponseWithPotential` маппит `*GetByPetIDResult` в DTO.
   - `DonorResponseToApplication` принимает `*petmodel.PotentialDonor` и использует `CompensationType`/`TaxiCompensation` из донорских настроек.
6. ✅ `SelectDonorHandler` (п. 6.2).
   - **Отличие от плана**: в транзакции использует `GetByID` (`fresh.DonorApplications` приходит загруженным, без отдельного `GetDonorResponsesByRequestID`).
   - Статус `accepted` устанавливается через `donorResponse.Accept()`, а не прямым присваиванием.
   - Пересчёт статуса донора делается через `petEnricher.Fetch([donorPet.ID])` + `Recalculate` **по его собственным данным** (НЕ по заявке реципиента — иначе `peekStatus` лепит донору `Recipient`/`BloodFound`).
   - `RecoveryPeriodMonths` для enrich берётся из `DonorPreference` владельца донора (не `0`).
   - Условия донации (`CompensationType`/`TaxiCompensation`) тоже берутся из `DonorPreference` владельца донора — реципиент их не задаёт (поля убраны из `SelectDonorBody`).
7. ✅ HTTP endpoint `POST /v1/blood-request/{req_id}/donor/select` (п. 8.2) + DTOs. Действие адресовано конкретной заявке (`req_id` в path совпадает с `path:"req_id"` в `commondto.BloodRequestIDPath`), донор передаётся в body. Путь выровнен с архитектурой остальных bloodrequest-эндпоинтов (singleton resource + action). Ownership-проверка реципиента добавлена и в `GetBloodRequestByPetID` через `callerUserID` в `GetByPetIDHandler.Handle`.
8. ✅ DI wiring (п. 9).
9. ⬜ Прогон `go test ./...` + линтер + пересборка OpenAPI + миграция БД.

Каждый шаг — отдельный коммит.

---

## 15. Что осталось (открытые задачи)

### Блокер деплоя
- [ ] **Миграция БД**: создать SQL/atlas-файл для `open_for_contact` в `donor_preferences`. Применить на dev/staging/prod.

### Желательно до мержа
- [ ] **Тесты §10.1-§10.5** (все нулевые сейчас). Минимум — unit-тесты для:
  - `GetByPetIDHandler` с моками (`bloodRepo`/`petRepo`/`petEnricher`) на happy path + forbidden + enrich-failure.
  - `SelectDonorHandler` happy path (через stub `TxManager` + моки репозиториев) — главное покрыть «accepted через доменный метод», forbidden, donor-not-found, donor-not-available.
  - Mapper `ToResponseWithPotential` + `DonorResponseToApplication` на nil/empty списках.
  - §10.4 HTTP-тест через `httptest` — отдельная задача.
- [ ] **CQRS Read/Write сплит**: единичные `select_donor.go`/`get_by_pet_id.go` уже используют `pet.PetReadRepository`, но остальные соседние cmd-хэндлеры (`accept_response`, `confirm_donation`, `create_req`, `reject_donation`) ещё на `pet.Repository`. Эта унификация лежит отдельно от текущей ветки.
- [x] ~~**Двойной `userRepo.GetByID` владельца донора**~~ — ✅ объединено в один вызов `{WithDonorPreference: true, WithIdentities: true}`, `donorUser` заменён на `donorOwner` для уведомления.
- [x] ~~**`recoveryMonths` без fallback**~~ — ✅ дефолт `2` (соответствует `DefaultDonorPreference()`), используется если у владельца нет `DonorPreference` или `RecoveryPeriodMonths == 0`.
- [x] ~~**`callerUserID == ""` в `GetByPetIDHandler`**~~ — ✅ ужесточено до `Unauthorized` вместо тихого пропуска.
- [x] ~~**Полировка: `DonorResponseToSelected`**~~ — ✅ поле `WarnFactors` маппера оставлено как есть; это мистический нонсенс, но правка лежит отдельно от текущей ветки (см. ниже).
- [ ] **`DonorResponseToSelected` читает пустые `WarnFactors`** — мистика: свежесозданный `DonorResponse` приходит без варн-факторов. Либо убрать поле из маппера, либо обогащать в handler-е. Полировка.
- [ ] **`events.DonorData.UserName` фактически содержит имя питомца** — путаница имён в существующем event-структе (мимоходом, не в этом PR).

### Документация
- [ ] OpenAPI перегенерирован (`go run ./cmd/api/main.go openapi`) — ✅ локально; проверить что CI тоже гоняет.
- [ ] TS-клиент обновлён (`task generate-ts`) — ✅ локально.

### Полировка кода (некритично)
- [ ] Маппер `DonorResponseToSelected` читает `resp.WarnFactors`, но свежесозданный `DonorResponse` приходит из `CreateDonorResponse` без варн-факторов. Поле всегда пустое — стоит либо убрать его из маппера, либо обогащать на стороне cmd-хэндлера.
- [ ] `events.DonorData.UserName` фактически содержит имя питомца (`donorPet.Name`) — существующая путаница имён (не в этом PR, мимоходом).
