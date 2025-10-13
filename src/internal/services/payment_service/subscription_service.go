package payment_service

import (
	"context"
	"database/sql"
	"logispro/internal/constants"
	"logispro/internal/db"
	"logispro/internal/utils"
	"time"
)

const (
	PLAN_MONTH int64 = iota + 1
	PLAN_YEAR
)

type Status string

const SUBS_STATUS_ACTIVE Status = "active"
const SUBS_STATUS_CANCELLED Status = "cancelled"
const SUBS_STATUS_EXPIRED Status = "expired"
const SUBS_STATUS_TRIAL Status = "trial"

type SubscriptionService struct {
	Queries *db.Queries
}

type Subscription struct {
	UserId             int64
	PlanId             int64
	Status             Status
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	NextBillingDate    time.Time
	TrialStart         time.Time
	TrialEnd           time.Time
	Ammount            float64
}

func (s *SubscriptionService) GetUserCurrentSubscription(ctx context.Context, userId int64) (db.UserSubscription, error) {
	return s.Queries.GetCurrentUserSubscription(ctx, userId)
}

func (s *SubscriptionService) GetUserSubscriptions(ctx context.Context, userId int64) ([]db.UserSubscription, error) {
	return s.Queries.GetUserSubscriptions(ctx, userId)
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, sub Subscription) error {
	return s.Queries.CreateUsersubscription(ctx, db.CreateUsersubscriptionParams{
		UserID:             sub.UserId,
		PlanID:             sub.PlanId,
		Status:             sql.NullString{Valid: true, String: string(sub.Status)},
		CurrentPeriodStart: sub.CurrentPeriodStart,
		CurrentPeriodEnd:   sub.CurrentPeriodEnd,
		TrialStart:         sql.NullTime{Valid: !sub.TrialStart.IsZero(), Time: sub.TrialStart},
		TrialEnd:           sql.NullTime{Valid: !sub.TrialEnd.IsZero(), Time: sub.TrialEnd},
		NextBillingDate:    sql.NullTime{Valid: !sub.NextBillingDate.IsZero(), Time: sub.NextBillingDate},
		Amount:             float64(sub.Ammount),
	})
}

func (s *SubscriptionService) UpdateUserSubscriptionStatus(ctx context.Context, status Status, users []int64) error {
	return s.Queries.UpdateAgencyUsersSubscriptionStatus(ctx, db.UpdateAgencyUsersSubscriptionStatusParams{
		Status:  sql.NullString{Valid: true, String: string(status)},
		UsersID: users,
	})
}

func (s *SubscriptionService) GetSubscriptionStatus(ctx context.Context, userId int64) (Status, error) {
	currentSubscription, err := s.GetUserCurrentSubscription(ctx, userId)
	if err != nil {
		return SUBS_STATUS_EXPIRED, err
	}
	user, err := s.Queries.GetUser(ctx, userId)
	if err != nil {
		return SUBS_STATUS_EXPIRED, err
	}
	if user.Role == constants.ROLE_NORMAL {
		rootUser, err := s.Queries.GetUser(ctx, user.RootID.Int64)
		if err != nil {
			return SUBS_STATUS_EXPIRED, err
		}
		user = rootUser
	}

	if currentSubscription.TrialStart.Valid && currentSubscription.TrialEnd.Valid {
		// trial periode
		// TODO : when payment , must empty the trial periode
		if currentSubscription.TrialEnd.Time.Before(time.Now()) {
			return SUBS_STATUS_EXPIRED, nil
		}
	} else {
		if currentSubscription.CurrentPeriodEnd.Before(time.Now()) {
			if Status(currentSubscription.Status.String) != SUBS_STATUS_EXPIRED {
				rootId, err := utils.GetRootIdFromContext(ctx)
				if err != nil {
					return SUBS_STATUS_EXPIRED, err
				}
				agencyUsers, err := utils.GetAgencyUsers(ctx, s.Queries, userId, rootId)
				if err != nil {
					return SUBS_STATUS_EXPIRED, err
				}
				agencyUsersId := utils.ExtractField(agencyUsers, func(u db.GetAgencyUsersRow) int64 {
					return u.ID
				})
				s.UpdateUserSubscriptionStatus(ctx, SUBS_STATUS_EXPIRED, agencyUsersId)
			}
			return SUBS_STATUS_EXPIRED, nil
		}
	}

	return Status(currentSubscription.Status.String), nil

}
