package controller

import (
	"encoding/json"
	"html/template"
	"mvc/pkg/database"
	"mvc/pkg/types"
	"net/http"
	"strconv"
)

func AddDummyBookData() {
	// if the number of books in the table is more than 5, do nothing
	var count int64
	database.DB.Model(&types.Book{}).Count(&count)
	if count > 5 {
		return
	}

	books := []types.Book{
		{Title: "The Great Adventure", Author: "John Doe", Genre: "Adventure", Language: "English", Summary: "An epic journey through unknown lands.", Count: 5},
		{Title: "Mystery of the Lost Island", Author: "Jane Smith", Genre: "Mystery", Language: "English", Summary: "A thrilling mystery set on a deserted island.", Count: 3},
		{Title: "Science and the Future", Author: "Albert Newton", Genre: "Science", Language: "English", Summary: "Exploring the possibilities of future technologies.", Count: 4},
		{Title: "Romance in Paris", Author: "Emily Rose", Genre: "Romance", Language: "English", Summary: "A love story set in the city of lights.", Count: 7},
		{Title: "The Art of War", Author: "Sun Tzu", Genre: "Strategy", Language: "Chinese", Summary: "Ancient Chinese military treatise.", Count: 2},
		{Title: "Cooking Made Easy", Author: "Gordon Ramsay", Genre: "Cooking", Language: "English", Summary: "Simple recipes for everyday meals.", Count: 6},
		{Title: "Space Explorers", Author: "Isaac Asimov", Genre: "Science Fiction", Language: "English", Summary: "A journey through the cosmos with a team of explorers.", Count: 5},
		{Title: "The Haunted Mansion", Author: "Edgar Poe", Genre: "Horror", Language: "English", Summary: "A spooky tale of a haunted house.", Count: 4},
		{Title: "History of the World", Author: "William Durant", Genre: "History", Language: "English", Summary: "A comprehensive overview of world history.", Count: 8},
		{Title: "The Zen Mind", Author: "Shunryu Suzuki", Genre: "Philosophy", Language: "Japanese", Summary: "Teachings on Zen meditation and mindfulness.", Count: 3},
	}

	database.DB.Create(&books)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var books []types.APIBook
		database.DB.Model(&types.Book{}).Find(&books)

		json.NewEncoder(w).Encode(books)
	}
}

func ManageBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/manage_books.html"))

		var books []types.APIBook
		database.DB.Model(&types.Book{}).Find(&books)

		t.Execute(w, books)
		return
	}
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		count_s := r.FormValue("count")
		count, err := strconv.ParseInt(count_s, 10, 64)
		if err != nil || count <= 0 {
			http.Error(w, "count should be a positive integer!", http.StatusInternalServerError)
			return
		}

		if count > 10000 {
			http.Error(w, "can't add too many books at once!", http.StatusInternalServerError)
		}

		book := types.Book{
			Title:    r.FormValue("name"),
			Author:   r.FormValue("author"),
			Genre:    r.FormValue("genre"),
			Language: r.FormValue("language"),
			Summary:  r.FormValue("summary"),
			Count:    int(count),
		}

		database.DB.Create(&book)
	}
}

func RemoveBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		id_s := r.FormValue("id")
		id, err := strconv.ParseInt(id_s, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid book id!", http.StatusInternalServerError)
			return
		}

		database.DB.Delete(&types.Book{}, id)
	}
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}

		id_s := r.FormValue("id")
		id, err := strconv.ParseInt(id_s, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid book id!", http.StatusInternalServerError)
			return
		}

		count_s := r.FormValue("count")
		title := r.FormValue("title")
		author := r.FormValue("author")
		genre := r.FormValue("genre")
		language := r.FormValue("language")
		summary := r.FormValue("summary")
		count, err := strconv.ParseInt(count_s, 10, 64)
		if err != nil || count < 0 {
			http.Error(w, "count should be a positive integer!", http.StatusInternalServerError)
			return
		}

		if count > 10000 {
			http.Error(w, "can't add too many books at once!", http.StatusInternalServerError)
		}

		if count == 0 {
			database.DB.Delete(&types.Book{}, id)
			return
		}

		book := types.Book{
			Title:    title,
			Author:   author,
			Genre:    genre,
			Language: language,
			Summary:  summary,
			Count:    int(count),
		}

		err = database.DB.Model(&types.Book{}).Where("id = ?", id).Updates(book).Error
		if err != nil {
			http.Error(w, "Some error occured", http.StatusInternalServerError)
			return
		}
	}
}
