package quizzes

import (
	"ekb-edu/src/api/middleware"
	"ekb-edu/src/database/storage"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetQuizzes(c *fiber.Ctx) error {
	lessonID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid Lesson ID"})
	}

	var quizzes []storage.EeQuiz
	if err := storage.DB.Where("lesson_id = ?", lessonID).Find(&quizzes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", err.Error())})
	}

	return c.JSON(&quizzes)
}

func CreateQuiz(c *fiber.Ctx) error {
	quiz := storage.EeQuiz{}

	lessonID, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "lesson ID is required"})
	}

	if err := c.BodyParser(&quiz); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse quiz data"})
	}

	quiz.LessonID = uint(lessonID)
	if err := storage.DB.Create(&quiz).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", err.Error())})
	}

	return c.JSON(&quiz)
}

// QuizQuestion model without the correct answer
type EeQuizQuestionWithoutAnswer struct {
	QuestionID   uint      `gorm:"primary_key" json:"question_id"`
	QuizID       uint      `gorm:"type:integer" json:"quiz_id"`
	QuestionText string    `gorm:"type:text;not null" json:"question_text"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EeQuizWithQuestions struct {
	Quiz      storage.EeQuiz                `json:"quiz"`
	Questions []EeQuizQuestionWithoutAnswer `json:"questions"`
}

func getQuiz(c *fiber.Ctx) error {
	quizID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	// Проверка на наличие ID урока для обновления.
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "quiz ID is required"})
	}

	var quizInfo EeQuizWithQuestions
	if err := storage.DB.Where("quiz_id = ?", quizID).First(&quizInfo.Quiz).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", err.Error())})
	}

	if err := storage.DB.Where("quiz_id = ?", quizID).Find(&quizInfo.Questions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", err.Error())})
	}

	return c.JSON(&quizInfo)
}

func addQuestion(c *fiber.Ctx) error {
	quizID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	// Проверка на наличие ID урока для обновления.
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "quiz ID is required"})
	}

	var question storage.EeQuizQuestion

	if err := c.BodyParser(&question); err != nil {
		// В случае ошибки разбора возвращаем HTTP статус 400 (Bad Request).
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse quiz data"})
	}

	question.QuizID = uint(quizID)

	if err := storage.DB.Create(&question).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", err.Error())})
	}

	return c.JSON(&question)
}

func answerQuestion(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	quizID, err := c.ParamsInt("quiz_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid quiz ID"})
	}

	questionID, err := c.ParamsInt("question_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid question ID"})
	}

	// Парсим тело запроса для получения ответа
	var answer storage.EeQuizAnswer
	if err := c.BodyParser(&answer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse the answer"})
	}

	answer.QuestionID = uint(questionID)
	answer.UserID = userID

	// Проверяем, является ли пользователь владельцем курса
	var courseOwnerCount int64
	storage.DB.Table("ee_course_owners").
		Joins("JOIN ee_courses ON ee_course_owners.course_id = ee_courses.course_id").
		Joins("JOIN ee_course_sections ON ee_courses.course_id = ee_course_sections.course_id").
		Joins("JOIN ee_lessons ON ee_course_sections.section_id = ee_lessons.section_id").
		Joins("JOIN ee_quizzes ON ee_lessons.lesson_id = ee_quizzes.lesson_id").
		Where("ee_course_owners.user_id = ? AND ee_quizzes.quiz_id = ?", userID, quizID).
		Count(&courseOwnerCount)

	if courseOwnerCount == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "user is not the owner of the course"})
	}

	var question storage.EeQuizQuestion
	if err := storage.DB.Where("question_id = ?", questionID).First(&question).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", err.Error())})
	}

	answer.IsCorrect = answer.AnswerText == question.CorrectAnswer

	// Если пользователь является владельцем, добавляем ответ
	if err := storage.DB.Create(&answer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", err.Error())})
	}

	return c.Status(fiber.StatusOK).JSON(&answer)
}

func RegisterService(app fiber.Router) {
	g := app.Group("/quizzes", middleware.TokenRequired)
	{
		g.Get("/:id", getQuiz)
		g.Post("/:quiz_id/:question_id", answerQuestion)

		admin := g.Group("/", middleware.AdminRequired)
		{
			admin.Post("/:id", addQuestion)
		}
	}
}
