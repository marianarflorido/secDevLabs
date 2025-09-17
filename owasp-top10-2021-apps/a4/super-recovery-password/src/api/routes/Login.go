package routes

import (
	"api/database"
	"api/services"
	"api/types"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"time"
)

var failedAttempts = make(map[string]int)
var lockoutUntil = make(map[string]time.Time)


const maxAttempts = 5
const lockoutDuration = 15 * time.Minute

func Login(c echo.Context) (err error) {
	u := new(types.UserLogin)
	if err = c.Bind(u); err != nil {
		return
	}
	u.Login = strings.ToLower(u.Login)

	user := types.UserLogin{
		Login:    u.Login,
		Password: u.Password,
	}

	if until, ok := lockoutUntil[u.Login]; ok {
        if time.Now().Before(until) {
            return c.JSON(http.StatusTooManyRequests, map[string]string{
                "error": "Conta temporariamente bloqueada. Tente novamente mais tarde.",
            })
        }
		delete(lockoutUntil, u.Login)
		delete(failedAttempts, u.Login)	
	
	}

	success, err := database.LoginUser(user.Login, user.Password)
	if !success || err != nil {
		
		failedAttempts[u.Login]++
		if failedAttempts[u.Login] >= maxAttempts {
			lockoutUntil[u.Login] = time.Now().Add(lockoutDuration)
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Muitas tentativas falhas. Conta bloqueada por 15 minutos.",
			})
		}

		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Bad Credentials! User not authorized.",
		})
	}

	delete(failedAttempts, u.Login)
	delete(lockoutUntil, u.Login)

	token, err := services.GenerateJwt(u.Login, false)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Error to generate token.",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
