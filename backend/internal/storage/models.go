package storage

type Pet struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Weight       float64 `json:"weight"`
	Age          int     `json:"age"`
	AgeWriteDate string  `json:"age_write_date"`
	Type         string  `json:"type"`
	Anamnesis    string  `json:"anamnesis"`
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
