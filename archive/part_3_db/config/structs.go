package config

// структура списка оценок
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

// структура периодов
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

// структура описывает данные полученные в JSON
type StudentPerformance struct {
	Payload []Payload `json:"payload"`
}

// структуры описывают личные данные студента
type Profile struct {
	ID              int    `json:"id"`
	Type            string `json:"type"`
	UserID          int    `json:"user_id"`
	SchoolID        int    `json:"school_id"`
	SchoolName      string `json:"school_name"`
	SchoolShortname string `json:"school_shortname"`
	OrganizationID  string `json:"organization_id"`
}

type Student struct {
	ID          int       `json:"id"`
	GUID        string    `json:"guid"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	MiddleName  string    `json:"middle_name"`
	DateOfBirth string    `json:"date_of_birth"`
	Sex         string    `json:"sex"`
	PersonID    string    `json:"person_id"`
	Profiles    []Profile `json:"profiles"`
}
