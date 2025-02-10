CREATE TABLE loan (
    loan_id TEXT, 
    borrower_id TEXT, 
    principle_amount INT, 
    roi DECIMAL, 
    rate DECIMAL, 
    agreement_link TEXT, 
    version INT
    );

CREATE TABLE loan_state (
    loan_id TEXT, 
    loan_state TEXT, 
    approve_emp_id TEXT,
    approve_proof TEXT,
    approve_time TIMESTAMP,
    disburse_emp_id TEXT,
    disburse_proof TEXT,
    disburse_time TEXT,
    version INT
    );

CREATE TABLE loan_investment (
    loan_id TEXT, 
    investor_id TEXT, 
    amount INT, 
    version INT
    );

CREATE TABLE loan_amount (
    loan_id TEXT, 
    required_amount INT, 
    invested_amount INT, 
    version INT
    );