package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    _ "github.com/go-sql-driver/mysql" // Importar el driver de MySQL
)

// Estructura para la serie
type Serie struct {
    ID               int    `json:"id"`
    Title            string `json:"title"`
    Status           string `json:"status"`
    LastEpisodeWatched int    `json:"lastEpisodeWatched"`
    TotalEpisodes    int    `json:"totalEpisodes"`
    Ranking          int    `json:"ranking"`
}

// Estructura para el mensaje en JSON
type Message struct {
    Message string `json:"message"`
}

// Variable global para la conexión a la base de datos
var db *sql.DB

// Función para inicializar la base de datos
func initDB() {
    var err error
    db, err = sql.Open("mysql", "user:userpassword@tcp(localhost:3306)/series_tracker")
    if err != nil {
        log.Fatal(err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Conexión exitosa a la base de datos")

    // Verificar si ya existen datos en la tabla
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM series").Scan(&count)
    if err != nil {
        log.Fatal("Error al verificar datos en la base de datos:", err)
    }

    if count == 0 {
        fmt.Println("Insertando datos iniciales...")
        insertInitialData()
    } else {
        fmt.Println("La base de datos ya tiene datos, no es necesario insertar.")
    }
}

// Función para insertar datos iniciales
func insertInitialData() {
    initialSeries := []Serie{
        {Title: "Breaking Bad", Status: "Completed", LastEpisodeWatched: 62, TotalEpisodes: 62, Ranking: 10},
        {Title: "Attack on Titan", Status: "Ongoing", LastEpisodeWatched: 87, TotalEpisodes: 87, Ranking: 9},
        {Title: "Stranger Things", Status: "Ongoing", LastEpisodeWatched: 34, TotalEpisodes: 34, Ranking: 8},
        {Title: "Game of Thrones", Status: "Completed", LastEpisodeWatched: 73, TotalEpisodes: 73, Ranking: 7},
    }

    query := "INSERT INTO series (title, status, last_episode_watched, total_episodes, ranking) VALUES (?, ?, ?, ?, ?)"
    for _, serie := range initialSeries {
        _, err := db.Exec(query, serie.Title, serie.Status, serie.LastEpisodeWatched, serie.TotalEpisodes, serie.Ranking)
        if err != nil {
            log.Println("Error insertando datos iniciales:", err)
        }
    }
    fmt.Println("Datos iniciales insertados correctamente.")
}

func main() {
    initDB() // Llamar a la función de inicialización de la base de datos

    r := mux.NewRouter()

    // Ruta principal ("/") que retorna un mensaje en JSON
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json") // Asegura que la respuesta sea en JSON
        message := Message{"Hello, World!"}                // Mensaje
        json.NewEncoder(w).Encode(message)                  // Codifica el mensaje como JSON y lo envía al cliente
    })

    // Definir rutas para las solicitudes
    r.HandleFunc("/api/series", getAllSeries).Methods("GET")
    r.HandleFunc("/api/series/{id}", getSerieByID).Methods("GET")
    r.HandleFunc("/api/series", crearSerie).Methods("POST")
    r.HandleFunc("/api/series/{id}", updateSerie).Methods("PUT")
    r.HandleFunc("/api/series/{id}", deleteSerie).Methods("DELETE")

	//r.HandleFunc("/api/series/{id}/upvote", upvoteSerie).Methods("PATCH")
    //r.HandleFunc("/api/series/{id}/downvote", downvoteSerie).Methods("PATCH")

    // Configurar CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"}, // Origen permitido (tu frontend)
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Métodos permitidos
        AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Cabeceras permitidas
        AllowCredentials: true,
    })

    // Aplicar el middleware CORS
    handler := c.Handler(r)

    // Iniciar el servidor
    log.Println("Servidor corriendo en http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}

// getAllSeries maneja las solicitudes GET para obtener todas las series
func getAllSeries(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        rows, err := db.Query("SELECT id, title, status, last_episode_watched, total_episodes, ranking FROM series")
        if err != nil {
            http.Error(w, "Error al obtener las series", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var series []Serie
        for rows.Next() {
            var serie Serie
            if err := rows.Scan(&serie.ID, &serie.Title, &serie.Status, &serie.LastEpisodeWatched, &serie.TotalEpisodes, &serie.Ranking); err != nil {
                http.Error(w, "Error al leer los datos de la base de datos", http.StatusInternalServerError)
                return
            }
            series = append(series, serie)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(series)
    } else {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}

// crearSerie maneja las solicitudes POST para crear una nueva serie
func crearSerie(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        var nuevaSerie Serie

        err := json.NewDecoder(r.Body).Decode(&nuevaSerie)
        if err != nil {
            http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
            return
        }

        query := "INSERT INTO series (title, status, last_episode_watched, total_episodes, ranking) VALUES (?, ?, ?, ?, ?)"
        result, err := db.Exec(query, nuevaSerie.Title, nuevaSerie.Status, nuevaSerie.LastEpisodeWatched, nuevaSerie.TotalEpisodes, nuevaSerie.Ranking)
        if err != nil {
            http.Error(w, "Error al crear la serie", http.StatusInternalServerError)
            return
        }

        id, err := result.LastInsertId()
        if err != nil {
            http.Error(w, "Error al obtener el ID de la nueva serie", http.StatusInternalServerError)
            return
        }

        nuevaSerie.ID = int(id)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(nuevaSerie)
    } else {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}

// deleteSerie maneja las solicitudes DELETE para eliminar una serie por ID
func deleteSerie(w http.ResponseWriter, r *http.Request) {
    if r.Method == "DELETE" {
        vars := mux.Vars(r)
        id := vars["id"]

        query := "DELETE FROM series WHERE id = ?"
        result, err := db.Exec(query, id)
        if err != nil {
            http.Error(w, "Error al eliminar la serie", http.StatusInternalServerError)
            return
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil || rowsAffected == 0 {
            http.Error(w, "Serie no encontrada", http.StatusNotFound)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    } else {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}

// getSerieByID maneja las solicitudes GET para obtener una serie por ID
func getSerieByID(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        vars := mux.Vars(r)
        id := vars["id"]

        row := db.QueryRow("SELECT id, title, status, last_episode_watched, total_episodes, ranking FROM series WHERE id = ?", id)

        var serie Serie
        if err := row.Scan(&serie.ID, &serie.Title, &serie.Status, &serie.LastEpisodeWatched, &serie.TotalEpisodes, &serie.Ranking); err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "Serie no encontrada", http.StatusNotFound)
            } else {
                http.Error(w, "Error al obtener la serie", http.StatusInternalServerError)
            }
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(serie)
    } else {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}

// updateSerie maneja las solicitudes PUT para actualizar una serie existente
func updateSerie(w http.ResponseWriter, r *http.Request) {
    if r.Method == "PUT" {
        vars := mux.Vars(r)
        id := vars["id"]

        var serie Serie
        err := json.NewDecoder(r.Body).Decode(&serie)
        if err != nil {
            http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
            return
        }

        query := "UPDATE series SET title = ?, status = ?, last_episode_watched = ?, total_episodes = ?, ranking = ? WHERE id = ?"
        _, err = db.Exec(query, serie.Title, serie.Status, serie.LastEpisodeWatched, serie.TotalEpisodes, serie.Ranking, id)
        if err != nil {
            http.Error(w, "Error al actualizar la serie", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(serie)
    } else {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
    }
}
