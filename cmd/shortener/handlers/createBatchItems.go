package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
)

// CreateBatchUrls creates for multiple url request shorten url
func (h *Handler) CreateBatchUrls(w http.ResponseWriter, r *http.Request) {
	// unmarshall request items
	var storageItems []storage.StorageItemWithCorrelationID
	if err := json.NewDecoder(r.Body).Decode(&storageItems); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	// convert to interface (why is this needed?)
	var optionItems []storage.StorageItemOptionsInterface
	for i := range storageItems {
		err := storage.ValidateStorageItemWithCorrelationIDRequest(&storageItems[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		optionItems = append(optionItems, &storageItems[i])
	}

	// send to storage interface
	if err := h.Storage.MultiStoreURL(&optionItems); err != nil {
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
