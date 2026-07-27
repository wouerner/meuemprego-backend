package database

import (
	"log"

	"github.com/wouerner/runter-backend/internal/config"
	"github.com/wouerner/runter-backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Conexão com PostgreSQL estabelecida com sucesso!")

	// Executa AutoMigrate para criar/atualizar as tabelas
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, err
	}

	log.Println("Migrações do banco de dados executadas com sucesso!")

	return db, nil
}
