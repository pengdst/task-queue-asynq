package web

type EmailDeliveryPayload struct {
	UserId     int    `json:"user_id" validate:"required"`
	TemplateId string `json:"template_id" validate:"required"`
}
