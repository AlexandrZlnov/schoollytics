package service

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Структура для десериализации JSON
// type ApiResponse struct {
// 	ID int `json:"id"`
// }

// Get запрос по установленному URL для получения JSON
// сохраняется в сыром виде файл get_response.json
// и в структурированном через структуру ApiResponse
func GetRequest(url string) {
	//клиент с поддержкой cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Printf("Ошибка создания cookie jar: %v\n", err)
		return
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	// Делаем запрос с кастомными заголовками
	req, _ := http.NewRequest("GET", "https://stats.mos.ru/handler/handler.js?time=1743089604714", nil)
	//req.Header.Add("Authorization", "Bearer token")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Сервер вернул ошибку: %s\n", resp.Status)
		return
	}

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Printf("Ошибка при чтении ответа: %v\n", err)
	// 	return
	// }

	// записывает JSON в файл
	// if err = os.WriteFile("get_response.json", body, 0644); err != nil {
	// 	fmt.Printf("Ошибка при сохранении файла: %v\n", err)
	// 	return
	// }

	// вариант 1 десериализации JSON в заранее известную структуру
	// для этого варианта нужно включить структуру ApiResponse в начале
	// var result ApiResponse
	// if err = json.Unmarshal(body, &result); err != nil {
	// 	fmt.Printf("Ошибка десериализации JSON по Get запросу: %v\n", err)
	// 	return

	// вариант 2 десериализации JSON с неизвестными параметрами
	// выдает ошибку если в начале !
	// var result map[string]interface{}
	// var result []interface{}
	// if err = json.Unmarshal(body, &result); err != nil {
	// 	fmt.Printf("Ошибка десериализации JSON в структуру по Get запросу: %v\n", err)
	// 	return
	// }

	cookies := resp.Cookies()
	//fmt.Println(result)
	fmt.Println(cookies)

	// Анализируем полученные cookies
	for _, cookie := range jar.Cookies(req.URL) {
		fmt.Printf("Cookie: %s=%s (Domain: %s, Path: %s)\n",
			cookie.Name, cookie.Value, cookie.Domain, cookie.Path)
	}

}
