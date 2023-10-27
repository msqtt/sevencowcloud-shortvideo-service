package sample

import (
	"fmt"
	"math/rand"
	"time"
)

var genderMap = []string{"male", "female", "unknown"}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s.%s", RandomStr(rand.Intn(10)), RandomStr(10), RandomStr(5))
}

func RandomGender() string {
	return genderMap[rand.Intn(3)]
}

func RandomBirthDate() time.Time{
	days := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	months := []time.Month{time.January, time.February, time.March, time.April, time.May,
		time.June, time.July, time.August, time.September, time.October, time.November, time.December}

	year := int(RandomInt(1980, 2023))
	month := int(rand.Intn(12))
	day := 0	
	add := 0
	if month == 1 && (year % 4 == 0 && year % 100 != 0 || year % 400 == 0 ){
		add = 1
	}
	day = int(RandomInt(1, int64(days[month]+add)))

	return time.Date(year, months[month], day, 0, 0, 0, 0, time.Now().Location())
}
