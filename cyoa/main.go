package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultTemp))
}

var tpl *template.Template

var defaultTemp = `<!DOCTYPE html>
    <head>
		<meta charset="utf-8">
        <title>
            Create Your Own Adventure
        </title>
    </head>
	<body>
		<h1>{{.Title}}</h1>
		{{range .Paragraphs}}
		<p>{{.}}</p>
		{{end}}
		<ul>
			{{range .Options}}
			<li><a href="/{{.Arc}}">{{.Text}}</a></li>
			{{end}}
		</ul>
	</body>`

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    Option   `json:"options"`
}

type Option []struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func main() {

	storyJsonFile := flag.String("story", "story.json", "json file containing the whole story")

	flag.Parse()

	f, err := os.Open(*storyJsonFile)

	if err != nil {
		panic(err)
	}

	d := json.NewDecoder(f)

	var story Story

	if err := d.Decode(&story); err != nil {
		panic(err)
	}

	h := NewHandler(story)

	log.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", h))

}

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			panic(err)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}
