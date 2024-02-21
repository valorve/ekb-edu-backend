package courses

import (
	"ekb-edu/src/api/middleware"
	"ekb-edu/src/database/storage"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func getCourses(c *fiber.Ctx) error {
	var courses []storage.EeCourse

	if err := storage.DB.Find(&courses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", err.Error())})
	}

	return c.JSON(courses)
}

func getCourse(c *fiber.Ctx) error {
	var course storage.EeCourse

	result := storage.DB.
		Where("course_id = ?", c.Params("id")).
		First(&course)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.JSON(&course)
}

func addCourse(c *fiber.Ctx) error {
	courseInfo := storage.EeCourse{}

	if err := c.BodyParser(&courseInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse course data"})
	}

	result := storage.DB.Create(&courseInfo)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.SendStatus(fiber.StatusOK)
}

func getSection(c *fiber.Ctx) error {
	var section storage.EeCourseSection

	result := storage.DB.
		Where("section_id = ?", c.Params("id")).
		First(&section)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.JSON(&section)
}

func addSection(c *fiber.Ctx) error {
	courseSection := storage.EeCourseSection{}

	if err := c.BodyParser(&courseSection); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse section data"})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course id"})
	}

	courseSection.CourseID = uint(id)

	result := storage.DB.Create(&courseSection)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.JSON(&courseSection)
}

func addSections(c *fiber.Ctx) error {
	courseSection := []*storage.EeCourseSection{}

	if err := c.BodyParser(&courseSection); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse sections data"})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course id"})
	}

	for _, course := range courseSection {
		course.CourseID = uint(id)
	}

	result := storage.DB.Create(&courseSection)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.SendStatus(fiber.StatusOK)
}

func getCourseSections(c *fiber.Ctx) error {
	var sections []storage.EeCourseSection

	result := storage.DB.
		Where("course_id = ?", c.Params("id")).
		Find(&sections)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.JSON(sections)
}

func linkUser(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("user_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid User ID",
		})
	}

	courseID, err := strconv.ParseUint(c.Params("course_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course id"})
	}

	courseOwner := storage.EeCourseOwner{
		UserID:   uint(userID),
		CourseID: uint(courseID),
	}

	result := storage.DB.Create(&courseOwner)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.SendStatus(fiber.StatusOK)
}

func getMyCourses(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	var courses []storage.EeCourse

	if err := storage.DB.Where("instructor_id = ?", userID).Find(&courses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", err.Error())})
	}

	return c.JSON(courses)
}

func deleteCourse(c *fiber.Ctx) error {
	courseID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course id"})
	}

	course := storage.EeCourse{
		CourseID: uint(courseID),
	}

	result := storage.DB.Delete(&course)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.SendStatus(fiber.StatusOK)
}

func updateCourse(c *fiber.Ctx) error {
	courseInfo := storage.EeCourse{}

	if err := c.BodyParser(&courseInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse course data"})
	}

	courseID, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course id"})
	}

	result := storage.DB.Model(&storage.EeCourse{}).Where("course_id = ?", courseID).Updates(&courseInfo)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.SendStatus(fiber.StatusOK)
}

func updateSection(c *fiber.Ctx) error {
	sectionInfo := storage.EeCourseSection{}

	if err := c.BodyParser(&sectionInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse section data",
		})
	}

	sectionID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid section id"})
	}

	result := storage.DB.Model(&storage.EeCourseSection{}).Where("section_id = ?", sectionID).Updates(sectionInfo)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.SendStatus(fiber.StatusOK)
}

func deleteSection(c *fiber.Ctx) error {
	sectionID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid section id"})
	}

	section := storage.EeCourseSection{
		SectionID: uint(sectionID),
	}

	result := storage.DB.Delete(&section)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.SendStatus(fiber.StatusOK)
}

func getLessonsBySection(c *fiber.Ctx) error {
	sectionIDParam := c.Params("id")
	sectionID, err := strconv.Atoi(sectionIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid section ID",
		})
	}

	var lessons []storage.EeLesson
	result := storage.DB.Where("section_id = ?", sectionID).Find(&lessons)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %s", result.Error.Error())})
	}

	return c.JSON(lessons)
}

func RegisterService(app fiber.Router) {
	courses := app.Group("/courses")
	{
		courses.Get("/", getCourses)
		courses.Get("/:id", getCourse)

		courses.Get("/:id/sections", getCourseSections)
		courses.Get("/sections/:id", getSection)

		admin := courses.Group("/", middleware.TokenRequired, middleware.AdminRequired)
		{
			admin.Post("/", addCourse)
			admin.Patch("/:id", updateCourse)

			admin.Post("/:id/section", addSection)
			admin.Post("/:id/sections", addSections)

			admin.Post("/link/:course_id/:user_id", linkUser)

			admin.Delete("/:id", deleteCourse)

			{
				admin.Patch("/sections/:id", updateSection)
				admin.Delete("/sections/:id", deleteSection)
				admin.Get("/sections/:id/lessons", getLessonsBySection)
			}
		}
	}

	app.Get("/my_courses", middleware.TokenRequired, middleware.AdminRequired, getMyCourses)
}
