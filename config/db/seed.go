package db

import (
	"chatvis-chat/internal/models"
	"chatvis-chat/internal/pkg"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// SeedDefaultAdmin verifica si ya existe un usuario admin por defecto,
// si no existe lo crea con la contraseña hasheada.
func SeedDefaultAdmin() {
	const defaultEmail = "admin@chatvist.com"

	var existing models.Usuarios
	err := DB.Where("email = ?", defaultEmail).First(&existing).Error

	if err == nil {
		log.Println("Seed: El usuario admin por defecto ya existe, omitiendo creación.")
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Seed: Error al buscar usuario admin: %v", err)
		return
	}

	hashedPassword, err := pkg.HashPassword("Admin2026!")
	if err != nil {
		log.Printf("Seed: Error al hashear contraseña del admin: %v", err)
		return
	}

	admin := models.Usuarios{
		Nombre:   "Administrador",
		Apodo:    "admin",
		Email:    defaultEmail,
		Password: hashedPassword,
		Fecha:    time.Now(),
		IsLlm:    false,
		IsAdmin:  true,
		IsActive: true,
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Printf("Seed: Error al crear usuario admin por defecto: %v", err)
		return
	}

	log.Printf("Seed: Usuario admin creado correctamente (ID: %d, Email: %s)", admin.Id, admin.Email)
}
