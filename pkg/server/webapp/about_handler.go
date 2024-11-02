package webapp

import "net/http"

func (r *WebAppRouter) aboutHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "about.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
