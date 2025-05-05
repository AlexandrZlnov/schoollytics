package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AlexandrZlnov/schoollytics/internal/domain"
)

// сохраняет информацию о школьнике если его еще нет БД
func SaveStudent(db *sql.DB, student domain.Students, schoolID *int) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer tx.Rollback()

	var studentID int

	err = tx.QueryRow(`
        INSERT INTO students (user_id, profile_id, guid, first_name, last_name, middle_name, phone_number, 
authentication_token, person_id, pswrd_change_required, regional_auth, date_of_birth, sex, school_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		  ON CONFLICT (user_id) DO UPDATE SET
            regional_auth = EXCLUDED.regional_auth,
            date_of_birth = EXCLUDED.date_of_birth,
            phone_number = EXCLUDED.phone_number,
			school_id = $14
			RETURNING id`,
		student.UserID,
		student.Profile[0].ExternalID,
		student.GuID,
		student.FirstName,
		student.LastName,
		student.MiddleName,
		student.PhoneNumber,
		student.AuthenticationToken,
		student.PersonID,
		student.PswrdChangeRequired,
		student.RegionalAuth,
		student.DateOfBirht,
		student.Sex,
		schoolID,
	).Scan(&studentID)

	if err != nil {
		return 0, fmt.Errorf("ошибка при вставке/обновлении школьника: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("не удалось зафиксировать транзакцию: %v", err)
	}

	log.Println("Данные школьника успешно обновлены")
	return studentID, nil
}
