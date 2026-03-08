package user

type AppValidator interface {
	ValidateHash(initData string) (providerID int64)
	ValidateBySecret(id int64, secret string) (providerID int64)
}

type TokenGenerator interface {
	Generate(entityID, role string) (token string)
	Validate(token string) (entityID string, role string, err error)
}
