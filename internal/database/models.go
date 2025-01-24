// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
)

type User struct {
	ID                   interface{}
	Email                string
	HashedPassword       string
	RiskPreference       string
	PlanType             string
	StripeCustomerID     sql.NullString
	StripeSubscriptionID sql.NullString
	ScholarshipFlag      sql.NullInt64
	CreatedAt            sql.NullTime
	UpdatedAt            sql.NullTime
}
