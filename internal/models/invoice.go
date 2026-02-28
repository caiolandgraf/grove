package models

import (
	"time"

	"github.com/caiolandgraf/grove/internal/app"
	"github.com/caiolandgraf/grove/internal/database"
	"gorm.io/gorm"
)

type Invoice struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Invoice) TableName() string { return "invoices" }

// Invoices returns a repository scoped to the Invoice model.
func Invoices() *database.Repository[Invoice] {
	return database.New[Invoice](app.DB)
}
