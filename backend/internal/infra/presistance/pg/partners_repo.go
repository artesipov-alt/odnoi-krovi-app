package pg

import (
	"context"
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	partnermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/partner/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/partner"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/domainmapper"
)

type EntPartnerRepository struct {
	db *ent.Client
}

func (r *EntPartnerRepository) client() *ent.Client {
	return r.db
}

func NewEntPartnerRepository(db *ent.Client) *EntPartnerRepository {
	return &EntPartnerRepository{db: db}
}

func (r *EntPartnerRepository) CreatePartner(ctx context.Context, name, apiKey, role, status, description string) (*partnermodel.Partner, error) {
	p, err := r.client().Partner.Create().
		SetName(name).
		SetAPIKey(apiKey).
		SetRole(partner.Role(role)).
		SetStatus(partner.Status(status)).
		SetNillableDescription(&description).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return EntPartnerToModel(p), nil
}

func (r *EntPartnerRepository) GetByAPIKey(ctx context.Context, apiKey string) (*partnermodel.Partner, error) {
	p, err := r.client().Partner.Query().
		Where(partner.APIKeyEQ(apiKey)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return EntPartnerToModel(p), nil
}

func (r *EntPartnerRepository) GetByID(ctx context.Context, id string) (*partnermodel.Partner, error) {
	p, err := r.client().Partner.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return EntPartnerToModel(p), nil
}

func (r *EntPartnerRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.client().Partner.UpdateOneID(id).
		SetStatus(partner.Status(status)).
		Exec(ctx)
}

func (r *EntPartnerRepository) UpdateLastUsedAt(ctx context.Context, id string, lastUsedAt time.Time) error {
	return r.client().Partner.UpdateOneID(id).
		SetLastUsedAt(lastUsedAt).
		Exec(ctx)
}

func (r *EntPartnerRepository) Delete(ctx context.Context, id string) error {
	return r.client().Partner.DeleteOneID(id).Exec(ctx)
}

func (r *EntPartnerRepository) ExistsByAPIKey(ctx context.Context, apiKey string) (bool, error) {
	count, err := r.client().Partner.Query().
		Where(partner.APIKeyEQ(apiKey)).
		Count(ctx)
	return count > 0, err
}

func (r *EntPartnerRepository) GetPartnerIdentities(ctx context.Context, partnerID string) ([]*authmodel.Identity, error) {
	identities, err := r.client().Partner.Query().
		Where(partner.ID(partnerID)).
		QueryPartnerIdentities().
		All(ctx)
	if err != nil {
		return nil, err
	}
	var result []*authmodel.Identity
	for _, i := range identities {
		result = append(result, domainmapper.EntIdentityToModel(i))
	}
	return result, nil
}

func EntPartnerToModel(p *ent.Partner) *partnermodel.Partner {
	model := &partnermodel.Partner{
		ID:          p.ID,
		Name:        p.Name,
		APIKey:      p.APIKey,
		Role:        string(p.Role),
		Status:      string(p.Status),
		Description: &p.Description,
		LastUsedAt:  &p.LastUsedAt,
	}
	return model
}
