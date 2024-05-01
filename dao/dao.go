package dao

import (
	"context"
	"database/sql"
	"log"

	"github.com/myorn/gepard-m/constants"
)

func SaveTxAndUpdateBalance(
	ctx context.Context,
	jsonedReq []byte,
	amount,
	state,
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
	_, err = dbSession.ExecContext(ctx, `DO $$
BEGIN
    UPDATE deposit SET balance =
    CASE
    WHEN ? = 'deposit' THEN balance + ?
    ELSE balance - ?;

    INSERT INTO messages (message, tx_id)
    VALUES (?,?);
END $$;`,
		state, amount, amount, jsonedReq, txID)
	if err != nil {
		return err
	}

	return nil
}

type Message struct {
	ID        string
	State     string
	Amount    uint64
	TxID      uint64
	CreatedAt string
}

func GetLast10OddMessages(db *sql.DB) ([]Message, error) {
	// Prepare the query
	query := `
        SELECT id, message->'state', message->'amount', tx_id, created_at
        FROM messages
        ORDER BY created_at DESC
	WHERE deleted_at IS NULL
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

func DeleteMessageByTxID(db *sql.DB, txID uint64) error {
	// Prepare the update statement
	stmt, err := db.Prepare("UPDATE messages SET deleted_at = NOW() WHERE tx_id = ?")
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

func UpdateAddBalance(db *sql.DB, balance int64) error {
	// Prepare the update statement
	stmt, err := db.Prepare("UPDATE deposit SET balance = ?")
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
