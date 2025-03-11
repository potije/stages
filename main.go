package main

import (
	"html/template"
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Stage struct {
	Seconds int
	Delay int
	Name string
}

type MeditationData struct {
    Stages  int
    Minutes int
    Seconds int
    List []Stage
}

func handler(w http.ResponseWriter, r *http.Request) {
	    fmt.Fprintf(w, "Method: %s\n", r.Method)
	    fmt.Fprintf(w, "URL: %s\n", r.URL)
	    fmt.Fprintf(w, "Header: %v\n", r.Header)
	    fmt.Fprintf(w, "Content-Type: %s\n", r.Header.Get("Content-Type"))

	    // Read and close the request body
	    body, _ := ioutil.ReadAll(r.Body)
	    r.Body.Close()
	    fmt.Fprintf(w, "Body: %s\n", body)

	    // Parse and access form data
	    r.ParseForm()
	    fmt.Fprintf(w, "Form: %v\n", r.Form)
	    fmt.Fprintf(w, "Form value 'name': %s\n", r.Form.Get("name"))
	}

func main() {
	// Serve static files from the "static" directory for the gong
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

	tmpl := template.Must(template.ParseFiles("templates/index.html"))


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/test", handler)

	http.HandleFunc("/begin", func(w http.ResponseWriter, r*http.Request) {
		tmplm := template.Must(template.ParseFiles("templates/meditation.html"))
		err := r.ParseForm()
		if err != nil {
        	http.Error(w, "Failed to parse form", http.StatusBadRequest)
        	return
    	}
    	stages, err := strconv.Atoi(r.FormValue("stages"))
    	minutes, err := strconv.Atoi(r.FormValue("minutes"))
    	if err != nil {
        	http.Error(w, "Failed to cast stages/minutes to int", http.StatusBadRequest)
        	return
    	}
    	list := make([]Stage, stages)

    	for i := 0; i < stages; i++ {
    		list[i] = Stage{
    			Seconds: minutes*60, 
    			Delay: 5+i*minutes*60, 
    			Name: fmt.Sprintf("Stage %d", i+1),
    		}
    	}
    	// Access form values
    	data := MeditationData {
        	Stages:  stages,
        	Minutes: minutes,
        	Seconds: minutes*60,
        	List: list,
    	}
    	fmt.Println(data.List)
    	tmplm.Execute(w, data)
	})

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}