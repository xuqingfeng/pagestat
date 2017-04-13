// Package vars contains all variables used in other packages
package vars

const (
	Topic = "pagestat"
)

type Task struct {
	Url  string `json:"url"`
	Cron string `json:"cron"` // 1h, 2m, 30s ...
}
