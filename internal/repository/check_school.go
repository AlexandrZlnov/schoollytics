package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AlexandrZlnov/schoollytics/internal/domain"
)

func CheckSchool(db *sql.DB, school domain.Students) (int, error) {
	var schoolID int
	fmt.Println("========ПРОВЕРИМ ДЛИНУ PROFILE:", len(school.Profile))

	if len(school.Profile) > 0 {
		profile := school.Profile[0]

		err := db.QueryRow(`SELECT id FROM schools WHERE organization_id = $1`,
			profile.OrganizationID).Scan(&schoolID)
		fmt.Println("========SCHOOLID------:", schoolID, err)

		if err == nil {
			_, err = SaveSchool(db, school)
			if err != nil {
				log.Fatalf("ошибка обновления данных школы: %v", err)
			} else {
				log.Printf("Данные школы обновлены")
			}
		}

		if err == sql.ErrNoRows {
			// школы нет в БД
			// сохраняем данные по школе
			log.Println("=== ВЫЗЫВАЕМ SaveSchool ===")
			schoolID, err = SaveSchool(db, school)
			if err != nil {
				log.Fatalf("ошибка сохранения школы: %v", err)
			}
		}

	}
	return schoolID, nil
}
