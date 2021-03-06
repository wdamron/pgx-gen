package example

import (
	"github.com/satori/go.uuid"
	"github.com/wdamron/pgx"
)

type Point struct {
	X  []string           `pgx:"name:x;type:varchar[]"`
	Y  *int64             `pgx:"name:y;type:int4"`
	Z  *pgx.NullInt32     `pgx:"name:z;type:integer"`
	H  *map[string]string `pgx:"name:h;type:hstore"`
	H2 pgx.Hstore         `pgx:"name:h2;type:hstore"`
	u  string             `pgx:"name:id;type:uuid"`
	u2 *uuid.UUID         `pgx:"name:id2;type:uuid"`
	j  *string            `pgx:"name:j;type:json"`
	j2 map[string]int     `pgx:"name:j2;type:json"`
	j3 []byte             `pgx:"name:j3;type:json"`
}
