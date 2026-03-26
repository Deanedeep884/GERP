-- internal/mdm/schema.ddl
-- Cloud Spanner DDL (Strict DDD Isolation)

CREATE TABLE GlobalEntities (
    ID STRING(36) NOT NULL,
    LegalName STRING(255) NOT NULL,
    TaxID STRING(100) NOT NULL,
    CountryCode STRING(10) NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (ID);

CREATE TABLE EntityMappings (
    GlobalEntityID STRING(36) NOT NULL,
    Domain STRING(50) NOT NULL,
    LocalID STRING(36) NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (GlobalEntityID, Domain, LocalID),
  INTERLEAVE IN PARENT GlobalEntities ON DELETE CASCADE;
