package actions

import (
	"testing"
	"time"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestGetTimeFromFile(t *testing.T) {
	cases := []struct {
		name     string
		file     string
		isErr    bool
		expected time.Time
	}{
		{
			"case1 OK",
			"../test/testdata/20230810_025005.txt",
			false,
			time.Date(2023, 8, 10, 2, 50, 5, 0, time.Now().Location()),
		},
		{
			"case2 OK",
			"../test/testdata/20230811_123456.txt",
			false,
			time.Date(2023, 8, 11, 12, 34, 56, 0, time.Now().Location()),
		},
		{
			"case3 NG",
			"../test/testdata/nonformatted.txt",
			true,
			time.Time{},
		},
		{
			"case4 NG",
			"../test/testdata/nonexisted.txt",
			true,
			time.Time{},
		},
		{
			"case5 NG",
			"",
			true,
			time.Time{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := getTimeFromFile(c.file)
			if err != nil {
				if !c.isErr {
					t.Fatal(err.Error())
				}
			}
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestGetDelay(t *testing.T) {
	cases := []struct {
		name     string
		sTime    time.Time
		cTime    time.Time
		expected time.Duration
	}{
		{
			"case1 delay 0",
			time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Duration(0),
		},
		{
			"case2 delay +1m1s",
			time.Date(2023, 1, 2, 3, 3, 4, 0, time.UTC),
			time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Duration(time.Minute + time.Second),
		},
		{
			"case3 delay -10m20s",
			time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Date(2023, 1, 2, 2, 53, 45, 0, time.UTC),
			time.Duration(10*time.Minute + 20*time.Second),
		},
		{
			"case4 delay +11h22m33s",
			time.Date(2023, 1, 1, 15, 41, 32, 0, time.UTC),
			time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Duration(11*time.Hour + 22*time.Minute + 33*time.Second),
		},
		{
			"case5 delay +3d4h5m6s",
			time.Date(2022, 12, 29, 22, 58, 59, 0, time.UTC),
			time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC),
			time.Duration(76*time.Hour + 5*time.Minute + 6*time.Second),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := getDelay(c.sTime, c.cTime)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestIntegration(t *testing.T) {
	cases := []struct {
		name           string
		opts           DelayOptions
		expectedStatus checkers.Status
		expectedString string
	}{
		{
			"case1 OK",
			DelayOptions{
				"../test/testdata/20230810_025005.txt",
				300,
				600,
				time.Date(2023, 8, 10, 2, 51, 23, 0, time.Now().Location()),
			},
			checkers.OK,
			"delay = 1m18s, serverTime = 2023/08/10 02:50:05, clientTime = 2023/08/10 02:51:23",
		},
		{
			"case2 WARNING",
			DelayOptions{
				"../test/testdata/20230810_025005.txt",
				120,
				240,
				time.Date(2023, 8, 10, 2, 52, 34, 0, time.Now().Location()),
			},
			checkers.WARNING,
			"delay = 2m29s, serverTime = 2023/08/10 02:50:05, clientTime = 2023/08/10 02:52:34",
		},
		{
			"case3 CRITICAL",
			DelayOptions{
				"../test/testdata/20230810_025005.txt",
				60,
				180,
				time.Date(2023, 8, 10, 2, 53, 43, 0, time.Now().Location()),
			},
			checkers.CRITICAL,
			"delay = 3m38s, serverTime = 2023/08/10 02:50:05, clientTime = 2023/08/10 02:53:43",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expectedStatus, c.opts.Check().Status)
			assert.Equal(t, c.expectedString, c.opts.Check().Message)
		})
	}
}
