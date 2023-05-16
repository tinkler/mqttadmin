package db

import (
	"context"
	"net/http"

	"gorm.io/gorm"
)

func WrapGorm() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), dbKey, DB())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WrapGormTx() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var dbInst *gorm.DB
			if old := r.Context().Value(dbKey); old != nil {
				dbInst = old.(*gorm.DB)
			} else {
				dbInst = DB()
			}
			tx := dbInst.Begin()
			ctx := context.WithValue(r.Context(), dbKey, tx)
			next.ServeHTTP(w, r.WithContext(ctx))
			if r.Context().Err() != nil {
				tx.Rollback()
				return
			}
			tx.Commit()
			// TODO: handle error
		})
	}
}

func WithValue(ctx context.Context, db *gorm.DB) context.Context {
	if db == nil {
		return ctx
	}
	return context.WithValue(ctx, dbKey, db)
}
