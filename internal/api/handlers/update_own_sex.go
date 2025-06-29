package handlers

import (
	"context"

	"github.com/chains-lab/elector-cab-svc/internal/api/responses"
	svc "github.com/chains-lab/proto-storage/gen/go/electorcab"
	"github.com/google/uuid"
)

func (s Service) UpdateOwnSex(ctx context.Context, req *svc.UpdateOwnSexRequest) (*svc.Biography, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	bio, err := s.app.UpdateSex(ctx, meta.InitiatorID, req.Sex)
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("failed to update user sex")

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.Biography(bio), nil
}
