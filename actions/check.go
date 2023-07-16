package actions

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/mackerelio/checkers"
)

const (
	fileTimeFormat   = "20060102_150405"
	resultTimeFormat = "2006/01/02 15:04:05"
)

func readDateFromFile(file string) (time.Time, error) {
	f, err := os.Open(file)
	if err != nil {
		return time.Time{}, err
	}

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	date := scanner.Text()

	serverTime, err := time.ParseInLocation(fileTimeFormat, date, time.Now().Location())
	if err != nil {
		return time.Time{}, err
	}
	return serverTime, nil
}

func checkTimeDiff(serverTime, clientTime time.Time) time.Duration {
	timeDiff := serverTime.Sub(clientTime)
	absTimeDiff := math.Abs(float64(timeDiff))
	absDiffDuration := time.Duration(absTimeDiff).Truncate(time.Second)
	return absDiffDuration
}

func checkDelay(diff time.Duration, warning int64, critical int64, serverTime time.Time, clientTime time.Time) *checkers.Checker {
	warn := time.Duration(warning) * time.Second
	crit := time.Duration(critical) * time.Second
	checkStatus := checkers.UNKNOWN

	if warn != 0 && diff < warn {
		checkStatus = checkers.OK
	}
	if warn == 0 && diff < crit {
		checkStatus = checkers.OK
	}
	if warn != 0 && diff >= warn {
		checkStatus = checkers.WARNING
	}
	if crit != 0 && diff >= crit {
		checkStatus = checkers.CRITICAL
	}

	msg := fmt.Sprintf("diff = %s, serverTime = %s, clientTime = %s", diff, serverTime.Format(resultTimeFormat), clientTime.Format(resultTimeFormat))

	return checkers.NewChecker(checkStatus, msg)
}

func Run(file string, warning int64, critical int64) *checkers.Checker {
	serverTime, err := readDateFromFile(file)
	if err != nil {
		return checkers.Unknown(err.Error())
	}
	clientTime := time.Now()

	diff := checkTimeDiff(serverTime, clientTime)

	result := checkDelay(diff, warning, critical, serverTime, clientTime)

	return result
}
