package cron

import (
	"database/sql"
	"log"
	"time"

	"github.com/myorn/gepard-m/dao"
)

const jobSleepDuration = 1

func RunCancelJob(dbSession *sql.DB, ch chan any) {
	for {
		time.Sleep(time.Minute * jobSleepDuration)

		// get 10 odd operations
		msgs, err := dao.GetLast10OddMessages(dbSession)
		if err != nil {
			log.Println(err)
		}

		var amountOverall int64

		for i := range msgs {
			if msgs[i].State == "deposit" {
				amountOverall -= int64(msgs[i].Amount) // this could potentially overflow and I can fix it, but do I really need to?
			}

			amountOverall += int64(msgs[i].Amount)
		}

		// update balance and cancel them
		if err := dao.UpdateAddBalance(dbSession, amountOverall); err != nil {
			log.Println(err)
		}

		for i := range msgs {
			if err := dao.DeleteMessageByTxID(dbSession, msgs[i].TxID); err != nil {
				log.Println(err)
			}
		}
	}
}
