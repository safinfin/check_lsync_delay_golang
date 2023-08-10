package actions

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/mackerelio/checkers"
)

type DelayOptions struct {
	File     string
	Warning  int
	Critical int
	Time     time.Time
}

const (
	fileTimeFormat   string = "20060102_150405"
	outputTimeFormat string = "2006/01/02 15:04:05"
)

func getTimeFromFile(file string) (time.Time, error) {
	f, err := os.Open(file)
	if err != nil {
		return time.Time{}, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	textData := scanner.Text()

	fileTime, err := time.ParseInLocation(fileTimeFormat, textData, time.Now().Location())
	if err != nil {
		return time.Time{}, err
	}

	return fileTime, nil
}

func getDelay(sTime, cTime time.Time) time.Duration {
	delay := cTime.Sub(sTime)
	absDelay := time.Duration(math.Abs(float64(delay)))
	return absDelay.Truncate(time.Second)
}

func (o *DelayOptions) Check() *checkers.Checker {
	checkStatus := checkers.OK
	sTime, err := getTimeFromFile(o.File)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	cTime := o.Time
	delay := getDelay(sTime, cTime)

	warning := time.Duration(o.Warning) * time.Second
	critical := time.Duration(o.Critical) * time.Second

	if warning != 0 && warning <= delay {
		checkStatus = checkers.WARNING
	}
	if critical != 0 && critical <= delay {
		checkStatus = checkers.CRITICAL
	}

	msg := fmt.Sprintf("delay = %s, serverTime = %s, clientTime = %s", delay, sTime.Format(outputTimeFormat), cTime.Format(outputTimeFormat))

	return checkers.NewChecker(checkStatus, msg)
}
