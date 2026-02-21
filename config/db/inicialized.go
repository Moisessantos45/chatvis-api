package db

import "chatvis-chat/internal/models"

func Init() error {

	if err := DB.AutoMigrate(models.Models...); err != nil {
		return err
	}

	return nil
}
