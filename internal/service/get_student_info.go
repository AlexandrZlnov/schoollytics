package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AlexandrZlnov/schoollytics/internal/domain"
)

// направляет GET запрос по адресу https://school.mos.ru/api/ej/acl/v1/sessions
// и получает JSON с ID пользователя и студента
// ID студента будет использован для получения JSON с отметками
func GetStudentInfo(token string) (int, error) {
	// Создаем http слиент
	client := &http.Client{}

	requestBody := map[string]string{
		"auth_token": token,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания тела запроса: %v", err)
	}

	// Создаем новый запрос
	url := "https://school.mos.ru/api/ej/acl/v1/sessions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer((jsonBody)))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}

	// Добавляем заголовки
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("x-mes-subsystem", "familyweb")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return 00, fmt.Errorf("request failed: %v, Статус: %v", err, resp.StatusCode)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("server returned error in resp.StatCode: %s, body: %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка чтения ответа сервера: %v", err)
	}

	// Сохраняем сырой JSON в файл c ощими даннымиы студента
	err = os.WriteFile("../student_info.json", body, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to save response to file: %v", err)
	}

	var student domain.StudentInfo

	if err := json.Unmarshal(body, &student); err != nil {
		return 0, fmt.Errorf("ошибка ошибка десериализации ID студента: %v", err)
	}

	fmt.Println("----------->>>>>>>>>>", student)
	return student.Profiles[0].StudentID, nil
}

// func parseStudentData(jsonData []byte) (*domain.StudentInfo, error) {
// 	var student domain.StudentInfo
// 	err := json.Unmarshal(jsonData, &student)
// 	if err != nil {
// 		return nil, fmt.Errorf("error parsing JSON: %v", err)
// 	}
// 	return &student, nil
// }
