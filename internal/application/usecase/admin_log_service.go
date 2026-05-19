package usecase

import adminlogdto "github.com/Mahoura-shop/Backend/internal/application/dto/admin_log"

type AdminLogService interface {
	CreateLog(adminID uint, method, path, ip string, statusCode int) error
	GetLogs() ([]adminlogdto.AdminActivityLogDTO, error)
}
