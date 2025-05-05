package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"

	//"encoding/json"
	//"fmt"
	//"io"
	//"net/http"
	//"net/http/cookiejar"
	//"os"
	"time"
)

// сгенерировать Хэш
func generateURLHash(original string) string {
	hash := sha256.Sum256([]byte(original))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:7] // Берем первые 7 символов
}

// Unix время с милисекундами
func unixTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func makeAuthenticatedRequest(url, bearerToken string) error {
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
	req.Header.Add("Authorization", "Bearer "+bearerToken)
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
	err = os.WriteFile("response.json", body, 0644)
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
	fmt.Println("\nJSON response preview:")
	switch v := jsonData.(type) {
	case map[string]interface{}:
		fmt.Println("Response is JSON object")
		i := 0
		for key, val := range v {
			fmt.Printf("  %s: %v\n", key, val)
			i++
			if i >= 3 {
				break
			}
		}
	case []interface{}:
		fmt.Println("Response is JSON array")
		for i := 0; i < len(v) && i < 3; i++ {
			fmt.Printf("  [%d]: %v\n", i, v[i])
		}
	default:
		fmt.Printf("Unknown JSON type: %T\n", v)
	}

	return nil
}

func main1() {
	// fmt.Println("Хэш сумма в кодировке SHA256:", generateURLHash("https://mos.ru"))
	// fmt.Println("Время в формате Unix с милисекундами:", unixTime())

	url := "https://school.mos.ru/api/family/web/v1/subject_marks?student_id=10205869"
	token := "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxMTkyODU0Iiwic2NwIjoib3BlbmlkIHByb2ZpbGUiLCJuYmYiOjE3NDI5Mjg1NDUsIm1zaCI6IjQxMGQwMDZjLTQ3Y2QtNDFlNC04OWU4LWU3YzdlN2Y3NGJkYyIsImF0aCI6InN1ZGlyIiwiaXNzIjoiaHR0cHM6XC9cL3NjaG9vbC5tb3MucnUiLCJybHMiOiJ7MjpbMjE6MjpbXSwzMTo0OltdLDQxOjE6W10sNTM0OjQ4OltdLDE4MjoxNjpbXSw1MjY6NDQ6W11dfSIsImV4cCI6MTc0Mzc5MjU0NSwiaWF0IjoxNzQyOTI4NTQ1LCJqdGkiOiI2OTJkZGU2OC02Y2Q0LTQ4ODgtYWRkZS02ODZjZDQ1ODg4NDAiLCJyb2wiOiIiLCJzc28iOiJkZjUyYzE2Ny0zM2QzLTQwNmItOGU4ZC0yZmJlNDEyMDY3YTcifQ.Zk55eohFo2iVw5gh_4KXAfMPKy4tGn5RIOlZQww8Y-jgHncjnRG3IyxUly8NExtL9rLgZNU96ZY7tTC6PHICjdv7fasG24nUAt7STcjtklBpVrkE0OMEny1kO_oS2fq89fd2U4VtC2D001S4foUXPN5vts4XbNbIXTTC-x1Gt3GKU3LBzc0oGDbofCMX01FG-aXM2BJBL6wprw3RAXgYK4sz4ckRvTlVzzL6HGt_tG5r_eRvWheUz0aGsvMt9BTVBZ9PPEGJ756cxZ0xnaAAD6hJy_BK1SpChTThbzLj5Yq7RarZlR8H0HAVD7vwfuB0nECW9oHONxb88CdAt-ROfQ"
	// GetRequest(url)

	makeAuthenticatedRequest(url, token)
}
