package pg_test

import (
	"testing"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/analytics/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/pg"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawQueryRepo_GetPortalStats(t *testing.T) {
	godotenv.Load("../.env")

	ctx := t.Context()

	env := config.GetEnv("ENV", "local")
	_, rawdb, err := config.ConnectEnt(config.NewEntConfig(env))
	require.NoError(t, err, "подключение к БД")
	t.Cleanup(func() {
		if closeErr := rawdb.Close(); closeErr != nil {
			t.Logf("ошибка при закрытии rawdb: %v", closeErr)
		}
	})

	// Act
	repo := pg.NewRawQueryRepository(rawdb)
	handler := query.NewPortalStatsHandler(repo)
	stats, err := handler.Handle(ctx)

	// Assert
	require.NoError(t, err, "запрос статистики портала")
	require.NotNil(t, stats, "результат не должен быть nil")

	// Проверяем, что все поля имеют осмысленные значения
	assert.GreaterOrEqual(t, stats.TotalUsers, int64(0), "TotalUsers")
	assert.GreaterOrEqual(t, stats.UsersWithPhone, int64(0), "UsersWithPhone")
	assert.GreaterOrEqual(t, stats.PhoneConversionPercent, 0.0, "PhoneConversionPercent")
	assert.GreaterOrEqual(t, stats.VerifiedUsers, int64(0), "VerifiedUsers")
	assert.GreaterOrEqual(t, stats.UnverifiedUsers, int64(0), "UnverifiedUsers")
	assert.GreaterOrEqual(t, stats.TotalPets, int64(0), "TotalPets")
	assert.GreaterOrEqual(t, stats.ActiveBloodRequests, int64(0), "ActiveBloodRequests")
	assert.GreaterOrEqual(t, stats.TotalDonations, int64(0), "TotalDonations")
	assert.GreaterOrEqual(t, stats.CompletedDonations, int64(0), "CompletedDonations")

	// Логическая инвариантность: компоненты <= целое
	assert.LessOrEqual(t, stats.UsersWithPhone, stats.TotalUsers, "UsersWithPhone <= TotalUsers")
	assert.Equal(t, stats.VerifiedUsers+stats.UnverifiedUsers, stats.TotalUsers,
		"VerifiedUsers + UnverifiedUsers == TotalUsers")
	assert.LessOrEqual(t, stats.CompletedDonations, stats.TotalDonations,
		"CompletedDonations <= TotalDonations")
	assert.LessOrEqual(t, stats.PhoneConversionPercent, 100.0, "PhoneConversionPercent <= 100%%")

	t.Logf("PortalStats: %+v", stats)
}
