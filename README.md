# Loan Engine

Loan Engine is a service which helps in maintaining the state for a loan application.

A loan application goes through multiple states:
- Proposed
- Approved
- Invested
- Disbursed

This service has 4 tables in the `loan_engine`:
- Loan: Has the details about the loan.
- Loan State: Has the details about the current state of the loan.
- Loan Investment: Has details about the current investors and their investment amount in a particular loan.
- Loan Amount: Has details about the total investment of a loan.

This service implements the following APIs:
- POST /v1/loan/propose : To create a new loan
  - required params as JSON: 
    - borrower_id string
    - principle_amount int
    - roi float64
    - rate float64
  - response as JSON:
    - message string (SUCCESS | Error message)
- POST /v1/loan/invest :
  - required params as JSON: 
    - loan_id string
    - investor_id string
    - amount int
  - response as JSON:
    - message string (SUCCESS | Error message)
- POST /v1/loan/approve
  - required params as multipart form: 
    - loan_id string
    - date timestamp
    - photo file
  - response as JSON:
    - message string (SUCCESS | Error message)
- POST /v1/loan/disburse 
  - required params as multipart form: 
    - loan_id string
    - date timestamp
    - photo file
  - response as JSON:
    - message string (SUCCESS | Error message)
- GET /v1/loan/state/:id
  - required params in URL:
    - id string (Loan ID)
  - response as JSON 
    - loan_id string
    - loan_state string
  - error response as JSON:
    - message string (error message)
- GET /v1/loan/list/*state
  - optional parameter in URL:
    - state string (proposed | approved | invested | disbursed)
  - required parameter as query string:
    - limit int
    - offset int
  - response as JSON:
    - array of:
      - loan_id string
      - principle_amount int
      - rate float64
      - roi float64
  - error response as JSON:
    - message string (error message)

To run the migrations use the following command from the root folder:
```
migrate -path migrations/ -database "postgresql://<username>:<pass>@<host>:<port>/loan_engine?sslmode=disable" -verbose up
```

To run the code use the following command from the root folder:
```
go run main.go
```

PS: Currently there are no unit tests written for this code, but as we use interface and dependecy injection for all the layers it shouldn't be much complicated to write the tests by mocking the required function calls.