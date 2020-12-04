package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Materia struct {
	Alumno       string
	Materia      []string
	Calificacion []string
}

type AdminMaterias struct {
	Materias []Materia
}

var misMaterias AdminMaterias

func (materias *AdminMaterias) AgregarCalificacion(materia Materia) {
	flag := false
	for mat, _ := range materias.Materias {
		if materias.Materias[mat].Alumno == materia.Alumno {
			flag = true
			materias.Materias[mat].Materia = append(materias.Materias[mat].Materia, materia.Materia[0])
			materias.Materias[mat].Calificacion = append(materias.Materias[mat].Calificacion, materia.Calificacion[0])
		}
	}
	if flag == false {
		materias.Materias = append(materias.Materias, materia)
	}
}

func (materias *AdminMaterias) promedioalumno(name string) float64 {
	flag := false
	var suma float64
	var con float64
	for mat, _ := range materias.Materias {
		if materias.Materias[mat].Alumno == name {
			flag = true
			for cal, _ := range materias.Materias[mat].Calificacion {
				if s, err := strconv.ParseFloat(materias.Materias[mat].Calificacion[cal], 64); err == nil {
					suma = suma + s
					fmt.Println("suma ", suma)
					con = con + 1.0
				}
			}
		}
	}
	if flag == false {
		return 0.0
	} else {
		return suma / con
	}
}

func (materias *AdminMaterias) promedioamateria(name string) float64 {
	flag := false
	var suma float64
	var con float64
	for mat, _ := range materias.Materias {
		for mate, _ := range materias.Materias[mat].Materia {
			if materias.Materias[mat].Materia[mate] == name {
				flag = true

				if s, err := strconv.ParseFloat(materias.Materias[mat].Calificacion[mate], 64); err == nil {
					suma = suma + s
					con = con + 1.0
				}
			}
		}
	}
	if flag == false {
		return 0.0
	} else {
		return suma / con
	}
}

func (materias *AdminMaterias) promediogrl() float64 {
	flag := false
	var suma float64
	var con float64
	for mat, _ := range materias.Materias {
		flag = true
		for cal, _ := range materias.Materias[mat].Calificacion {
			if s, err := strconv.ParseFloat(materias.Materias[mat].Calificacion[cal], 64); err == nil {
				suma = suma + s
				con = con + 1.0
			}
		}
	}
	if flag == false {
		return 0.0
	} else {
		return suma / con
	}
}

func cargarHtml(a string) string {
	html, _ := ioutil.ReadFile(a)

	return string(html)
}

func form(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("form.html"),
	)
}

func formprommateria(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("formprommate.html"),
	)
}

func formpromalum(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("promedioalumno.html"),
	)
}

func formpromgrl(res http.ResponseWriter, req *http.Request) {
	prom := misMaterias.promediogrl()
	fmt.Println(prom)
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("promediogrl.html"),
		prom,
	)
}

func materias(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	var materias []string
	var calificacion []string
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		materias = append(materias, req.FormValue("materia"))
		calificacion = append(calificacion, req.FormValue("calificacion"))
		materia_1 := Materia{Alumno: req.FormValue("alumno"), Materia: materias, Calificacion: calificacion}
		misMaterias.AgregarCalificacion(materia_1)
		fmt.Println(misMaterias)
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("respuesta.html"),
			materia_1.Alumno,
		)
	}
}

func promedioalum(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		alumno := req.FormValue("alumno")
		prom := misMaterias.promedioalumno(alumno)
		fmt.Println(prom)

		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("respuesta1.html"),
			prom,
		)
	}
}

func promediomate(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		mate := req.FormValue("materia")
		prom := misMaterias.promedioamateria(mate)
		fmt.Println(prom)

		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("respuesta2.html"),
			prom,
		)
	}
}

func main() {
	http.HandleFunc("/form", form)
	http.HandleFunc("/materias", materias)
	http.HandleFunc("/formpromalum", formpromalum)
	http.HandleFunc("/promedioalum", promedioalum)
	http.HandleFunc("/formprommate", formprommateria)
	http.HandleFunc("/promediomate", promediomate)
	http.HandleFunc("/formpromgrl", formpromgrl)
	fmt.Println("Corriendo servirdor de materias...")
	http.ListenAndServe(":9000", nil)
}
