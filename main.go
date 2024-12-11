package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type DataItem struct {
	ID    int
	Name  string
	Value float64
}

type DataPage struct {
	Items      []DataItem
	PageNum    int
	TotalPages int
}

var data []DataItem

func init() {
	// Simulate a large dataset
	for i := 1; i <= 100; i++ {
		data = append(data, DataItem{
			ID:    i,
			Name:  "Item " + strconv.Itoa(i),
			Value: float64(i) * 1.1,
		})
	}
}

func renderPage(w http.ResponseWriter, r *http.Request) {
	// Parse page number from URL
	pageStr := r.URL.Path[len("/page/"):]
	pageNum, err := strconv.Atoi(pageStr)
	if err != nil || pageNum < 1 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	itemsPerPage := 10
	totalPages := (len(data) + itemsPerPage - 1) / itemsPerPage

	if pageNum > totalPages {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	start := (pageNum - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > len(data) {
		end = len(data)
	}

	tmpl, err := template.New("page.html").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}).ParseFiles("page.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := DataPage{
		Items:      data[start:end],
		PageNum:    pageNum,
		TotalPages: totalPages,
	}

	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/page/", renderPage)
	log.Println("Server is running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
