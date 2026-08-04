# 🏗️ Go DDD Bible — Clean Architecture + DDD + Hexagonal на Go

Практическое руководство по построению backend-приложений на Go с использованием Clean Architecture, Domain-Driven Design и Hexagonal Architecture (Ports & Adapters). Все примеры на идиоматичном Go.

---

## Содержание

1. [Правило зависимостей](#1-правило-зависимостей)
2. [Структура директорий](#2-структура-директорий)
3. [Domain Layer — Бизнес-логика](#3-domain-layer--бизнес-логика)
4. [Application Layer — Сценарии использования](#4-application-layer--сценарии-использования)
5. [Infrastructure Layer — Адаптеры](#5-infrastructure-layer--адаптеры)
6. [Composition Root — Сборка зависимостей](#6-composition-root--сборка-зависимостей)
7. [CQRS — Разделение команд и запросов](#7-cqrs--разделение-команд-и-запросов)
8. [Domain Events и Outbox](#8-domain-events-и-outbox)
9. [Тестирование](#9-тестирование)
10. [Анти-паттерны](#10-анти-паттерны)
11. [Когда применять, а когда нет](#11-когда-применять-а-когда-нет)

---

## 1. Правило зависимостей

> **Зависимости направлены ТОЛЬКО внутрь.** Внешние слои зависят от внутренних, никогда наоборот.

```
Infrastructure → Application → Domain
   (адаптеры)    (use cases)   (ядро)
```

**Нарушения, которые нужно отлавливать:**

- ❌ Domain импортирует `database/sql`, `net/http`, ORM
- ❌ Контроллеры напрямую вызывают репозитории, минуя application-слой
- ❌ Сущности зависят от сервисов приложения

**Лакмусовая проверка:** если ты можешь запустить всю доменную логику из тестов без базы данных и HTTP — границы корректны.

---

## 2. Структура директорий

```
internal/
├── domain/                    # Ядро — бизнес-логика (НУЛЕВЫЕ внешние зависимости)
│   ├── order/                 # Агрегат
│   │   ├── order.go           # Агрегат-корень (entity)
│   │   ├── order_item.go      # Дочерняя сущность
│   │   ├── value_objects.go   # OrderID, Money, OrderStatus
│   │   ├── events.go          # Доменные события
│   │   ├── repository.go      # Интерфейс репозитория (driven port)
│   │   ├── service.go         # Доменный сервис
│   │   └── errors.go          # Ошибки домена
│   ├── customer/
│   │   └── ...
│   └── shared/
│       ├── money.go           # Общий Value Object
│       └── errors.go          # Базовые ошибки
│
├── application/               # Сценарии использования (оркестрация)
│   ├── order/
│   │   ├── place_order/
│   │   │   ├── command.go     # PlaceOrderCommand (DTO)
│   │   │   ├── handler.go     # PlaceOrderHandler
│   │   │   └── port.go        # PlaceOrderUseCase (driver port)
│   │   ├── get_order/
│   │   │   ├── query.go       # GetOrderQuery
│   │   │   ├── handler.go     # GetOrderHandler
│   │   │   └── dto.go         # OrderDTO
│   │   └── ship_order/
│   │       └── ...
│   └── shared/
│       ├── unit_of_work.go    # UnitOfWork (driven port)
│       └── event_publisher.go # EventPublisher (driven port)
│
├── infrastructure/            # Адаптеры (внешние системы)
│   ├── persistence/
│   │   └── postgres/
│   │       ├── order_repository.go  # Реализация репозитория
│   │       ├── unit_of_work.go      # Реализация UoW
│   │       └── mapper.go            # Domain ↔ DB маппинг
│   ├── messaging/
│   │   └── rabbitmq/
│   │       └── event_publisher.go   # Реализация паблишера
│   ├── http/
│   │   └── handler/
│   │       └── order_handler.go     # HTTP-контроллер (driver adapter)
│   └── config/
│       └── di.go                    # Composition Root / DI
│
└── main.go                    # Точка входа
```

---

## 3. Domain Layer — Бизнес-логика

### 3.1 Value Object

Неизменяемый объект, определяемый **значениями своих атрибутов**, а не идентичностью.

```go
// domain/shared/money.go
package shared

import (
    "errors"
    "fmt"
)

type Money struct {
    amount   int64  // храним в копейках/центах
    currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
    if amount < 0 {
        return Money{}, errors.New("money: amount cannot be negative")
    }
    if currency == "" {
        return Money{}, errors.New("money: currency is required")
    }
    return Money{amount: amount, currency: currency}, nil
}

func MustMoney(amount int64, currency string) Money {
    m, err := NewMoney(amount, currency)
    if err != nil {
        panic(err)
    }
    return m
}

func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, fmt.Errorf("money: currency mismatch: %s vs %s", m.currency, other.currency)
    }
    return NewMoney(m.amount+other.amount, m.currency)
}

func (m Money) Multiply(factor int64) Money {
    return MustMoney(m.amount*factor, m.currency)
}

func (m Money) Amount() int64    { return m.amount }
func (m Money) Currency() string { return m.currency }

func (m Money) Equals(other Money) bool {
    return m.amount == other.amount && m.currency == other.currency
}
```

```go
// domain/order/value_objects.go
package order

import "github.com/google/uuid"

type OrderID struct{ value string }

func NewOrderID() OrderID         { return OrderID{uuid.New().String()} }
func OrderIDFrom(s string) OrderID { return OrderID{s} }
func (id OrderID) String() string  { return id.value }
func (id OrderID) Equals(other OrderID) bool { return id.value == other.value }

type CustomerID struct{ value string }
func CustomerIDFrom(s string) CustomerID { return CustomerID{s} }
func (id CustomerID) String() string      { return id.value }

type ProductID struct{ value string }
func ProductIDFrom(s string) ProductID { return ProductID{s} }
func (id ProductID) String() string     { return id.value }

type OrderStatus int

const (
    StatusDraft     OrderStatus = iota // 0 — неизвестно/черновик
    StatusConfirmed                    // 1
    StatusShipped                      // 2
    StatusCancelled                    // 3
)

func (s OrderStatus) String() string {
    switch s {
    case StatusDraft:     return "draft"
    case StatusConfirmed: return "confirmed"
    case StatusShipped:   return "shipped"
    case StatusCancelled: return "cancelled"
    default:              return "unknown"
    }
}
```

### 3.2 Entity

Объект с **уникальной идентичностью**, которая сохраняется во времени.

```go
// domain/order/order_item.go
package order

import "odnoi-krovi-app/internal/domain/shared"

type OrderItem struct {
    id        string
    productID ProductID
    quantity  int
    unitPrice shared.Money
}

func NewOrderItem(productID ProductID, quantity int, unitPrice shared.Money) (*OrderItem, error) {
    if quantity <= 0 {
        return nil, ErrInvalidQuantity
    }
    return &OrderItem{
        id:        uuid.New().String(),
        productID: productID,
        quantity:  quantity,
        unitPrice: unitPrice,
    }, nil
}

func (oi *OrderItem) IncreaseQuantity(amount int) error {
    if amount <= 0 {
        return ErrInvalidQuantity
    }
    oi.quantity += amount
    return nil
}

func (oi *OrderItem) Subtotal() shared.Money {
    return oi.unitPrice.Multiply(int64(oi.quantity))
}

func (oi *OrderItem) ProductID() ProductID { return oi.productID }
func (oi *OrderItem) Quantity() int         { return oi.quantity }
```

### 3.3 Aggregate Root

Кластер сущностей и value objects с **единой границей консистентности**. Все изменения проходят только через корень агрегата.

```go
// domain/order/order.go
package order

import (
    "time"
    "odnoi-krovi-app/internal/domain/shared"
)

// Order — агрегат-корень
type Order struct {
    id         OrderID
    customerID CustomerID
    items      []*OrderItem
    status     OrderStatus
    createdAt  time.Time
    events     []DomainEvent // накопленные доменные события
}

// NewOrder — фабричный метод (конструктор агрегата)
func NewOrder(customerID CustomerID) *Order {
    o := &Order{
        id:         NewOrderID(),
        customerID: customerID,
        items:      make([]*OrderItem, 0),
        status:     StatusDraft,
        createdAt:  time.Now(),
    }
    o.addEvent(OrderCreated{OrderID: o.id, CustomerID: customerID})
    return o
}

// --- Поведенческие методы (бизнес-операции) ---

func (o *Order) AddItem(productID ProductID, quantity int, price shared.Money) error {
    if o.status == StatusCancelled || o.status == StatusShipped {
        return ErrCannotModifyOrder
    }
    if quantity <= 0 {
        return ErrInvalidQuantity
    }

    // Если товар уже есть — увеличиваем количество
    for _, item := range o.items {
        if item.productID.Equals(productID) {
            return item.IncreaseQuantity(quantity)
        }
    }

    item, err := NewOrderItem(productID, quantity, price)
    if err != nil {
        return err
    }
    o.items = append(o.items, item)
    return nil
}

func (o *Order) Confirm() error {
    if o.status != StatusDraft {
        return ErrInvalidStateTransition
    }
    if len(o.items) == 0 {
        return ErrEmptyOrder
    }

    o.status = StatusConfirmed
    o.addEvent(OrderConfirmed{OrderID: o.id, Total: o.Total()})
    return nil
}

func (o *Order) Ship(trackingNumber string) error {
    if o.status != StatusConfirmed {
        return ErrInvalidStateTransition
    }

    o.status = StatusShipped
    o.addEvent(OrderShipped{OrderID: o.id, TrackingNumber: trackingNumber})
    return nil
}

func (o *Order) Cancel(reason string) error {
    if o.status == StatusShipped {
        return ErrCannotCancelShippedOrder
    }

    o.status = StatusCancelled
    o.addEvent(OrderCancelled{OrderID: o.id, Reason: reason})
    return nil
}

// --- Запросы (не меняют состояние) ---

func (o *Order) Total() shared.Money {
    total := shared.MustMoney(0, "USD")
    for _, item := range o.items {
        total, _ = total.Add(item.Subtotal())
    }
    return total
}

func (o *Order) ItemCount() int {
    count := 0
    for _, item := range o.items {
        count += item.Quantity()
    }
    return count
}

// --- Доступ к полям ---

func (o *Order) ID() OrderID           { return o.id }
func (o *Order) CustomerID() CustomerID { return o.customerID }
func (o *Order) Status() OrderStatus    { return o.status }
func (o *Order) Items() []*OrderItem    { return o.items }

// --- Доменные события ---

func (o *Order) Events() []DomainEvent { return o.events }
func (o *Order) ClearEvents()          { o.events = nil }

func (o *Order) addEvent(event DomainEvent) {
    o.events = append(o.events, event)
}
```

### 3.4 Доменные события

```go
// domain/order/events.go
package order

import (
    "time"
    "odnoi-krovi-app/internal/domain/shared"
)

type DomainEvent interface {
    EventType() string
    OccurredAt() time.Time
}

type OrderCreated struct {
    OrderID    OrderID
    CustomerID CustomerID
    occurredAt time.Time
}

func (e OrderCreated) EventType() string   { return "order.created" }
func (e OrderCreated) OccurredAt() time.Time { return e.occurredAt }

type OrderConfirmed struct {
    OrderID    OrderID
    Total      shared.Money
    occurredAt time.Time
}

func (e OrderConfirmed) EventType() string   { return "order.confirmed" }
func (e OrderConfirmed) OccurredAt() time.Time { return e.occurredAt }

type OrderShipped struct {
    OrderID        OrderID
    TrackingNumber string
    occurredAt     time.Time
}

func (e OrderShipped) EventType() string    { return "order.shipped" }
func (e OrderShipped) OccurredAt() time.Time { return e.occurredAt }

type OrderCancelled struct {
    OrderID    OrderID
    Reason     string
    occurredAt time.Time
}

func (e OrderCancelled) EventType() string    { return "order.cancelled" }
func (e OrderCancelled) OccurredAt() time.Time { return e.occurredAt }
```

### 3.5 Интерфейс репозитория (Driven Port)

```go
// domain/order/repository.go
package order

import "context"

type Repository interface {
    FindByID(ctx context.Context, id OrderID) (*Order, error)
    Save(ctx context.Context, order *Order) error
    Delete(ctx context.Context, id OrderID) error
}
```

### 3.6 Доменные ошибки

```go
// domain/order/errors.go
package order

import "errors"

var (
    ErrEmptyOrder              = errors.New("order: cannot confirm empty order")
    ErrInvalidQuantity         = errors.New("order: quantity must be positive")
    ErrCannotModifyOrder       = errors.New("order: cannot modify order in current state")
    ErrInvalidStateTransition  = errors.New("order: invalid state transition")
    ErrCannotCancelShippedOrder = errors.New("order: cannot cancel shipped order")
)
```

### 3.7 Доменный сервис

Используется когда логика не принадлежит одной сущности, а затрагивает несколько.

```go
// domain/order/service.go
package order

import "odnoi-krovi-app/internal/domain/shared"

// PricingService — доменный сервис расчёта скидки
type PricingService struct{}

func (PricingService) CalculateDiscount(order *Order, isVIP bool) shared.Money {
    discount := shared.MustMoney(0, "USD")

    if order.ItemCount() > 10 {
        fivePercent, _ := order.Total().Multiply(5)
        discount, _ = discount.Add(fivePercent.Divide(100))
    }

    if isVIP {
        tenPercent, _ := order.Total().Multiply(10)
        vipDiscount, _ := discount.Add(tenPercent.Divide(100))
        discount = vipDiscount
    }

    // Максимальная скидка 20%
    maxDiscount, _ := order.Total().Multiply(20)
    maxDiscount = maxDiscount.Divide(100)

    if discount.Amount() > maxDiscount.Amount() {
        return maxDiscount
    }
    return discount
}
```

---

## 4. Application Layer — Сценарии использования

### 4.1 Driver Port (интерфейс use case)

```go
// application/order/place_order/port.go
package place_order

import (
    "context"
    "odnoi-krovi-app/internal/domain/order"
)

type UseCase interface {
    Execute(ctx context.Context, cmd Command) (order.OrderID, error)
}
```

### 4.2 Command DTO

```go
// application/order/place_order/command.go
package place_order

type OrderItemCommand struct {
    ProductID string
    Quantity  int
}

type Command struct {
    CustomerID string
    Items      []OrderItemCommand
}
```

### 4.3 Handler (реализация use case)

```go
// application/order/place_order/handler.go
package place_order

import (
    "context"
    "fmt"

    "odnoi-krovi-app/internal/domain/order"
    "odnoi-krovi-app/internal/domain/product"
    "odnoi-krovi-app/internal/application/shared"
)

type Handler struct {
    orderRepo   order.Repository
    productRepo product.Repository
    uow         shared.UnitOfWork
    eventPub    shared.EventPublisher
}

func NewHandler(
    orderRepo order.Repository,
    productRepo product.Repository,
    uow shared.UnitOfWork,
    eventPub shared.EventPublisher,
) *Handler {
    return &Handler{
        orderRepo:   orderRepo,
        productRepo: productRepo,
        uow:         uow,
        eventPub:    eventPub,
    }
}

func (h *Handler) Execute(ctx context.Context, cmd Command) (order.OrderID, error) {
    // Создаём агрегат
    o := order.NewOrder(order.CustomerIDFrom(cmd.CustomerID))

    // Добавляем товары
    for _, item := range cmd.Items {
        p, err := h.productRepo.FindByID(ctx, order.ProductIDFrom(item.ProductID))
        if err != nil {
            return order.OrderID{}, fmt.Errorf("finding product %s: %w", item.ProductID, err)
        }
        if p == nil {
            return order.OrderID{}, fmt.Errorf("product %s: not found", item.ProductID)
        }

        if err := o.AddItem(p.ID(), item.Quantity, p.Price()); err != nil {
            return order.OrderID{}, fmt.Errorf("adding item: %w", err)
        }
    }

    // Сохраняем агрегат
    if err := h.orderRepo.Save(ctx, o); err != nil {
        return order.OrderID{}, fmt.Errorf("saving order: %w", err)
    }

    // Публикуем доменные события
    for _, event := range o.Events() {
        if err := h.eventPub.Publish(ctx, event); err != nil {
            return order.OrderID{}, fmt.Errorf("publishing event: %w", err)
        }
    }

    return o.ID(), nil
}
```

### 4.4 Driven Ports (интерфейсы для инфраструктуры)

```go
// application/shared/unit_of_work.go
package shared

import "context"

type UnitOfWork interface {
    Begin(ctx context.Context) (context.Context, error)
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
}
```

```go
// application/shared/event_publisher.go
package shared

import "context"

type EventPublisher interface {
    Publish(ctx context.Context, event any) error
}
```

### 4.5 Query и DTO (Read Model)

```go
// application/order/get_order/query.go
package get_order

type Query struct {
    OrderID string
}

// application/order/get_order/dto.go
package get_order

type OrderItemDTO struct {
    ProductID string  `json:"product_id"`
    Quantity  int     `json:"quantity"`
    UnitPrice float64 `json:"unit_price"`
    Subtotal  float64 `json:"subtotal"`
}

type OrderDTO struct {
    ID         string         `json:"id"`
    CustomerID string         `json:"customer_id"`
    Status     string         `json:"status"`
    Items      []OrderItemDTO `json:"items"`
    Total      float64        `json:"total"`
    CreatedAt  string         `json:"created_at"`
}
```

```go
// application/order/get_order/handler.go
package get_order

import (
    "context"
    "fmt"
    "odnoi-krovi-app/internal/domain/order"
)

type ReadModel interface {
    FindByID(ctx context.Context, id string) (*OrderDTO, error)
}

type Handler struct {
    readModel ReadModel
}

func NewHandler(readModel ReadModel) *Handler {
    return &Handler{readModel: readModel}
}

func (h *Handler) Execute(ctx context.Context, q Query) (*OrderDTO, error) {
    dto, err := h.readModel.FindByID(ctx, q.OrderID)
    if err != nil {
        return nil, fmt.Errorf("finding order %s: %w", q.OrderID, err)
    }
    return dto, nil
}
```

---

## 5. Infrastructure Layer — Адаптеры

### 5.1 Реализация репозитория (Driven Adapter)

```go
// infrastructure/persistence/postgres/order_repository.go
package postgres

import (
    "context"
    "database/sql"
    "fmt"

    "odnoi-krovi-app/internal/domain/order"
)

type OrderRepository struct {
    db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

// compile-time проверка: реализует ли интерфейс
var _ order.Repository = (*OrderRepository)(nil)

func (r *OrderRepository) FindByID(ctx context.Context, id order.OrderID) (*order.Order, error) {
    // SELECT из БД и маппинг в доменный объект
    row := r.db.QueryRowContext(ctx,
        "SELECT id, customer_id, status, created_at FROM orders WHERE id = $1",
        id.String(),
    )

    var (
        orderID    string
        customerID string
        status     int
        createdAt  string
    )
    if err := row.Scan(&orderID, &customerID, &status, &createdAt); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("scanning order: %w", err)
    }

    // ... загрузка items, маппинг в доменный Order
    // Возвращаем восстановленный агрегат
    return nil, nil // упрощено
}

func (r *OrderRepository) Save(ctx context.Context, o *order.Order) error {
    // UPSERT в БД
    _, err := r.db.ExecContext(ctx,
        `INSERT INTO orders (id, customer_id, status, created_at)
         VALUES ($1, $2, $3, $4)
         ON CONFLICT (id) DO UPDATE SET status = $3`,
        o.ID().String(), o.CustomerID().String(), int(o.Status()), "now()",
    )
    if err != nil {
        return fmt.Errorf("saving order: %w", err)
    }
    return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id order.OrderID) error {
    _, err := r.db.ExecContext(ctx, "DELETE FROM orders WHERE id = $1", id.String())
    return err
}
```

### 5.2 HTTP-контроллер (Driver Adapter)

```go
// infrastructure/http/handler/order_handler.go
package handler

import (
    "encoding/json"
    "net/http"

    "odnoi-krovi-app/internal/application/order/place_order"
    "odnoi-krovi-app/internal/application/order/get_order"
)

type OrderHandler struct {
    placeOrder place_order.UseCase
    getOrder   *get_order.Handler
}

func NewOrderHandler(placeOrder place_order.UseCase, getOrder *get_order.Handler) *OrderHandler {
    return &OrderHandler{placeOrder: placeOrder, getOrder: getOrder}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        CustomerID string `json:"customer_id"`
        Items      []struct {
            ProductID string `json:"product_id"`
            Quantity  int    `json:"quantity"`
        } `json:"items"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    cmd := place_order.Command{
        CustomerID: req.CustomerID,
        Items:      make([]place_order.OrderItemCommand, len(req.Items)),
    }
    for i, item := range req.Items {
        cmd.Items[i] = place_order.OrderItemCommand{
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
        }
    }

    orderID, err := h.placeOrder.Execute(r.Context(), cmd)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"id": orderID.String()})
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
    orderID := r.PathValue("id")

    dto, err := h.getOrder.Execute(r.Context(), get_order.Query{OrderID: orderID})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if dto == nil {
        http.Error(w, "order not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(dto)
}
```

---

## 6. Composition Root — Сборка зависимостей

Все зависимости собираются в **одном месте** — точке входа приложения.

```go
// infrastructure/config/di.go
package config

import (
    "database/sql"

    "odnoi-krovi-app/internal/application/order/place_order"
    "odnoi-krovi-app/internal/application/order/get_order"
    "odnoi-krovi-app/internal/infrastructure/persistence/postgres"
    httphandler "odnoi-krovi-app/internal/infrastructure/http/handler"
)

type Container struct {
    PlaceOrderHandler *place_order.Handler
    GetOrderHandler   *get_order.Handler
    OrderHTTPHandler  *httphandler.OrderHandler
}

func NewContainer(db *sql.DB) *Container {
    // Driven adapters
    orderRepo := postgres.NewOrderRepository(db)
    productRepo := postgres.NewProductRepository(db)
    eventPub := postgres.NewOutboxEventPublisher(db)

    // Application handlers
    placeOrderHandler := place_order.NewHandler(orderRepo, productRepo, nil, eventPub)
    getOrderHandler := get_order.NewHandler(postgres.NewOrderReadModel(db))

    // Driver adapters
    orderHTTPHandler := httphandler.NewOrderHandler(placeOrderHandler, getOrderHandler)

    return &Container{
        PlaceOrderHandler: placeOrderHandler,
        GetOrderHandler:   getOrderHandler,
        OrderHTTPHandler:  orderHTTPHandler,
    }
}
```

```go
// main.go
package main

import (
    "database/sql"
    "log"
    "net/http"

    "odnoi-krovi-app/internal/infrastructure/config"
    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "postgres://localhost:5432/mydb?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    container := config.NewContainer(db)

    mux := http.NewServeMux()
    mux.HandleFunc("POST /orders", container.OrderHTTPHandler.Create)
    mux.HandleFunc("GET /orders/{id}", container.OrderHTTPHandler.Get)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

---

## 7. CQRS — Разделение команд и запросов

### Когда применять

| Применяй CQRS | Пропусти CQRS |
|---|---|
| Чтение и запись имеют **радикально разные** нагрузки | Простой CRUD |
| Сложные запросы, которые не ложатся на доменную модель | Чтение/запись похожи |
| Разные команды работают над read/write сторонами | Маленькая команда |
| Используется Event Sourcing | «На всякий случай» |

### Упрощённый CQRS (начни с этого)

Одна БД, разные модели для чтения и записи:

```go
// Write side — через агрегат
func (h *PlaceOrderHandler) Execute(ctx context.Context, cmd Command) (OrderID, error) {
    o := order.NewOrder(...)
    // ... бизнес-логика
    h.orderRepo.Save(ctx, o)
    return o.ID(), nil
}

// Read side — через оптимизированную read-модель
func (h *GetOrderHandler) Execute(ctx context.Context, q Query) (*OrderDTO, error) {
    // Прямой SQL-запрос с JOIN'ами, денормализованный
    return h.readModel.FindByID(ctx, q.OrderID)
}
```

### Read Model (проекция)

```go
// infrastructure/persistence/postgres/order_read_model.go
package postgres

type OrderReadModel struct {
    db *sql.DB
}

func (m *OrderReadModel) FindByID(ctx context.Context, id string) (*get_order.OrderDTO, error) {
    row := m.db.QueryRowContext(ctx, `
        SELECT o.id, o.customer_id, o.status, o.created_at,
               oi.product_id, oi.quantity, oi.unit_price
        FROM orders o
        LEFT JOIN order_items oi ON o.id = oi.order_id
        WHERE o.id = $1
    `, id)
    // ... маппинг в DTO
    return nil, nil
}
```

---

## 8. Domain Events и Outbox

### Проблема двойной записи

```go
// ❌ Плохо — crash между Save и Publish теряет события
h.orderRepo.Save(ctx, o)    // запись в БД
h.eventPub.Publish(events)  // публикация в брокер
// 💥 crash здесь — события потеряны!
```

### Outbox Pattern

События пишутся в outbox-таблицу **в той же транзакции**, что и агрегат.

```go
// infrastructure/persistence/postgres/outbox_event_publisher.go
package postgres

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
)

type OutboxEventPublisher struct {
    db *sql.DB
}

func NewOutboxEventPublisher(db *sql.DB) *OutboxEventPublisher {
    return &OutboxEventPublisher{db: db}
}

// SaveEvents сохраняет события в outbox в рамках переданной транзакции
func (p *OutboxEventPublisher) SaveEvents(ctx context.Context, tx *sql.Tx, events []order.DomainEvent) error {
    for _, event := range events {
        payload, err := json.Marshal(event)
        if err != nil {
            return fmt.Errorf("marshaling event: %w", err)
        }

        _, err = tx.ExecContext(ctx,
            `INSERT INTO outbox (event_type, payload, created_at)
             VALUES ($1, $2, NOW())`,
            event.EventType(), payload,
        )
        if err != nil {
            return fmt.Errorf("inserting outbox: %w", err)
        }
    }
    return nil
}

// Publish читает непроцессированные события и отправляет их в брокер
func (p *OutboxEventPublisher) Publish(ctx context.Context, event any) error {
    // В реальном коде: читаем из outbox, отправляем в RabbitMQ/Kafka,
    // помечаем как обработанное
    return nil
}
```

### Использование в handler'е

```go
func (h *Handler) Execute(ctx context.Context, cmd Command) (order.OrderID, error) {
    o := order.NewOrder(...)
    // ... бизнес-логика

    // Начинаем транзакцию
    tx, err := h.db.BeginTx(ctx, nil)
    if err != nil {
        return order.OrderID{}, err
    }
    defer tx.Rollback()

    // Сохраняем агрегат
    if err := h.orderRepo.Save(ctx, tx, o); err != nil {
        return order.OrderID{}, err
    }

    // Сохраняем события в outbox (в той же транзакции!)
    if err := h.outbox.SaveEvents(ctx, tx, o.Events()); err != nil {
        return order.OrderID{}, err
    }

    // Коммитим — агрегат и события атомарны
    if err := tx.Commit(); err != nil {
        return order.OrderID{}, err
    }

    return o.ID(), nil
}
```

---

## 9. Тестирование

### 9.1 Тесты доменного слоя (без моков!)

```go
// domain/order/order_test.go
package order_test

import (
    "testing"

    "odnoi-krovi-app/internal/domain/order"
    "odnoi-krovi-app/internal/domain/shared"
)

func TestOrder_AddItem(t *testing.T) {
    t.Run("adds item to empty order", func(t *testing.T) {
        o := order.NewOrder(order.CustomerIDFrom("cust-1"))
        price := shared.MustMoney(1000, "USD") // $10.00

        err := o.AddItem(order.ProductIDFrom("prod-1"), 2, price)

        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if len(o.Items()) != 1 {
            t.Errorf("expected 1 item, got %d", len(o.Items()))
        }
        if o.Items()[0].Quantity() != 2 {
            t.Errorf("expected quantity 2, got %d", o.Items()[0].Quantity())
        }
    })

    t.Run("increases quantity for existing product", func(t *testing.T) {
        o := order.NewOrder(order.CustomerIDFrom("cust-1"))
        price := shared.MustMoney(1000, "USD")

        o.AddItem(order.ProductIDFrom("prod-1"), 2, price)
        o.AddItem(order.ProductIDFrom("prod-1"), 3, price)

        if o.Items()[0].Quantity() != 5 {
            t.Errorf("expected quantity 5, got %d", o.Items()[0].Quantity())
        }
    })

    t.Run("rejects zero quantity", func(t *testing.T) {
        o := order.NewOrder(order.CustomerIDFrom("cust-1"))
        price := shared.MustMoney(1000, "USD")

        err := o.AddItem(order.ProductIDFrom("prod-1"), 0, price)

        if err == nil {
            t.Fatal("expected error, got nil")
        }
    })

    t.Run("rejects modification of cancelled order", func(t *testing.T) {
        o := order.NewOrder(order.CustomerIDFrom("cust-1"))
        price := shared.MustMoney(1000, "USD")
        o.AddItem(order.ProductIDFrom("prod-1"), 1, price)
        o.Cancel("test")

        err := o.AddItem(order.ProductIDFrom("prod-2"), 1, price)

        if err == nil {
            t.Fatal("expected error, got nil")
        }
    })
}

func TestOrder_Confirm(t *testing.T) {
    t.Run("confirms order with items", func(t *testing.T) {
        o := order.NewOrder(order.CustomerIDFrom("cust-1"))
        price := shared.MustMoney(1000, "USD")
        o.AddItem(order.ProductIDFrom("prod-1"), 1, price)

        err := o.Confirm()

        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if o.Status() != order.StatusConfirmed {
            t.Errorf("expected confirmed, got %s", o.Status())
        }
    })

    t.Run("rejects confirmation of empty order", func(t *testing.T) {
        o := order.NewOrder(order.CustomerIDFrom("cust-1"))

        err := o.Confirm()

        if err == nil {
            t.Fatal("expected error, got nil")
        }
    })
}

func TestOrder_Events(t *testing.T) {
    t.Run("emits OrderCreated on creation", func(t *testing.T) {
        o := order.NewOrder(order.CustomerIDFrom("cust-1"))

        events := o.Events()
        if len(events) != 1 {
            t.Fatalf("expected 1 event, got %d", len(events))
        }
        if events[0].EventType() != "order.created" {
            t.Errorf("expected 'order.created', got %s", events[0].EventType())
        }
    })
}
```

### 9.2 Тесты application-слоя (с моками)

```go
// application/order/place_order/handler_test.go
package place_order_test

import (
    "context"
    "testing"

    "odnoi-krovi-app/internal/application/order/place_order"
    "odnoi-krovi-app/internal/domain/order"
    "odnoi-krovi-app/internal/domain/shared"
)

// Mock репозитория
type mockOrderRepo struct {
    saved *order.Order
}

func (m *mockOrderRepo) FindByID(ctx context.Context, id order.OrderID) (*order.Order, error) {
    return nil, nil
}

func (m *mockOrderRepo) Save(ctx context.Context, o *order.Order) error {
    m.saved = o
    return nil
}

func (m *mockOrderRepo) Delete(ctx context.Context, id order.OrderID) error {
    return nil
}

// Mock паблишера событий
type mockEventPublisher struct {
    published []any
}

func (m *mockEventPublisher) Publish(ctx context.Context, event any) error {
    m.published = append(m.published, event)
    return nil
}

func TestHandler_Execute(t *testing.T) {
    t.Run("creates order and publishes events", func(t *testing.T) {
        orderRepo := &mockOrderRepo{}
        productRepo := &mockProductRepo{
            products: map[string]*product.Product{
                "prod-1": product.NewProduct(
                    product.ProductIDFrom("prod-1"),
                    "Test Product",
                    shared.MustMoney(1000, "USD"),
                ),
            },
        }
        eventPub := &mockEventPublisher{}

        handler := place_order.NewHandler(orderRepo, productRepo, nil, eventPub)

        cmd := place_order.Command{
            CustomerID: "cust-1",
            Items: []place_order.OrderItemCommand{
                {ProductID: "prod-1", Quantity: 2},
            },
        }

        orderID, err := handler.Execute(context.Background(), cmd)

        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if orderID.String() == "" {
            t.Error("expected non-empty order ID")
        }
        if orderRepo.saved == nil {
            t.Fatal("expected order to be saved")
        }
        if len(eventPub.published) == 0 {
            t.Error("expected events to be published")
        }
    })
}
```

### 9.3 Архитектурные тесты

```go
// Проверка, что domain не зависит от infrastructure
// Используй go-arch-lint или кастомную проверку в CI
```

```bash
# Проверка через go vet и кастомный анализатор
go vet ./internal/domain/...
```

---

## 10. Анти-паттерны

| Анти-паттерн | Проблема | Исправление |
|---|---|---|
| **Анемичная доменная модель** | Сущности — просто структуры с данными, логика в сервисах | Перенести поведение В сущности |
| **Репозиторий на каждую таблицу** | `OrderItemRepository` | Один репозиторий на АГРЕГАТ |
| **Утечка инфраструктуры** | Domain импортирует `database/sql` | Domain имеет НОЛЬ внешних зависимостей |
| **God Aggregate** | Слишком много сущностей, долгие транзакции | Разделить на меньшие агрегаты |
| **Пропуск Application-слоя** | Контроллеры напрямую вызывают репозитории | Всё через use case handler'ы |
| **CRUD-мышление** | Моделируем данные, а не поведение | Моделируем бизнес-операции |
| **Преждевременный CQRS** | Сложность без необходимости | Начни с простого, усложняй по мере необходимости |
| **Кросс-агрегатные транзакции** | Несколько агрегатов в одной транзакции | Доменные события + eventual consistency |

### Пример: анемичная модель → богатая модель

```go
// ❌ Анемичная модель — логика снаружи
type Order struct {
    ID     string
    Status string
    Items  []OrderItem
}

// Где-то в сервисе:
if order.Status == "draft" {
    order.Status = "confirmed"
}

// ✅ Богатая модель — логика внутри
func (o *Order) Confirm() error {
    if o.status != StatusDraft {
        return ErrInvalidStateTransition
    }
    if len(o.items) == 0 {
        return ErrEmptyOrder
    }
    o.status = StatusConfirmed
    o.addEvent(OrderConfirmed{OrderID: o.id})
    return nil
}
```

---

## 11. Когда применять, а когда нет

### ✅ Применяй Clean Architecture + DDD + Hexagonal когда:

- Сложная бизнес-логика с множеством правил
- Долгоживущая система (годы поддержки)
- Команда 5+ разработчиков
- Несколько точек входа (API, CLI, события, cron)
- Нужна возможность замены инфраструктуры (БД, брокер)
- Требуется высокое покрытие тестами

### ❌ Пропусти когда:

- Простой CRUD (большинство приложений)
- Прототип / MVP / одноразовый код
- Маленькая команда (1-2 разработчика)
- Короткоживущий проект
- Тривиальная бизнес-логика

### Лестница сложности (начинай с простого)

```
Уровень 1: Простая слоистая архитектура (Controller → Service → Repository)
   ↓ Когда бизнес-правила усложняются
Уровень 2: Доменная модель (сущности с поведением)
   ↓ Когда нужно несколько точек входа
Уровень 3: Hexagonal (Ports & Adapters)
   ↓ Когда чтение/запись радикально расходятся
Уровень 4: CQRS (раздельные read/write модели)
   ↓ Когда нужен полный аудит / темпоральные запросы
Уровень 5: Event Sourcing (храним события, вычисляем состояние)
```

**Не перескакивай уровни.** Каждый уровень добавляет сложность. Переходи на следующий только когда текущий доказал свою недостаточность.

---

## Быстрые правила

1. **Зависимости только внутрь:** `Infrastructure → Application → Domain`
2. **Domain — чистая бизнес-логика:** ни одного импорта из `database/sql`, `net/http`, ORM
3. **Один репозиторий на агрегат:** не на таблицу, не на сущность
4. **Поведение в сущностях:** `order.Confirm()`, а не `order.Status = "confirmed"`
5. **Интерфейсы определяет потребитель:** репозиторий объявляется в domain, реализуется в infrastructure
6. **Принимай интерфейсы, возвращай структуры:** конструкторы возвращают конкретные типы
7. **Composition Root в одном месте:** все зависимости собираются в `main.go` или `di.go`
8. **Тестируй domain без моков:** доменный слой не имеет зависимостей — тесты летают