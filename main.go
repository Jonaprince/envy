package envy

import "database/sql"

type Envy struct {
	*sql.DB
}
