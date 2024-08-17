package router

import (
	"fmt"
	"gian/apply"
	"gian/jobs"
	"gian/middleware"
	"gian/saver"
	"gian/user"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	// Writing a simple response to the client
	fmt.Fprintf(w, "Hello, World!")
}

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Create a new CORS middleware handler
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "software", "source-type", "email", "role", "name", "company", "platform", "code", "user_id", "uuid", "job_id"},
		AllowCredentials: true,
	})

	r.Use(corsHandler.Handler)

	// Define your routes here
	// POST APIS

	r.Get("/", homePage)

	// USER ROUTES LATER WE MOVE THESE TO DIFFERENT SERVICE
	r.Post("/signup", middleware.CorsMiddleware(user.CreateUser))
	r.Post("/login", middleware.CorsMiddleware(user.GetUser))
	r.Post("/request-otp", middleware.CorsMiddleware(user.GetOtp))
	r.Post("/verify-otp", middleware.CorsMiddleware(user.VerifyOtp))
	r.Post("/create-job", middleware.CorsMiddleware(jobs.CreateJob))

	r.Get("/all-jobs", middleware.CorsMiddleware(jobs.GetJobsList))

	r.Post("/saver", middleware.CorsMiddleware(saver.SaveJob))
	r.Get("/saver", middleware.CorsMiddleware(saver.GetSavedJobs))
	r.Post("/unsave", middleware.CorsMiddleware(saver.RemoveSavedJob))

	r.Post("/apply", middleware.CorsMiddleware(apply.SaveApplication))
	r.Get("/apply", middleware.CorsMiddleware(apply.GetApplications))

	r.Post("/myjobs", middleware.CorsMiddleware(jobs.UpdateJobStatus))
	r.Get("/myjobs", middleware.CorsMiddleware(jobs.GetJobs))

	return r
}
