package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/elector-cab-svc/internal/app/ape"
	"github.com/chains-lab/elector-cab-svc/internal/app/domain"
	"github.com/chains-lab/elector-cab-svc/internal/app/enums"
	"github.com/chains-lab/elector-cab-svc/internal/app/models"
	"github.com/chains-lab/elector-cab-svc/internal/dbx"
	"github.com/google/uuid"
)

type BiographiesQ interface {
	New() dbx.BiographiesQ

	Insert(ctx context.Context, input dbx.BioModel) error
	Update(ctx context.Context, input dbx.UpdateBioInput) error
	Select(ctx context.Context) ([]dbx.BioModel, error)
	Get(ctx context.Context) (dbx.BioModel, error)
	Delete(ctx context.Context) error

	FilterUserID(userID uuid.UUID) dbx.BiographiesQ

	Count(ctx context.Context) (int, error)
	Page(limit, offset uint64) dbx.BiographiesQ
}

type Biographies struct {
	queries BiographiesQ
}

func NewBiographies(db *sql.DB) (Biographies, error) {
	return Biographies{
		queries: dbx.NewBiographies(db),
	}, nil
}

func (b Biographies) Create(ctx context.Context, userID uuid.UUID) error {
	if err := b.queries.Insert(ctx, dbx.BioModel{
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

func (b Biographies) GetByUserID(ctx context.Context, userID uuid.UUID) (models.Biography, error) {
	bio, err := b.queries.New().FilterUserID(userID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Biography{}, ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		default:
			return models.Biography{}, ape.ErrorInternal(err)
		}
	}

	return BioFromDb(bio), nil
}

func (b Biographies) UpdateSex(ctx context.Context, userID uuid.UUID, sex string) error {
	if err := enums.ValidateSex(sex); err != nil {
		return ape.ErrorPropertyIsNotValid(err)
	}

	now := time.Now().UTC()

	bio, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if bio.SexUpdatedAt != nil {
		last := *bio.SexUpdatedAt

		return domain.ValidateUpdateProperty(last, 365*24*time.Hour)
	}

	if err = b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		Sex:          &sex,
		SexUpdatedAt: &now,
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

func (b Biographies) UpdateBirthday(ctx context.Context, userID uuid.UUID, birthday time.Time) error {
	bio, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if bio.Birthday != nil {
		return ape.ErrorPropertyUpdateNotAllowed(fmt.Errorf("birthday is already set you can do it once")) //TODO: add error
	}

	if err = b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		Birthday: &birthday,
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func (b Biographies) SetNationality(ctx context.Context, userID uuid.UUID, nationality string) error {
	//TODO validate nationality from other api
	if err := enums.ValidateNationality(nationality); err != nil {
		return ape.ErrorPropertyIsNotValid(err)
	}

	bio, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	if bio.NationalityUpdatedAt != nil {
		last := *bio.NationalityUpdatedAt

		return domain.ValidateUpdateProperty(last, 365*24*time.Hour)
	}

	if err = b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		Nationality:          &nationality,
		NationalityUpdatedAt: &now,
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

func (b Biographies) SetPrimaryLanguage(ctx context.Context, userID uuid.UUID, primaryLanguage string) error {
	//TODO validate primaryLanguage from other api
	if err := enums.ValidateLanguage(primaryLanguage); err != nil {
		return ape.ErrorPropertyIsNotValid(err)
	}

	bio, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	if bio.PrimaryLanguageUpdatedAt != nil {
		last := *bio.PrimaryLanguageUpdatedAt

		return domain.ValidateUpdateProperty(last, 365*24*time.Hour)
	}

	if err = b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		PrimaryLanguage:          &primaryLanguage,
		PrimaryLanguageUpdatedAt: &now,
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func (b Biographies) UpdateResidence(ctx context.Context, userID uuid.UUID, country string, city string) error {
	//TODO validate country and city from other api
	err := enums.ValidateResidence(city, country)
	if err != nil {
		return ape.ErrorPropertyIsNotValid(err)
	}

	bio, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	if bio.ResidenceUpdatedAt != nil {
		last := *bio.ResidenceUpdatedAt

		return domain.ValidateUpdateProperty(last, 365*24*time.Hour)
	}

	if err = b.queries.New().FilterUserID(userID).Update(ctx, dbx.UpdateBioInput{
		City:               &city,
		Country:            &country,
		ResidenceUpdatedAt: &now,
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

type AdminBioUpdate struct {
	Birthday        *time.Time
	Sex             *string
	Nationality     *string
	PrimaryLanguage *string
	City            *string
	Country         *string
}

func (b Biographies) AdminUpdateBio(ctx context.Context, userID uuid.UUID, input AdminBioUpdate) error {
	_, err := b.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	var dbInput dbx.UpdateBioInput

	if input.Birthday != nil {
		dbInput.Birthday = input.Birthday
	}

	if input.Sex != nil {
		if err := enums.ValidateSex(*input.Sex); err != nil {
			return ape.ErrorPropertyIsNotValid(err)
		}

		dbInput.Sex = input.Sex
		dbInput.SexUpdatedAt = &now
	}

	if input.City != nil && input.Country != nil {
		//TODO implement this functionality
		if err = enums.ValidateResidence(*input.City, *input.Country); err != nil {
			return ape.ErrorPropertyIsNotValid(err)
		}

		dbInput.City = input.City
		dbInput.Country = input.Country
		dbInput.ResidenceUpdatedAt = &now
	}
	if input.Nationality != nil {
		//TODO implement this functionality
		if err = enums.ValidateNationality(*input.Nationality); err != nil {
			return ape.ErrorPropertyIsNotValid(err)
		}

		dbInput.Nationality = input.Nationality
		dbInput.NationalityUpdatedAt = &now
	}
	if input.PrimaryLanguage != nil {
		//TODO implement this functionality
		if err = enums.ValidateLanguage(*input.PrimaryLanguage); err != nil {
			return ape.ErrorPropertyIsNotValid(err)
		}

		dbInput.PrimaryLanguage = input.PrimaryLanguage
		dbInput.PrimaryLanguageUpdatedAt = &now
	}

	if err = b.queries.New().FilterUserID(userID).Update(ctx, dbInput); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ape.ErrorCabinetForUserDoesNotExist(err, userID.String())
		default:
			return ape.ErrorInternal(err) //TODO
		}
	}

	return nil
}

func BioFromDb(input dbx.BioModel) models.Biography {
	return models.Biography{
		UserID:          input.UserID,
		Sex:             input.Sex,
		Birthday:        input.Birthday,
		Nationality:     input.Nationality,
		PrimaryLanguage: input.PrimaryLanguage,
		City:            input.City,
		Country:         input.Country,

		SexUpdatedAt:             input.SexUpdatedAt,
		NationalityUpdatedAt:     input.NationalityUpdatedAt,
		PrimaryLanguageUpdatedAt: input.PrimaryLanguageUpdatedAt,
		ResidenceUpdatedAt:       input.ResidenceUpdatedAt,
	}
}
