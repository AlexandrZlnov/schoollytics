package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
)

func MakeAuthRequest(studentId int, token string) error {
	url := fmt.Sprintf("https://school.mos.ru/api/family/web/v1/subject_marks?student_id=%d", studentId)

	// Создаем cookie jar для хранения cookies между запросами
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("failed to create cookie jar: %v", err)
	}

	// Создаем HTTP клиент с поддержкой cookies
	client := &http.Client{
		Jar: jar,
	}

	// Создаем новый запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Добавляем заголовки
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("x-mes-subsystem", "familyweb")
	req.Header.Add("Accept", "application/json")

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned error: %s, body: %s", resp.Status, string(body))
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Сохраняем сырой JSON в файл
	err = os.WriteFile("../response.json", body, 0644)
	if err != nil {
		return fmt.Errorf("failed to save response to file: %v", err)
	}

	// Парсим JSON в интерфейс (работает с любой структурой JSON)
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Выводим информацию о cookies
	fmt.Println("Received cookies:")
	for i, cookie := range resp.Cookies() {
		fmt.Printf("[%d] %s=%s (Domain: %s, Path: %s, Expires: %v, Secure: %t, HttpOnly: %t)\n",
			i+1, cookie.Name, cookie.Value, cookie.Domain, cookie.Path, cookie.Expires, cookie.Secure, cookie.HttpOnly)
	}

	// Выводим информацию о JSON (первые несколько полей)
	// fmt.Println("\nJSON response preview:")
	// switch v := jsonData.(type) {
	// case map[string]interface{}:
	// 	fmt.Println("Response is JSON object")
	// 	i := 0
	// 	for key, val := range v {
	// 		fmt.Printf("  %s: %v\n", key, val)
	// 		i++
	// 		if i >= 3 {
	// 			break
	// 		}
	// 	}
	// case []interface{}:
	// 	fmt.Println("Response is JSON array")
	// 	for i := 0; i < len(v) && i < 3; i++ {
	// 		fmt.Printf("  [%d]: %v\n", i, v[i])
	// 	}
	// default:
	// 	fmt.Printf("Unknown JSON type: %T\n", v)
	// }

	return nil
}
