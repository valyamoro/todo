package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/valyamoro/TODO/internal/domain"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Tasks interface {
	Create(task domain.Task) (domain.Task, error)
	GetByID(id int64) (domain.Task, error)
	GetAll() ([]domain.Task, error)
	Delete(id int64) (domain.Task, error)
	Update(id int64, inp domain.UpdateTaskInput) (domain.Task, error)
}

type Handler struct { 
	tasksService Tasks
}

func NewHandler(tasks Tasks) *Handler {
	return &Handler{
		tasksService: tasks,
	}
}

func (h *Handler) InitRouter(logger *zap.Logger) *gin.Engine {
	r := gin.Default()
	r.Use(LoggingMiddleware(logger))

	tasks := r.Group("/tasks")
	{
		tasks.POST("", h.createTask)
		tasks.GET("", h.getAllTasks)
		tasks.GET("/:id", h.getTaskByID)
		tasks.DELETE("/:id", h.deleteTask)
		tasks.PUT("/:id", h.updateTask)
	}
	
	return r
}

func (h *Handler) getTaskByID(c *gin.Context) {
	id, err := getIdFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	task, err := h.tasksService.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return 
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) createTask(c *gin.Context) {
	var task domain.Task 
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	_, err := h.tasksService.Create(task); 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	c.JSON(http.StatusCreated, task)
}

func (h *Handler) deleteTask(c *gin.Context) {
	id, err := getIdFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	_, err = h.tasksService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getAllTasks(c *gin.Context) {
	tasks, err := h.tasksService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) updateTask(c *gin.Context) {
	id, err := getIdFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	var inp domain.UpdateTaskInput 
	if err := c.ShouldBindJSON(&inp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	_, err = h.tasksService.Update(id, inp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	c.Status(http.StatusOK)
}

func getIdFromRequest(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err 
	}

	if id == 0 {
		return 0, errors.New("id cant be 0")
	}

	return id, nil 
}
