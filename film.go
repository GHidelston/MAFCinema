package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Film struct {
	Judul     string `json:"Judul"`
	Sinopsis  string `json:"Sinopsis"`
	Rating    string `json:"Rating"`
	Tahun     string `json:"Tahun"`
	Genre     string `json:"Genre"`
	Sutradara string `json:"Sutradara"`
}

var movies []Film

// ================= LOAD DATA =================
func loadMovies() error {
	file, err := os.Open("film.json")
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &movies)
}

// ================= HOME =================
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, "Welcome to MAF Cinema API")
	fmt.Fprintln(w, "Available endpoints:")
	fmt.Fprintln(w, "GET  /film")
	fmt.Fprintln(w, "GET  /film/{keyword}")
	fmt.Fprintln(w, "POST /film")
}

// ================= FILM API =================
func filmHandler(w http.ResponseWriter, r *http.Request) {
	// === CORS ===
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Ambil keyword dari URL
	pathInput := strings.Trim(strings.TrimPrefix(r.URL.Path, "/film"), "/")

	switch r.Method {

	// ===== GET =====
	case http.MethodGet:

		// Jika ada keyword (/film/nolan)
		if pathInput != "" {
			query := strings.ToLower(pathInput)
			var hasil []Film

			for _, m := range movies {
				if strings.Contains(strings.ToLower(m.Judul), query) ||
					strings.Contains(strings.ToLower(m.Sinopsis), query) ||
					strings.Contains(strings.ToLower(m.Rating), query) ||
					strings.Contains(strings.ToLower(m.Tahun), query) ||
					strings.Contains(strings.ToLower(m.Genre), query) ||
					strings.Contains(strings.ToLower(m.Sutradara), query) {
					hasil = append(hasil, m)
				}
			}

			if len(hasil) == 0 {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Film tidak ditemukan",
				})
				return
			}

			json.NewEncoder(w).Encode(hasil)
			return
		}

		// Jika hanya /film
		json.NewEncoder(w).Encode(movies)

	// ===== POST =====
	case http.MethodPost:
		var newFilm Film

		if err := json.NewDecoder(r.Body).Decode(&newFilm); err != nil {
			http.Error(w, "JSON tidak valid", http.StatusBadRequest)
			return
		}

		movies = append(movies, newFilm)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newFilm)

	// ===== METHOD LAIN =====
	default:
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
	}
}

// ================= MAIN =================
func main() {
	if err := loadMovies(); err != nil {
		fmt.Println("Gagal memuat film.json:", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/film", filmHandler)
	http.HandleFunc("/film/", filmHandler)

	fmt.Println("Server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
