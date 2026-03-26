-- internal/finance/schema.ddl
-- Cloud Spanner DDL (Zero Cross-Domain Foreign Keys)

CREATE TABLE Accounts (
    ID STRING(36) NOT NULL,
    Name STRING(255) NOT NULL,
    Type STRING(50) NOT NULL,
    AccountOwnerID STRING(36), -- Soft link to HCM/MDM UUIDs
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (ID);

CREATE TABLE LedgerEntries (
    ID STRING(36) NOT NULL,
    TransactionID STRING(36), -- Soft link to Pipeline/Temporal Sagas
    Description STRING(MAX),
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (ID);

CREATE TABLE LineItems (
    LedgerEntryID STRING(36) NOT NULL,
    LineItemID STRING(36) NOT NULL,
    AccountID STRING(36) NOT NULL,
    AmountCents INT64 NOT NULL,
    CustomerID STRING(36), -- Soft link to MDM/Revenue UUIDs
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (LedgerEntryID, LineItemID),
  INTERLEAVE IN PARENT LedgerEntries ON DELETE CASCADE;
