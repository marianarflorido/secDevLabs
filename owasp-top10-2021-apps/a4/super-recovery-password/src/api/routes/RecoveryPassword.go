package routes

import (
	"api/database"
	"api/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// RecoveryPasswordHandler substitui perguntas de segurança por token via e-mail
func RecoveryPasswordHandler(c echo.Context) error {
	type Request struct {
		Login string `json:"login" form:"login" query:"login"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Parâmetros inválidos"})
	}

	req.Login = strings.ToLower(req.Login)

	// 1️⃣ Pega e-mail do usuário
	user, err := database.UserEmail(req.Login)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Usuário não encontrado"})
	}

	// 2️⃣ Gera token e salva no banco
	token, err := database.GenerateRecoveryToken(user.Login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Erro ao gerar token"})
	}

	// 3️⃣ Envia e-mail de recuperação
	if err := utils.SendRecoveryEmail(user.Email, token); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Erro ao enviar e-mail"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "E-mail de recuperação enviado. Verifique sua caixa de entrada.",
	})
}
