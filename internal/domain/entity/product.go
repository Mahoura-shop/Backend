package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Product struct {
	database.Model
	Name          string    `gorm:"type:varchar(50);uniqueIndex"`
	Slug          string    `gorm:"type:varchar(50);uniqueIndex"`
	Price         float64   `gorm:"type:decimal(10,2);not null"`
	CurrencyID    uint
	Currency      Currency  `gorm:"foreignKey:CurrencyID"`
	IRRPrice      uint      `gorm:"type:int;not null"`
	ConsumerPrice uint      `gorm:"type:int"`
	Step1Percent  float64   `gorm:"type:decimal(10,2);check:step1_percent >= -100"`
	Step2Percent  float64   `gorm:"type:decimal(10,2);check:step2_percent >= -100"`
	Step3Percent  float64   `gorm:"type:decimal(10,2);check:step3_percent >= -100"`
	Step4Percent  float64   `gorm:"type:decimal(10,2);check:step4_percent >= -100"`
	Step1Price    uint      `gorm:"type:int"`
	Step2Price    uint      `gorm:"type:int"`
	Step3Price    uint      `gorm:"type:int"`
	Step4Price    uint      `gorm:"type:int"`
	Step1Origin   bool      `gorm:"default:false"`
	Step2Origin   bool      `gorm:"default:false"`
	Step3Origin   bool      `gorm:"default:false"`
	Step4Origin   bool      `gorm:"default:false"`
	Quantity      uint      `gorm:"default:0;not null"`
	QuantityType  string    `gorm:"default:'pieces';not null"`
	Priority      uint      `gorm:"default:0;index" validate:"min=0"`
	MinOrder      uint      `gorm:"default:1" validate:"min=1"`
	CategoryID    *uint
	Category      *Category `gorm:"foreignKey:CategoryID"`
	BrandID       *uint
	Brand         *Brand    `gorm:"foreignKey:BrandID"`
	Description   string    `gorm:"type:text"`
	Offer         string    `gorm:"type:text"`
	IsActive      bool      `gorm:"default:true"`
	IsNew         bool      `gorm:"default:true"`
	ProductPic    string    `gorm:"type:varchar(255);default:null"`
	Size          string    `gorm:"type:varchar(255);default:null"`
}