/*
Часть 2
Парсинг сырого json файла и приведение к нужной структуре
получение данных из JSON и печать данных с оценками по предмету
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Mark struct {
	ID                      int64  `json:"id"`
	Value                   string `json:"value"`
	Values                  any    `json:"values"`
	Comment                 string `json:"comment"`
	Weight                  int    `json:"weight"`
	PointDate               any    `json:"point_date"`
	ControlFormName         string `json:"control_form_name"`
	CommentExists           bool   `json:"comment_exists"`
	CreatedAt               any    `json:"created_at"`
	UpdatedAt               any    `json:"updated_at"`
	Criteria                any    `json:"criteria"`
	Date                    string `json:"date"`
	IsPoint                 bool   `json:"is_point"`
	IsExam                  bool   `json:"is_exam"`
	OriginalGradeSystemType string `json:"original_grade_system_type"`
}

type Period struct {
	Start      string `json:"start"`
	End        string `json:"end"`
	Title      string `json:"title"`
	Dynamic    string `json:"dynamic"`
	Value      string `json:"value"`
	Marks      []Mark `json:"marks"`
	Count      int    `json:"count"`
	Target     any    `json:"target"`
	FixedValue string `json:"fixed_value"`
	StartISO   string `json:"start_iso"`
	EndISO     string `json:"end_iso"`
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
	var studentData StudentPerformance

	makeStructureJSON(&studentData)

	currentPeriod := calculatePeriod()
	if currentPeriod == 0 {
		log.Fatal("ошибка определения текущего учебного периода")
	}
	fmt.Println("Текущий учебный период: ", currentPeriod)

	printMarks(currentPeriod, &studentData)

}

// Функция берет сырой JSON и формирует структурированный в новый файл
func makeStructureJSON(studentData *StudentPerformance) {

	//читаем сырой JSON
	rawData, err := os.ReadFile("studentPerformance.json")
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	//десериализация сырого JSON в структуру - парсинг JSON
	if err := json.Unmarshal(rawData, &studentData); err != nil {
		log.Fatal("ошибка парсинга JSON файла:", err)
	}

	//преобразует сурую структуру с данными из JSON в упорядоченный список с отступами и переносами
	prettyJSON, err := json.MarshalIndent(studentData, "", "   ")
	if err != nil {
		log.Fatal("ошибка форматированя JSON данных:", err)
	}

	//записываем упорядоченный список данных из JSON в файл
	if err := os.WriteFile("prettyStudentPerformance.json", prettyJSON, 0644); err != nil {
		log.Fatal("ошибка записи в файл отформатированных JSON данных:", err)
	}

	prettyJSON, err = os.ReadFile("prettyStudentPerformance.json")
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	if err := json.Unmarshal(prettyJSON, &studentData); err != nil {
		log.Fatal("ошибка парсинга JSON файла:", err)
	}
}

// функция определяет текущий учебный период по текущей дате
// если текущая дата выпала на период каникул
// текущим учебным периодом будет принят предшествующий каникулам
func calculatePeriod() int {
	//now := time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC)
	//now := time.Date(2023, time.December, 15, 0, 0, 0, 0, time.UTC)
	//now := time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)
	// now := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	//now := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)

	now := time.Now()
	year := now.Year()
	var holidays = true

	var Periods = []struct {
		TrimNumber int
		Start      time.Time
		End        time.Time
		CrossYear  bool
	}{
		{
			1,
			time.Date(year, time.September, 01, 0, 0, 0, 0, time.UTC),
			time.Date(year, time.October, 25, 0, 0, 0, 0, time.UTC),
			false,
		},
		{
			2,
			time.Date(year, time.November, 05, 0, 0, 0, 0, time.UTC),
			time.Date(year, time.February, 14, 0, 0, 0, 0, time.UTC),
			true,
		},
		{
			3,
			time.Date(year, time.February, 02, 0, 0, 0, 0, time.UTC),
			time.Date(year, time.May, 30, 0, 0, 0, 0, time.UTC),
			false,
		},
	}

	for _, period := range Periods {
		switch period.CrossYear {
		case true:
			if now.After(period.Start) || now.Equal(period.Start) ||
				now.Before(period.End) || now.Equal(period.End) {
				//fmt.Println("Period number: ", period.TrimNumber)
				holidays = false
				return period.TrimNumber
			}
		case false:
			if (now.After(period.Start) || now.Equal(period.Start)) &&
				(now.Before(period.End) || now.Equal(period.End)) {
				//fmt.Println("Period number: ", period.TrimNumber)
				holidays = false
				return period.TrimNumber
			}

		}
	}

	// обработка периодов выпавших на каникулы
	// период будет равен предшествующему текущим каникулам
	if holidays {
		if now.After(Periods[2].End) && now.Before(Periods[0].Start) {
			return Periods[2].TrimNumber
		} else if now.After(Periods[0].End) && now.Before(Periods[1].Start) {
			return Periods[0].TrimNumber
		} else if now.After(Periods[1].End) && now.Before(Periods[2].Start) {
			return Periods[1].TrimNumber
		}

	}
	return 0
}

func printMarks(currentPeriod int, studentData *StudentPerformance) {
	count := 0
	for j := 0; j < currentPeriod; j++ {
		fmt.Printf("Оценки - %s, Триместр -  %d.\n", studentData.Payload[7].SubjectName, j+1)
		count = 0
		for i := 0; i < studentData.Payload[7].Periods[j].Count; i++ {
			//fmt.Print(studentData.Payload[7].Periods[j].Count, " ")
			if count%10 == 0 {
				fmt.Printf("\n")
			}
			fmt.Printf("%v  ", studentData.Payload[7].Periods[j].Marks[i].Value)
			count++
		}
		fmt.Printf("\n\n")
	}
}
