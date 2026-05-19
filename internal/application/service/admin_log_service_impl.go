package service

import (
	adminlogdto "github.com/Mahoura-shop/Backend/internal/application/dto/admin_log"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type AdminLogService struct {
	repo domainPostgres.AdminActivityLogRepository
	db   database.Database
}

type AdminLogServiceDeps struct {
	Repo domainPostgres.AdminActivityLogRepository
	DB   database.Database
}

func NewAdminLogService(deps AdminLogServiceDeps) *AdminLogService {
	return &AdminLogService{repo: deps.Repo, db: deps.DB}
}

func (s *AdminLogService) CreateLog(adminID uint, method, path, ip string, statusCode int) error {
	log := entity.AdminActivityLog{
		AdminID:    adminID,
		Method:     method,
		Path:       path,
		IPAddress:  ip,
		StatusCode: statusCode,
	}
	return s.repo.Create(s.db, log)
}

func (s *AdminLogService) GetLogs() ([]adminlogdto.AdminActivityLogDTO, error) {
	logs, err := s.repo.GetAll(s.db)
	if err != nil {
		return nil, err
	}
	result := make([]adminlogdto.AdminActivityLogDTO, len(logs))
	for i, l := range logs {
		result[i] = adminlogdto.AdminActivityLogDTO{
			ID:         l.ID,
			AdminID:    l.AdminID,
			AdminPhone: l.Admin.Phone,
			AdminName:  l.Admin.FirstName + " " + l.Admin.LastName,
			Method:     l.Method,
			Path:       l.Path,
			IPAddress:  l.IPAddress,
			StatusCode: l.StatusCode,
			CreatedAt:  l.CreatedAt,
		}
	}
	return result, nil
}
