package routes

import (
	"api/database"
	"net/http"

	"github.com/labstack/echo"
)

// ResetPasswordHandler reseta a senha usando o token enviado por e-mail
func ResetPasswordHandler(c echo.Context) error {
	type Request struct {
		Token          string `json:"token" form:"token" query:"token"`
		NewPassword    string `json:"newPassword" form:"newPassword" query:"newPassword"`
		RepeatPassword string `json:"repeatPassword" form:"repeatPassword" query:"repeatPassword"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Parâmetros inválidos"})
	}

	// Confirma se as senhas coincidem
	if req.NewPassword != req.RepeatPassword {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "As senhas não coincidem"})
	}

	// Tenta resetar a senha usando o token
	if err := database.ResetPasswordByToken(req.Token, req.NewPassword); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Senha alterada com sucesso"})
}
