package main

import (
	"encoding/json"
	"log"
	"log/slog"
	_ "mock/docs"
	"net/http"
	"strings"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

var counter int

// @title Mock API
// @version 1.0
// @description This is a sample mock.

// @contact.name API Support
// @contact.email smagulmyrzakhmet@gmail.com

// @host localhost:3031
// @BasePath /mock

func main() {
	http.HandleFunc("GET /mock/info", getLyrics)
	http.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	log.Println("starting on :3031")
	if err := http.ListenAndServe(":3031", nil); err != nil {
		log.Fatalln("bruh", err)
	}
}

type Description struct {
	Description string `json:"description"`
}

type LyricsAPIResponse struct {
	Release *time.Time `json:"releaseDate"`
	Link    string     `json:"link"`
	Text    string     `json:"text"`
}

// @Summary Get song info
// @Description Получить данные о музыке
// @Tags Songs
// @Produce  json
// @Param group query string true "Group name"
// @Param song query string true "Song name"
// @Success 200 {object} LyricsAPIResponse
// @Failure 400 {object} Description
// @Failure 500 {object} Description
// @Router /info [get]
func getLyrics(w http.ResponseWriter, r *http.Request) {
	defer func() {
		counter++
	}()

	song := r.URL.Query().Get("song")
	group := r.URL.Query().Get("group")
	if strings.TrimSpace(song) == "" || strings.TrimSpace(group) == "" {
		slog.Error("bad request")
		defaultResponse(w, http.StatusBadRequest)
		return
	}
	log.Printf("song: \"%s\"; group: \"%s\"", song, group)

	day, err := time.Parse("02.01.2006", "16.07.2006")
	if err != nil {
		log.Println("bruh error", err)
		defaultResponse(w, 500)
		return
	}

	switch counter % 3 {
	case 2:
		defaultResponse(w, 500)
		return
	default:
		resp := &LyricsAPIResponse{
			Release: &day,
			Text:    "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			Link:    "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(resp)
	}
}

func defaultResponse(w http.ResponseWriter, statusCode int) {
	errResp := Description{
		Description: http.StatusText(statusCode),
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errResp)
}
