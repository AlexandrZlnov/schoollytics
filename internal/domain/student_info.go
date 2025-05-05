package domain

type StudentInfo struct {
	UserID       int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MiddleName   string `json:"Middle_name"`
	PphoneNumber string `json:"phone_number"`
	RegionalAuth string `json:"regional_auth"`
	DateOfBirth  string `json:"date_of_birth"`
	Sex          string `json:"sex"`
	Profiles     []struct {
		StudentID       int    `json:"id"`
		Type            string `json:"type"`
		SchoolId        int    `json:"school_id"`
		SchoolShortname string `json:"school_shortname"`
		SchoolName      string `json:"school_name"`
		OrganizationID  string `json:"organization_id"`
	} `json:"profiles"`
	Token string `json:"authentication_token"`
}

// Дублирующая структура
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
