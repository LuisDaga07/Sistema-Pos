package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pos-saas/restaurant-pos/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary      Registrar restaurante
// @Description  Crea un nuevo restaurante con usuario admin
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  service.RegisterInput  true  "Datos de registro"
// @Success      201   {object}  service.AuthResponse
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var input service.RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "datos inválidos: " + err.Error()})
		return
	}

	resp, err := c.authService.Register(ctx.Request.Context(), input)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary      Iniciar sesión
// @Description  Autentica usuario y devuelve JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  service.LoginInput  true  "Credenciales"
// @Success      200   {object}  service.AuthResponse
// @Failure      401   {object}  map[string]string
// @Router       /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var input service.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "datos inválidos: " + err.Error()})
		return
	}

	resp, err := c.authService.Login(ctx.Request.Context(), input)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
