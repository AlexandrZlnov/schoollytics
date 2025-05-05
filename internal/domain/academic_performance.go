package domain

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
