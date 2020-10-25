package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	engine2 "github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type Parametro struct {
	Velocidad float64 `json:"velocidad,omitempty"`
	Distancia float64 `json:"distancia,omitempty"`
}

type Respuesta struct {
	Dentro      int    `json:"dentro"`
	Fuera       int    `json:"fuera"`
	Rendimiento string `json:"rendimiento"`
}

func (rp *Respuesta) String() string {
	return fmt.Sprintf("El rendimiento esta %s. Dentro = %d  Fuera = %d", rp.Rendimiento, rp.Dentro, rp.Fuera)
}

var Parametros = &Parametro{
	Velocidad: 0,
	Distancia: 0,
}

func calcularRendimientoEndPoint(writer http.ResponseWriter, request *http.Request) {

	/*var parametro Parametro
	_ = json.NewDecoder(request.Body).Decode(&parametro)
	Parametros = &Parametro{
		Velocidad: parametro.Velocidad,
		Distancia: parametro.Distancia,
	}
	var response = CalcularRendimientos()
	type allResults []Respuesta
	var rosos = allResults{
		{
			Dentro:      response.Dentro,
			Fuera:       response.Fuera,
			Rendimiento: response.Rendimiento,
		},
	}
	json.NewEncoder(writer).Encode(rosos)*/

	var newParametro Parametro
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(writer, "Datos invalidos")
	}

	json.Unmarshal(reqBody, &newParametro)
	Parametros = &Parametro{
		Velocidad: newParametro.Velocidad,
		Distancia: newParametro.Distancia,
	}

	var response = CalcularRendimientos()
	type allResults []Respuesta
	var result = allResults{
		{
			Dentro:      response.Dentro,
			Fuera:       response.Fuera,
			Rendimiento: response.Rendimiento,
		},
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(result)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/CalcularRendimiento", calcularRendimientoEndPoint).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", router))
}

func CalcularRendimientos() (res *Respuesta) {

	respuesta := &Respuesta{}

	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	err := rb.BuildRuleFromResource("Calcular Rendimiento", "0.0.1", pkg.NewFileResource("RulesRendimiento.grl"))
	if err != nil {
		panic(err)
	}

	engine := engine2.NewGruleEngine()

	kb := lib.NewKnowledgeBaseInstance("Calcular Rendimiento", "0.0.1")

	dctx := ast.NewDataContext()
	dctx.Add("Respuesta", respuesta)
	dctx.Add("Parametro", Parametros)
	err = engine.Execute(dctx, kb)
	if err != nil {
		panic(err)
	}

	res = respuesta
	return
	//fmt.Println(respuesta.String())
}
