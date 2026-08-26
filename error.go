package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

const (
	DuplicateKeyErrorCode        = pq.ErrorCode("23505")
	ForeignKeyViolationErrorCode = pq.ErrorCode("23503")
)

var (
	ErrBeginTransactionFailed     = errors.New("begin transaction failed")
	ErrCommitTransactionFailed    = errors.New("commit transaction failed")
	ErrRollbackTransactionFailed  = errors.New("revert transaction failed")
	ErrCommitWithoutTransaction   = errors.New("invalid transaction, can not perform commit without transaction")
	ErrRollbackWithoutTransaction = errors.New("invalid transaction, can not perform rollback without transaction")
	ErrNoPgConnection             = errors.New("postgres connection is not established")
)

func IsErrorCode(err error, errCode pq.ErrorCode) bool {
	var pgxErr *pgconn.PgError
	ok := errors.As(err, &pgxErr)
	if ok {
		return pq.ErrorCode(pgxErr.Code) == errCode
	}
	return false
}

func TryCastErrorToPgError(err error) any {
	if err == nil {
		return nil
	}
	if pqErr, ok := errors.AsType[*pq.Error](err); ok {
		return pqErr
	}
	if pgxErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgxErr
	}
	return err.Error()
}
