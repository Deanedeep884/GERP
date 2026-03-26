-- internal/hcm/schema.ddl
-- Cloud Spanner DDL (Strict DDD Isolation)

CREATE TABLE Employees (
    ID STRING(36) NOT NULL,
    FirstName STRING(255) NOT NULL,
    LastName STRING(255) NOT NULL,
    Role STRING(100) NOT NULL,
    Email STRING(255) NOT NULL,
    IsActive BOOL NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (ID);

CREATE TABLE PayrollRuns (
    EmployeeID STRING(36) NOT NULL,
    ID STRING(36) NOT NULL,
    GrossPay INT64 NOT NULL,
    NetPay INT64 NOT NULL,
    FinanceRefID STRING(36), -- Soft Link traversing across domains to internal/finance. 
    PayPeriodEnd TIMESTAMP NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (EmployeeID, ID),
  INTERLEAVE IN PARENT Employees ON DELETE CASCADE;
