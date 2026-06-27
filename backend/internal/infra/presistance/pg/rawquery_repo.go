package pg

import (
	"context"
	"database/sql"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/analytics/query"
)

type RawQueryRepo struct {
	db *sql.DB
}

func NewRawQueryRepository(client *sql.DB) *RawQueryRepo {
	return &RawQueryRepo{
		db: client,
	}
}

func (r *RawQueryRepo) GetPortalStats(ctx context.Context) (*query.PortalStats, error) {
	newquery := `
		select
			count(*) as total_users,

			count(phone) as users_with_phone,

			round(
				count(phone)::numeric / nullif(count(*), 0) * 100,
				2
			) as phone_conversion_percent,

			count(*) filter (where verified = true) as verified_users,

			count(*) filter (where verified = false) as unverified_users,

			(
				select count(*)
				from pets
			) as total_pets,

			(
				select count(*)
				from blood_requests
				where status = 'active'
			) as active_blood_requests,

			(
				select count(*)
				from donor_responses
			) as total_donations,

			(
				select count(*)
				from donor_responses
				where status = 'completed'
						and is_confirmed = true
			) as completed_donations

		from users;
	`

	row := r.db.QueryRowContext(ctx, newquery)

	var s query.PortalStats

	err := row.Scan(
		&s.TotalUsers,
		&s.UsersWithPhone,
		&s.PhoneConversionPercent,
		&s.VerifiedUsers,
		&s.UnverifiedUsers,
		&s.TotalPets,
		&s.ActiveBloodRequests,
		&s.TotalDonations,
		&s.CompletedDonations,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
