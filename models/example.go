// Package models holds the service's domain types. Each type embeds
// orm.Model and uses `orm.Typed[T]` for queries — that is the ONE way to
// read/write per-(org,user) SQLite data.
//
// The skeleton ships one Example type so the layout is visible. Replace
// with your domain.
package models

// Example is the skeleton's placeholder domain type. Swap the body out
// for your real fields; keep the embedding pattern intact.
type Example struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
}
