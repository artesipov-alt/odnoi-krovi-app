package job

const queryDonorNotAccepted = `
		SELECT
            dr.id,
            recipient_pet.name,
            recipient_pet.blood_group,
            br.blood_volume_needed,
            MAX(CASE WHEN i.provider = 'telegram_bot' THEN i.provider_user_id END) AS telegram_id,
            MAX(CASE WHEN i.provider = 'max_bot'      THEN i.provider_user_id END) AS max_id
        FROM donor_responses dr
        JOIN blood_requests br       ON br.id = dr.request_id
        JOIN pets donor_pet          ON donor_pet.id = dr.donor_id
        JOIN users donor_user        ON donor_user.id = donor_pet.user_id
        LEFT JOIN user_identities i       ON i.user_id = donor_user.id
        JOIN pets recipient_pet      ON recipient_pet.id = br.pet_id
        WHERE dr.status = 'pending'
          AND dr.created_at < NOW() - INTERVAL '1 hour'
          AND br.status = 'active'
        GROUP BY dr.id, recipient_pet.name, recipient_pet.blood_group, br.blood_volume_needed
`

const queryDonorWaiting = `
		SELECT
		    dr.id,
		    donor_pet.name,
		    donor_pet.blood_group,
		    MAX(CASE WHEN i.provider = 'telegram_bot' THEN i.provider_user_id END) AS telegram_id,
		    MAX(CASE WHEN i.provider = 'max_bot'      THEN i.provider_user_id END) AS max_id
		FROM donor_responses dr
		JOIN blood_requests br      ON br.id = dr.request_id
		JOIN pets donor_pet         ON donor_pet.id = dr.donor_id
		JOIN pets recipient_pet     ON recipient_pet.id = br.pet_id
		JOIN users recipient_user   ON recipient_user.id = recipient_pet.user_id
		LEFT JOIN user_identities i ON i.user_id = recipient_user.id
		WHERE dr.status = 'pending'
		  AND dr.created_at < NOW() - INTERVAL '30 minutes'
		  AND br.status = 'active'
		GROUP BY dr.id, donor_pet.name, donor_pet.blood_group
`

const queryRecipientInactive6h = `
		SELECT
		    br.id,
		    MAX(CASE WHEN i.provider = 'telegram_bot' THEN i.provider_user_id END) AS telegram_id,
		    MAX(CASE WHEN i.provider = 'max_bot'      THEN i.provider_user_id END) AS max_id
		FROM blood_requests br
		JOIN pets recipient_pet     ON recipient_pet.id = br.pet_id
		JOIN users recipient_user   ON recipient_user.id = recipient_pet.user_id
		LEFT JOIN user_identities i ON i.user_id = recipient_user.id
		WHERE br.status = 'active'
					  AND COALESCE(recipient_user.last_seen_at, recipient_user.created_at) < NOW() - INTERVAL '6 hours'
					  AND EXISTS (
					      SELECT 1 FROM donor_responses dr
					      WHERE dr.request_id = br.id
					        AND dr.status = 'pending'
					  )
					GROUP BY br.id
`
const queryRecipientInactive12h = `
		SELECT
		    br.id,
		    MAX(CASE WHEN i.provider = 'telegram_bot' THEN i.provider_user_id END) AS telegram_id,
		    MAX(CASE WHEN i.provider = 'max_bot'      THEN i.provider_user_id END) AS max_id
		FROM blood_requests br
		JOIN pets recipient_pet     ON recipient_pet.id = br.pet_id
		JOIN users recipient_user   ON recipient_user.id = recipient_pet.user_id
		LEFT JOIN user_identities i ON i.user_id = recipient_user.id
		WHERE br.status = 'active'
					  AND COALESCE(recipient_user.last_seen_at, recipient_user.created_at) < NOW() - INTERVAL '12 hours'
					  AND EXISTS (
					      SELECT 1 FROM donor_responses dr
					      WHERE dr.request_id = br.id
					        AND dr.status = 'pending'
					  )
					GROUP BY br.id
`
const queryRecipientEmptyShowcase24h = `
		SELECT
		    br.id,
		    recipient_pet.name,
		    recipient_pet.blood_group,
		    br.blood_volume_needed,
		    MAX(CASE WHEN i.provider = 'telegram_bot' THEN i.provider_user_id END) AS telegram_id,
		    MAX(CASE WHEN i.provider = 'max_bot'      THEN i.provider_user_id END) AS max_id
		FROM blood_requests br
		JOIN pets recipient_pet     ON recipient_pet.id = br.pet_id
		JOIN users recipient_user   ON recipient_user.id = recipient_pet.user_id
		LEFT JOIN user_identities i ON i.user_id = recipient_user.id
		WHERE br.status = 'active'
					  AND COALESCE(recipient_user.last_seen_at, recipient_user.created_at) < NOW() - INTERVAL '24 hours'
					  AND NOT EXISTS (
					      SELECT 1 FROM donor_responses dr
					      WHERE dr.request_id = br.id
					        AND dr.status = 'pending'
					  )
					GROUP BY br.id, recipient_pet.name, recipient_pet.blood_group, br.blood_volume_needed
`
const queryRecipientEmptyShowcase48h = `
		SELECT
		    br.id,
		    MAX(CASE WHEN i.provider = 'telegram_bot' THEN i.provider_user_id END) AS telegram_id,
		    MAX(CASE WHEN i.provider = 'max_bot'      THEN i.provider_user_id END) AS max_id
		FROM blood_requests br
		JOIN pets recipient_pet     ON recipient_pet.id = br.pet_id
		JOIN users recipient_user   ON recipient_user.id = recipient_pet.user_id
		LEFT JOIN user_identities i ON i.user_id = recipient_user.id
		WHERE br.status = 'active'
					  AND COALESCE(recipient_user.last_seen_at, recipient_user.created_at) < NOW() - INTERVAL '48 hours'
					  AND NOT EXISTS (
				      SELECT 1 FROM donor_responses dr
				      WHERE dr.request_id = br.id
				        AND dr.status = 'pending'
				  )
				GROUP BY br.id
`

