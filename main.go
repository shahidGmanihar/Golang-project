package main

import (
	"database/sql"
	//"strconv"
	"crypto/sha1"
	"fmt"

	"strings"

	//"github.com/mailru/easyjson"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	//"strconv"

	"text/template"

	"github.com/360EntSecGroup-Skylar/excelize"
	//"github.com/joho/sqltocsv"
	"github.com/tealeg/xlsx"

	"github.com/gorilla/securecookie"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id            int    `xml:"id"`
	Name          string `xml:"name"`
	Birthdate     string `xml:"birthday"`
	Address       string `xml:"address"`
	Pancard       string `xml:"pancard"`
	Familydetails string `xml:"familydetails"`
}

//Cookie Authentication
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var filename string

func getUserName(request *http.Request) (userName string) {

	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name1"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {

	value := map[string]string{
		"name1": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	println(userName)
	clearSession(response)
	userNamee := getUserName(request)
	println("userName", userNamee)
	http.Redirect(response, request, "/", 302)
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "admin"
	dbName := "userdb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Login", nil)
}

func Loginn(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	name1 := r.FormValue("name")
	y := r.FormValue("pancard")
	selDB, err := db.Query("SELECT * FROM employee where name=? and pancard=?", name1, y)

	println("Login Done", "name:", name1, "y:", y)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name string
		var birthdate string
		var address string
		var pancard string
		var familydetails string
		err = selDB.Scan(&id, &name, &birthdate, &address, &pancard, &familydetails)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Birthdate = birthdate
		emp.Address = address
		emp.Pancard = pancard
		emp.Familydetails = familydetails

		res = append(res, emp)

	}
	println(emp.Name, emp.Pancard)

	if emp.Name == name1 && emp.Pancard == y {
		setSession(name1, w)
		http.Redirect(w, r, "/index", 301)
	} else {
		http.Redirect(w, r, "/", 301)
	}

	defer db.Close()
}

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	userName := getUserName(r)
	println(userName)
	selDB, err := db.Query("SELECT * FROM employee ORDER BY id DESC")

	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name string
		var birthdate string
		var address string
		var pancard string
		var familydetails string
		err = selDB.Scan(&id, &name, &birthdate, &address, &pancard, &familydetails)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Birthdate = birthdate
		emp.Address = address
		emp.Pancard = pancard
		emp.Familydetails = familydetails

		res = append(res, emp)
	}
	//tmpl.ExecuteTemplate(w, "Index", res)
	if userName != "" {
		tmpl.ExecuteTemplate(w, "Index", res)
	} else {
		http.Redirect(w, r, "/", 302)

	}

	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name string
		var birthdate string
		var address string
		var pancard string
		var familydetails string
		err = selDB.Scan(&id, &name, &birthdate, &address, &pancard, &familydetails)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Birthdate = birthdate
		emp.Address = address
		emp.Pancard = pancard
		emp.Familydetails = familydetails
	}

	tmpl.ExecuteTemplate(w, "Show", emp)

	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name string
		var birthdate string
		var address string
		var pancard string
		var familydetails string
		err = selDB.Scan(&id, &name, &birthdate, &address, &pancard, &familydetails)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Birthdate = birthdate
		emp.Address = address
		emp.Pancard = pancard
		emp.Familydetails = familydetails
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	println("Add check")

	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		birthdate := r.FormValue("birthdate")
		address := r.FormValue("address1") + ", " + r.FormValue("address2") + ", " + r.FormValue("stateSelect") + ", " + r.FormValue("citySelect") + ", " + r.FormValue("pincode")
		pancard := r.FormValue("pancard")
		familydetails := "Relation: " + r.FormValue("relationship") + ", Name: " + r.FormValue("name1") + ", DOB: " + r.FormValue("dob") + ", Address: " + r.FormValue("address3") + ", Occupation: " + r.FormValue("occupation") + "Gender :" + r.FormValue("gender")
		println("Add Done")
		insForm, err := db.Prepare("INSERT INTO Employee(name, birthdate, address, pancard, familydetails) Values(?,?,?,?,?)")

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, birthdate, address, pancard, familydetails)
		log.Println("INSERT: Name: " + name + " | Birthdate: " + birthdate + " |Address: " + address + " |Pancard:" + pancard + " |Familydetails:" + familydetails)
	}
	defer db.Close()
	http.Redirect(w, r, "/index", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		birthdate := r.FormValue("birthdate")
		address := r.FormValue("address")
		pancard := r.FormValue("pancard")
		familydetails := r.FormValue("familydetails")

		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Employee SET name=?, birthdate=?, address=?, pancard=?, familydetails=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, birthdate, address, pancard, familydetails, id)
		log.Println("UPDATE: Name: " + name + " | Birthdate: " + birthdate + " |Address: " + address + " |Pancard:" + pancard + " |Familydetails:" + familydetails)
	}
	defer db.Close()
	http.Redirect(w, r, "/index", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/index", 301)
}

//Excel file upload handling

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleUpload(w, r)
		return
	}
	tmpl.ExecuteTemplate(w, "Fileupload", nil)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) //10MB
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename = path.Base(fileHeader.Filename)
	println(filename)

	dest, err := os.Create(filename)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	if _, err = io.Copy(dest, file); err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

//EXcel Read and Export

