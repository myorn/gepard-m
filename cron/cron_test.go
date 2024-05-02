package cron_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/myorn/gepard-m/constants"
	"github.com/myorn/gepard-m/cron"
	"github.com/myorn/gepard-m/dao"
	"github.com/myorn/gepard-m/db"

	"github.com/stretchr/testify/assert"
)

func contextWithDB() context.Context {
	session := db.InitDB()
	db.Flush(session)
	db.Migrate(session)
	return context.WithValue(context.Background(), constants.DBSession, session)
}

type TxData struct {
	Source string `json:"source"`
	State  string `json:"state"`
	Amount string `json:"amount"`
	TxId   string `json:"tx_id"`
}

func generateRandomTxData() ([]byte, int, string) {
	rand.Seed(time.Now().UnixNano())
	amount := rand.Intn(1000) + 1                                   // Random amount between 1 and 1000
	states := []string{"deposit", "deposit", "deposit", "withdraw"} // State can be "deposit" or "withdraw"
	sources := []string{"game", "payment", "service"}               // Source can be "game", "payment", or "service"
	state := states[rand.Intn(len(states))]                         // Randomly select a state
	source := sources[rand.Intn(len(sources))]                      // Randomly select a source
	txID := strconv.FormatInt(time.Now().UnixNano(), 10)            // Use current time as transaction ID

	data := TxData{
		Amount: strconv.Itoa(amount),
		State:  state,
		Source: source,
		TxId:   txID,
	}

	dataBytes, _ := json.Marshal(data)

	if state == "withdraw" {
		amount = -amount
	}

	return dataBytes, amount, txID
}

func TestCancel10OddMessagesHappyPath(t *testing.T) {
	ctx := contextWithDB()
	dbSession := ctx.Value(constants.DBSession).(*sql.DB)

	var overallAmount int
	for i := 0; i < 20; i++ {
		data, amount, txID := generateRandomTxData()
		assert.NoError(t, dao.SaveTx(ctx, data, txID))

		overallAmount += amount
	}

	assert.NoError(t, dao.AddBalance(dbSession, int64(overallAmount)))

	err := cron.Cancel10OddMessages(dbSession)
	assert.NoError(t, err)

	balance, err := dao.GetBalance(ctx)
	assert.NoError(t, err)

	sumOfLiveMessages, err := dao.GetSumOfLiveMessages(ctx)
	assert.NoError(t, err)
	assert.Equal(t, sumOfLiveMessages, balance)
}
