package lessons

import (
	"ekb-edu/src/api/middleware"
	"ekb-edu/src/database/storage"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Struct to aggregate lesson with course and section info
type LessonWithCourseAndSection struct {
	storage.EeLesson        // Embedding EeLesson to include all its fields
	CourseTitle      string `json:"course_title"`
	CourseID         uint   `json:"course_id"`
	SectionTitle     string `json:"section_title"`
	SectionID        uint   `json:"section_id"`
}

func getLessons(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	// Сначала находим все курсы, принадлежащие пользователю
	var ownedCourses []storage.EeCourseOwner
	storage.DB.Where("user_id = ?", userID).Find(&ownedCourses)
	var courseIDs []uint
	for _, co := range ownedCourses {
		courseIDs = append(courseIDs, co.CourseID)
	}

	// Теперь находим все секции этих курсов
	var sections []storage.EeCourseSection
	storage.DB.Where("course_id IN ?", courseIDs).Find(&sections)
	var sectionIDs []uint
	for _, section := range sections {
		sectionIDs = append(sectionIDs, section.SectionID)
	}

	// И, наконец, находим все уроки для этих секций
	var lessons []storage.EeLesson
	storage.DB.Where("section_id IN ?", sectionIDs).Find(&lessons)

	// Создаем словарь для курсов и секций для удобства доступа
	courseMap := make(map[uint]storage.EeCourse)
	sectionMap := make(map[uint]storage.EeCourseSection)
	for _, course := range ownedCourses {
		var c storage.EeCourse
		storage.DB.First(&c, course.CourseID)
		courseMap[c.CourseID] = c
	}
	for _, section := range sections {
		sectionMap[section.SectionID] = section
	}

	// Агрегируем уроки с информацией о курсе и секции
	var result []LessonWithCourseAndSection
	for _, lesson := range lessons {
		section := sectionMap[lesson.SectionID]
		course := courseMap[section.CourseID]
		result = append(result, LessonWithCourseAndSection{
			EeLesson:     lesson,
			CourseTitle:  course.Title,
			CourseID:     course.CourseID,
			SectionTitle: section.Title,
			SectionID:    section.SectionID,
		})
	}

	// Возвращаем агрегированные уроки как JSON
	return c.JSON(result)
}

func addLesson(c *fiber.Ctx) error {
	lesson := storage.EeLesson{}

	// Разбор JSON тела запроса в структуру EeLesson.
	if err := c.BodyParser(&lesson); err != nil {
		// В случае ошибки разбора возвращаем HTTP статус 400 (Bad Request).
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse lesson data"})
	}

	result := storage.DB.Create(&lesson)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Database error: %s", result.Error.Error())})
	}

	return c.JSON(&lesson)
}

func getLesson(c *fiber.Ctx) error {
	var lesson storage.EeLesson

	result := storage.DB.
		Where("lesson_id = ?", c.Params("id")).
		First(&lesson)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Database error: %s", result.Error.Error())})
	}

	return c.JSON(&lesson)
}

func updateLesson(c *fiber.Ctx) error {
	lessonInfo := storage.EeLesson{}

	// Разбор JSON тела запроса в структуру EeLesson.
	if err := c.BodyParser(&lessonInfo); err != nil {
		// В случае ошибки разбора возвращаем HTTP статус 400 (Bad Request).
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse lesson data"})
	}

	lessonID, err := strconv.ParseUint(c.Params("id"), 10, 32)

	// Проверка на наличие ID урока для обновления.
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Lesson ID is required"})
	}

	// Обновление урока в базе данных.
	result := storage.DB.Model(&storage.EeLesson{}).Where("lesson_id = ?", lessonID).Updates(lessonInfo)
	if result.Error != nil {
		// В случае ошибки обновления возвращаем HTTP статус 500 (Internal Server Error).
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update lesson"})
	}

	// Возвращение HTTP статуса 200 (OK), если обновление прошло успешно.
	return c.SendStatus(fiber.StatusOK)
}

func deleteLesson(c *fiber.Ctx) error {
	lessonID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Lesson ID"})
	}

	lesson := storage.EeLesson{
		LessonID: uint(lessonID),
	}

	result := storage.DB.Delete(&lesson)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Database error: %s", result.Error.Error())})
	}

	return c.SendStatus(fiber.StatusOK)
}

func RegisterService(app fiber.Router) {
	g := app.Group("/lessons", middleware.TokenRequired)
	{
		g.Get("/", getLessons)
		g.Get("/:id", getLesson)

		admin := g.Group("/", middleware.AdminRequired)
		{
			admin.Post("/", addLesson)
			admin.Patch("/:id", updateLesson)
			admin.Delete("/:id", deleteLesson)
		}
	}
}