// queryAccepted12h — заявки в active/reserved_full с accepted-откликом старше 12ч, донация не подтверждена.
// Возвращает данные и реципиента, и донора в одной строке.
const queryAccepted12h = `
		SELECT
		    dr.id AS response_id,
		    br.id AS request_id,
		    recipient_pet.name AS recipient_pet_name,
		    recipient_pet.blood_group AS recipient_blood_group,
		    donor_pet.name AS donor_pet_name,
		    donor_pet.blood_group AS donor_blood_group,
		    dr.amount,
		    MAX(CASE WHEN ri.provider = 'telegram_bot' THEN ri.provider_user_id END) AS recipient_telegram_id,
		    MAX(CASE WHEN ri.provider = 'max_bot'      THEN ri.provider_user_id END) AS recipient_max_id,
		    MAX(CASE WHEN di.provider = 'telegram_bot' THEN di.provider_user_id END) AS donor_telegram_id,
		    MAX(CASE WHEN di.provider = 'max_bot'      THEN di.provider_user_id END) AS donor_max_id
		FROM donor_responses dr
		JOIN blood_requests br      ON br.id = dr.request_id
		JOIN pets recipient_pet     ON recipient_pet.id = br.pet_id
		JOIN users recipient_user   ON recipient_user.id = recipient_pet.user_id
		LEFT JOIN user_identities ri ON ri.user_id = recipient_user.id
		JOIN pets donor_pet          ON donor_pet.id = dr.donor_id
		JOIN users donor_user        ON donor_user.id = donor_pet.user_id
		LEFT JOIN user_identities di ON di.user_id = donor_user.id
		WHERE dr.status = 'accepted'
		  AND dr.is_confirmed = false
		  AND dr.updated_at < NOW() - INTERVAL '12 hours'
		  AND br.status IN ('active', 'reserved_full')
		GROUP BY dr.id, br.id, recipient_pet.name, recipient_pet.blood_group, donor_pet.name, donor_pet.blood_group, dr.amount
`

// queryAccepted24h — заявки в active/reserved_full с accepted-откликом старше 24ч, донация не подтверждена.
// Возвращает данные только реципиента (уведомление для донора на 24ч не предусмотрено).
const queryAccepted24h = `
		SELECT
		    dr.id AS response_id,
		    br.id AS request_id,
		    recipient_pet.name AS recipient_pet_name,
		    recipient_pet.blood_group AS recipient_blood_group,
		    MAX(CASE WHEN ri.provider = 'telegram_bot' THEN ri.provider_user_id END) AS recipient_telegram_id,
		    MAX(CASE WHEN ri.provider = 'max_bot'      THEN ri.provider_user_id END) AS recipient_max_id
		FROM donor_responses dr
		JOIN blood_requests br      ON br.id = dr.request_id
		JOIN pets recipient_pet     ON recipient_pet.id = br.pet_id
		JOIN users recipient_user   ON recipient_user.id = recipient_pet.user_id
		LEFT JOIN user_identities ri ON ri.user_id = recipient_user.id
		WHERE dr.status = 'accepted'
		  AND dr.is_confirmed = false
		  AND dr.updated_at < NOW() - INTERVAL '24 hours'
		  AND br.status IN ('active', 'reserved_full')
		GROUP BY dr.id, br.id, recipient_pet.name, recipient_pet.blood_group
`
