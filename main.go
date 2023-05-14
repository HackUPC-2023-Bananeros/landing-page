package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
)

func main() {
	db, err := sql.Open("sqlite3", "./data2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS requests (	ip_address TEXT, seat TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.RawQuery
		re := regexp.MustCompile(`seat=(\w+)`)

		resultado := re.FindStringSubmatch(param)
		if len(resultado) > 1 {
			valor := resultado[1]

			// Imprimimos el valor encontrado
			// Obtener la dirección IP del solicitante
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)

			// Guardar la dirección IP en la base de datos SQLite
			_, err := db.Exec("INSERT INTO requests (ip_address, seat) VALUES (?, ?)", ip, valor)
			if err != nil {
				log.Fatal(err)
				os.Exit(-1)
			}
			// Leer el contenido del archivo "test.txt"
			fileContent, err := ioutil.ReadFile("index.html")
			if err != nil {
				http.Error(w, "Error al leer el archivo", http.StatusInternalServerError)
				os.Exit(-1)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			w.Write(fileContent)
		}
	})
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {

		// Leer el contenido del archivo "test.txt"
		fileContent, err := ioutil.ReadFile("test.txt")
		if err != nil {
			http.Error(w, "Error al leer el archivo", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=test.txt")
		w.Write(fileContent)

	})
	// Iniciar el servidor en el puerto 8080
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
