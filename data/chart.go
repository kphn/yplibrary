package data

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/pkg/errors"
)

// GenChart .
func GenChart(t time.Time) error {

	csvFile := fileName(t)
	f, err := os.Open(csvFile)
	if err != nil {
		return err
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return errors.WithStack(err)
	}

	times := make([]string, len(records))
	totalNums := make([]float64, len(records))
	currentNums := make([]float64, len(records))
	for i, record := range records {

		if len(record) != 3 {
			continue
		}

		timeUnix, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		totalNum, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			continue
		}

		currentNum, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			continue
		}

		times[i] = time.Unix(timeUnix, 0).Format("01-02 3:04PM")
		totalNums[i] = totalNum
		currentNums[i] = currentNum
	}

	line := charts.NewLine()
	line.SetGlobalOptions(charts.TitleOpts{Title: "杨浦图书馆日流量"})
	line.AddXAxis(times).AddYAxis("日流总量", totalNums).
		AddYAxis("当前人数", currentNums)

	htmlFileName := t.Format("20060102.html")
	f, err = os.Create(htmlFileName)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	return errors.WithStack(line.Render(f))
}
