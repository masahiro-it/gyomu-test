package main


import (
	"strconv"
	_"fmt"
	"database/sql"
	_"github.com/mattn/go-sqlite3"

	"log"
	"net/http"
	"text/template"
)


func index(w http.ResponseWriter, rq *http.Request, tmp *template.Template) {
	er := tmp.Execute(w, nil)
	if er != nil {
		log.Fatal(er)
	}

}

func gyomu(w http.ResponseWriter, rq *http.Request, tmp *template.Template) {
	msg := "必要項目を入力してください。"

	if rq.Method == "POST" {
		nm := rq.PostFormValue("name")

		con, er := sql.Open("sqlite3", "data.sqlite3")
		if er != nil {
			panic(er)
		}
		defer con.Close()

		q := "insert into process_info (name) values (?)"
		con.Exec(q, nm)
	}

	item := struct {
		Message string
	} {
		Message: msg,
	}

	er := tmp.Execute(w, item)
	if er != nil {
		log.Fatal(er)
	}


}

func show(w http.ResponseWriter, rq *http.Request, tmp *template.Template) {
	item := struct {
		Title string
		Names []string
	} {
		Title: "レコード表示",
	}

	con, er := sql.Open("sqlite3", "data.sqlite3")
	if er != nil {
		panic(er)
	}
	defer con.Close()

	q := "select * from process_info"
	rs, er := con.Query(q)
	if er != nil {
		panic(er)
	}
	for rs.Next() {
		var pi Process_info
		er := rs.Scan(&pi.ID, &pi.Name)
		if er != nil {
			panic(er)
		}
		item.Names = append(item.Names, pi.Name)
	}

	er = tmp.Execute(w, item)
	if er != nil {
		log.Fatal(er)
	}

}


func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		tf, er := template.ParseFiles("templates/index.html")
		if er != nil {
		
			tf, _ = template.New("index").Parse("<html><body><h1>NO TEMPLATE.</h1></body></html>")
			
		}
		index(w, rq, tf)
	})

	http.HandleFunc("/gyomu", func(w http.ResponseWriter, rq *http.Request) {
		gtf, ger := template.ParseFiles("templates/gyomu.html")
		if ger != nil {
		
			gtf, _ = template.New("index").Parse("<html><body><h1>NO TEMPLATE.</h1></body></html>")
			
		}
		gyomu(w, rq, gtf)
	})

	http.HandleFunc("/show", func(w http.ResponseWriter, rq *http.Request) {
		stf, ser := template.ParseFiles("templates/show.html")
		if ser != nil {
		
			stf, _ = template.New("index").Parse("<html><body><h1>NO TEMPLATE.</h1></body></html>")
			
		}
		show(w, rq, stf)
	})

	http.ListenAndServe("", nil)


}




// Process_info構造体を用意する
type Process_info struct {
	ID int
	Name string
}

// Strファンクションはstring値を取得する。
func (p *Process_info) Str() string {
	return "<\"" + strconv.Itoa(p.ID) + ":" + p.Name + "\" " + ">"
}
