package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Task struct {
	ID       int
	DETAIL   string
	ASSIGNEE string
	DEADLINE string
	STATUS   string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "usbw"
	dbName := "db_task"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM task ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	task := Task{}
	res := []Task{}
	for selDB.Next() {
		var id int
		var detail, assignee, deadline, status string
		err = selDB.Scan(&id, &detail, &assignee, &deadline, &status)
		if err != nil {
			panic(err.Error())
		}
		task.ID = id
		task.DETAIL = detail
		task.ASSIGNEE = assignee
		task.DEADLINE = deadline
		task.STATUS = status
		res = append(res, task)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func View(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM task WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	task := Task{}
	for selDB.Next() {
		var id int
		var detail, assignee, deadline, status string
		err = selDB.Scan(&id, &detail, &assignee, &deadline, &status)
		if err != nil {
			panic(err.Error())
		}
		task.ID = id
		task.DETAIL = detail
		task.ASSIGNEE = assignee
		task.DEADLINE = deadline
		task.STATUS = status
	}
	tmpl.ExecuteTemplate(w, "View", task)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM task WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	task := Task{}
	for selDB.Next() {
		var id int
		var detail, assignee, deadline, status string
		err = selDB.Scan(&id, &detail, &assignee, &deadline, &status)
		if err != nil {
			panic(err.Error())
		}
		task.ID = id
		task.DETAIL = detail
		task.ASSIGNEE = assignee
		task.DEADLINE = deadline
		task.STATUS = status
	}
	tmpl.ExecuteTemplate(w, "Edit", task)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		detail := r.FormValue("detail")
		assignee := r.FormValue("assignee")
		deadline := r.FormValue("deadline")
		insForm, err := db.Prepare("INSERT INTO task(detail, assignee, deadline) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(detail, assignee, deadline)
		log.Println("INSERT: Detail: " + detail + " | Assignee: " + assignee + " | Deadline: " + deadline)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		detail := r.FormValue("detail")
		assignee := r.FormValue("assignee")
		deadline := r.FormValue("deadline")
		id := r.FormValue("id")
		insForm, err := db.Prepare("UPDATE task SET detail=?, assignee=?, deadline=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(detail, assignee, deadline, id)
		log.Println("UPDATE: Detail: " + detail + " | Assignee: " + assignee + " | Deadline: " + deadline)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func UpdateStatusDone(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	id := r.URL.Query().Get("id")
	insForm, err := db.Prepare("UPDATE task SET status='done' WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(id)
	log.Println("UPDATE: Status: done")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	task := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM task WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(task)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/view", View)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/update-status", UpdateStatusDone)
	http.HandleFunc("/delete", Delete)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.ListenAndServe(":8080", nil)
}
