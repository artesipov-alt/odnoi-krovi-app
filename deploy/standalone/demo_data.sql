-- Демонстрационные данные standalone-экземпляра ПО «Одной Крови».
--
-- Применяется одноразовым seed-контейнером после старта backend,
-- потому что схему БД создает Ent (авто-миграция при старте).
-- Идемпотентно: повторный запуск безопасен (ON CONFLICT DO NOTHING).
--
-- Демо-пользователь USR-DEMO00001 соответствует входу через
-- http://localhost (frontend автоматически выполняет signin
-- через /v1/auth/signin/service с providerId 11111111).

BEGIN;

-- Партнер с API-ключом, который frontend использует на localhost.
-- Роль clinic входит в enum ролей пользователей (SetRole(partner.Role)).
INSERT INTO partners (id, name, api_key, role, status, description, created_at, updated_at)
VALUES (
    'PRT-DEMO00001',
    'Демо-стенд',
    'usvc_eebbbdaac43ddb3883ccdbd33d6e912547ed8e20dbdc44f3',
    'clinic',
    'active',
    'Партнер демо-стенда: обеспечивает вход через localhost без Telegram/Max',
    now(),
    now()
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO users (
    id, full_name, phone, verified, verified_at, location_id,
    consent_pd, allow_geo, role, origin_source,
    on_boarding, photo_urls, priority_search_count, last_seen_at,
    created_at, updated_at
) VALUES
    (
        'USR-DEMO00001', 'Демо-пользователь', '+79000000001', true, now(), '45',
        true, true, 'user', 'service',
        '[]'::jsonb, '[]'::jsonb, 0, now(),
        now(), now()
    ),
    (
        'USR-DEMO00002', 'Анна Смирнова', '+79000000002', true, now(), '45',
        true, true, 'user', 'self',
        '[]'::jsonb, '[]'::jsonb, 0, now(),
        now(), now()
    ),
    (
        'USR-DEMO00003', 'Иван Петров', '+79000000003', true, now(), '45',
        true, true, 'user', 'self',
        '[]'::jsonb, '[]'::jsonb, 0, now(),
        now(), now()
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_identities (
    id, user_id, partner_id, provider, provider_user_id, metadata,
    created_at, updated_at
) VALUES
    (
        'IDN-DEMO00001', 'USR-DEMO00001', 'PRT-DEMO00001', 'service', '11111111', NULL,
        now(), now()
    ),
    (
        'IDN-DEMO00002', 'USR-DEMO00002', NULL, 'telegram_bot', '700000002', NULL,
        now(), now()
    ),
    (
        'IDN-DEMO00003', 'USR-DEMO00003', NULL, 'telegram_bot', '700000003', NULL,
        now(), now()
    )
ON CONFLICT (id) DO NOTHING;

-- Настройки донора: без них список заявок для донора возвращает
-- «Настройки донора не заполнены» (recipient-list).
INSERT INTO donor_preferences (
    id, user_id, preferred_location_ids, recovery_period_months,
    compensation_type, taxi_compensation, notification_frequency,
    open_for_contact, created_at, updated_at
) VALUES
    (
        'DPR-DEMO00001', 'USR-DEMO00001', '["45"]'::jsonb, 2,
        'free', false, 'immediately',
        true, now(), now()
    ),
    (
        'DPR-DEMO00003', 'USR-DEMO00003', '["45"]'::jsonb, 3,
        'free', false, 'immediately',
        true, now(), now()
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO pet_healths (
    id, health_status, last_donation, transfused, medications, surgical_interventions,
    created_at, updated_at
) VALUES
    ('PHL-DEMO00001', 'healthy', NULL, false, '', '', now(), now()),
    ('PHL-DEMO00002', 'healthy', NULL, false, '', '', now(), now()),
    ('PHL-DEMO00003', 'healthy', NULL, false, '', '', now(), now()),
    ('PHL-DEMO00004', 'healthy', NULL, false, '', '', now(), now()),
    ('PHL-DEMO00005', 'healthy', NULL, false, '', '', now(), now())
ON CONFLICT (id) DO NOTHING;

INSERT INTO pet_treatments (
    id, rabies_vaccination_date, infection_vaccination_date,
    ectoparasite_treatment_date, deworming_date,
    created_at, updated_at
) VALUES
    (
        'PTR-DEMO00001', now() - interval '4 months', now() - interval '4 months',
        now() - interval '1 month', now() - interval '1 month',
        now(), now()
    ),
    (
        'PTR-DEMO00002', now() - interval '3 months', now() - interval '3 months',
        now() - interval '20 days', now() - interval '20 days',
        now(), now()
    ),
    (
        'PTR-DEMO00003', now() - interval '6 months', now() - interval '6 months',
        now() - interval '2 months', now() - interval '2 months',
        now(), now()
    ),
    (
        'PTR-DEMO00004', now() - interval '5 months', now() - interval '5 months',
        now() - interval '1 month', now() - interval '1 month',
        now(), now()
    ),
    (
        'PTR-DEMO00005', now() - interval '2 months', now() - interval '2 months',
        now() - interval '20 days', now() - interval '20 days',
        now(), now()
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO pets (
    id, name, type, weight_kg, gender, birth_date, breed_id, user_id,
    health_id, treatment_id,
    blood_group, living_condition, reproductive_status,
    photo_urls, is_profile_lock, created_at, updated_at
) VALUES
    (
        'PET-DEMO00001', 'Рекс', 'dog', 28.5, 'male', now() - interval '5 years', 'LBR_RTR', 'USR-DEMO00001',
        'PHL-DEMO00001', 'PTR-DEMO00001',
        'DEA 1+', 'home', 'not_sterilized',
        '["pets/2026/PET-DEMO00001/photos/1.jpg"]'::jsonb, false, now(), now()
    ),
    (
        'PET-DEMO00002', 'Барсик', 'cat', 5.2, 'male', now() - interval '4 years', 'MYN_KN', 'USR-DEMO00001',
        'PHL-DEMO00002', 'PTR-DEMO00002',
        'A', 'home', 'sterilized',
        '["pets/2026/PET-DEMO00002/photos/1.jpg"]'::jsonb, false, now(), now()
    ),
    (
        'PET-DEMO00003', 'Бим', 'dog', 18.0, 'male', now() - interval '3 years', 'LBR_RTR', 'USR-DEMO00002',
        'PHL-DEMO00003', 'PTR-DEMO00003',
        'DEA 1-', 'home', 'sterilized',
        '["pets/2026/PET-DEMO00003/photos/1.jpg"]'::jsonb, false, now(), now()
    ),
    (
        'PET-DEMO00004', 'Луна', 'dog', 22.0, 'female', now() - interval '6 years', 'LBR_RTR', 'USR-DEMO00003',
        'PHL-DEMO00004', 'PTR-DEMO00004',
        'DEA 1+', 'home', 'sterilized',
        '["pets/2026/PET-DEMO00004/photos/1.jpg"]'::jsonb, false, now(), now()
    ),
    (
        'PET-DEMO00005', 'Мурка', 'cat', 4.0, 'female', now() - interval '2 years', 'SIA', 'USR-DEMO00002',
        'PHL-DEMO00005', 'PTR-DEMO00005',
        'B', 'home', 'sterilized',
        '["pets/2026/PET-DEMO00005/photos/1.jpg"]'::jsonb, false, now(), now()
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO blood_requests (
    id, pet_id, blood_volume_needed, regions, small_pets_notify_allowed,
    status, description, photo_urls, blood_group_names, blood_component_ids,
    on_boarding, priority_search, include_unknown_blood_group,
    created_at, updated_at
) VALUES
    (
        'BLS-DEMO00001', 'PET-DEMO00003', 450, '["45"]'::jsonb, true,
        'active',
        'Срочно нужна кровь для собаки после полостной операции. Клиника в Москве, забор возможен ежедневно.',
        '[]'::jsonb, '["DEA 1+"]'::jsonb, '["BLC-2"]'::jsonb,
        '[]'::jsonb, false, false,
        now(), now()
    ),
    (
        'BLS-DEMO00002', 'PET-DEMO00005', 0, '["45"]'::jsonb, true,
        'active',
        'Ищем донора-кота с группой крови A для переливания в московской клинике.',
        '[]'::jsonb, '["A"]'::jsonb, '["BLC-3"]'::jsonb,
        '[]'::jsonb, false, false,
        now(), now()
    )
ON CONFLICT (id) DO NOTHING;

COMMIT;
