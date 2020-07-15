package kasset

import (
	"net/http"

	"git.kanosolution.net/kano/kaos"
)

var (
	Event kaos.EventHub
	Topic string

	e error
)

func (ae *AssetEngine) HttpViewer(w http.ResponseWriter, r *http.Request) {
	assetID := r.URL.Query().Get("id")
	_ = r.URL.Query().Get("t") == "yes"
	dl := r.URL.Query().Get("t") == "dl"

	if Event == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("EventHub is not initialized"))
	}

	ast := new(Asset)
	if e = Event.Publish(Topic, assetID, ast); e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
		return
	}

	content, e := ae.fs.Read(ast.URI)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
		return
	}

	if dl {
		w.Header().Set("Content-Disposition", "attachment; filename=\""+ast.OriginalFileName+"\"")
	}
	w.Header().Set("Content-Type", ast.ContentType)
	w.Write(content)
}
