package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
)

// CreateBatchUrls creates for multiple url request shorten url
func (h *Handler) CreateBatchUrls(w http.ResponseWriter, r *http.Request) {
	// unmarshall request items
	var storageItems []storage.StorageItemWithCorrelationId
	if err := json.NewDecoder(r.Body).Decode(&storageItems); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	// convert to interface (why?)
	var optionItems []storage.StorageItemOptionsInterface
	for i := range storageItems {
		optionItems = append(optionItems, &storageItems[i])
	}

	// send to storage interface
	if err := h.Storage.MultiStoreUrl(&optionItems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// updates short url format
	for ix := range storageItems {
		storageItems[ix].Shorten = h.Conf.HostLocation + "/" + storageItems[ix].Shorten
	}

	// marshal
	response, err := json.Marshal(storageItems)
	if err != nil {
		http.Error(w, "cant generate response", http.StatusInternalServerError)
		return
	}

	// send
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
