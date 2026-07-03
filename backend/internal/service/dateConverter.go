package service

import "time"

const layout = "20060102"
const fromLayout = "2006-01-02T15:04:05.000Z"

func ConvertDate(dateStr string) (string, error) {
	date, err := time.Parse(fromLayout, dateStr)
	if err != nil {
		return "", err
	}
	result := date.Format(layout)

	return result, nil
}
func CalcAge(dateStart string, age int) (int, error) {
	now := time.Now()
	date, err := time.Parse(layout, dateStart)
	if err != nil {
		return 0, err
	}
	diff := now.Year() - date.Year()

	return age + diff, nil

}
