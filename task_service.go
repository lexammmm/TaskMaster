package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

var tasks = []Task{}

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
    router.HandleFunc("/api/tasks", createTask).Methods("POST")
    router.HandleFunc("/api/tasks", getTasks).Methods("GET")
    router.HandleFunc("/api/tasks/{id}", getTask).Methods("GET")
    router.HandleFunc("/api/tasks/{id}", updateTask).Methods("PUT")
    router.HandleFunc("/api/tasks/{id}", deleteTask).Methods("DELETE")

    log.Printf("Server starting on port %s", port)
    http.ListenAndServe(":"+port, router)
}

func createTask(w http.ResponseWriter, r *http.Request) {
    var task Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    tasks = append(tasks, task)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(task)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    for _, task := range tasks {
        if task.ID == vars["id"] {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(task)
            return
        }
    }
    http.NotFound(w, r)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    for i, task := range tasks {
        if task.ID == vars["id"] {
            var updatedTask Task
            err := json.NewDecoder(r.Body).Decode(&updatedTask)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            tasks[i] = updatedTask
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(updatedTask)
            return
        }
    }
    http.NotFound(w, r)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    for i, task := range tasks {
        if task.ID == vars["id"] {
            tasks = append(tasks[:i], tasks[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    http.NotFound(w, r)
}