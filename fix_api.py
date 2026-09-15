import sys
content = open('backend/internal/api/api.go').read()
replacement = '''if r.Method != http.MethodGet {
		http.Error(w, {"error":"method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	if s.store == nil || s.store.PG == nil {
		http.Error(w, {"error":"internal server error: database offline"}, http.StatusInternalServerError)
		return
	}'''

content = content.replace('''if r.Method != http.MethodGet {
		http.Error(w, {"error":"method not allowed"}, http.StatusMethodNotAllowed)
		return
	}''', replacement)

open('backend/internal/api/api.go', 'w').write(content)
