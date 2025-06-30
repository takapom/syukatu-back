package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// authMiddleware はリクエストヘッダーから JWT を検証し、userID をコンテキストにセットします
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := ParseJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// registerHandler は新規ユーザー登録を行うハンドラを返します
func registerHandler(db *gorm.DB) gin.HandlerFunc {
	type req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	return func(c *gin.Context) {
		var body req
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pwHash, err := HashPassword(body.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
			return
		}
		if _, err := createUser(db, body.Email, pwHash); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"email": body.Email})
	}
}

// loginHandler はログイン認証を行い、JWT を返すハンドラを返します
func loginHandler(db *gorm.DB) gin.HandlerFunc {
	type req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	return func(c *gin.Context) {
		var body req
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, err := getUserByEmail(db, body.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if err := CheckPassword(u.Password, body.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		token, err := GenerateJWT(u.ID)
		if err != nil {
			log.Printf("[login] GenerateJWT error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token, "user_id": u.ID})
	}
}

// createCompanyListHandler は新規 CompanyList 作成のハンドラ
func createCompanyListHandler(db *gorm.DB) gin.HandlerFunc {
	type req struct {
		Company    string `json:"company" binding:"required"`
		Occupation string `json:"occupation"`
		Member     int    `json:"member" binding:"required"`
		Selection  string `json:"selection"`
		Intern     bool   `json:"intern"`
	}
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		var body req
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cl, err := createCompanyList(
			db,
			userID,
			body.Company,
			body.Occupation,
			body.Member,
			body.Selection,
			body.Intern,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, cl)
	}
}

// listCompanyListsHandler はユーザーの CompanyList 一覧を返すハンドラ
func listCompanyListsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		list, err := listCompanyLists(db, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

// updateCompanyListHandler は既存 CompanyList 更新のハンドラ
func updateCompanyListHandler(db *gorm.DB) gin.HandlerFunc {
	type req struct {
		Company    string `json:"company" binding:"required"`
		Occupation string `json:"occupation"`
		Member     int    `json:"member" binding:"required"`
		Selection  string `json:"selection"`
		Intern     bool   `json:"intern"`
	}
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		var body req
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		if err := updateCompanyList(
			db,
			id,
			userID,
			body.Company,
			body.Occupation,
			body.Member,
			body.Selection,
			body.Intern,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// deleteCompanyListHandler は既存 CompanyList 削除のハンドラ
func deleteCompanyListHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		if err := deleteCompanyList(db, id, userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}


//インターンシップ作成handler処理
func createInternshipHandler(db *gorm.DB) gin.HandlerFunc {
    type req struct {
        Title       string `json:"title" binding:"required"`
        Company     string `json:"company" binding:"required"`
        Dailystart  int    `json:"dailystart" binding:"required"`
        Dailyfinish int    `json:"dailyfinish" binding:"required"`
        Content     string `json:"content"`
        Selection   string `json:"selection" binding:"required"`
        Joined      bool   `json:"joined"`
    }
    return func(c *gin.Context) {
        userID := c.GetUint("userID")
        var body req
        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        i, err := createInternship(
            db,
            userID,
            body.Title,
            body.Company,
            body.Dailystart,
            body.Dailyfinish,
            body.Content,
            body.Selection,
            body.Joined,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, i)
    }
}

//インターンシップ更新handler処理
func updateInternshipHandler(db *gorm.DB) gin.HandlerFunc {
    type req struct {
        Title       string `json:"title" binding:"required"`
        Company     string `json:"company" binding:"required"`
        Dailystart  int    `json:"dailystart" binding:"required"`
        Dailyfinish int    `json:"dailyfinish" binding:"required"`
        Content     string `json:"content"`
        Selection   string `json:"selection" binding:"required"`
        Joined      bool   `json:"joined"`
	}
    
    return func(c *gin.Context) {
        userID := c.GetUint("userID")
        var body req
        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        var id uint
        fmt.Sscanf(c.Param("id"), "%d", &id)
        if err := updateInternship(
            db,
            id,
            userID,
            body.Title,
            body.Company,
            body.Dailystart,
            body.Dailyfinish,
            body.Content,
            body.Selection,
            body.Joined,
        ); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.Status(http.StatusNoContent)
    }
}


//インターンシップ一覧handler処理
func listInternshipsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")  // コンテキストからuserIDを取得
		list, err := listInternships(db, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

//インターンシップ削除handler処理
func deleteInternshipHandler(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		userID := c.GetUint("userID")
		var id uint
		fmt.Sscanf(c.Param("id"), "%d", &id)
		//repositoryのdeleteInternshipを呼び出し処理
		if err := deleteInternship(db, id, userID); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// 掲示板関連のハンドラー

// 投稿作成ハンドラー
func createPostHandler(db *gorm.DB) gin.HandlerFunc {
	type req struct {
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content" binding:"required"`
		DisplayName string `json:"display_name" binding:"required,max=20"`
	}
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		var body req
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		post := &Post{
			Title:       body.Title,
			Content:     body.Content,
			DisplayName: body.DisplayName,
			UserID:      userID,
		}
		
		if err := createPost(db, post); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// デバッグ用ログ
		log.Printf("Created post: %+v", post)
		
		c.JSON(http.StatusCreated, post)
	}
}

// 投稿一覧取得ハンドラー
func getPostsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ページネーションパラメータ
		var limit = 20
		var offset = 0
		
		if l := c.Query("limit"); l != "" {
			fmt.Sscanf(l, "%d", &limit)
		}
		if o := c.Query("offset"); o != "" {
			fmt.Sscanf(o, "%d", &offset)
		}
		
		posts, err := getPosts(db, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// 各投稿に対していいね状態を追加
		userID := c.GetUint("userID")
		type PostResponse struct {
			ID           uint   `json:"ID"`
			Title        string `json:"title"`
			Content      string `json:"content"`
			DisplayName  string `json:"display_name"`
			LikeCount    int    `json:"like_count"`
			CommentCount int    `json:"comment_count"`
			UserID       uint   `json:"user_id"`
			CreatedAt    time.Time `json:"CreatedAt"`
			IsLiked      bool   `json:"is_liked"`
		}
		
		var response []PostResponse
		for _, post := range posts {
			isLiked, _ := checkUserLiked(db, post.ID, userID)
			response = append(response, PostResponse{
				ID:           post.ID,
				Title:        post.Title,
				Content:      post.Content,
				DisplayName:  post.DisplayName,
				LikeCount:    post.LikeCount,
				CommentCount: post.CommentCount,
				UserID:       post.UserID,
				CreatedAt:    post.CreatedAt,
				IsLiked:      isLiked,
			})
		}
		
		c.JSON(http.StatusOK, response)
	}
}

// 投稿詳細取得ハンドラー
func getPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var postID uint
		fmt.Sscanf(c.Param("id"), "%d", &postID)
		
		post, err := getPost(db, postID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		
		// いいね状態を確認
		userID := c.GetUint("userID")
		isLiked, _ := checkUserLiked(db, postID, userID)
		
		// コメント一覧を取得
		comments, _ := getComments(db, postID)
		
		c.JSON(http.StatusOK, gin.H{
			"post":     post,
			"is_liked": isLiked,
			"comments": comments,
		})
	}
}

// 投稿削除ハンドラー
func deletePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		var postID uint
		fmt.Sscanf(c.Param("id"), "%d", &postID)
		
		if err := deletePost(db, postID, userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// いいね追加ハンドラー
func likePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		var postID uint
		fmt.Sscanf(c.Param("id"), "%d", &postID)
		
		like := &Like{
			PostID: postID,
			UserID: userID,
		}
		
		if err := createLike(db, like); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	}
}

// いいね削除ハンドラー
func unlikePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		var postID uint
		fmt.Sscanf(c.Param("id"), "%d", &postID)
		
		if err := deleteLike(db, postID, userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// コメント作成ハンドラー
func createCommentHandler(db *gorm.DB) gin.HandlerFunc {
	type req struct {
		Content     string `json:"content" binding:"required"`
		DisplayName string `json:"display_name" binding:"required,max=20"`
	}
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		var postID uint
		fmt.Sscanf(c.Param("id"), "%d", &postID)
		
		var body req
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		comment := &Comment{
			Content:     body.Content,
			DisplayName: body.DisplayName,
			PostID:      postID,
			UserID:      userID,
		}
		
		if err := createComment(db, comment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, comment)
	}
}
