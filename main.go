package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type PR struct {
	Reviewer    string
	PRLink      string
	Description string
}

type Main struct {
	PageTitle string
	PRs       []PR
}

func main() {

	http.HandleFunc("/", handler)
	http.HandleFunc("/generate", generateHandler)

	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	file, _ := ioutil.ReadFile("PRs.json")
	var prs struct {
		PR []PR
	}

	_ = json.Unmarshal([]byte(file), &prs)
	data := Main{
		PageTitle: "PR Generator",
		PRs:       prs.PR,
	}

	tmpl.Execute(w, data)
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	reviewer := fmt.Sprintf("@%v", r.FormValue("reviewer"))
	prLink := fmt.Sprintf("`%v`", r.FormValue("prlink"))
	description := fmt.Sprintf("```\n%v\n```", r.FormValue("description"))

	data := PR{
		Reviewer:    reviewer,
		PRLink:      prLink,
		Description: description,
	}
	savePR(data)
	t, _ := template.ParseFiles("result.html")
	t.Execute(w, data)
}

func savePR(pr PR) {
	file, _ := ioutil.ReadFile("PRs.json")
	var prs struct {
		PR []PR
	}

	_ = json.Unmarshal([]byte(file), &prs)
	fmt.Println(prs)
	prs.PR = append(prs.PR, pr)
	fmt.Println("prs", prs)

	newPRs, _ := json.MarshalIndent(prs, "", "")
	_ = ioutil.WriteFile("PRs.json", newPRs, 0644)
}
