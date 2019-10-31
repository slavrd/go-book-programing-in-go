package main

import (
	"fmt"
	"log"
	"math/cmplx"
	"net/http"
	"strconv"
)

const (
	pageTop = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>Statistics</title>
<body><h3>Quadratic equation solver</h3>
<p>Computes equations in the form of ax<sup>2</sup> + bx + c = 0</p>`
	form = `    <form action="/" method="POST">
<input type="text" name="coefa" size="3">x<sup>2</sup> + 
<input type="text" name="coefb" size="3">x + 
<input type="text" name="coefc" size="3"> = 0 ->        
<input type="submit" value="Calculate">
</form></br>`
	pageBottom = `</body></html>`
	anError    = `<p class="error">%s</p>`
)

type quadratic struct {
	a float64
	b float64
	c float64
}

func main() {
	http.HandleFunc("/", homePage)
	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm() // Must be called before writing response
	fmt.Fprint(writer, pageTop, form)
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	} else {
		if equation, message, ok := processRequest(request); ok {
			solution1, solution2 := resolveEquation(equation)
			fmt.Fprint(writer, fmtSolution(equation, solution1, solution2))
		} else if message != "" {
			fmt.Fprintf(writer, anError, message)
		}
	}
	fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) (quadratic, string, bool) {
	var eq quadratic
	log.Printf("%v", request.Form)
	if len(request.Form) == 0 {
		return eq, "", false
	}
	var err error
	eq.a, err = strconv.ParseFloat(request.Form["coefa"][0], 64)
	if err != nil {
		return quadratic{}, "invalid 'a' coefficient", false
	}

	eq.b, err = strconv.ParseFloat(request.Form["coefb"][0], 64)
	if err != nil {
		return quadratic{}, "invalid 'b' coefficient", false
	}

	eq.c, err = strconv.ParseFloat(request.Form["coefc"][0], 64)
	if err != nil {
		return quadratic{}, "invalid 'c' coefficient", false
	}

	log.Printf("%v", eq)
	return eq, "", true
}

func resolveEquation(eq quadratic) (complex128, complex128) {
	cmplA := complex(eq.a, 0)
	cmplB := complex(eq.b, 0)
	cmplC := complex(eq.c, 0)
	r1 := (-1*cmplB + cmplx.Sqrt(cmplx.Pow(cmplB, 2)-4*cmplA*cmplC)) / 2 * cmplA
	r2 := (-1*cmplB - cmplx.Sqrt(cmplx.Pow(cmplB, 2)-4*cmplA*cmplC)) / 2 * cmplA
	log.Printf("%v", 2.0+cmplx.Sqrt(complex(-4, 0)))
	return r1, r2
}

func fmtSolution(eq quadratic, s1, s2 complex128) string {
	return fmt.Sprintf(`<div>%fx<sup>2</sup> %+fx %+f = 0 -> x1 = %v, x2 = %v</div>`,
		eq.a, eq.b, eq.c, s1, s2)
}
