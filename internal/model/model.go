package model

type User struct {
	Email               string `json:"email" gorm:"size:255;unique;not null"`
	Username            string `json:"username" gorm:"size:255;unique;not null"`
	PassHash            string `json:"-" gorm:"size:255;not null"`
	VerificationCode    string `json:"verification_code" gorm:"size:255"`
	ID                  uint   `json:"id" gorm:"primary_key;unique;not null"`
	RefreshTokenVersion uint   `json:"-"`
	AmountOfBookmarks   uint   `json:"amount_of_bookmarks"`
	IsVerified          bool   `json:"-" gorm:"default:false"`
	IsPremium           bool   `json:"is_premium" gorm:"default:false"`
}

type Bookmark struct {
	Title    string `json:"title" gorm:"size:128"`
	URL      string `json:"url" gorm:"size:255"`
	IconURL  string `json:"icon_url" gorm:"size:255"`
	ID       uint   `json:"id" gorm:"primaryKey;not null;unique"`
	UserID   uint   `json:"user_id"`
	ShowText bool   `json:"show_text" gorm:"default:false"`
}
