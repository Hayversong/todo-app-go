package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
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
			salvarTarefas()
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
		salvarTarefas()
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
		salvarTarefas()
	}

	http.Redirect(
		w,
		r,
		"/",
		http.StatusSeeOther,
	)
}

func salvarTarefas() {

	arquivo, err :=
		os.Create(
			"tarefas.json",
		)

	if err != nil {

		return
	}

	defer arquivo.Close()

	json.NewEncoder(
		arquivo,
	).Encode(
		tarefas,
	)
}

func carregarTarefas() {

	arquivo, err :=
		os.Open(
			"tarefas.json",
		)

	if err != nil {

		return
	}

	defer arquivo.Close()

	json.NewDecoder(
		arquivo,
	).Decode(
		&tarefas,
	)
}

func main() {

	carregarTarefas()

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
