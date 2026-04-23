package query

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type GetAssignedBonusesResult struct {
	TotalPriority int
	Food          []*bonusmodel.Bonus
	Preparation   []*bonusmodel.Bonus
	Other         []*bonusmodel.Bonus
}

type AssignedBonusesHandler struct {
	userRepo  user.Repository
	bonusRepo bonus.Repository
}

func NewAssignedBonusesHandler(userRepo user.Repository, bonusRepo bonus.Repository) *AssignedBonusesHandler {
	return &AssignedBonusesHandler{
		userRepo:  userRepo,
		bonusRepo: bonusRepo,
	}
}

func (h *AssignedBonusesHandler) Handle(ctx context.Context, userID string) (*GetAssignedBonusesResult, error) {
	user, err := h.userRepo.GetByID(ctx, userID, user.UserPreloadOptions{})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	bonuses, err := h.bonusRepo.GetAssignedBonuses(ctx, userID)
	if err != nil {
		return nil, err
	}

	var food []*bonusmodel.Bonus
	var preparation []*bonusmodel.Bonus
	var other []*bonusmodel.Bonus

	for _, b := range bonuses {
		switch b.Category {
		case bonusmodel.CategoryFood:
			food = append(food, b)
		case bonusmodel.CategoryPreparation:
			preparation = append(preparation, b)
		default:
			other = append(other, b)
		}
	}

	result := &GetAssignedBonusesResult{
		TotalPriority: user.PrioritySearchCount,
		Food:          food,
		Preparation:   preparation,
		Other:         other,
	}

	return result, nil
}
