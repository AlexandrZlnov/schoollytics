package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/AlexandrZlnov/schoollytics/internal/domain"
)

/*
Блок функций которые сохраняют оцеки школьника в БД
вызывает вспомогательные функции:
- SaveSubject - добавляет предмет если его небыло
- SavePeriod - добавляет период(тримест) если его небыло
- SaveMark - сохраняет оценки
*/
func SaveGrades(studentID int, db *sql.DB, studentGrades *domain.StudentPerformance) error {
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
func SaveMark(db *sql.DB, studentID int, subjectID int, periodID int, mark *domain.Mark) error {
	_, err := db.Exec(`INSERT INTO grades (student_id, subject_id, period_id, external_id, value, weight, control_form_name, 
	date, original_grade_system_type)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		studentID, subjectID, periodID, mark.ID, mark.Value, mark.Weight, mark.ControlFormName,
		mark.Date, mark.OriginalGradeSystemType)

	return err
}

// сохраняет период обучения - триместр, если его там нет
// и возврщает его ID
func SavePeriod(db *sql.DB, period *domain.Period) (int, error) {
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