func ExcelRead(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	emp := Employee{}
	res := []Employee{}

	//var id int
	var name string
	var birthdate string
	var address string
	var pancard string
	var familydetails string

	f, err := excelize.OpenFile(filename)
	println("filensme :", filename)

	if err != nil {
		println(err.Error())
		return
	}

	rows := f.GetRows("Sheet1")
	println("Rows: ", rows)
	var colCell string
	var d string
	println("Debugging")
	println("FileName", filename)
	//Mycode

	//file, err := xlsx.OpenFile("AddUser.xlsx")
	file, err := xlsx.OpenFile(filename)
	if err != nil {
		panic(err.Error())
	}
	for i := 1; i < len(rows); i++ {
		//id, _ = strconv.Atoi(file.Sheets[0].Rows[i].Cells[0].Value)
		name = file.Sheets[0].Rows[i].Cells[1].Value
		birthdate = (file.Sheets[0].Rows[i].Cells[2].Value)
		address = file.Sheets[0].Rows[i].Cells[3].Value
		pancard = file.Sheets[0].Rows[i].Cells[4].Value
		familydetails = file.Sheets[0].Rows[i].Cells[5].Value

		//emp.Id = id
		emp.Name = name
		//dt, _ := time.Parse("2006-01-02", birthdate)
		emp.Birthdate = birthdate

		//emp.Birthdate = dt.Format("2006-01-02")
		println("Birthdate :", emp.Birthdate)
		emp.Address = address
		emp.Pancard = pancard
		emp.Familydetails = familydetails

		println("Add Done")
		insForm, err := db.Prepare("INSERT INTO Employee(name, birthdate, address, pancard, familydetails) Values(?,?,?,?,?)")
		println(insForm)
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(emp.Name, emp.Birthdate, emp.Address, emp.Pancard, emp.Familydetails)
		println("MyCell", insForm)
		res = append(res, emp)

	}

	for _, row := range rows {
		for _, colCell = range row {
			d = colCell
			//print(colCell, "\t")

			print("d:", d, "\t")

		}
		http.Redirect(w, r, "Index", 301)

		println()
	}

	defer db.Close()

	//log.Println("INSERT Name: " + name + " |Salary: " + salary + " |City:" + city)
}

func ExcelWrite(w http.ResponseWriter, r *http.Request) {

	f, err := excelize.OpenFile("emp1.xlsx")
	//println("filename :", filename)

	if err != nil {
		println(err.Error())
		return
	}

	//file, err := xlsx.OpenFile("AddUser.xlsx")
	file, err := xlsx.OpenFile("MyExcel.xlsx")
	if err != nil {
		panic(err.Error())
	}

	db := dbConn()
	selDB, err := db.Query("SELECT * FROM employee ")

	println("selDB", selDB)
	//r :=[]rune(selDB)

	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}

	for selDB.Next() {
		var id int
		var name string
		var birthdate string
		var address string
		var pancard string
		var familydetails string
		err = selDB.Scan(&id, &name, &birthdate, &address, &pancard, &familydetails)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Birthdate = birthdate
		emp.Address = address
		emp.Pancard = pancard
		emp.Familydetails = familydetails
		println("Name ::::", emp.Name)

		res = append(res, emp)
		println("res: ", len(res))
		i := len(res)

		cap := cap(res)
		println(i, "capacity of res", cap)
		//a:=strconv.Itoa(id)

		f.SetCellValue("sheet1", "A", res)

	}
	println("res", res)

	/*rows := f.GetRows("Sheet1")
	println("Rows: ", rows)
	println("Debugging")
	for _, record := range rows {
		f.SetCellValue("sheet1", "A", record)
		println(record)
	}*/
	i := len(res)

	for _, m := range res {
		println("m", m.Name)
		//file.Sheets[1].Rows[i].Cells[0].Value = strconv.Itoa(emp.Id)

		file.Sheets[1].Rows[i].Cells[1].Value = m.Name
		sheet := file.Sheets[1].Rows[i].Cells[1].Value
		println("sheet :::::", sheet)
		(file.Sheets[1].Rows[i].Cells[2].Value) = emp.Birthdate
		file.Sheets[1].Rows[i].Cells[3].Value = emp.Address
		file.Sheets[1].Rows[i].Cells[4].Value = emp.Pancard
		file.Sheets[1].Rows[i].Cells[5].Value = emp.Familydetails

	}

	defer db.Close()
}

func photogallery(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		mf, fh, err := req.FormFile("nf")
		if err != nil {
			println(err)
		}
		defer mf.Close()
		// create sha for file name
		ext := strings.Split(fh.Filename, ".")[1]
		h := sha1.New()
		io.Copy(h, mf)
		fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
		// create new file
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		path := filepath.Join(wd, "gallery", fname)
		nf, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer nf.Close()
		mf.Seek(0, 0)
		io.Copy(nf, mf)
	}

	file, err := os.Open("gallery")
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	tmpl.ExecuteTemplate(w, "photo.gohtml", list)
}

func main() {

	log.Println("Server started on: http://localhost:9000")
	http.HandleFunc("/photo", photogallery)

	http.Handle("/gallery/", http.StripPrefix("/gallery", http.FileServer(http.Dir("./gallery"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", Login)
	http.HandleFunc("/loginn", Loginn)
	http.HandleFunc("/index", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/new1", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/add", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/del", Delete)
	http.HandleFunc("/handle", handleUpload)
	http.HandleFunc("/bulkupload", home)
	http.HandleFunc("/rowins", ExcelRead)
	http.HandleFunc("/expotoexcel", ExcelWrite)

	http.ListenAndServe(":9000", nil)
}
