package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

type Task struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Project     string   `json:"project"`
    AssignedTo  []string `json:"assignedTo"`
}

var taskList = []Task{}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

    router := mux.NewRouter()
    router.Use(loggingMiddleware)
    setupRoutes(router)

    log.Printf("Server starting on port %s", port)
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatalf("Error starting server: %s", err)
    }
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
    newTask.ID = uuid.New().String()
    taskList = append(taskList, newTask)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(newTask)
}

func listTasksHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(taskList)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    for _, task := range taskList {
        if task.ID == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(task)
            return
        }
    }
    http.NotFound(w, r)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    for i, task := range taskList {
        if task.ID == id {
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
    id := mux.Vars(r)["id"]
    for i, task := range taskList {
        if task.ID == id {
            taskList = append(taskList[:i], taskList[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    http.NotFound(w, r)
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.RequestURI)
        next.ServeHTTP(w, r)
    })
}