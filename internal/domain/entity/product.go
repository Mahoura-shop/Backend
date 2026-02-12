package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Product struct {
	database.Model
	Name          string    `gorm:"type:varchar(50);uniqueIndex"`
	Slug          string    `gorm:"type:varchar(50);uniqueIndex"`
	Price         float64   `gorm:"type:decimal(10,2);not null"`
	CurrencyCode  string    `gorm:"type:varchar(5);default:'IRR';not null"`
	IRRPrice      float64   `gorm:"type:decimal(10,2);not null"`
	ConsumerPrice float64   `gorm:"type:decimal(10,2)"`
	Step1Percent  float64   `gorm:"type:decimal(10,2);check:step1_percent >= -100 AND step1_percent <= 100"`
	Step2Percent  float64   `gorm:"type:decimal(10,2);check:step2_percent >= -100 AND step2_percent <= 100"`
	Step3Percent  float64   `gorm:"type:decimal(10,2);check:step3_percent >= -100 AND step3_percent <= 100"`
	Step1Price    float64   `gorm:"type:decimal(10,2)"`
	Step2Price    float64   `gorm:"type:decimal(10,2)"`
	Step3Price    float64   `gorm:"type:decimal(10,2)"`
	Quantity      uint      `gorm:"default:0;not null"`
	QuantityType  string    `gorm:"default:'pieces';not null"`
	Priority      uint      `gorm:"default:0;index" validate:"min=0"`
	MinOrder      uint      `gorm:"default:1" validate:"min=1"`
	CategoryID    *uint
	Category      *Category `gorm:"foreignKey:CategoryID"`
	BrandID       *uint
	Brand         *Brand    `gorm:"foreignKey:BrandID"`
	Description   string    `gorm:"type:text"`
	IsActive      bool      `gorm:"default:true"`
	IsNew         bool      `gorm:"default:true"`
	ProductPic    string    `gorm:"type:varchar(255);default:null"`
}