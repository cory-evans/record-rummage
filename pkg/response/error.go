package response

import "net/http"

func Error(w http.ResponseWriter, err error, code ...int) {

	if len(code) > 0 {
		http.Error(w, err.Error(), code[0])
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
