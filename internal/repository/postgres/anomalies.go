package postgres

import "fmt"

const (
	insertWalletAffectedRows  = 1
	updateAddressAffectedRows = 1
)

type InsertAnomalyError struct {
	Insert               bool
	RowsAffected         int64
	ExpectedRowsAffected int64
}

func (e *InsertAnomalyError) Error() string {
	return fmt.Sprintf("insertion anomaly: expected insert operation true (%d row), got insert %t (%d rows)",
		e.ExpectedRowsAffected, e.Insert, e.RowsAffected)
}

type UpdateAnomalyError struct {
	Update               bool
	RowsAffected         int64
	ExpectedRowsAffected int64
}

func (e *UpdateAnomalyError) Error() string {
	return fmt.Sprintf("insertion anomaly: expected insert operation true (%d row), got insert %t (%d rows)",
		e.ExpectedRowsAffected, e.Update, e.RowsAffected)
}
