package calendar_service

import (
	"context"
	"logispro/internal/db"
	"logispro/internal/web/requests"
)

type CalendarService struct {
	Queries *db.Queries
}

func (s *CalendarService) Create(ctx context.Context, req requests.CreateCalendarEventRequest) (db.CalendarEvent, error) {
	return s.Queries.CreateCalendar(ctx, db.CreateCalendarParams{
		UserID:  req.UserId,
		Title:   req.Title,
		Content: req.Content,
		ForDate: req.ForDate,
	})
}
func (s *CalendarService) Update(ctx context.Context, req requests.CreateCalendarEventRequest) {}
func (s *CalendarService) Delete(ctx context.Context, id int64, userId int64) error {
	return s.Queries.DeleteCalendar(ctx, db.DeleteCalendarParams{
		UserID: userId,
		ID:     id,
	})
}
func (s *CalendarService) All(ctx context.Context, userId int64) ([]db.CalendarEvent, error) {
	return s.Queries.GetUserCalendars(ctx, userId)
}
