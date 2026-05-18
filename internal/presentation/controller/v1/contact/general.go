package contact

import (
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralContactController struct {
	contactService usecase.ContactMessageService
}

func NewGeneralContactController(contactService usecase.ContactMessageService) *GeneralContactController {
	return &GeneralContactController{contactService: contactService}
}

func (c *GeneralContactController) Submit(ctx *gin.Context) {
	type params struct {
		Name    string `json:"name" validate:"required"`
		Email   string `json:"email" validate:"required,email"`
		Subject string `json:"subject" validate:"required"`
		Message string `json:"message" validate:"required,min=10"`
	}
	p := controller.Validated[params](ctx)

	if err := c.contactService.Submit(p.Name, p.Email, p.Subject, p.Message); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "پیام شما با موفقیت ارسال شد", nil)
}
