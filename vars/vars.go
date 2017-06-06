// Package vars contains all variables used in other packages
package vars

const (
	Channel = "pagestat"
)

type Task struct {
	UUID string `json:"uuid"`
	Url  string `json:"url"`
	Cron string `json:"cron"` // 1h, 2m, 30s ...
}
