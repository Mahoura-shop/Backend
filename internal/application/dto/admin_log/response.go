package adminlogdto

import "time"

type AdminActivityLogDTO struct {
	ID         uint      `json:"id"`
	AdminID    uint      `json:"adminID"`
	AdminPhone string    `json:"adminPhone"`
	AdminName  string    `json:"adminName"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	IPAddress  string    `json:"ipAddress"`
	StatusCode int       `json:"statusCode"`
	CreatedAt  time.Time `json:"createdAt"`
}
