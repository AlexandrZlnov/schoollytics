package main

import (
	//"fmt"
	//"database/sql"
	"log"
	//"time"

	"github.com/AlexandrZlnov/schoollytics/internal/repository"
	"github.com/AlexandrZlnov/schoollytics/internal/service"
)

func main() {
	// проходим аутентификацию и получаем токен
	token := service.Authentication()

	// получаем ID школьника для формирования GET запроса
	// в котором получим JSON с данными оценок
	studentId, err := service.GetStudentInfo(token)
	if err != nil {
		log.Fatal(err)
	}

	// формируем url для Get запроса по которому получим JSON с даннымы оценок
	// url := fmt.Sprintf("https://school.mos.ru/api/family/web/v1/subject_marks?student_id=%d", studentId)
	if err := service.MakeAuthRequest(studentId, token); err != nil {
		log.Fatal("Ошибка аутентифицированного запроса: ", err)
	}

	// инициируем базу данных
	db, err := repository.InitDB()
	if err != nil {
		log.Fatalf("ошибка подключения к PostgreSql: %v", err)
	}
	defer db.Close()

	// db.SetMaxOpenConns(25)
	// db.SetMaxIdleConns(25)
	// db.SetConnMaxLifetime(5 * time.Minute)

	if err := service.ProcessingJSON(db); err != nil {
		log.Fatal("Ошибка обработки JSON: ", err)
	}
}
