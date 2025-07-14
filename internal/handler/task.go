package handler

import (
	"go-task-service/internal/repository"
	"go-task-service/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type TaskHandler struct {
	svc  *service.TaskService
	repo *repository.TaskRepo
}

func New(svc *service.TaskService, repo *repository.TaskRepo) *TaskHandler {
	return &TaskHandler{svc, repo}
}

// @Summary Запуск задач в работу
// @Description Отправляет задачи на выполнение
// @Accept json
// @Produce json
// @Param data body ProcessRequest true "Параметры обработки"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /process [post]
func (h *TaskHandler) ProcessHandler(c *fiber.Ctx) error {
	req := new(ProcessRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body parameters"})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tasksInWork := h.svc.ProcessTasks(req.Limit, req.MinDelay, req.MaxDelay, req.SuccessRate)
	// у нас может быть меньше задач в пуле - чем запрошено в запросе
	return c.JSON(
		fiber.Map{"status": "processing", "count": tasksInWork},
	)
}

// @Summary Список задач
// @Description Получает список задач с пагинацией и фильтрацией
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы"
// @Param pageSize query int false "Размер страницы"
// @Param status query string false "Фильтр по статусу (NEW, PROCESSING, PROCESSED, FAILED)"
// @Success 200 {object} handler.TaskListResponse
// @Router /tasks [get]
func (h *TaskHandler) GetTasksHandler(c *fiber.Ctx) error {
	var query TaskQueryParams
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверные параметры запроса"})
	}

	// Значения по умолчанию
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	tasks, total, err := h.repo.GetTasks(query.Page, query.PageSize, query.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка при получении задач"})
	}

	return c.JSON(TaskListResponse{
		Data:     tasks,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	})
}
