package tempo

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	CapitalizableTime = "How much time did you spend developing, designing or testing software (in hours): "
	PtoTime           = "How much time did you spend with PTO (vacation or sick) (in hours): "
	MeetingTime       = "How much time did you spend in meetings (in hours): "
)

func requestTimeInput() (devTime, ptoTime, meetingTime int) {
	fmt.Printf("Answer the following questions to the best of your ability and estimate how you spent your time this week.\n")
	fmt.Printf("We will ask about 3 things: development, PTO, and meetings.\n")
	fmt.Printf("(For the moment, this cannot exceed a total of 40 hours)\n\n")

	devTime = getTime(CapitalizableTime)
	ptoTime = getTime(PtoTime)
	meetingTime = getTime(MeetingTime)

	totalHoursThisWeek := meetingTime + ptoTime + devTime

	print(fmt.Sprintf("Total hours this week: %s\n", strconv.Itoa(totalHoursThisWeek)))
	return
}

func getTime(printString string) int {
	fmt.Print(printString)
	var timeInput string
	if _, err := fmt.Scan(&timeInput); err != nil {
		log.Fatal(err)
	}
	return stringToInt(timeInput)
}

func stringToInt(input string) int {
	var convertedValue, err = strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	return convertedValue
}

func requestDayOfWeek() time.Time {
	mondayOfThisWeek := determineWeekforTimeSheet()

	fmt.Printf("Would you like to fill out time for %s (Y/N)? ", mondayOfThisWeek.Format(time.DateOnly))
	var confirmTime string
	if _, err := fmt.Scan(&confirmTime); err != nil {
		log.Fatal(err)
	}

	if confirmTime == "y" || confirmTime == "Y" {
		return mondayOfThisWeek
	}

	fmt.Printf("\nHow many weeks back would you like to fill out (ex. 1 means last week): ")
	var timeInput string
	if _, err := fmt.Scan(&timeInput); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Now we are filling out a timesheet for %s\n", mondayOfThisWeek.AddDate(0, 0, -7*stringToInt(timeInput)).Format(time.DateOnly))

	return mondayOfThisWeek.AddDate(0, 0, -7*stringToInt(timeInput))
}

func determineWeekforTimeSheet() time.Time {
	currentDay := time.Now()
	dayOfTheWeek := int(currentDay.Weekday())
	var distanceToMonday int
	if dayOfTheWeek-1 == -1 {
		distanceToMonday = -6
	} else {
		distanceToMonday = -(dayOfTheWeek - 1)
	}

	monday := currentDay.AddDate(0, 0, distanceToMonday)
	print(fmt.Sprintf("This will fill out the timesheet for the week of %s\n\n", monday.Format(time.DateOnly)))
	return monday
}
