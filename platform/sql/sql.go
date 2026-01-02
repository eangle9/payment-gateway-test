package sql

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

func Int32OrNullpntr(n *int32) sql.NullInt32 {
	if n != nil {
		return sql.NullInt32{Int32: *n, Valid: true}
	}
	return sql.NullInt32{Int32: 0, Valid: false}
}

// PositiveInt32OrNull returns properly configured sql.NullInt64
func PositiveInt32OrNullpntr(n *int32) sql.NullInt32 {
	if n != nil && *n > 0 {
		return sql.NullInt32{Int32: *n, Valid: true}
	}
	return sql.NullInt32{Int32: 0, Valid: false}
}

// Int64OrNullFromPntr returns properly configured sql.NullInt64
func Int64OrNullPntr(n *int64) sql.NullInt64 {
	if n != nil {
		return sql.NullInt64{Int64: *n, Valid: true}
	}
	return sql.NullInt64{Int64: 0, Valid: false}
}

// PositiveInt64OrNull returns properly configured sql.NullInt64
func PositiveInt64OrNullpntr(n *int64) sql.NullInt64 {
	if n != nil && *n > 0 {
		return sql.NullInt64{Int64: *n, Valid: true}
	}
	return sql.NullInt64{Int64: 0, Valid: false}
}

func Float64OrNullpntr(n *float64) sql.NullFloat64 {
	if n != nil {
		return sql.NullFloat64{Float64: *n, Valid: true}
	}
	return sql.NullFloat64{Float64: 0, Valid: false}
}

// PositiveFloatOrNull returns properly configured sql.NullFloat64 for a positive number
func PositiveFloat64OrNullpntr(n *float64) sql.NullFloat64 {
	if n != nil && *n > 0.0 {
		return sql.NullFloat64{Float64: *n, Valid: true}
	}
	return sql.NullFloat64{Float64: 0.0, Valid: false}
}

// StringOrNullFromPntr returns properly configured sql.NullString
func StringOrNullPntr(str *string) sql.NullString {
	if str != nil {
		return sql.NullString{String: *str, Valid: true}
	}
	return sql.NullString{String: "", Valid: false}
}

// DecimalOrNullFromPntr returns properly configured sql.NullString
func DecimalOrNullPntr(n *decimal.Decimal) decimal.NullDecimal {
	if n != nil {
		return decimal.NullDecimal{
			Decimal: *n,
			Valid:   true}
	}
	return decimal.NullDecimal{
		Decimal: decimal.Decimal{},
		Valid:   false,
	}
}

// decimal.NewFromInt(0)
// PositiveDecimalOrNull returns properly configured sql.NullString
func PositiveDecimalOrNullpntr(n *decimal.Decimal) decimal.NullDecimal {
	if n != nil && (*n).GreaterThan(decimal.NewFromInt(0)) {
		return decimal.NullDecimal{
			Decimal: *n,
			Valid:   true,
		}
	}
	return decimal.NullDecimal{
		Decimal: decimal.NewFromInt(0),
		Valid:   false}
}

// TimeOrNull returns properly configured pq.TimeNull
func TimeOrNullpntr(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Time: time.Time{}, Valid: false}
}

// UUIDOrNull returns properly configured uuid.NullUUID
func UUIDOrNullpntr(t *uuid.UUID) uuid.NullUUID {
	if t != nil {
		return uuid.NullUUID{UUID: *t, Valid: true}
	}
	return uuid.NullUUID{UUID: uuid.UUID{}, Valid: false}
}

// BoolOrNullFromPntr returns properly configured uuid.NullUUID
func BoolOrNullPntr(t *bool) sql.NullBool {
	if t != nil {
		return sql.NullBool{Bool: *t, Valid: true}
	}
	return sql.NullBool{Bool: false, Valid: false}
}
func MapJSONOrNullpntr(t []byte) pgtype.JSON {
	if t != nil {
		return pgtype.JSON{
			Bytes:  t,
			Status: pgtype.Present,
		}
	}
	return pgtype.JSON{
		Status: pgtype.Null,
	}
}
