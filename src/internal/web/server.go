package web

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"logispro/internal/config"
	"logispro/internal/constants"
	"logispro/internal/services/payment_service"
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

type ApiHandlerFunc func(w http.ResponseWriter, r *http.Request) response_types.Responder

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
							"error": fmt.Sprintf("failed to parse multipart form: %s", err),
						},
					})
					return
				}
			}
		}

		resp := handler(w, r)

		// Handle FileResponse separately - it manages its own headers
		if fileResp, ok := resp.(response_types.FileResponse); ok {
			// Only set Content-Disposition, let ServeFile handle Content-Type
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileResp.Name))

			// Handle errors for file responses
			if resp.GetError() != nil {
				if errors.Is(resp.GetError(), sql.ErrNoRows) {
					http.NotFound(w, r)
				} else {
					http.Error(w, "File not found", resp.GetStatusCode())
				}
				return
			}

			// Serve the file - this handles all headers and status codes
			http.ServeFile(w, r, fmt.Sprintf("uploads/%s", fileResp.Name))
			return
		}

		// Handle all other response types (JSON responses)
		w.Header().Set("Content-Type", "application/json")

		if errors.Is(resp.GetError(), sql.ErrNoRows) {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(resp.GetStatusCode())
		}

		if resp.GetError() != nil {
			if errors.Is(resp.GetError(), sql.ErrNoRows) {
				json.NewEncoder(w).Encode(map[string]map[string]any{
					"data": {
						"error": "not found",
					},
				})
			} else {
				if resp.GetStatusCode() == http.StatusInternalServerError {
					s.Logger.Error("error", slog.Any("content", resp.GetError()))
					json.NewEncoder(w).Encode(map[string]map[string]any{
						"data": {
							"error": "internal server error",
						},
					})
				} else {
					json.NewEncoder(w).Encode(map[string]map[string]any{
						"data": {
							"error": resp.GetError(),
						},
					})
				}
			}
		} else {
			switch v := resp.(type) {
			case response_types.ApiResponse:
				json.NewEncoder(w).Encode(map[string]any{
					"data": v.Content,
				})
			default:
				fmt.Println("Unknown response type")
			}
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
	mux.Handle("PATCH /user", AuthMiddleware(OwnerGuardMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.UserController.UpdateUserHandler), s))))
	mux.Handle("POST /contact", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ContactController.CreateContactHandler), s)))
	mux.Handle("GET /contact", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ContactController.GetContactsHandler), s)))
	mux.Handle("GET /contacts", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ContactController.GetContactsListHandler), s)))
	mux.Handle("GET /count-contacts", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ContactController.CountContactsHandler), s)))
	mux.Handle("GET /contact/{id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ContactController.GetContactHandler), s)))
	mux.Handle("DELETE /contact/{id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ContactController.DeleteContactHandler), s)))
	mux.Handle("POST /building", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.CreateBuildingHandler), s)))
	mux.Handle("GET /building", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetBuildingsHandler), s)))
	mux.Handle("GET /buildings-statistics", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetBuildingsStatisticsHandler), s)))
	mux.Handle("GET /buildings-gain", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetBuildingsGainHandler), s)))
	mux.Handle("GET /building/{id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetBuildingHandler), s)))
	mux.Handle("PATCH /building/{id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.UpdateBuildingHandler), s)))
	mux.Handle("DELETE /building/{id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.DeleteBuildingHandler), s)))
	mux.Handle("POST /building/{id}/images", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.CreateBuildingImagesHandler), s)))
	mux.Handle("DELETE /building/{id}/images/{imageId}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.DeleteBuildingImageHandler), s)))
	mux.Handle("POST /building/{id}/documents", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.CreateBuildingDocumentsHandler), s)))
	mux.Handle("DELETE /building/{id}/documents/{documentId}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.DeleteBuildingDocumentHandler), s)))
	mux.Handle("POST /building-vue/{id}", s.makeHttpHandler(s.Controller.BuildingController.AddVueHandler))
	mux.Handle("POST /sms", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.SmsController.CreateSmsHandler), s)))
	mux.Handle("GET /agency-users", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.UserController.GetAgencyUsers), s)))
	mux.Handle("POST /task", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.TaskController.CreateTaskHandler), s)))
	mux.Handle("GET /tasks", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.TaskController.GetTasksHandler), s)))
	mux.Handle("PATCH /task/{id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.TaskController.UpdateTaskHandler), s)))
	mux.Handle("POST /report", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ReportController.CreateReportHandler), s)))
	mux.Handle("PATCH /report/{id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ReportController.UpdateReportHandler), s)))
	mux.Handle("DELETE /report/{id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ReportController.DeleteReportHandler), s)))
	mux.Handle("GET /reports", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ReportController.GetReportsHandler), s)))
	mux.Handle("POST /calendar_events", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.CalendarController.CreateCalendarEventHandler), s)))
	mux.Handle("DELETE /calendar_events/{id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.CalendarController.DeleteCalendarEventHandler), s)))
	mux.Handle("GET /calendar_events", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.CalendarController.GetCalendarEventsHandler), s)))
	mux.Handle("GET /get-building-recommendations/{building_id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.RecommenderController.GetForBuildingHandler), s)))
	mux.Handle("GET /get-contact-recommendations/{contact_id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.RecommenderController.GetForContactsHandler), s)))
	mux.Handle("GET /buildings-daira-distributions", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetDairaDistributionHandler), s)))
	mux.Handle("GET /buildings-map", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.BuildingController.GetMapHandler), s)))
	mux.Handle("POST /share/{id}/{type}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.ShareableController.Share), s)))
	mux.Handle("POST /document", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.DocumentController.CreateDocumentHandler), s)))
	mux.Handle("GET /documents", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.DocumentController.GetDocumentsHandler), s)))
	mux.Handle("DELETE /document/{id}", AuthMiddleware(SubsribedMiddleware(s.makeHttpHandler(s.Controller.DocumentController.DeleteDocumentHandler), s)))
	mux.Handle("GET /d/{token}", s.makeHttpHandler(s.Controller.DocumentController.DownloadDocumentHandler))
	mux.Handle("GET /b/{token}", s.makeHttpHandler(s.Controller.BuildingController.GetSharedBuilding))
	mux.Handle("POST /payment-link", AuthMiddleware(s.makeHttpHandler(s.Controller.SubscriptionController.CreateCheckoutLink)))
	mux.Handle("POST /subscription-cancel", AuthMiddleware(s.makeHttpHandler(s.Controller.SubscriptionController.Cancel)))
	mux.Handle("POST /chargily-webhhook", s.makeHttpHandler(s.Controller.SubscriptionController.ChargilyWebhook))

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

func SubsribedMiddleware(next http.Handler, s *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, ok := ctx.Value(constants.UserIDContextKey).(int64)
		if !ok {
			http.Error(w, "failed to get user id from context", http.StatusUnauthorized)
			return
		}
		subscriptionStatus, err := s.Controller.SubscriptionController.GetStatus(ctx, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if subscriptionStatus != payment_service.SUBS_STATUS_ACTIVE && subscriptionStatus != payment_service.SUBS_STATUS_CANCELLED {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]string{
				"data": "subscription ended",
			})
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
