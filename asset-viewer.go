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

func Viewer(w http.ResponseWriter, r *http.Request) {
	assetID := r.URL.Query().Get("assetid")
	if Event == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("EventHub is not initialized"))
	}

	ad := new(AssetData)
	if e = Event.Publish(Topic, assetID, ad); e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+ad.Asset.OriginalFileName+"\"")
	w.Header().Set("Content-Type", ad.Asset.ContentType)
	w.Write(ad.Content)
}
