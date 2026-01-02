package sql

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

// Int32OrNull returns properly configured sql.NullInt64
func Int32OrNull(n int32) sql.NullInt32 {
	if n != 0 {
		return sql.NullInt32{Int32: n, Valid: true}
	}
	return sql.NullInt32{Int32: 0, Valid: false}
}

// PositiveInt32OrNull returns properly configured sql.NullInt64
func PositiveInt32OrNull(n int32) sql.NullInt32 {
	if n > 0 {
		return sql.NullInt32{Int32: n, Valid: true}
	}
	return sql.NullInt32{Int32: 0, Valid: false}
}

// Int64OrNull returns properly configured sql.NullInt64
func Int64OrNull(n int64) sql.NullInt64 {
	if n != 0 {
		return sql.NullInt64{Int64: n, Valid: true}
	}
	return sql.NullInt64{Int64: 0, Valid: false}
}

// PositiveInt64OrNull returns properly configured sql.NullInt64
func PositiveInt64OrNull(n int64) sql.NullInt64 {
	if n > 0 {
		return sql.NullInt64{Int64: n, Valid: true}
	}
	return sql.NullInt64{Int64: 0, Valid: false}
}

// IntOrNull returns properly configured sql.NullInt64
// The same with Int64OrNull
func IntOrNull(n int64) sql.NullInt64 {
	if n != 0 {
		return sql.NullInt64{Int64: n, Valid: true}
	}
	return sql.NullInt64{Int64: 0, Valid: false}
}

// PositiveIntOrNull returns properly configured sql.NullInt64 for a positive number
// The same with PositiveInt64OrNull
func PositiveIntOrNull(n int64) sql.NullInt64 {
	if n > 0 {
		return sql.NullInt64{Int64: n, Valid: true}
	}
	return sql.NullInt64{Int64: 0, Valid: false}
}

// FloatOrNull returns properly configured sql.NullFloat64
func Float64OrNull(n float64) sql.NullFloat64 {
	if n != 0.0 {
		return sql.NullFloat64{Float64: n, Valid: true}
	}
	return sql.NullFloat64{Float64: 0, Valid: false}
}

// PositiveFloatOrNull returns properly configured sql.NullFloat64 for a positive number
func PositiveFloat64OrNull(n float64) sql.NullFloat64 {
	if n > 0.0 {
		return sql.NullFloat64{Float64: n, Valid: true}
	}
	return sql.NullFloat64{Float64: 0.0, Valid: false}
}

// StringOrNull returns properly configured sql.NullString
func StringOrNull(str string) sql.NullString {
	if str != "" {
		return sql.NullString{String: str, Valid: true}
	}
	return sql.NullString{String: "", Valid: false}
}

// DecimalOrNull returns properly configured sql.NullString
func DecimalOrNull(n decimal.Decimal) decimal.NullDecimal {
	if !n.IsZero() {
		return decimal.NullDecimal{
			Decimal: n,
			Valid:   true,
		}

	}
	return decimal.NullDecimal{
		Decimal: decimal.NewFromInt(0),
		Valid:   false}
}

// PositiveDecimalOrNull returns properly configured sql.NullString
func PositiveDecimalOrNull(n decimal.Decimal) decimal.NullDecimal {
	if n.GreaterThan(decimal.NewFromInt(0)) {
		return decimal.NullDecimal{
			Decimal: n,
			Valid:   true,
		}
	}
	return decimal.NullDecimal{
		Decimal: decimal.NewFromInt(0),
		Valid:   false}
}

// TimeOrNull returns properly configured pq.TimeNull
func TimeOrNull(t time.Time) sql.NullTime {
	if !t.IsZero() {
		return sql.NullTime{Time: t, Valid: true}
	}
	return sql.NullTime{Time: time.Time{}, Valid: false}
}

// UUIDOrNull returns properly configured uuid.NullUUID
func UUIDOrNull(t uuid.UUID) uuid.NullUUID {
	if t != uuid.Nil {
		return uuid.NullUUID{UUID: t, Valid: true}
	}
	return uuid.NullUUID{UUID: uuid.UUID{}, Valid: false}
}

// BoolOrNull returns properly configured uuid.NullUUID
func BoolOrNull(t bool) sql.NullBool {
	if t {
		return sql.NullBool{Bool: t, Valid: true}
	}
	return sql.NullBool{Bool: t, Valid: false}
}

func MapJSONOrNull(t []byte) pgtype.JSON {
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
