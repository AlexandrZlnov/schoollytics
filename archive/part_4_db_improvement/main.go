/*
Часть 4
Улучшеная работа с БД
Доработаны функции добавления данных о Школьнике Школе Предметах Периодах и Успеваемости
Доработаны структуры данных ученика и успеваемости
*/

package main

import (
	"database/sql"
	//"fmt"

	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// структура Школа
type Profiles struct {
	SchoolID       int    `json:"school_id"`
	Name           string `json:"school_name"`
	Shortname      string `json:"school_shortname"`
	OrganizationID string `json:"organization_id"`
	ExternalID     int    `json:"id"`
}

// структура Студент
type Students struct {
	UserID              int        `json:"id"`
	Profile             []Profiles `json:"profiles"`
	GuID                string     `json:"guid"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	MiddleName          string     `json:"middle_name"`
	PhoneNumber         string     `json:"phone_number"`
	AuthenticationToken string     `json:"authentication_token"`
	PersonID            string     `json:"person_id"`
	PswrdChangeRequired bool       `json:"password_change_required"`
	RegionalAuth        string     `json:"regional_auth"`
	DateOfBirht         string     `json:"date_of_birth"`
	Sex                 string     `json:"sex"`
}

// структура Успеваемость школьника - оценки
type Mark struct {
	ID                      int64  `json:"id"`                         // id оценки - 2447848392
	Value                   string `json:"value"`                      // оценка
	Values                  any    `json:"values"`                     // ????????
	Comment                 string `json:"comment"`                    // ????????
	Weight                  int    `json:"weight"`                     // вес (коэффициент) оценки
	PointDate               any    `json:"point_date"`                 // ???????
	ControlFormName         string `json:"control_form_name"`          // форма оцениваемой работы - Творческая работа
	CommentExists           bool   `json:"comment_exists"`             // ??????
	CreatedAt               any    `json:"created_at"`                 // ??????
	UpdatedAt               any    `json:"updated_at"`                 // ??????
	Criteria                any    `json:"criteria"`                   // ??????
	Date                    string `json:"date"`                       // дата оценки
	IsPoint                 bool   `json:"is_point"`                   // ??????
	IsExam                  bool   `json:"is_exam"`                    // ??????
	OriginalGradeSystemType string `json:"original_grade_system_type"` // ???? вероятно тип системы оценок к примеру пятибальная - "five"
}

type Period struct {
	Start      string `json:"start"`       // начало учебного периода
	End        string `json:"end"`         // конец учебного периода
	Title      string `json:"title"`       // номер триместра
	Dynamic    string `json:"dynamic"`     // динамика - NONE
	Value      string `json:"value"`       // оценка - 5.00
	Marks      []Mark `json:"marks"`       // слайс оценок
	Count      int    `json:"count"`       // количество оценок в этом триместре
	Target     any    `json:"target"`      // не понятно, возможно цель - null
	FixedValue string `json:"fixed_value"` // не понятно, возможно оценка за тримест - 5
	StartISO   string `json:"start_iso"`   // дата начала триместра
	EndISO     string `json:"end_iso"`     // дата окончания триместра
}

type Payload struct {
	Average      string   `json:"average"`
	Dynamic      string   `json:"dynamic"`
	Periods      []Period `json:"periods"`
	SubjectName  string   `json:"subject_name"`
	SubjectID    int      `json:"subject_id"`
	AverageByAll string   `json:"average_by_all"`
	YearMark     any      `json:"year_mark"`
}

type StudentPerformance struct {
	Payload []Payload `json:"payload"`
}

func main() {
	log.Println("=== ПРОГРАММА ЗАПУЩЕНА ===")
	// инициируем базу данных
	db, err := InitDB()
	if err != nil {
		log.Fatalf("ошибка подключения к PostgreSql: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

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
	var school Students
	if err := json.Unmarshal(jsonDataStudent, &school); err != nil {
		log.Fatalf("ошибка десериализации json: %v", err)
	}

	// десериализуем JSON с данными оценок  в структуру
	var studentGrades StudentPerformance
	err = json.Unmarshal(jsonDataMarks, &studentGrades)
	if err != nil {
		log.Fatalf("ошибка десериализации json: %v", err)
	}

	// проверяем наличие школы в БД
	var schoolID int
	schoolID, err = checkSchool(db, school)
	if err != nil {
		log.Fatalf("ошибка проверки наличих школы в БД: %v", err)
	}

	var studentID int
	// сохраняем данные по ученику в БД
	// и получаем его ID
	studentID, err = SaveStudent(db, school, &schoolID)
	if err != nil {
		log.Fatalf("ошибка сохранения школьника: %v", err)
	}

	// БЛОК ВСТАВКИ ДАННЫХ ОБ ОЦЕНКАХ В БД ------------------->
	err = SaveGrades(studentID, db, &studentGrades)
	if err != nil {
		log.Fatal("ошибка сохранения оценок в БД: %w", err)
	}
}

// функция сохраняет оцеки школьника
// вызывает вспомогательные функции:
// SaveSubject - проверяет предмет на вкллюченность в таблицу
// если нет сохраняет новый и возвращает его ID
// если предмет уже есть - возвращает его ID
func SaveGrades(studentID int, db *sql.DB, studentGrades *StudentPerformance) error {
	for _, payload := range studentGrades.Payload {
		// добавим предмет если его еще нет и получим его ID
		// тут и далее payload.SubjectID - внешний subjectID полученный от МЭШ
		// subjectID - значение из нашей базы данных
		subjectID, err := SaveSubject(db, payload.SubjectID, payload.SubjectName)
		if err != nil {
			log.Fatalf("ошибка сохранения предмета: %s, в БД: %v", payload.SubjectName, err)
			continue
		}
		log.Printf("ID %s - %d", payload.SubjectName, subjectID)

		// добавим периоды
		for _, period := range payload.Periods {
			periodID, err := SavePeriod(db, &period)
			if err != nil {
				log.Fatalf("ошибка сохранения периода: %s, в БД: %v", payload.SubjectName, err)
				continue
			}
			log.Printf("Период ID - %d c %s по %s", periodID, period.Start, period.End)

			// добавим оценки
			for _, mark := range period.Marks {
				err := SaveMark(db, studentID, subjectID, periodID, &mark)
				if err != nil {
					log.Fatalf("ошибка сохранения оценки в БД: Студент - %d, Предмет - %d, Оценка - %s, ID оценки - %d. Ошибка: %v",
						studentID, subjectID, mark.Value, mark.ID, err)
					continue
				}
			}

		}

	}

	return nil
}

// сохраняет оценки
func SaveMark(db *sql.DB, studentID int, subjectID int, periodID int, mark *Mark) error {
	_, err := db.Exec(`INSERT INTO grades (student_id, subject_id, period_id, external_id, value, weight, control_form_name, 
	date, original_grade_system_type)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		studentID, subjectID, periodID, mark.ID, mark.Value, mark.Weight, mark.ControlFormName,
		mark.Date, mark.OriginalGradeSystemType)

	return err
}

