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
	Velocidad          float64 `json:"Velocidad,omitempty"`
	Aceleracion        float64 `json:"Aceleracion,omitempty"`
	TemperaturaCompost float64 `json:"TemperaturaCompost,omitempty"`
	HumedadCompost     float64 `json:"HumedadCompost,omitempty"`
	PresionCompost     float64 `json:"PresionCompost,omitempty"`
	TemperaturaRuta    float64 `json:"TemperaturaRuta,omitempty"`
	HumedadRuta        float64 `json:"HumedadRuta,omitempty"`
	PresionRuta        float64 `json:"PresionRuta,omitempty"`
	TipoVia            string  `json:"TipoVia,omitempty"`
	Autonomia          float64 `json:"Autonomia,omitempty"`
}

type Respuesta struct {
	Rendimiento string `json:"Rendimiento"`
	Codigo      int    `json:"Codigo"`
	Regla       string `json:"Regla"`
}

func (rp *Respuesta) String() string {
	return fmt.Sprintf("El Rendimiento esta %s. Dentro = %d  Fuera = %d", rp.Rendimiento, rp.Codigo, rp.Regla)
}

var Parametros = &Parametro{
	Velocidad:          0,
	Aceleracion:        0,
	TemperaturaCompost: 0,
	HumedadCompost:     0,
	PresionCompost:     0,
	TemperaturaRuta:    0,
	HumedadRuta:        0,
	PresionRuta:        0,
	TipoVia:            "",
	Autonomia:          0,
}

func calcularRendimientoEndPoint(writer http.ResponseWriter, request *http.Request) {

	var newParametro Parametro
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(writer, "Datos invalidos")
	}

	json.Unmarshal(reqBody, &newParametro)
	Parametros = &Parametro{
		Velocidad:          newParametro.Velocidad,
		Aceleracion:        newParametro.Aceleracion,
		TemperaturaCompost: newParametro.TemperaturaCompost,
		HumedadCompost:     newParametro.HumedadCompost,
		PresionCompost:     newParametro.PresionCompost,
		TemperaturaRuta:    newParametro.TemperaturaRuta,
		HumedadRuta:        newParametro.HumedadRuta,
		PresionRuta:        newParametro.PresionRuta,
		TipoVia:            newParametro.TipoVia,
		Autonomia:          newParametro.Autonomia,
	}

	var response = CalcularRendimientos()
	type allResults []Respuesta

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(response)
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
}
