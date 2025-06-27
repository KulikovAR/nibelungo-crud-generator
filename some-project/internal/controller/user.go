package controller

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/domain"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/usecase"
)

type userController struct {
	useCase usecase.userUseCase
}

func NewuserController(useCase usecase.userUseCase) *userController {
	return &userController{useCase: useCase}
}

// Createuser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.user true "user object"
// @Success 201 {object} domain.user
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [post]
func (c *userController) Create(ctx *gin.Context) {
	var entity domain.user
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.useCase.Create(ctx, &entity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, entity)
}

// Getuser godoc
// @Summary Get a user by ID
// @Description Get a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "user ID"
// @Success 200 {object} domain.user
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [get]
func (c *userController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	entity, err := c.useCase.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if entity == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, entity)
}

// Updateuser godoc
// @Summary Update a user
// @Description Update a user with the input payload
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "user ID"
// @Param user body domain.user true "user object"
// @Success 200 {object} domain.user
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [put]
func (c *userController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var entity domain.user
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity.ID = id
	if err := c.useCase.Update(ctx, &entity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entity)
}

// Deleteuser godoc
// @Summary Delete a user
// @Description Delete a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "user ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [delete]
func (c *userController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := c.useCase.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Listuser godoc
// @Summary List all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} domain.user
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func (c *userController) List(ctx *gin.Context) {
	entities, err := c.useCase.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entities)
}
