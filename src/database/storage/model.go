package storage

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User model
type EeUser struct {
	UserID       uint      `gorm:"primary_key" json:"user_id"`
	Username     string    `gorm:"type:varchar(255);unique;not null" json:"username"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:char(64);not null" json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Course model
type EeCourse struct {
	CourseID     uint           `gorm:"primary_key" json:"course_id"`
	Title        string         `gorm:"type:varchar(255);not null" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	Meta         datatypes.JSON `gorm:"type:text" json:"meta"`
	InstructorID uint           `gorm:"type:integer" json:"instructor_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// Course owner model
type EeCourseOwner struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UserID    uint      `gorm:"type:integer" json:"user_id"`
	CourseID  uint      `gorm:"type:integer" json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CourseSection model
type EeCourseSection struct {
	SectionID uint      `gorm:"primary_key" json:"section_id"`
	CourseID  uint      `gorm:"type:integer" json:"course_id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Order     int       `gorm:"type:integer;not null" json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Lesson model
type EeLesson struct {
	LessonID    uint      `gorm:"primary_key" json:"lesson_id"`
	SectionID   uint      `gorm:"type:integer" json:"section_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	ContentText string    `gorm:"type:text" json:"content_text"`
	VideoID     uint      `gorm:"type:integer" json:"video_id"` // Optional, relation will be established separately
	Order       int       `gorm:"type:integer;not null" json:"order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Video model
type EeVideo struct {
	VideoID   uint      `gorm:"primary_key" json:"video_id"`
	LessonID  uint      `gorm:"type:integer;unique" json:"lesson_id"`
	URL       string    `gorm:"type:varchar(255);not null" json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Quiz model
type EeQuiz struct {
	QuizID    uint      `gorm:"primary_key" json:"quiz_id"`
	LessonID  uint      `gorm:"type:integer" json:"lesson_id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// QuizQuestion model
type EeQuizQuestion struct {
	QuestionID    uint      `gorm:"primary_key" json:"question_id"`
	QuizID        uint      `gorm:"type:integer" json:"quiz_id"`
	QuestionText  string    `gorm:"type:text;not null" json:"question_text"`
	CorrectAnswer string    `gorm:"type:text;not null" json:"correct_answer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// QuizAnswer model
type EeQuizAnswer struct {
	AnswerID   uint      `gorm:"primary_key" json:"answer_id"`
	QuestionID uint      `gorm:"not null;index:idx_user_question,unique" json:"question_id"`
	UserID     uint      `gorm:"not null;index:idx_user_question,unique" json:"user_id"`
	AnswerText string    `gorm:"type:text;not null" json:"answer_text"`
	IsCorrect  bool      `gorm:"not null" json:"is_correct"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
