package contact

import (
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminContactController struct {
	contactService usecase.ContactMessageService
}

func NewAdminContactController(contactService usecase.ContactMessageService) *AdminContactController {
	return &AdminContactController{contactService: contactService}
}

type contactMessageResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
	CreatedAt string `json:"createdAt"`
}

func (c *AdminContactController) GetAll(ctx *gin.Context) {
	msgs, err := c.contactService.GetAll()
	if err != nil {
		panic(err)
	}
	response := make([]contactMessageResponse, 0, len(msgs))
	for _, m := range msgs {
		response = append(response, contactMessageResponse{
			ID:        m.ID,
			Name:      m.Name,
			Email:     m.Email,
			Subject:   m.Subject,
			Message:   m.Message,
			CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	controller.Response(ctx, 200, "", response)
}
