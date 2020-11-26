package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func scheduleHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		getHandler(writer)
		break
	case http.MethodPut:
		addOrUpdateHandler(writer, request)
		break
	case http.MethodDelete:
		deleteHandler(writer, request)
		break
	}
}

func getHandler(writer http.ResponseWriter) {
	tasks, err := GetTasksFromDB(db)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusOK)
	responseData, err := json.Marshal(tasks)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	writer.Write(responseData)
}

func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	var task Task
	err := json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		fmt.Println("Error while decoding json")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(task.ID)
	err = DeleteTask(db, task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	UpdateSchedule(db, scheduler)
}

func addOrUpdateHandler(writer http.ResponseWriter, request *http.Request) {
	var settings Task
	err := json.NewDecoder(request.Body).Decode(&settings)
	if err != nil {
		fmt.Println("Error while decoding json")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = SaveOrUpdate(db, settings)
	if err != nil {
		fmt.Println("Error while saving settings")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	UpdateSchedule(db, scheduler)
	writer.WriteHeader(http.StatusOK)
}
