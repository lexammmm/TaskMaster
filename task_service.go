package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

var taskList = []Task{}

type Task struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Project     string   `json:"project"`
    AssignedTo  []string `json:"assignedTo"`
}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

    router := mux.NewRouter()
    setupRoutes(router)

    log.Printf("Server starting on port %s", port)
    http.ListenAndServe(":"+port, router)
}

func setupRoutes(router *mux.Router) {
    router.HandleFunc("/api/tasks", createTaskHandler).Methods("POST")
    router.HandleFunc("/api/tasks", listTasksHandler).Methods("GET")
    router.HandleFunc("/api/tasks/{id}", getTaskHandler).Methods("GET")
    router.HandleFunc("/api/tasks/{id}", updateTaskHandler).Methods("PUT")
    router.HandleFunc("/api/tasks/{id}", deleteTaskHandler).Methods("DELETE")
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
    var newTask Task
    err := json.NewDecoder(r.Body).Decode(&newTask)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    taskList = append(taskList, newTask)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(newTask)
}

func listTasksHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(taskList)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    for _, task := range taskList {
        if task.ID == vars["id"] {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(task)
            return
        }
    }
    http.NotFound(w, r)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    for i, task := range taskList {
        if task.ID == vars["id"] {
            var updatedTask Task
            err := json.NewDecoder(r.Body).Decode(&updatedTask)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            updatedTask.ID = task.ID
            taskList[i] = updatedTask
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(updatedTask)
            return
        }
    }
    http.NotFound(w, r)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    for i, task := range taskList {
        if task.ID == vars["id"] {
            taskList = append(taskList[:i], taskList[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    http.NotFound(w, r)
}