// сохраняет период обучения - триместр, если его там нет
// и возврщает его ID
func SavePeriod(db *sql.DB, period *Period) (int, error) {
	var periodID int
	// проверка наличия периода в БД

	// добавил targetJSON типа []byte т.к. в некоторых target приходит
	// переменная типа map[string]interface {} и возникает ошибка при
	// попытке добавить в БД с полем типа jsonb
	targetJSON, err := json.Marshal(period.Target)
	if err != nil {
		return 0, fmt.Errorf("ошибка сериализации данных из Target: %v", err)
	}

	err = db.QueryRow(`SELECT id FROM periods WHERE start_date = $1 AND end_date = $2`,
		period.Start, period.End).Scan(&periodID)
	if err != nil || err == sql.ErrNoRows {
		log.Printf("Периода с таким началом: %s и концом: %s нет в базе данных, нужно добавить", period.Start, period.End)
	} else {
		log.Printf("Периода с таким началом: %s и концом: %s уже существует", period.Start, period.End)
		return periodID, nil
	}

	// Если периода нет, вставляем его
	err = db.QueryRow(`INSERT INTO periods (start_date, end_date, title, dynamic, value, 
	count, target, fixed_value, start_iso, end_iso)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`,
		period.Start,
		period.End,
		period.Title,
		period.Dynamic,
		period.Value,
		period.Count,
		targetJSON,
		period.FixedValue,
		period.StartISO,
		period.EndISO).Scan(&periodID)
	if err != nil {
		return 0, fmt.Errorf("oшибка добавления нового периода в таблицу: %w", err)
	}

	return periodID, nil

}

