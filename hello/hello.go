package hello

import (
  "html/template"
  "net/http"
  "time"
  "fmt"
)

func init() {
  http.HandleFunc("/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
  param := r.FormValue("date")
  var date time.Time
  if param == "" {
    date = time.Now()
  } else {
    d, err := time.Parse("2006-01-02", param)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    date = d
  }
  render(w, r, date)
}

type DateAndUrl struct {
  Date string
  Url string
}
type Params struct {
  Urls []DateAndUrl
  Date string
  PrevDate string
  NextDate string
}

func render(w http.ResponseWriter, r *http.Request, date time.Time) {
  urls := make([]DateAndUrl, 7)
  for i := range urls {
    d := addDays(date, -i)
    urls[i] = DateAndUrl{ d.Format("2006-01-02"), urlForDate(d) }
  }
  err := garfieldTemplate.Execute(w, Params{
    urls,
    date.Format("2006-01-02"),
    addDays(date, -7).Format("2006-01-02"),
    addDays(date, 7).Format("2006-01-02") })
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func addDays(date time.Time, days int) time.Time {
  return date.Add(time.Duration(24 * days) * time.Hour)
}

func urlForDate(date time.Time) string {
  year, month, day := date.Date()
  var ext string
  if date.Weekday() == time.Sunday {
    ext = "jpg"
  } else {
    ext = "gif"
  }
  return fmt.Sprintf("http://picayune.uclick.com/comics/ga/%4d/ga%02d%02d%02d.%s",
    year, year % 100, month, day, ext)
}

var garfieldTemplate = template.Must(template.New("garfield").Parse(garfieldTemplateHTML))
const garfieldTemplateHTML = `
<html>
  <head>
    <style>
      .comic { width: 100% }
      .date { padding-top: 1em; display: block; font-size: 200%; }
    </style>
  </head>
  <body>
    <a href="/?date={{.PrevDate}}">Prev week</a>
    {{.Date}}
    <a href="/?date={{.NextDate}}">Next week</a>
    {{range .Urls}}
      <div>
        <span class="date">{{.Date}}</span>
        <img class="comic" src="{{.Url}}">
      </div>
    {{end}}
  </body>
</html>
`
