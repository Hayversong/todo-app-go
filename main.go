package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type Tarefa struct {
	Titulo    string
	Concluida bool
}

var tarefas []Tarefa

func home(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		tarefa := r.FormValue("tarefa")

		if tarefa != "" {

			tarefas = append(
				tarefas,
				Tarefa{
					Titulo: tarefa,
				},
			)
		}

		http.Redirect(
			w,
			r,
			"/",
			http.StatusSeeOther,
		)

		return
	}

	tmpl, err := template.ParseFiles(
		"static/index.html",
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	tmpl.Execute(
		w,
		tarefas,
	)
}

func concluir(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if id >= 0 && id < len(tarefas) {

		tarefas[id].Concluida =
			!tarefas[id].Concluida
	}

	http.Redirect(
		w,
		r,
		"/",
		http.StatusSeeOther,
	)
}

func excluir(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if id >= 0 &&
		id < len(tarefas) {

		tarefas =
			append(
				tarefas[:id],
				tarefas[id+1:]...,
			)
	}

	http.Redirect(
		w,
		r,
		"/",
		http.StatusSeeOther,
	)
}

func main() {

	http.HandleFunc("/", home)

	http.HandleFunc(
		"/concluir",
		concluir,
	)

	http.HandleFunc(
		"/excluir",
		excluir,
	)

	http.Handle(
		"/style.css",
		http.FileServer(
			http.Dir("./static"),
		),
	)

	http.ListenAndServe(
		":8080",
		nil,
	)
}
