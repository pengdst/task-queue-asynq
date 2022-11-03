package web

type EmailDeliveryPayload struct {
	UserId     int    `json:"user_id"`
	TemplateId string `json:"template_id"`
}
