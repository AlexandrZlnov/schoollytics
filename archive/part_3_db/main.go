/*
Часть 3
Работа с базой данных
Инициация БД с загрузкой конфигурации из ENV
Проверка наличия таблиц
Создание таблиц
Сохранение данных
Сохранение данных о студенте
Миграции
*/
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"test3/config"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// структура с данными из файла .env
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
}

type Profile struct {
	ID              int    `json:"id"`
	Type            string `json:"type"`
	SchoolID        int    `json:"school_id"`
	SchoolName      string `json:"school_name"`
	SchoolShortname string `json:"school_shortname"`
}

type Student struct {
	ID          int       `json:"id"`
	GUID        string    `json:"guid"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	MiddleName  string    `json:"middle_name"`
	DateOfBirth string    `json:"date_of_birth"` // или используйте time.Time
	Sex         string    `json:"sex"`
	PersonID    string    `json:"person_id"`
	Profiles    []Profile `json:"profiles"`
	// Другие поля по необходимости
}

func main() {
	// загружаем конфигурацию из .env
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	// создаем подключение к базе данных с учетом конфигурационных данных
	// host port user password dbname sslmode
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("could not ping database: %v", err)
	}
	defer db.Close()

	// применяем миграции
	// ######### добавил для эксперимента В дальнейшем можно убрать
	if err := runMigrations(db); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
	log.Println("Миграция применена")

	// проверка наличия таблиц
	if err := createDBTables(db); err != nil {
		log.Fatalf("Ошибка при работе с таблицей students: %v", err)
	}

	// читаем prettyJSON файл
	// jsonData, err := os.ReadFile("prettyStudentPerformance.json")
	// if err != nil {
	// 	log.Fatalf("Ошибка чтения prettyJSON file: %v", err)
	// }
	jsonData, err := os.ReadFile("student_info.json")
	if err != nil {
		log.Fatalf("Ошибка чтения prettyJSON file: %v", err)
	}

	var student config.Student

	if err := json.Unmarshal(jsonData, &student); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Println("----->>>>> Данные студента в структуре ------>>", student)
	// работа с данными оценок
	// var payload config.StudentPerformance
	// if err := json.Unmarshal(jsonData, &payload); err != nil {
	// 	log.Fatalf("Failed to parse JSON: %v", err)
	// }

	// if err := saveData(db, payload, 1, "4Ж"); err != nil {
	// 	log.Fatalf("ошибка функции saveData: %v", err)
	// }

	log.Println("Дошли до SaveStudent")
	SaveStudentFromJSONFile(db, &student, "student_info.json")

}

// функция сохраняет школу
func SaveStudentFromJSONFile(db *sql.DB, student *config.Student, filename string) error {

	// Начало транзакции
	// tx, err := db.Begin()
	// if err != nil {
	// 	return fmt.Errorf("could not begin transaction: %v", err)
	// }
	// defer tx.Rollback()

	// // 1. Сохранение школы (если есть профили)
	// var schoolID int
	// profile := student.Profiles[0]
	// if len(student.Profiles) > 0 {
	// 	log.Println("Воши в условие len>0")
	// 	err = tx.QueryRow(`
	// 	INSERT INTO schools (organization_id, name, short_name)
	// 	VALUES ($1, $2, $3)
	// 	RETURNING id`,
	// 		profile.OrganizationID, profile.SchoolName, profile.SchoolShortname,
	// 	).Scan(&schoolID)
	// } else {
	// 	log.Println("Воши в условие else")
	// 	// Если школа найдена, обновляем ее данные
	// 	_, err = tx.Exec(`
	// 	UPDATE schools
	// 	SET name = $1, short_name = $2
	// 	WHERE organization_id = $3`,
	// 		profile.SchoolName, profile.SchoolShortname, profile.OrganizationID,
	// 	)
	// }

	// if err != nil {
	// 	return fmt.Errorf("error saving school: %v", err)
	// }

	// // Завершение транзакции
	// err = tx.Commit()
	// if err != nil {
	// 	return fmt.Errorf("could not commit transaction: %v", err)
	// }

	// return nil
	// Проверяем, есть ли профили у студента
	if len(student.Profiles) == 0 {
		return fmt.Errorf("student has no profiles")
	}

	profile := student.Profiles[0] // Берём первый профиль

	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback() // Откатываем, если будет ошибка

	// Используем UPSERT (вставка или обновление)
	_, err = tx.Exec(`
        INSERT INTO schools (organization_id, name, short_name)
        VALUES ($1, $2, $3)
        ON CONFLICT (organization_id) DO UPDATE SET
            name = EXCLUDED.name,
            short_name = EXCLUDED.short_name`,
		profile.OrganizationID, profile.SchoolName, profile.SchoolShortname)

	if err != nil {
		return fmt.Errorf("error saving school: %v", err)
	}

	// Завершаем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	log.Printf("School saved successfully: %s (org_id: %s)",
		profile.SchoolName, profile.OrganizationID)
	return nil

}

