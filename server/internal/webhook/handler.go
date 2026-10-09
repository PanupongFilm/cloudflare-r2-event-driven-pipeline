package webhook

import(
	
	"server/internal/middleware"
	"server/internal/r2"

	"github.com/gofiber/fiber/v3"
	_package "server/package"
)

type Handler struct{
	service       Service
	webhookSecret string
}


func NewHandler(service Service, webhookSecret string) *Handler{
	return &Handler{
		service:       service,
		webhookSecret: webhookSecret,
	}
}



func (h *Handler) UploadTrigger(c fiber.Ctx) error{
	
	receivedSecret := c.Get("Webhook-Secret")

	if receivedSecret == "" {
		return _package.SendJSON(
			c,
			fiber.StatusUnauthorized,
			_package.ErrorResponse(
				"Webhook secret is required",
				"UNAUTHORIZED",
				"Missing Webhook-Secret header",
			),
		)
	}

	if receivedSecret != h.webhookSecret {
		return _package.SendJSON(
			c,
			fiber.StatusUnauthorized,
			_package.ErrorResponse(
				"Invalid webhook secret",
				"UNAUTHORIZED",
				"Webhook secret does not match",
			),
		)
	}

	var payload r2.CreateUploadLogRequest
	if err := c.Bind().JSON(&payload); err != nil{
		return _package.SendJSON(
			c,
			fiber.StatusBadRequest,
			_package.ErrorResponse(
				"Invalid request body",
				"BODY_INVALID",
				err.Error(),
			),
		)
	}

	
	if payload.AccountID == "" || payload.ObjectKey == ""{
		return _package.SendJSON(
			c,
			fiber.StatusBadRequest,
			_package.ErrorResponse(
				"Missing required fields",
				"VALIDATION_ERROR",
				"account_id and object_key are required",
			),
		)
	}

	uploadLog, err := h.service.AddUploadLog(c.Context(), payload)
	if err != nil {
		return _package.SendJSON(
			c,
			fiber.StatusInternalServerError,
			_package.ErrorResponse(
				"Failed to create upload log",
				"SERVICE_ERROR",
				err.Error(),
			),
		)
	}


	return _package.SendJSON(
		c,
		fiber.StatusOK,
		_package.SuccessResponse(
			"Webhook received successfully",
			uploadLog,
		),
	)
}



func (h *Handler) RegisterRoutes(router fiber.Router) {

	router.Post("/webhook/upload-trigger", middleware.ReqEmptyCheck, h.UploadTrigger)
}