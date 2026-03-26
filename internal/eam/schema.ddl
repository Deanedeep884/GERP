-- internal/eam/schema.ddl
-- Cloud Spanner DDL (Strict DDD Isolation)

CREATE TABLE Assets (
    ID STRING(36) NOT NULL,
    Name STRING(255) NOT NULL,
    Type STRING(50) NOT NULL,
    Status STRING(50) NOT NULL,
    FinanceAssetID STRING(36), -- Soft Link crossing domains to internal/finance
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (ID);

CREATE TABLE MaintenanceLogs (
    AssetID STRING(36) NOT NULL,
    ID STRING(36) NOT NULL,
    TechnicianID STRING(36), -- Soft Link crossing domains to internal/hcm
    Description STRING(MAX),
    Cost INT64 NOT NULL,
    CompletedAt TIMESTAMP NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (AssetID, ID),
  INTERLEAVE IN PARENT Assets ON DELETE CASCADE;
