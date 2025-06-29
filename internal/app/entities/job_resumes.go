package entities

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
	"github.com/chains-lab/elector-cab-svc/internal/app/domain"
	"github.com/chains-lab/elector-cab-svc/internal/app/enums"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/elector-cab-svc/internal/dbx"
	"github.com/google/uuid"
)

type JobResumeQ interface {
	New() dbx.JobResumesQ

	Insert(ctx context.Context, input dbx.JobResumeModel) error
	Update(ctx context.Context, input dbx.UpdateJobInput) error
	Select(ctx context.Context) ([]dbx.JobResumeModel, error)
	Get(ctx context.Context) (dbx.JobResumeModel, error)
	Delete(ctx context.Context) error

	FilterUserID(userID uuid.UUID) dbx.JobResumesQ

	Page(limit, offset uint64) dbx.JobResumesQ
	Count(ctx context.Context) (int, error)
}

type JobResumes struct {
	queries JobResumeQ
}

func NewJobResumes(db *sql.DB) (JobResumes, error) {
	return JobResumes{
		queries: dbx.NewJobs(db),
	}, nil
}

func (j JobResumes) Create(ctx context.Context, userID uuid.UUID) error {
	if err := j.queries.Insert(ctx, dbx.JobResumeModel{
		UserID: userID,
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorCabinetForUserAlreadyExists(err, userID.String())
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

func (j JobResumes) Get(ctx context.Context, userID uuid.UUID) (models.JobResume, error) {
	job, err := j.queries.FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.JobResume{}, ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		default:
			return models.JobResume{}, ape.ErrorInternal(err)
		}
	}

	return JobFromDb(job), nil
}

func (j JobResumes) UpdateDegree(ctx context.Context, userID uuid.UUID, degree string) error {
	if err := enums.ValidateDegree(degree); err != nil {
		return ape.ErrorPropertyIsNotValid(err)
	}

	job, err := j.Get(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	if job.DegreeUpdatedAt != nil {
		last := *job.DegreeUpdatedAt

		return domain.ValidateUpdateProperty(last, 365*24*time.Hour)
	}

	if err = j.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateJobInput{
		Degree:          &degree,
		DegreeUpdatedAt: &now,
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

func (j JobResumes) UpdateIndustry(ctx context.Context, userID uuid.UUID, industry string) error {
	if err := enums.ValidateIndustry(industry); err != nil {
		return ape.ErrorPropertyIsNotValid(err)
	}

	job, err := j.Get(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	if job.IndustryUpdatedAt != nil {
		last := *job.IndustryUpdatedAt

		return domain.ValidateUpdateProperty(last, 365*24*time.Hour)
	}

	if err = j.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateJobInput{
		Industry:          &industry,
		IndustryUpdatedAt: &now,
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

func (j JobResumes) UpdateIncome(ctx context.Context, userID uuid.UUID, income string) error {
	if err := enums.ValidateIncome(income); err != nil {
		return ape.ErrorPropertyIsNotValid(err)
	}

	job, err := j.Get(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	if job.IncomeUpdatedAt != nil {
		last := *job.IncomeUpdatedAt

		return domain.ValidateUpdateProperty(last, 365*24*time.Hour)
	}

	if err = j.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateJobInput{
		Income:          &income,
		IncomeUpdatedAt: &now,
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

type AdminJobUpdate struct {
	Degree   *string `json:"degree"`
	Industry *string `json:"industry"`
	Income   *string `json:"income"`
}

func (j JobResumes) AdminUpdate(ctx context.Context, userID uuid.UUID, input AdminJobUpdate) error {
	_, err := j.Get(ctx, userID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()

	var dbInput dbx.UpdateJobInput

	if input.Degree != nil {
		if err = enums.ValidateDegree(*input.Degree); err != nil {
			return ape.ErrorPropertyIsNotValid(err)
		}

		dbInput.Degree = input.Degree
		dbInput.DegreeUpdatedAt = &now
	}

	if input.Industry != nil {
		if err = enums.ValidateIndustry(*input.Industry); err != nil {
			return ape.ErrorPropertyIsNotValid(err)
		}

		dbInput.Industry = input.Industry
		dbInput.IndustryUpdatedAt = &now

	}

	if input.Income != nil {
		if err = enums.ValidateIncome(*input.Income); err != nil {
			return ape.ErrorPropertyIsNotValid(err)
		}

		dbInput.Income = input.Income
		dbInput.IncomeUpdatedAt = &now

	}

	if err := j.queries.New().FilterUserID(userID).Update(ctx, dbInput); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		default:
			return ape.ErrorInternal(err)
		}
	}

	return nil
}

func JobFromDb(model dbx.JobResumeModel) models.JobResume {
	return models.JobResume{
		UserID:   model.UserID,
		Degree:   model.Degree,
		Industry: model.Industry,
		Income:   model.Income,

		DegreeUpdatedAt:   model.DegreeUpdatedAt,
		IndustryUpdatedAt: model.IndustryUpdatedAt,
		IncomeUpdatedAt:   model.IncomeUpdatedAt,
	}
}
