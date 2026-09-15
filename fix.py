import re

content = open('backend/internal/api/api.go').read()

def inject_check(func_name):
    global content
    pattern = rf"(func \(s \*Server\) {func_name}\(w http\.ResponseWriter, r \*http\.Request\) {{\n\t(.*)\n\t\t(.*)\n\t\treturn\n\t}})"
    
    replacement = r"\1\n\tif s.store == nil || s.store.PG == nil {\n\t\thttp.Error(w, {\"error\":\"internal server error: database offline\"}, http.StatusInternalServerError)\n\t\treturn\n\t}"
    content = re.sub(pattern, replacement, content)

inject_check("handleEvents")
inject_check("handleIncidents")
inject_check("handleDecisions")
# handleInvestigations calls handleDecisions, wait, let's check. No, it doesn't.
inject_check("handleInvestigations")

open('backend/internal/api/api.go', 'w').write(content)
