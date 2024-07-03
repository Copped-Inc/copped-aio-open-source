package traiding

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func predictIntra(d map[string]Series, t time.Time, symbol string) (float64, float64, string, error) {

	current := t.Format("2006-01-02")
	o := d[current].open()

	parse := getDate(current)
	promptArray := []float64{
		o / 500,
		getLast(&parse, d),
		getIntraDay(parse, d),
		getLast(&parse, d),
		getIntraDay(parse, d),
		getLast(&parse, d),
		getIntraDay(parse, d),
		getLast(&parse, d),
		getIntraDay(parse, d),
		getLast(&parse, d),
		getIntraDay(parse, d),
		(float64(t.Year()) - 2010) / 13,
		float64(t.Month()) / 12,
		float64(t.Day()) / float64(time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()),
		float64(t.Weekday()) / 5,
	}

	prompt := fmt.Sprintf("%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f",
		promptArray[14],
		promptArray[13],
		promptArray[12],
		promptArray[11],
		promptArray[10],
		promptArray[9],
		promptArray[8],
		promptArray[7],
		promptArray[6],
		promptArray[5],
		promptArray[4],
		promptArray[3],
		promptArray[2],
		promptArray[1],
		promptArray[0],
	)

	cmd := exec.Command("./ai-model/ai-model-cli", "predict", prompt, "./ai-model/"+symbol+"-intra.json")

	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0., 0., prompt, err
	}

	s := strings.Split(out.String()[1:len(out.String())-2], ", ")
	f1, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return 0., 0., prompt, err
	}

	f2, err := strconv.ParseFloat(s[1], 64)
	return f1, f2, prompt, err

}

func predictDiff(d map[string]Series, t time.Time, symbol string) (float64, float64, string, error) {

	current := t.Format("2006-01-02")
	c := d[current].close()

	parse := getDate(current)
	promptArray := []float64{
		c / 500,
		getIntraDay(parse, d),
		getLast(&parse, d),
		getIntraDay(parse, d),
		getLast(&parse, d),
		getIntraDay(parse, d),
		getLast(&parse, d),
		getIntraDay(parse, d),
		getLast(&parse, d),
		getIntraDay(parse, d),
		getLast(&parse, d),
		getIntraDay(parse, d),
		(float64(t.Year()) - 2010) / 13,
		float64(t.Month()) / 12,
		float64(t.Day()) / float64(time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()),
		float64(t.Weekday()) / 5,
	}

	prompt := fmt.Sprintf("%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f,%.17f",
		promptArray[15],
		promptArray[14],
		promptArray[13],
		promptArray[12],
		promptArray[11],
		promptArray[10],
		promptArray[9],
		promptArray[8],
		promptArray[7],
		promptArray[6],
		promptArray[5],
		promptArray[4],
		promptArray[3],
		promptArray[2],
		promptArray[1],
		promptArray[0],
	)

	cmd := exec.Command("./ai-model/ai-model-cli", "predict", prompt, "./ai-model/"+symbol+"-diff.json")

	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0., 0., "", err
	}

	s := strings.Split(out.String()[1:len(out.String())-2], ", ")
	f1, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return 0., 0., "", err
	}

	f2, err := strconv.ParseFloat(s[1], 64)
	return f1, f2, prompt, err

}

func getLast(date *time.Time, d map[string]Series) float64 {
	o := d[date.Format("2006-01-02")].open()
	for true {
		*date = date.AddDate(0, 0, -1)
		if d[date.Format("2006-01-02")].Open != "" {
			c := d[date.Format("2006-01-02")].close()
			return (o / c) / 2
		}
	}

	return 0.0
}

func getIntraDay(date time.Time, d map[string]Series) float64 {
	o := d[date.Format("2006-01-02")].open()
	c := d[date.Format("2006-01-02")].close()
	return (c / o) / 2
}

func (d Series) open() float64 {
	f, _ := strconv.ParseFloat(d.Open, 64)
	return f
}

func (d Series) close() float64 {
	f, _ := strconv.ParseFloat(d.Close, 64)
	return f
}

func getDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}
