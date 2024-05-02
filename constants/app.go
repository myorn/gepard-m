package constants

type contextKey string

const (
	ConfigFilePath            = "config.ini"
	DBSession      contextKey = "dbSession"

	// accepted sources
	SourceGame    = "game"
	SourcePayment = "payment"
	SourceService = "service"

	// accepted states
	StateDeposit  = "deposit"
	StateWithdraw = "withdraw"
)
