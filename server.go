package main

import (
	"github.com/drone/routes"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Information struct{
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	Id string `json:"id"`
	Password string `json:"password"`
	Mac string `json:"mac"`
}
type Auth struct{
	Status string `json:"status"`
}

type Info struct{
	Id string `json:"id"`
	Here string `json:"here"`
}

func main() {
	mux := routes.New()
	mux.Get("/profile/:Id/:password/:mac", GetProfile)	
	mux.Post("/profile", PostProfile)
	mux.Post("/here", PostProf)
	mux.Put("/profile/:mac", PutProfile)
	http.Handle("/", mux)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func PostProfile(w http.ResponseWriter, r *http.Request) {
	var temp Information
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	w.WriteHeader(201)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &temp); err != nil { 
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
	}
	}
	fname:=temp.FirstName
	lname:=temp.LastName
	id:=temp.Id
	pass:=temp.Password
	mac:=temp.Mac
	fmt.Println(fname,lname,id,pass,mac)
	db, err := sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/273")
	checkErr(err)
	stmt, err:=db.Query("INSERT INTO students(firstname,lastname,id,password,mac) VALUES(?,?,?,?,?)",fname,lname,id,pass,mac)
	checkErr(err)
	defer stmt.Close()
	defer db.Close()
	
}

func PostProf(w http.ResponseWriter, r *http.Request) {
	var temp1 Info
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	w.WriteHeader(201)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &temp1); err != nil { 
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
	}
	}
	here:=temp1.Here
	id:=temp1.Id
	fmt.Println(here,id)
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/273?charset=utf8")
    checkErr(err)
	stmt, err := db.Prepare("update students set here=? where id=?")
    checkErr(err)
	
	res, err := stmt.Exec(here, id)
    checkErr(err)
	
	affect, err := res.RowsAffected()
    checkErr(err)

    fmt.Println(affect)
	defer db.Close()
}


func PutProfile (w http.ResponseWriter, r *http.Request) {
	params:=r.URL.Query()
	mac:=params.Get(":mac")
	fmt.Println(mac)
	w.WriteHeader(204)
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/273?charset=utf8")
    checkErr(err)
	stmt, err := db.Prepare("update students set attendance=? where here=? and mac=?")
    checkErr(err)
	
	res, err := stmt.Exec("present", "true", mac)
    checkErr(err)
	
	affect, err := res.RowsAffected()
    checkErr(err)

    fmt.Println(affect)
	defer db.Close()
	
	
	
}

func GetProfile (w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":Id")
	pass:=params.Get(":password")
	mac:=params.Get(":mac")
	w.WriteHeader(200)
	fmt.Println(id)
	fmt.Println(pass)
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/273?charset=utf8")
    checkErr(err)
	stmt, err := db.Query("select * from students where id=? and password=? and mac=?",id,pass,mac)
    checkErr(err)
	
    var temp Auth
    status:="false"
     count :=0
     for stmt.Next() {
     	count=count+1
     }
     fmt.Println(count)
    if (count!=0){
    	status="true"
    }
	defer db.Close()
	fmt.Println(status)
	temp.Status=status
	
		
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(temp); err != nil {
			panic(err)
		}

}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

