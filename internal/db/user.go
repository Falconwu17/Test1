package db_

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	db := variables.DB

	err := db.QueryRow("SELECT id, name_, surname, email, password FROM users WHERE email = $1", email).Scan(
		&user.Id, &user.Name, &user.Surname, &user.Email, &user.Password,
	)
	if err != nil {
		log.Printf("Ошибка получения пользователя по email: %v", err)
		return user, err
	}
	return user, nil
}
func CreateUser(user *models.User) error {
	db := variables.DB
	hashePasswd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка при хэшировании пароля: %v", err)
	}
	err = db.QueryRow("INSERT INTO users (name_, surname , email , password ) VALUES ($1,$2,$3,$4 ) RETURNING id", user.Name, user.Surname, user.Email, string(hashePasswd)).Scan(&user.Id)
	if err != nil {
		log.Printf("Ошибка создания пользователя: %v", err)
		return err
	}
	return nil
}
func CheckUserExistsByEmail(email string) (bool, error) {
	var count int
	err := variables.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", email).Scan(&count)
	return count > 0, err
}
