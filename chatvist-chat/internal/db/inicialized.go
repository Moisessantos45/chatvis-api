package db

import "chatvis-chat/internal/app"

func Init() error {

	if err := DB.AutoMigrate(app.Models...); err != nil {
		return err
	}

	return nil
}
