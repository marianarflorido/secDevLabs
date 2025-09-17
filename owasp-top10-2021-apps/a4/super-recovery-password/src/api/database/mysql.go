package database

import (
	"api/types"
	"api/utils"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDatabaseConnection() (*sql.DB, error) {

	databaseEndpoint := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := sql.Open("mysql", databaseEndpoint)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(40)
	return db, nil
}

func ChangePassword(login string, password string, repeatPassword string) error {
	db, err := OpenDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	userExists, err := CheckUserExists(login)
	if !userExists {
		return err
	}

	passwordHashed, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.Query("UPDATE Users SET Password = ? WHERE Login = ?", passwordHashed, login)
	if err != nil {
		return err
	}

	return nil
}

func GenerateRecoveryToken(login string) (string, error) {
	db, err := OpenDatabaseConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Garante que usuário existe
	userExists, err := CheckUserExists(login)
	if !userExists {
		return "", errors.New("usuário não encontrado")
	}

	// Gera token
	token, err := utils.GenerateToken()
	if err != nil {
		return "", err
	}

	// Salva no banco
	_, err = db.Exec("UPDATE Users SET RecoveryToken = ? WHERE Login = ?", token, login)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUserByToken(token string) (types.RecoveryEmail, error) {
	db, err := OpenDatabaseConnection()
	if err != nil {
		return types.RecoveryEmail{}, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT Login, Email FROM Users WHERE RecoveryToken = ?", token)

	recoveryEmail := types.RecoveryEmail{}
	err = row.Scan(&recoveryEmail.Login, &recoveryEmail.Email)
	if err != nil {
		return types.RecoveryEmail{}, err
	}

	return recoveryEmail, nil
}

func ResetPasswordByToken(token, newPassword string) error {
	db, err := OpenDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Pega usuário pelo token
	row := db.QueryRow("SELECT Login FROM Users WHERE RecoveryToken = ?", token)
	var login string
	err = row.Scan(&login)
	if err != nil {
		return errors.New("token inválido ou expirado")
	}

	// Hash da nova senha
	passwordHashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Atualiza senha e limpa token
	_, err = db.Exec("UPDATE Users SET Password = ?, RecoveryToken = NULL WHERE Login = ?", passwordHashed, login)
	if err != nil {
		return err
	}

	return nil
}

func RegisterUser(login string, password string, email string) (bool, error) {

	db, err := OpenDatabaseConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	userExists, err := CheckUserExists(login)
	if userExists {
		return false, err
	}

	emailExists, err := CheckEmailExists(email)
	if emailExists {
		return false, err
	}

	passwordHashed, err := utils.HashPassword(password)
	if err != nil {
		return false, err
	}

	result, err := db.Query("INSERT INTO Users (Login, Password, Email) VALUES (?, ?, ?)", login, passwordHashed, email)
	if err != nil {
		return false, err
	}
	defer result.Close()

	return true, nil
}

func LoginUser(login string, password string) (bool, error) {
	db, err := OpenDatabaseConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	userExists, err := CheckUserExists(login)
	if !userExists {
		return false, err
	}

	result, err := db.Query("SELECT ID, Login, Password FROM Users WHERE Login = ?", login)
	if err != nil {
		return false, err
	}
	defer result.Close()

	loginUser := types.UserLogin{}

	for result.Next() {
		err := result.Scan(&loginUser.ID, &loginUser.Login, &loginUser.Password)
		if err != nil {
			return false, err
		}
	}

	err = CheckUserPassword(password, loginUser.Password)
	if err != nil {
		return false, err
	}

	return true, nil
}
func UserEmail(login string) (types.RecoveryEmail, error) {
	db, err := OpenDatabaseConnection()
	if err != nil {
		return types.RecoveryEmail{}, err
	}
	defer db.Close()

	userExists, err := CheckUserExists(login)
	if !userExists {
		return types.RecoveryEmail{}, err
	}

	row := db.QueryRow("SELECT Login, Email FROM Users WHERE Login = ?", login)

	recoveryEmail := types.RecoveryEmail{}
	err = row.Scan(&recoveryEmail.Login, &recoveryEmail.Email)
	if err != nil {
		return types.RecoveryEmail{}, err
	}

	return recoveryEmail, nil
}

func CheckUserExists(login string) (bool, error) {
	db, err := OpenDatabaseConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM Users WHERE Login = ?", login)
	if err != nil {
		return false, err
	}
	defer result.Close()

	if !result.Next() {
		err = errors.New("Error: User not found!")
		return false, err
	}
	return true, nil
}

func CheckEmailExists(email string) (bool, error) {
	db, err := OpenDatabaseConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM Users WHERE Email = ?", email)
	if err != nil {
		return false, err
	}
	defer result.Close()

	if !result.Next() {
		err = errors.New("Error: Email not found!")
		return false, err
	}
	return true, nil
}

func CheckUserPassword(password string, hashedPassword string) error {
	hashedPass := []byte(hashedPassword)

	err := utils.ComparePassword(password, hashedPass)
	if err != nil {
		return err
	}

	return nil
}

func initUsers(db *sql.DB) error {
	users := [9]string{
		"admin",
		"thiago",
		"teste",
		"joana",
		"roberto",
		"estag",
		"marleda",
		"junior",
		"maria",
	}
	for _, u := range users {
		hashPassword, err := utils.HashPassword(u)
		_, err = db.Query("INSERT INTO Users (Login, Password) VALUES (?, ?)", u, hashPassword)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitDatabase() error {

	dbConn, err := OpenDatabaseConnection()
	if err != nil {
		errOpenDBConnection := fmt.Sprintf("OpenDBConnection error: %s", err)
		return errors.New(errOpenDBConnection)
	}

	defer dbConn.Close()

	_, err = dbConn.Query("SELECT * FROM Users")

	if err != nil {
		queryCreate := fmt.Sprint(`
		CREATE TABLE Users (
			ID int NOT NULL AUTO_INCREMENT, 
			Login varchar(20), 
			Password varchar(80), 
			Email varchar(80),
			RecoveryToken varchar(200),
			PRIMARY KEY (ID)
		);`)
		_, err = dbConn.Exec(queryCreate)
	}

	if err != nil {
		errInitDB := fmt.Sprintf("InitDatabase error: %s", err)
		return errors.New(errInitDB)
	}

	err = initUsers(dbConn)
	if err != nil {
		return err
	}

	return nil
}
