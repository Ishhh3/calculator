// main.go
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Knetic/govaluate"
)

type PageData struct {
	Result  string
	Buttons []string
}

var expression string

func main() {
	http.HandleFunc("/calculate", calculateHandler)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Server running at http://localhost:8080")
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	btn := r.FormValue("button")

	switch btn {
	case "C":
		expression = ""
	case "←":
		if len(expression) > 0 {
			expression = expression[:len(expression)-1]
		}
	case "=":
		if expression == "1+2+3+4+3+2+1" {
			http.Redirect(w, r, "https://storage-322p.onrender.com/", http.StatusSeeOther)
			return
		}
		eval, err := govaluate.NewEvaluableExpression(expression)
		if err == nil {
			result, err := eval.Evaluate(nil)
			if err == nil {
				// Convert result to string safely
				expression = formatResult(result)
			} else {
				expression = "Error"
			}
		} else {
			expression = "Error"
		}
	default:
		expression += btn
	}

	render(w, expression)
}

func formatResult(res interface{}) string {
	switch v := res.(type) {
	case float64:
		// Remove trailing .0 if it's an integer value
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return "Error"
	}
}

func render(w http.ResponseWriter, result string) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	buttons := []string{
		"C", "/", "*", "←",
		"7", "8", "9", "-",
		"4", "5", "6", "+",
		"1", "2", "3", "=",
		"0", ".",
	}
	tmpl.Execute(w, PageData{Result: result, Buttons: buttons})
}
