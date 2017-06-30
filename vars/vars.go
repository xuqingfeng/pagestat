// Package vars contains all variables used in other packages
package vars

const (
	Name = "pagestat"

	Channel = "pagestat"

	APIBasePath = "/" + Name + "/v1"
)

type Task struct {
	UUID string `json:"uuid"`
	Url  string `json:"url"`
	Cron string `json:"cron"` // 1h, 2m, 30s ...
}

type Msg struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
