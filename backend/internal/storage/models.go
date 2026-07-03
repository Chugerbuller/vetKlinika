package storage

type Pet struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Weight       float64 `json:"weight"`
	Age          int     `json:"age"`
	AgeWriteDate string  `json:"age_write_date"`
	Type         string  `json:"type"`
	Allergies    string  `json:"allergies"`
}
type Person struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
}
type Types struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}
type Appointment struct {
	Id              int    `json:"id"`
	PetId           int    `json:"pet_id"`
	Date            string `json:"date"`
	Anamnesis       string `json:"anamnesis"`
	Disease         string `json:"disease"`
	Treatment       string `json:"treatment"`
	Recommendations string `json:"recommendations"`
}
type Document struct {
	Id            int    `json:"id"`
	PetId         int    `json:"pet_id"`
	Date          string `json:"date"`
	AppointmentId int    `json:"appointment_id"`
	Name          string `json:"name"`
	Content       []byte `json:"content"`
}