// загружает .env файл и возвращает заполненную структуру
// с данными из .env или дефолтными если файли или ключ отсутствует
func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки .env: %w", err)
	}
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "8080"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
		SSLMode:    getEnv("SSL_MODE", "disable"),
	}, nil
}

// вспомогательная функция к func loadConfig()
// присваивает значение из .env или дефолнтное
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// инициируем базу данных
func initDB(cfg *Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)

	fmt.Println("connStr in initDB ----> ", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к PostgreSql: %w", err)
	}

	// Настройки соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверка подключения к базе данных
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping DB: %w", err)
	}

	return db, nil
}

// применение миграция к БД
// ######### добавил для эксперимента В дальнейшем можно убрать
func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("создание драйвера миграций: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("инициализация миграций: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("применение миграций: %w", err)
	}

	return nil
}

// функция проверит наличие таблицы students в БД
// SELECT EXISTS всегда возвращает true & false поэтому в запросе db.QueryRow никогда не будет
// нелевого значения и ErrNoRows не возникнет
// Scan(&tableExist) перезаписывает true & false по результатам работы db.QueryRow
func createDBTables(db *sql.DB) error {
	tableExist := false

	log.Println("Начинаем проверку наличия таблицы в БД")
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'students'
		)
	`).Scan(&tableExist)
	log.Printf("Результат: tableExist=%v, err=%v", tableExist, err)

	if err != nil {
		return fmt.Errorf("ошибка проверки существования таблицы: %w", err)
	}

	if tableExist {
		log.Println("Таблица students уже существует")
		return nil
	}

	log.Println("Таблица students не создана")

	return nil

}

// func saveData(db *sql.DB, payload config.StudentPerformance, schoolID int, className string) error {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return fmt.Errorf("failed to begin transaction: %v", err)
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()
// 	log.Println("1-Создана тракзакция ")

// 	// -------------->>>>>>>>!!!!!!!!!!   Для примера - создаем одного студента
// 	studentID, err := insertStudent(tx, schoolID, className, "Иванов Иван Иванович", "12345")
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("Присвоены данные студента:", studentID)

// 	log.Println("4-приступаем range по payload.Payload")
// 	for _, subject := range payload.Payload {
// 		// Вставляем предмет
// 		if err := insertSubject(tx, subject.SubjectID, subject.SubjectName); err != nil {
// 			return err
// 		}

// 	}
// 	return nil
// }

// Вспомогательнаяфункция для вставки данных студента
// func insertStudent(tx *sql.Tx, schoolID int, className, fullName, externalID string) (int, error) {
// 	log.Println("2-перешлли в функцию insertStudent")
// 	var id int
// 	err := tx.QueryRow(`
// 		INSERT INTO students (school_id, class_name, full_name, external_id)
// 		VALUES ($1, $2, $3, $4)
// 		ON CONFLICT (external_id) DO UPDATE SET
// 			full_name = EXCLUDED.full_name,
// 			class_name = EXCLUDED.class_name
// 		RETURNING id`,
// 		schoolID, className, fullName, externalID).Scan(&id)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to insert student: %v", err)
// 	}
// 	log.Println("3-завершен INSERT Student")
// 	return id, nil
// }

// func insertSubject(tx *sql.Tx, subjectID int, subjectName string) error {
// 	_, err := tx.Exec(`
// 		INSERT INTO subjects (id, name)
// 		VALUES ($1, $2)
// 		ON CONFLICT (id) DO UPDATE SET
// 			name = EXCLUDED.name`,
// 		subjectID, subjectName)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert subject %d: %v", subjectID, err)
// 	}
// 	return nil
// }
