package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tmlamb/document-api/cmd/document-api/model"
	"github.com/tmlamb/document-api/cmd/document-api/repository"
)

// UploadDocument ...
func UploadDocument(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: uploadDocument")
	vars := mux.Vars(r)
	accountId, err := uuid.Parse(vars["accountId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}

	var document model.Document

	document.Filename = uuid.New()
	document.AccountID = accountId

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./"+document.Filename.String(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	newDocument, err := repository.SaveDocument(document)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	err = encoder.Encode(newDocument)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}

// DownloadDocument ...
func DownloadDocument(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: downloadDocument")
	vars := mux.Vars(r)
	accountId, err := uuid.Parse(vars["accountId"])
	if err != nil {
		log.Printf("Error parsing accountId: %v\n", err)
		w.WriteHeader(400)
		return
	}
	documentId, err := uuid.Parse(vars["documentId"])
	if err != nil {
		log.Printf("Error parsing documentId: %v\n", err)
		w.WriteHeader(400)
		return
	}

	document, err := repository.FindOneDocument(documentId, accountId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fmt.Sprintf("%x", document.DocumentID))
	w.Header().Set("Content-Type", "application/octet-stream")
	log.Println("Endpoint Hit: downloadDocument: " + document.Filename.String())
	http.ServeFile(w, r, "./"+document.Filename.String())
}

// ReturnDocumentsIndex ...
func ReturnDocumentsIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: returnDocuments")
	vars := mux.Vars(r)
	accountId, err := uuid.Parse(vars["accountId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	documents, err := repository.FindAllDocuments(accountId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	err = encoder.Encode(documents)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}

// UpdateDocument ...
func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: updateDocument")
	vars := mux.Vars(r)
	accountId, err := uuid.Parse(vars["accountId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	documentId, err := uuid.Parse(vars["documentId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var document model.Document
	err = json.Unmarshal(reqBody, &document)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	document.DocumentID = documentId
	document.AccountID = accountId

	err = repository.UpdateDocument(document)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}

// DeleteDocument ...
func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: deleteDocument")
	vars := mux.Vars(r)
	accountId, err := uuid.Parse(vars["accountId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	documentId, err := uuid.Parse(vars["documentId"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}

	err = repository.DeleteDocument(documentId, accountId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	// Just a reminder that the filesystem driven approach isn't cleaning up. Store these in postgres
	// err = os.Remove("./" + documentToDelete.Filename.String())
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(500)
	// 	return
	// }
}
