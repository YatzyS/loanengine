package constants

const (
	LOAN_TABLE            = "loan"
	LOAN_STATE_TABLE      = "loan_state"
	LOAN_INVESTMENT_TABLE = "loan_investment"
	LOAN_AMOUNT_TABLE     = "loan_amount"
)

// These variable need to come from config, but keeping it here as const for now.
const (
	DB_HOST = "localhost"
	DB_PORT = 5432
	DB_USER = "postgres"
	DB_PASS = "root"
	DB_NAME = "loan_engine"
)
