package dtx

import (
	"context"
	"database/sql"
	"net/http"

)

func main() {
	ctx := context.Background()
	db, _ := sql.Open("postgres", "localhost:5432")
}

type Controller struct {
	db SqlDBWrapper
}

func (c *Controller) apiEntryPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, tx, err := c.db.NewTxDistributed(ctx, nil)

}