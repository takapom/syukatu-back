package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"` // bcrypt でハッシュ化したものを保存
}

// 企業名
// 職種
// 従業員人数
// 選考状況
// 　インターンの有無
type CompanyList struct {
	gorm.Model
	Company    string `gorm:"not null"`
	Occupation string
	Member     int `gorm:"index;not null"`
	Selection  string
	Intern     bool
	UserID     uint `gorm:"index;not null"`
}

// 後々にインターンモデルも作成予定(モデル名Internship)
type Internship struct {
	gorm.Model
	Title       string
	Company     string
	Dailystart  int
	Dailyfinish int
	Content     string
	Selection   string
	Joined      bool
	UserID      uint `gorm:"index;not null"`
}

// 掲示板投稿モデル
type Post struct {
	gorm.Model
	Title        string `json:"title" gorm:"not null"`
	Content      string `json:"content" gorm:"not null"`
	DisplayName  string `json:"display_name" gorm:"not null;size:20"`
	LikeCount    int    `json:"like_count" gorm:"default:0"`
	CommentCount int    `json:"comment_count" gorm:"default:0"`
	UserID       uint   `json:"user_id" gorm:"index;not null"`
}

// コメントモデル
type Comment struct {
	gorm.Model
	Content     string `json:"content" gorm:"not null"`
	DisplayName string `json:"display_name" gorm:"not null;size:20"`
	PostID      uint   `json:"post_id" gorm:"index;not null"`
	UserID      uint   `json:"user_id" gorm:"index;not null"`
}

// いいねモデル
type Like struct {
	gorm.Model
	PostID uint `json:"post_id" gorm:"index;not null"`
	UserID uint `json:"user_id" gorm:"index;not null"`
}
