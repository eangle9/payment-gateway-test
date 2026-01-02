package persistencedb

import (
	"context"
	"pg/internal/constant/errors"
	"pg/internal/constant/model/db"
	"pg/internal/constant/model/dto"

	"go.uber.org/zap"
)

func (q *PersistenceDB) GenerateCompanyCredential(ctx context.Context,
	req dto.CreateCompanyToken) (err error) {
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return errors.ErrUnableToCreate.Wrap(err, "error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				q.log.Error(ctx, "error rolling back transaction", zap.Error(rbErr))
			}
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				q.log.Error(ctx, "error rolling back transaction", zap.Error(rbErr))
			}
		}
	}()

	tQ := q.WithTx(tx)
	if err = tQ.InActiveCompanyToken(ctx, req.CompanyID); err != nil {
		err = errors.ErrUnableToUpdate.Wrap(err, "error updating the merchant key")
		q.log.Error(ctx, "unable to update merchant key", zap.Error(err),
			zap.String("merchant-id", req.CompanyID.String()))
		return err
	}
	if _, err = tQ.CreateCompanyToken(ctx, db.CreateCompanyTokenParams{
		TokenID:   req.TokenID,
		CompanyID: req.CompanyID,
	}); err != nil {
		err = errors.ErrUnableToCreate.Wrap(err, "error creating the merchant key")
		q.log.Error(ctx, "unable to create merchant key",
			zap.Error(err), zap.String("merchant-id", req.CompanyID.String()),
			zap.String("token-id", req.TokenID.String()))
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return errors.ErrUnableToCreate.Wrap(err, "error committing transaction")
	}
	return nil
}
