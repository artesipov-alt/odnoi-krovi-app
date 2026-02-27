package repositories

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
)

// TxManager управляет жизненным циклом транзакции
type TxManager struct {
	client *ent.Client
}

func NewTxManager(client *ent.Client) *TxManager {
	return &TxManager{client: client}
}

// WithTx оборачивает функцию в транзакцию и прокидывает её через контекст
func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Откат при панике
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	// Выполняем бизнес-логику, передавая транзакцию в контексте
	if err := fn(ent.NewTxContext(ctx, tx)); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("rollback error: %v (original error: %w)", rerr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
