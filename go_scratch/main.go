package main

import ("fmt"
"log"
"os"
"net/http"
"github.com/joho/godotenv"
"github.com/go-chi/chi"
"github.com/go-chi/cors"
_ "github.com/lib/pq"
  "github.com/wagslane/rssagg/internal/database"
)


type apiConfig struct{
	DB *database.Queries
}

func main(){
	fmt.Println("cool")

	godotenv.Load(".env")
	port :=os.Getenv("PORT")
	if port ==""{
		log.Fatal("PORT is not found in the env")
	}

	dbUrl :=os.Getenv("DB_URL")
	if dbUrl ==""{
		log.Fatal("dbUrl is not found in the env")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	router :=chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router :=chi.NewRouter()
	v1Router.Post("/users", apiCfg.handlerUsersCreate)
	v1Router.HandleFunc("/ready",handlerReadiness)

	router.Mount("/v1",v1Router)
	srv :=&http.Server{
		Handler: router,
		Addr: ":"+port,
	}
	log.Printf(" Server starting on port %v", port)
	err :=srv.ListenAndServe()
	if err!=nil {
		log.Fatal(err)
	}

	fmt.Println(":Port", port)
}