// сохраняет предмет в базу данных если его там нет
// возвращает ID предмета
func SaveSubject(db *sql.DB, subjectID int, subjectName string) (int, error) {
	var existingID int
	// проверка наличия предмета в БД
	err := db.QueryRow(`SELECT id FROM subjects WHERE external_id = $1`, subjectID).Scan(&existingID)
	if err != nil {
		log.Printf("Предмета под названием %s нет в базе данных, нужно добавить", subjectName)
	} else {
		log.Printf("Предмет %s уже существует его ID - %d", subjectName, existingID)
		return existingID, nil
	}

	// Если предмета нет, вставляем его
	err = db.QueryRow(`INSERT INTO subjects (name, external_id)
	VALUES ($1, $2)
	ON CONFLICT (external_id) DO NOTHING
	RETURNING id`,
		subjectName,
		subjectID).Scan(&existingID)
	if err != nil {
		return 0, fmt.Errorf("oшибка добавления нового предмета в таблицу: %w", err)
	}

	return existingID, nil
}

// содзает подключение к базе данных
// возвращает *sql.DB или error
func InitDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=Zel408 dbname=test4 sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, fmt.Errorf("ошибка подключения к PostgreSql: %v", err)
	}

	if err := db.Ping(); err != nil {
		return db, fmt.Errorf("ошибка Ping DB: %v", err)
	}

	return db, nil
}

// проверяет наличие школы в таблице schools
// возвращает ID школы в случа ее наличия
func checkSchool(db *sql.DB, school Students) (int, error) {
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

// добавляет новую школу или
// вносит изменения в уже существующую
func SaveSchool(db *sql.DB, school Students) (int, error) {
	log.Printf("Пытаемся сохранить школу: ID=%d, Name=%s, OrgID=%s",
		school.Profile[0].SchoolID,
		school.Profile[0].Name,
		school.Profile[0].OrganizationID)

	var schoolID int

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer tx.Rollback()
	log.Println("---------> 0")
	err = tx.QueryRow(`
        INSERT INTO schools (school_id, name, shortname, organization_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (organization_id) DO UPDATE SET
            school_id = EXCLUDED.school_id,
            name = EXCLUDED.name,
            shortname = EXCLUDED.shortname 
			RETURNING id`,
		school.Profile[0].SchoolID,
		school.Profile[0].Name,
		school.Profile[0].Shortname,
		school.Profile[0].OrganizationID,
	).Scan(&schoolID)

	fmt.Println(schoolID)
	log.Println("---------> 1")

	if err != nil {
		return 0, fmt.Errorf("ошибка при вставке/обновлении школы: %v", err)
	}
	log.Println("---------> 2")
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("не удалось зафиксировать транзакцию: %v", err)
	}
	log.Println("---------> 3")
	log.Println("Данные школы успешно обновлены")
	return schoolID, nil
}

// добавляет нового школьника
// данные школьника обновятся даже если он уже есть в таблице students
// обновляются: телефон, регион, дата рождения и школа (в случае изменения ее ID)
func SaveStudent(db *sql.DB, student Students, schoolID *int) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer tx.Rollback()

	var studentID int

	err = tx.QueryRow(`
        INSERT INTO students (user_id, profile_id, guid, first_name, last_name, middle_name, phone_number, 
authentication_token, person_id, pswrd_change_required, regional_auth, date_of_birth, sex, school_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		  ON CONFLICT (user_id) DO UPDATE SET
            regional_auth = EXCLUDED.regional_auth,
            date_of_birth = EXCLUDED.date_of_birth,
            phone_number = EXCLUDED.phone_number,
			school_id = $14
			RETURNING id`,
		student.UserID,
		student.Profile[0].ExternalID,
		student.GuID,
		student.FirstName,
		student.LastName,
		student.MiddleName,
		student.PhoneNumber,
		student.AuthenticationToken,
		student.PersonID,
		student.PswrdChangeRequired,
		student.RegionalAuth,
		student.DateOfBirht,
		student.Sex,
		schoolID,
	).Scan(&studentID)

	if err != nil {
		return 0, fmt.Errorf("ошибка при вставке/обновлении школьника: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("не удалось зафиксировать транзакцию: %v", err)
	}

	log.Println("Данные школьника успешно обновлены")
	return studentID, nil
}
