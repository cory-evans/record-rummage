package response

import "net/http"

func Redirect(w http.ResponseWriter, r *http.Request, url string, code ...int) {
	if len(code) > 0 {
		http.Redirect(w, r, url, code[0])
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
