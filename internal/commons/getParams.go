package commons

import "github.com/jackc/pgx/v5/pgtype"

var DefaultSkip = pgtype.Int4{Int32: 0, Valid: true}
var DefaultTake = pgtype.Int4{Int32: 50, Valid: true}
