/*
Часть 1 - Аутентификация
Аутентификация получение токена и вополнение авторизированного запроса
для получения данных о школьнике
*/

package main

import (
	"fmt"
	"log"

	//"strings"
	//"github.com/playwright-community/playwright-go"
	"test/service"
)

func main() {
	// проходим аутентификацию и получаем токен
	token := service.Authentication()

	// получаем ID школьника для формирования GET запроса
	// в котором получим JSON с данными оценок
	// !!! дальнейшая работа идет в папке test2
	studentId, err := service.GetStudentInfo(token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("StudentID:", studentId)

	// блок функций который получает данные о студенре
	//takeStudentInformation

	// формируем url для Get запроса по которому получим JSON с даннымы оценок
	url := fmt.Sprintf("https://school.mos.ru/api/family/web/v1/subject_marks?student_id=%d", studentId)
	if err := service.MakeAuthRequest(url, token); err != nil {
		log.Fatal("Ошибка аутентифицированного запроса: ", err)
	}

}
