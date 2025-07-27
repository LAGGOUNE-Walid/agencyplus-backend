package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"logispro/internal/config"
	"logispro/internal/constants"
	"logispro/internal/shared/response_types"
	"logispro/internal/web/controllers"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Server struct {
	Domain     string
	Logger     *slog.Logger
	Controller controllers.Controller
}

type ApiHandlerFunc func(w http.ResponseWriter, r *http.Request) response_types.ApiResponse

func (s *Server) makeHttpHandler(handler ApiHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch || r.Method == http.MethodDelete {
			if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				r.Body = http.MaxBytesReader(w, r.Body, 500<<20) // 500mb
				err := r.ParseMultipartForm(100 << 20)           // 100 mb
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(400)

					json.NewEncoder(w).Encode(map[string]map[string]any{
						"data": {
							"error": fmt.Sprintf("failed to parse multipart form: %w", err),
						},
					})
					return
				}
			}
		}
		resp := handler(w, r)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)

		if resp.Error != nil {
			if resp.StatusCode == http.StatusInternalServerError {
				s.Logger.Error("error", slog.Any("content", resp.Error))
				json.NewEncoder(w).Encode(map[string]map[string]any{
					"data": {
						"error": "internal server error",
					},
				})
			} else {
				json.NewEncoder(w).Encode(map[string]map[string]any{
					"data": {
						"error": resp.Error.Error(),
					},
				})
			}
		} else {
			json.NewEncoder(w).Encode(map[string]any{
				"data": resp.Content,
			})
		}
	}

}
func (s *Server) LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			s.Logger.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote", r.RemoteAddr),
				slog.Duration("duration ms", time.Duration(time.Since(start).Milliseconds())),
			)
		})
	}
}

func RecoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					s := debug.Stack()
					logger.Error("handler panic recovered",
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
						slog.String("remote", r.RemoteAddr),
						slog.Any("error", rec),
						slog.Any("stack", string(s)),
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					json.NewEncoder(w).Encode(map[string]any{
						"data": map[string]string{
							"error": "internal server error",
						},
					})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func NewServer(domain string, logger *slog.Logger, mainController controllers.Controller) Server {
	return Server{
		Domain:     domain,
		Logger:     logger,
		Controller: mainController,
	}
}

func (s *Server) Run() {

	mux := http.NewServeMux()

	// Routes
	mux.Handle("POST /user", s.makeHttpHandler(s.Controller.UserController.CreateUserHandler))
	mux.Handle("POST /auth", s.makeHttpHandler(s.Controller.UserController.Auth))
	mux.Handle("PATCH /user", AuthMiddleware(OwnerGuardMiddleware(s.makeHttpHandler(s.Controller.UserController.UpdateUserHandler))))
	mux.Handle("POST /contact", AuthMiddleware(s.makeHttpHandler(s.Controller.ContactController.CreateContactHandler)))
	mux.Handle("GET /contact", AuthMiddleware(s.makeHttpHandler(s.Controller.ContactController.GetContactsHandler)))
	mux.Handle("GET /contacts", AuthMiddleware(s.makeHttpHandler(s.Controller.ContactController.GetContactsListHandler)))
	mux.Handle("GET /count-contacts", AuthMiddleware(s.makeHttpHandler(s.Controller.ContactController.CountContactsHandler)))
	mux.Handle("GET /contact/{id}", AuthMiddleware(s.makeHttpHandler(s.Controller.ContactController.GetContactHandler)))
	mux.Handle("DELETE /contact/{id}", AuthMiddleware(s.makeHttpHandler(s.Controller.ContactController.DeleteContactHandler)))
	mux.Handle("POST /building", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.CreateBuildingHandler)))
	mux.Handle("GET /building", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetBuildingsHandler)))
	mux.Handle("GET /buildings-statistics", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetBuildingsStatisticsHandler)))
	mux.Handle("GET /buildings-gain", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetBuildingsGainHandler)))
	mux.Handle("GET /building/{id}", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetBuildingHandler)))
	mux.Handle("PATCH /building/{id}", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.UpdateBuildingHandler)))
	mux.Handle("DELETE /building/{id}", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.DeleteBuildingHandler)))
	mux.Handle("POST /building/{id}/images", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.CreateBuildingImagesHandler)))
	mux.Handle("DELETE /building/{id}/images/{imageId}", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.DeleteBuildingImageHandler)))
	mux.Handle("POST /building/{id}/documents", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.CreateBuildingDocumentsHandler)))
	mux.Handle("DELETE /building/{id}/documents/{documentId}", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.DeleteBuildingDocumentHandler)))
	mux.Handle("POST /building-vue/{id}", s.makeHttpHandler(s.Controller.BuildingController.AddVueHandler))
	mux.Handle("POST /sms", AuthMiddleware(s.makeHttpHandler(s.Controller.SmsController.CreateSmsHandler)))
	mux.Handle("POST /task", AuthMiddleware(s.makeHttpHandler(s.Controller.TaskController.CreateTaskHandler)))
	mux.Handle("GET /tasks", AuthMiddleware(s.makeHttpHandler(s.Controller.TaskController.GetTasksHandler)))
	mux.Handle("PATCH /task/{id}", AuthMiddleware(s.makeHttpHandler(s.Controller.TaskController.UpdateTaskHandler)))
	mux.Handle("POST /report", AuthMiddleware(s.makeHttpHandler(s.Controller.ReportController.CreateReportHandler)))
	mux.Handle("PATCH /report/{id}", AuthMiddleware(s.makeHttpHandler(s.Controller.ReportController.UpdateReportHandler)))
	mux.Handle("DELETE /report/{id}", AuthMiddleware(s.makeHttpHandler(s.Controller.ReportController.DeleteReportHandler)))
	mux.Handle("GET /reports", AuthMiddleware(s.makeHttpHandler(s.Controller.ReportController.GetReportsHandler)))
	mux.Handle("POST /calendar_events", AuthMiddleware(s.makeHttpHandler(s.Controller.CalendarController.CreateCalendarEventHandler)))
	mux.Handle("DELETE /calendar_events/{id}", AuthMiddleware(s.makeHttpHandler(s.Controller.CalendarController.DeleteCalendarEventHandler)))
	mux.Handle("GET /calendar_events", AuthMiddleware(s.makeHttpHandler(s.Controller.CalendarController.GetCalendarEventsHandler)))
	mux.Handle("GET /get-building-recommendations/{building_id}", AuthMiddleware(s.makeHttpHandler(s.Controller.RecommenderController.GetForBuildingHandler)))
	mux.Handle("GET /get-contact-recommendations/{contact_id}", AuthMiddleware(s.makeHttpHandler(s.Controller.RecommenderController.GetForContactsHandler)))
	mux.Handle("GET /buildings-daira-distributions", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetDairaDistributionHandler)))
	mux.Handle("GET /buildings-map", AuthMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetMapHandler)))

	handler := RecoveryMiddleware(s.Logger)(s.LoggingMiddleware(s.Logger)(mux))
	s.Logger.Info("starting server ", slog.String("domain ", s.Domain))
	if err := http.ListenAndServe(s.Domain, handler); err != nil {
		s.Logger.Error("server failed", slog.String("error", err.Error()))
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token missing", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return config.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims[constants.UserIDContextKey].(float64) // float64 because of JSON
		if !ok {
			http.Error(w, "user_id missing in token", http.StatusUnauthorized)
			return
		}
		role, ok := claims[constants.UserRoleContextKey].(float64)
		if !ok {
			http.Error(w, "role missing in token", http.StatusUnauthorized)
			return
		}

		var rootIdPtr *int64
		if rawRootId, exists := claims[constants.UserRootContextKey]; exists && rawRootId != nil {
			if floatVal, ok := rawRootId.(float64); ok {
				id := int64(floatVal)
				rootIdPtr = &id
			} else {
				http.Error(w, "invalid root_id in token", http.StatusUnauthorized)
				return
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.UserRootContextKey, rootIdPtr)
		ctx = context.WithValue(ctx, constants.UserIDContextKey, int64(userID))
		ctx = context.WithValue(ctx, constants.UserRoleContextKey, int64(role))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OwnerGuardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		role, ok := ctx.Value(constants.UserRoleContextKey).(int64)
		if !ok {
			http.Error(w, "role missing in context", http.StatusUnauthorized)
			return
		}
		if role != constants.ROLE_OWENER {
			http.Error(w, "this endpoint is gaurded with owner role", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
