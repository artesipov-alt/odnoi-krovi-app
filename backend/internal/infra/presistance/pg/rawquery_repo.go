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
				where deleted_at is null
			) as total_pets,

			(
				select count(*)
				from blood_requests
				where deleted_at is null
				  and status = 'active'
			) as active_blood_requests,

			(
				select count(*)
				from donor_responses
				where deleted_at is null
			) as total_donations,

			(
				select count(*)
				from donor_responses
				where deleted_at is null
				  and status = 'completed'
				  and is_confirmed = true
			) as completed_donations,

			(
				select count(*)
				from blood_requests
				where deleted_at is null
			) as total_searches,

			(
				select coalesce(sum(blood_volume_needed), 0)
				from blood_requests
				where deleted_at is null
			) as total_search_volume,

			(
				select coalesce(sum(amount), 0)
				from donor_responses
				where deleted_at is null
				  and status = 'completed'
				  and is_confirmed = true
			) as total_donation_volume,

			-- cat stats
			(
				select count(*)
				from pets
				where deleted_at is null
				  and type = 'cat'
			) as cat_total_pets,

			(
				select count(*)
				from blood_requests br
				join pets p on br.pet_id = p.id
				where br.deleted_at is null
				  and p.type = 'cat'
				  and br.status = 'active'
			) as cat_active_blood_requests,

			(
				select count(*)
				from donor_responses dr
				join blood_requests br on dr.request_id = br.id
				join pets p on br.pet_id = p.id
				where dr.deleted_at is null
				  and p.type = 'cat'
			) as cat_total_donations,

			(
				select count(*)
				from donor_responses dr
				join blood_requests br on dr.request_id = br.id
				join pets p on br.pet_id = p.id
				where dr.deleted_at is null
				  and dr.status = 'completed'
				  and dr.is_confirmed = true
				  and p.type = 'cat'
			) as cat_completed_donations,

			(
				select count(*)
				from blood_requests br
				join pets p on br.pet_id = p.id
				where br.deleted_at is null
				  and p.type = 'cat'
			) as cat_searches,

			(
				select coalesce(sum(br.blood_volume_needed), 0)
				from blood_requests br
				join pets p on br.pet_id = p.id
				where br.deleted_at is null
				  and p.type = 'cat'
			) as cat_search_volume,

			(
				select coalesce(sum(dr.amount), 0)
				from donor_responses dr
				join blood_requests br on dr.request_id = br.id
				join pets p on br.pet_id = p.id
				where dr.deleted_at is null
				  and dr.status = 'completed'
				  and dr.is_confirmed = true
				  and p.type = 'cat'
			) as cat_donation_volume,

			-- dog stats
			(
				select count(*)
				from pets
				where deleted_at is null
				  and type = 'dog'
			) as dog_total_pets,

			(
				select count(*)
				from blood_requests br
				join pets p on br.pet_id = p.id
				where br.deleted_at is null
				  and p.type = 'dog'
				  and br.status = 'active'
			) as dog_active_blood_requests,

			(
				select count(*)
				from donor_responses dr
				join blood_requests br on dr.request_id = br.id
				join pets p on br.pet_id = p.id
				where dr.deleted_at is null
				  and p.type = 'dog'
			) as dog_total_donations,

			(
				select count(*)
				from donor_responses dr
				join blood_requests br on dr.request_id = br.id
				join pets p on br.pet_id = p.id
				where dr.deleted_at is null
				  and dr.status = 'completed'
				  and dr.is_confirmed = true
				  and p.type = 'dog'
			) as dog_completed_donations,

			(
				select count(*)
				from blood_requests br
				join pets p on br.pet_id = p.id
				where br.deleted_at is null
				  and p.type = 'dog'
			) as dog_searches,

			(
				select coalesce(sum(br.blood_volume_needed), 0)
				from blood_requests br
				join pets p on br.pet_id = p.id
				where br.deleted_at is null
				  and p.type = 'dog'
			) as dog_search_volume,

			(
				select coalesce(sum(dr.amount), 0)
				from donor_responses dr
				join blood_requests br on dr.request_id = br.id
				join pets p on br.pet_id = p.id
				where dr.deleted_at is null
				  and dr.status = 'completed'
				  and dr.is_confirmed = true
				  and p.type = 'dog'
			) as dog_donation_volume

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
		&s.TotalSearches,
		&s.TotalSearchVolume,
		&s.TotalDonationVolume,
		// cat stats
		&s.CatStats.TotalPets,
		&s.CatStats.ActiveBloodRequests,
		&s.CatStats.TotalDonations,
		&s.CatStats.CompletedDonations,
		&s.CatStats.Searches,
		&s.CatStats.SearchVolume,
		&s.CatStats.DonationVolume,
		// dog stats
		&s.DogStats.TotalPets,
		&s.DogStats.ActiveBloodRequests,
		&s.DogStats.TotalDonations,
		&s.DogStats.CompletedDonations,
		&s.DogStats.Searches,
		&s.DogStats.SearchVolume,
		&s.DogStats.DonationVolume,
	)

	if err != nil {
		return nil, err
	}

	s.CatStats.Type = "cat"
	s.DogStats.Type = "dog"

	return &s, nil
}
