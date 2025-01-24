package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
)

type Item struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Value    int    `json:"value"`
	Quantity int    `json:"quantity"`
}

var (
	inventory []Item
	mu        sync.Mutex
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "テンプレートの読み込みに失敗しました: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, map[string]interface{}{
		"Inventory": inventory,
	})
	if err != nil {
		http.Error(w, "テンプレートのレンダリングに失敗しました: "+err.Error(), http.StatusInternalServerError)
	}
}

func addItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		category := r.FormValue("category")
		valueStr := r.FormValue("value")
		quantityStr := r.FormValue("quantity")

		if name == "" || category == "" || valueStr == "" || quantityStr == "" {
			http.Error(w, "すべてのフィールドを入力してください", http.StatusBadRequest)
			return
		}

		value, err := strconv.Atoi(valueStr)
		if err != nil {
			http.Error(w, "価値（ゴールド）は数値で入力してください", http.StatusBadRequest)
			return
		}

		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			http.Error(w, "個数は数値で入力してください", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()
		item := Item{
			Name:     name,
			Category: category,
			Value:    value,
			Quantity: quantity,
		}
		inventory = append(inventory, item)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func saveJSON(w http.ResponseWriter, r *http.Request) {
	file, err := os.Create("inventory.json")
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(inventory)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loadJSON(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("inventory.json")
	if err != nil {
		http.Error(w, "Error loading file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&inventory); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetInventory(w http.ResponseWriter, r *http.Request) {
	inventory = []Item{}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		log.Println("ブラウザの自動起動に対応していません")
	}
	if err != nil {
		log.Printf("ブラウザの起動に失敗しました: %v\n", err)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add-item", addItem)
	http.HandleFunc("/save-json", saveJSON)
	http.HandleFunc("/load-json", loadJSON)
	http.HandleFunc("/reset", resetInventory)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	go func() {
		url := "http://localhost:8080"
		log.Printf("ブラウザでサーバーを起動します: %s\n", url)
		openBrowser(url)
	}()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("サーバーエラー:", err)
	}
}
