package main

import(
	"net/http"
	"database/sql"
	"log"
	//"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
)

type History struct {
    Id	int
    Url	string
	Title	string
	Visit_time	int
}

var historys [] History
var db_file = History_path // google chrome のhistoryファイルへのパス

func main() {
	insta_db, err := sql.Open("sqlite3", "./insta.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer insta_db.Close()

	_, err = insta_db.Exec(
		`CREATE TABLE IF NOT EXISTS "INSTA" ("ID" INTEGER PRIMARY KEY, "URL" VARCHAR(255), "TITLE" VARCHAR(255), "VISIT_TIME" INTEGER)`,
	  )
	  if err != nil {
		panic(err)
	  }
	cmd := "select * FROM INSTA order by VISIT_TIME desc"
	rows, err := insta_db.Query(cmd, 18)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var a History
		err := rows.Scan(
			&a.Id,
			&a.Url,
			&a.Title,
			&a.Visit_time,
		)
		if err != nil {
			log.Fatal(err)
		}
		historys = append(historys, a)
	}  
	r := gin.Default()
	r.LoadHTMLFiles("./template/index.html", "./template/list.html", "./template/detail1.html")
	
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"history":historys,
		})
	})
	r.GET("/list", func(c *gin.Context) {
		date1_s := c.Query("q1")
		date2_s := c.Query("q2")
		
		if date1_s != "" && date2_s != ""{
			layout := "2006-01-02"
			date1, _ := time.Parse(layout, date1_s)
			date2, _ := time.Parse(layout, date2_s)
			date_dis, _ := time.Parse(layout, "1601-01-01")
			date_dis_i := date_dis.UnixMicro()
			date1_i := date1.UnixMicro() - date_dis_i
			date2_i := date2.UnixMicro() - date_dis_i
			//fmt.Println(date1_i)
			var results [] History
			insta_db, err := sql.Open("sqlite3", "./insta.sqlite3")
			if err != nil {
				log.Fatal(err)
			}
			defer insta_db.Close()
			cmd := "select * FROM INSTA WHERE ? <= VISIT_TIME AND VISIT_TIME <= ?"
			rows, err1 := insta_db.Query(cmd, date1_i, date2_i)
			if err1 != nil {
				log.Fatal(err1)
			}
			defer rows.Close()
			for rows.Next() {
				var t History

				err := rows.Scan(&t.Id, &t.Url, &t.Title, &t.Visit_time)
				if err != nil {
					log.Fatal(err)
				}
				results = append(results, t)
			}
			c.HTML(http.StatusOK, "list.html", gin.H{
				"history":results,
				"date1":date1_s,
				"date2":date2_s,
			})
		}else if date1_s != "" {
			layout := "2006-01-02"
			date1, _ := time.Parse(layout, date1_s)
			date_dis, _ := time.Parse(layout, "1601-01-01")
			date_dis_i := date_dis.UnixMicro()
			date1_i := date1.UnixMicro() - date_dis_i
			//fmt.Println(date1_i)
			var results [] History
			insta_db, err := sql.Open("sqlite3", "./insta.sqlite3")
			if err != nil {
				log.Fatal(err)
			}
			defer insta_db.Close()
			cmd := "select * FROM INSTA WHERE ? <= VISIT_TIME"
			rows, err1 := insta_db.Query(cmd, date1_i)
			if err1 != nil {
				log.Fatal(err1)
			}
			defer rows.Close()
			for rows.Next() {
				var t History

				err := rows.Scan(&t.Id, &t.Url, &t.Title, &t.Visit_time)
				if err != nil {
					log.Fatal(err)
				}
				results = append(results, t)
			}
			c.HTML(http.StatusOK, "list.html", gin.H{
				"history":results,
				"date1":date1_s,
			})
		}else if date2_s != "" {
			layout := "2006-01-02"
			date2, _ := time.Parse(layout, date2_s)
			date_dis, _ := time.Parse(layout, "1601-01-01")
			date_dis_i := date_dis.UnixMicro()
			date2_i := date2.UnixMicro() - date_dis_i
			//fmt.Println(date1_i)
			var results [] History
			insta_db, err := sql.Open("sqlite3", "./insta.sqlite3")
			if err != nil {
				log.Fatal(err)
			}
			defer insta_db.Close()
			cmd := "select * FROM INSTA WHERE VISIT_TIME <= ?"
			rows, err1 := insta_db.Query(cmd, date2_i)
			if err1 != nil {
				log.Fatal(err1)
			}
			defer rows.Close()
			for rows.Next() {
				var t History

				err := rows.Scan(&t.Id, &t.Url, &t.Title, &t.Visit_time)
				if err != nil {
					log.Fatal(err)
				}
				results = append(results, t)
			}
			c.HTML(http.StatusOK, "list.html", gin.H{
				"history":results,
				"date2":date2_s,
			})
		}else{
			c.HTML(http.StatusOK, "list.html", gin.H{
				"history":historys,
			})
		}	
	})
	r.GET("/detail/:id", ShowHistory)
	r.GET("/reload", ReadHistory)
	r.Run()
}

func ShowHistory(c *gin.Context) {
	var b History
	i := c.Param("id")
	insta_db, err := sql.Open("sqlite3", "./insta.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer insta_db.Close()
	//fmt.Println(i)
	cmd := "select * FROM INSTA WHERE ID = ?"
	err1 := insta_db.QueryRow(cmd, i).Scan(
		&b.Id,
		&b.Url,
		&b.Title,
		&b.Visit_time,
	)
	if err1 != nil {
		log.Fatal(err1)
	}
	c.HTML(http.StatusOK, "detail1.html", gin.H{
		"history":b,
	})
}

func ReadHistory(c *gin.Context) {
	os.Remove("./insta.sqlite3")
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insta_db, err := sql.Open("sqlite3", "./insta.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer insta_db.Close()
	_, err = insta_db.Exec(
		`CREATE TABLE IF NOT EXISTS "INSTA" ("ID" INTEGER PRIMARY KEY, "URL" VARCHAR(255), "TITLE" VARCHAR(255), "VISIT_TIME" INTEGER)`,
	  )
	  if err != nil {
		panic(err)
	  }
	cmd := "select visits.id, urls.url, urls.title, visits.visit_time,visits.from_visit from visits inner join urls on visits.url = urls.id"
	rows, err := db.Query(cmd)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var url string
		var title string
		var visit_time int
		var from_visit string

		err := rows.Scan(&id, &url, &title, &visit_time, &from_visit)
		if err != nil {
			log.Fatal(err)
		}
		if contain("https://www.instagram.com/p/", url) {
			// fmt.Println(url)
			_, err := insta_db.Exec(
				`INSERT INTO INSTA (ID, URL, TITLE, VISIT_TIME) VALUES (?, ?, ?, ?)`,
				id,
				url,
				title,
				visit_time,
			  )
			  if err != nil {
				panic(err)
			  }
		}
	}
    c.HTML(http.StatusOK, "index.html", gin.H{
		"history":historys,
	})
}

func contain(a string, b string) bool {
	if len(a) <= len(b) {
		if a == b[:len(a)] {
			return true
		}
	}
	return false
}
