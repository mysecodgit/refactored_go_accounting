package store

import (
	"database/sql"
)

func NullInt64(v *int64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{
		Int64: *v,
		Valid: true,
	}
}


func NullFloat64(v *float64) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{
		Float64: *v,
		Valid: true,
	}
}