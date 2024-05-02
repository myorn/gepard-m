package dao

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/myorn/gepard-m/constants"
)

func SaveTx(
	ctx context.Context,
	jsonedReq []byte,
	txID string,
) error {
	dbSession := ctx.Value(constants.DBSession).(*sql.DB)
	if dbSession == nil {
		log.Fatal("no database connection found")
		return nil
	}

	// Connect to the database
	err := dbSession.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Insert data into the database
	_, err = dbSession.ExecContext(ctx,
		`INSERT INTO messages (message, tx_id) VALUES ($1,$2);`,
		jsonedReq, txID)
	if err != nil {
		return err
	}

	return nil
}

func AddBalanceFromTx(ctx context.Context, state, amount string) error {
	dbSession := ctx.Value(constants.DBSession).(*sql.DB)
	if dbSession == nil {
		log.Fatal("no database connection found")
		return nil
	}

	// Connect to the database
	err := dbSession.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	updateStr := "UPDATE deposit SET credit = credit"
	switch state {
	case "deposit":
		updateStr += " + $1;"
	case "withdraw":
		updateStr += " - $1;"
	}

	// Insert data into the database
	_, err = dbSession.ExecContext(ctx, updateStr, amount)
	if err != nil {
		return err
	}

	return nil
}

type Message struct {
	ID        string
	State     string
	Amount    string
	TxID      string
	CreatedAt string
}

func GetLast10OddMessages(db *sql.DB) ([]Message, error) {
	// Prepare the query
	query := `
        SELECT id, message->>'state', message->>'amount', tx_id, created_at
        FROM messages
	    WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT 20
    `

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the result set
	var messages []Message
	var i uint8

	for rows.Next() {
		i++

		if i%2 == 0 {
			continue
		}

		var message Message
		err := rows.Scan(&message.ID,
			&message.State,
			&message.Amount,
			&message.TxID,
			&message.CreatedAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func DeleteMessageByTxID(db *sql.DB, txID string) error {
	// Prepare the update statement
	stmt, err := db.Prepare("UPDATE messages SET deleted_at = NOW() WHERE tx_id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the update statement
	_, err = stmt.Exec(txID)
	if err != nil {
		return err
	}

	return nil
}

func AddBalance(db *sql.DB, balance int64) error {
	// Prepare the update statement
	stmt, err := db.Prepare("UPDATE deposit SET credit = credit + $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the update statement
	_, err = stmt.Exec(balance)
	if err != nil {
		return err
	}

	return nil
}

func GetBalance(ctx context.Context) (int64, error) {
	dbSession := ctx.Value(constants.DBSession).(*sql.DB)
	if dbSession == nil {
		return 0, nil
	}

	// Query the database for the credit balance
	var credit int64
	err := dbSession.QueryRowContext(ctx, "SELECT credit FROM deposit;").Scan(&credit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No rows in the result set
			return 0, nil
		}

		return 0, err
	}

	return credit, nil
}

func GetSumOfLiveMessages(ctx context.Context) (int64, error) {
	dbSession := ctx.Value(constants.DBSession).(*sql.DB)
	if dbSession == nil {
		log.Fatal("no database connection found")
		return 0, nil
	}

	// Query the database for the sum of amounts of non-deleted messages
	var sum int64
	err := dbSession.QueryRowContext(ctx,
		"SELECT SUM((message->>'amount')::bigint * (CASE WHEN message->>'state' = 'deposit' THEN 1 ELSE -1 END)) FROM messages WHERE deleted_at IS NULL").Scan(&sum)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows in the result set
			return 0, nil
		}
		return 0, err
	}

	return sum, nil
}
