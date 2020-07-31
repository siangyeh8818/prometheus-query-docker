package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	//prometheus "github.com/ryotarai/prometheus-query/client"
	//prometheus "github.com/ryotarai/prometheus-query/client"
	prometheus "github.com/siangyeh8818/prometheus-query-docker/internal/client"
)

func PrintResp(resp *prometheus.QueryRangeResponse, format string, outputfilename string) error {
	switch format {
	case "tsv":
		return PrintRespXSV(resp, "\t", outputfilename)
	case "csv":
		return PrintRespXSV(resp, ",", outputfilename)
	case "json":
		return PrintRespJSON(resp)
	}

	return fmt.Errorf("unknown format: %s", format)
}

func PrintRespJSON(resp *prometheus.QueryRangeResponse) error {
	type valueEntry struct {
		Metric map[string]string `json:"metric"`
		Value  float64           `json:"value"`
	}
	type timeEntry struct {
		Time   int64         `json:"time"`
		Values []*valueEntry `json:"values"`
	}
	entryByTime := map[int64]*timeEntry{}

	for _, r := range resp.Data.Result {
		for _, v := range r.Values {
			t := v.Time()
			u := t.Unix()
			e, ok := entryByTime[u]
			if !ok {
				e = &timeEntry{
					Time:   u,
					Values: []*valueEntry{},
				}
				entryByTime[u] = e
			}

			val, err := v.Value()
			if err != nil {
				return err
			}
			e.Values = append(e.Values, &valueEntry{
				Metric: r.Metric,
				Value:  val,
			})
		}
	}

	s := make([]*timeEntry, len(entryByTime))
	i := 0
	for _, e := range entryByTime {
		s[i] = e
		i++
	}

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}

func PrintRespXSV(resp *prometheus.QueryRangeResponse, delimiter string, outputfilename string) error {
	type valueByMetric map[string]float64

	valuesByTime := map[time.Time]valueByMetric{}
	metrics := []string{}

	for _, r := range resp.Data.Result {
		metric := stringMapToString(r.Metric, "|")
		for _, v := range r.Values {
			t := v.Time()
			//t.Format("2006-01-02 15:04:05")
			//fmt.Println(t)
			d, ok := valuesByTime[t]
			if !ok {
				d = valueByMetric{}
				valuesByTime[t] = d
			}
			var err error
			d[metric], err = v.Value()
			if err != nil {
				return err
			}
		}

		found := false
		for _, m := range metrics {
			if m == metric {
				found = true
			}
		}
		if !found {
			metrics = append(metrics, metric)
		}
	}

	type st struct {
		time time.Time
		v    valueByMetric
	}
	slice := make([]st, len(valuesByTime))
	i := 0
	for t, v := range valuesByTime {
		slice[i] = st{t, v}
		i++
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].time.Before(slice[j].time)
	})
	var result string
	// header
	result = "time" + delimiter + strings.Join(metrics, delimiter) + "\n"
	fmt.Printf("time%s%s\n", delimiter, strings.Join(metrics, delimiter))

	// print rows
	for _, s := range slice {
		values := make([]string, len(metrics))
		for i, m := range metrics {
			if v, ok := s.v[m]; ok {
				values[i] = fmt.Sprintf("%f", v)
			} else {
				values[i] = ""
			}
		}
		result = result + s.time.Format("2006-01-02 15:04:05") + delimiter + strings.Join(values, delimiter) + "\n"
		fmt.Printf("%s%s%s\n", s.time.Format("2006-01-02 15:04:05"), delimiter, strings.Join(values, delimiter))
		//fmt.Printf("%d%s%s\n", s.time.Unix(), delimiter, strings.Join(values, delimiter))
	}
	WriteWithIoutil(outputfilename, result)

	return nil
}

func WriteWithIoutil(name, content string) {
	data := []byte(content)
	if ioutil.WriteFile(name, data, 0644) == nil {
		fmt.Println("Success to export to file\n", content)
	}

}

/*
func tranformTimeForamt(file string) error {

	fi, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	defer fi.Close()

	line :=0
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if line !=0 {
			//fmt.Println(string(a))
			temparray := strings.Split(string(a), ",")
			fmt.Println(temparray[0])

		}

		//temparray := strings.Split(string(a), ",")

        line++
		//fmt.Println(string(a))
	}
	return nil
}

func WriteWithIoutil(name, content string) {
	data := []byte(content)
	if ioutil.WriteFile(name, data, 0644) == nil {
		fmt.Println("Success to export to file\n", content)
	}

}
*/
