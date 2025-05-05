package service

import (
	"fmt"
	"log"

	//"strings"
	"github.com/playwright-community/playwright-go"
	//"test/service"
)

// проходим аутентификацию и получаем токен в случае успеха
// !!!! логин и пароль нужно убрать в дальнейшем
func Authentication() string {
	// Установка драйвера (если ещё нет)
	// if err := playwright.Install(); err != nil {
	// 	log.Fatalf("Ошибка установки драйвера: %v", err)
	// }

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Ошибка запуска Playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("Ошибка запуска браузера: %v", err)
	}
	defer browser.Close()

	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("Ошибка создания контекста: %v", err)
	}
	defer context.Close()

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("Ошибка создания страницы: %v", err)
	}

	if _, err = page.Goto("https://school.mos.ru"); err != nil {
		log.Fatalf("Ошибка перехода: %v", err)
	}

	log.Println("Хром запущен")

	interButton := page.Locator("text=Войти")
	if err := interButton.Click(); err != nil {
		log.Fatal("Ошибка клика на странице регистрации в MЭШ:", err)
	}

	// Дожидаемся редиректа
	page.WaitForURL("https://login.mos.ru/sps/login/methods/**")

	// (логика входа и извлечения токена)
	// Заполнение логина и пароля
	loginInput := page.Locator("input[name='login']")
	if err := loginInput.Fill("Zelenovba"); err != nil {
		log.Fatal("Ошибка ввода логина на mos.ru:", err)
	}

	passwordInput := page.Locator("input[name='password']")
	if err := passwordInput.Fill("Zel220580"); err != nil {
		log.Fatal("Ошибка ввода пароля на mos.ru:", err)
	}

	// page.Click("button[type='submit']")
	submitButton := page.Locator("button[id='bind']")
	if err := submitButton.Click(); err != nil {
		log.Fatal("Ошибка клика по кнопке войти на mos.ru:", err)
	}

	log.Println("Регистрацию на mos.ru прошли")

	page.WaitForURL("https://school.mos.ru/diary/schedules/schedule/**")

	// Получение кук
	cookies, err := page.Context().Cookies()
	if err != nil {
		log.Fatal(err)
	}

	var token string

	for _, cookie := range cookies {
		if cookie.Name == "aupd_token" {
			token = cookie.Value
		}
	}

	fmt.Println("Токен:", token)

	return token
}
