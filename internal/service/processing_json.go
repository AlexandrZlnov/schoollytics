package service

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/AlexandrZlnov/schoollytics/internal/domain"
	"github.com/AlexandrZlnov/schoollytics/internal/repository"
)

func ProcessingJSON(db *sql.DB) error {

	// читаем файл JSON с данными студента
	jsonDataStudent, err := os.ReadFile("student_info.json")
	if err != nil {
		log.Fatalf("ошибка чтения файла: %v", err)
	}

	// читаем файл JSON с данными оценок
	jsonDataMarks, err := os.ReadFile("studentPerformance.json")
	if err != nil {
		log.Fatalf("ошибка чтения json: %v", err)
	}

	// десериализуем JSON с данными студента в структуру
	var school domain.Students
	if err := json.Unmarshal(jsonDataStudent, &school); err != nil {
		log.Fatalf("ошибка десериализации json: %v", err)
	}

	// десериализуем JSON с данными оценок  в структуру
	var studentGrades domain.StudentPerformance
	err = json.Unmarshal(jsonDataMarks, &studentGrades)
	if err != nil {
		log.Fatalf("ошибка десериализации json: %v", err)
	}

	// проверяем наличие школы в БД
	var schoolID int
	schoolID, err = repository.CheckSchool(db, school)
	if err != nil {
		log.Fatalf("ошибка проверки наличих школы в БД: %v", err)
	}

	var studentID int
	// сохраняем данные по ученику в БД
	// и получаем его ID
	studentID, err = repository.SaveStudent(db, school, &schoolID)
	if err != nil {
		log.Fatalf("ошибка сохранения школьника: %v", err)
	}

	// БЛОК ВСТАВКИ ДАННЫХ ОБ ОЦЕНКАХ В БД ------------------->
	err = repository.SaveGrades(studentID, db, &studentGrades)
	if err != nil {
		log.Fatal("ошибка сохранения оценок в БД: %w", err)
	}
	return nil
}
