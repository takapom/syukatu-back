// package main

// import (
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func openGormDB() (*gorm.DB, error) {
// 	db, err := gorm.Open(sqlite.Open("example.db"), &gorm.Config{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	// マイグレーション：User, Todo テーブルを自動作成／更新
// 	if err := db.AutoMigrate(&User{}, &CompanyList{}); err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }

package main

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openGormDB(config *Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if strings.HasPrefix(config.DatabaseURL, "postgres://") || strings.HasPrefix(config.DatabaseURL, "postgresql://") {
		// PostgreSQL for production
		db, err = gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	} else if strings.HasPrefix(config.DatabaseURL, "sqlite://") {
		// SQLite for local development
		dbPath := strings.TrimPrefix(config.DatabaseURL, "sqlite://")
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	} else {
		return nil, fmt.Errorf("unsupported database URL: %s", config.DatabaseURL)
	}

	if err != nil {
		return nil, err
	}

	// マイグレーション：User, CompanyList, Internship, Post, Comment, Like テーブルを自動作成／更新
	if err := db.AutoMigrate(&User{}, &CompanyList{}, &Internship{}, &Post{}, &Comment{}, &Like{}); err != nil {
		return nil, err
	}

	// SQLiteの場合のみ外部キー制約を有効化
	if strings.HasPrefix(config.DatabaseURL, "sqlite://") {
		if err := db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
			log.Printf("Failed to enable foreign keys: %v", err)
		}
	}

	return db, nil
}