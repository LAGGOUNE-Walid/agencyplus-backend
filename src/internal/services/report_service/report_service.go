package report_service

import (
	"context"
	"logispro/internal/db"
	"logispro/internal/web/requests"
)

type ReportService struct {
	Queries *db.Queries
}

func (s *ReportService) Create(req requests.CreateReportRequest, ctx context.Context) (report db.Report, err error) {
	return s.Queries.CreateReport(ctx, db.CreateReportParams{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	})
}
func (s *ReportService) Update(req requests.UpdateReportRequest, ctx context.Context) error {
	return s.Queries.UpdateReport(ctx, db.UpdateReportParams{
		ID:      req.ID,
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	})
}

func (s *ReportService) Delete(ctx context.Context, id int64, userId int64) error {
	return s.Queries.DeleteReport(ctx, db.DeleteReportParams{
		ID:     id,
		UserID: userId,
	})
}

func (s *ReportService) All(ctx context.Context, userId int64) ([]db.Report, error) {
	return s.Queries.GetUserReports(ctx, userId)
}
