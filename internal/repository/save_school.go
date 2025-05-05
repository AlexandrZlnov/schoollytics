package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AlexandrZlnov/schoollytics/internal/domain"
)

// сохраняет информацию о школе если ее еще нет в БД
func SaveSchool(db *sql.DB, school domain.Students) (int, error) {
	log.Printf("Пытаемся сохранить школу: ID=%d, Name=%s, OrgID=%s",
		school.Profile[0].SchoolID,
		school.Profile[0].Name,
		school.Profile[0].OrganizationID)

	var schoolID int

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer tx.Rollback()
	log.Println("---------> 0")
	err = tx.QueryRow(`
        INSERT INTO schools (school_id, name, shortname, organization_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (organization_id) DO UPDATE SET
            school_id = EXCLUDED.school_id,
            name = EXCLUDED.name,
            shortname = EXCLUDED.shortname 
			RETURNING id`,
		school.Profile[0].SchoolID,
		school.Profile[0].Name,
		school.Profile[0].Shortname,
		school.Profile[0].OrganizationID,
	).Scan(&schoolID)

	fmt.Println(schoolID)
	log.Println("---------> 1")

	if err != nil {
		return 0, fmt.Errorf("ошибка при вставке/обновлении школы: %v", err)
	}
	log.Println("---------> 2")
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("не удалось зафиксировать транзакцию: %v", err)
	}
	log.Println("---------> 3")
	log.Println("Данные школы успешно обновлены")
	return schoolID, nil
}
