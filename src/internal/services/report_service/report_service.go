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
func (s *ReportService) Update(req requests.UpdateReportRequest, ctx context.Context, rootId *int64, agencyUsers []int64) error {
	if rootId == nil {
		return s.Queries.UpdateReportByMaster(ctx, db.UpdateReportByMasterParams{
			ID:      req.ID,
			UsersID: agencyUsers,
			Title:   req.Title,
			Content: req.Content,
		})
	}
	return s.Queries.UpdateReport(ctx, db.UpdateReportParams{
		ID:      req.ID,
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	})
}

func (s *ReportService) Delete(ctx context.Context, id int64, userId int64, rootId *int64, agencyUsers []int64) error {
	if rootId == nil {
		s.Queries.DeleteReportByMaster(ctx, db.DeleteReportByMasterParams{
			ID:      id,
			UsersID: agencyUsers,
		})
	}
	return s.Queries.DeleteReport(ctx, db.DeleteReportParams{
		ID:     id,
		UserID: userId,
	})
}

func (s *ReportService) All(ctx context.Context, userId int64, rootId *int64, agencyUsers []int64) ([]db.Report, error) {
	if rootId == nil {
		return s.Queries.GetUserMasterReports(ctx, agencyUsers)
	}
	return s.Queries.GetUserReports(ctx, userId)
}
