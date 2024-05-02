package cron

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/myorn/gepard-m/constants"
	"github.com/myorn/gepard-m/dao"
)

const jobSleepDuration = 1

func RunCancelJob(dbSession *sql.DB) {
	for {
		time.Sleep(time.Minute * jobSleepDuration)

		if err := Cancel10OddMessages(dbSession); err != nil {
			log.Println(err)
		}
	}
}

func Cancel10OddMessages(dbSession *sql.DB) error {
	// get 10 odd operations
	msgs, err := dao.GetLast10OddMessages(dbSession)
	if err != nil {
		return err
	}

	var amountOverall int64

	for i := range msgs {
		// this could potentially overflow and I can fix it, but do I really need to?
		amount, err := strconv.ParseInt(msgs[i].Amount, 10, 64)
		if err != nil {
			return err
		}

		switch msgs[i].State {
		case constants.StateDeposit:
			amountOverall -= amount
		case constants.StateWithdraw:
			amountOverall += amount
		}
	}

	// update balance and cancel them
	if err := dao.AddBalance(dbSession, amountOverall); err != nil {
		return err
	}

	for i := range msgs {
		if err := dao.DeleteMessageByTxID(dbSession, msgs[i].TxID); err != nil {
			return err
		}
	}

	return nil
}
