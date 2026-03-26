-- internal/scm/schema.ddl
-- Cloud Spanner DDL (Strict DDD Isolation)

CREATE TABLE Products (
    ID STRING(36) NOT NULL,
    SKU STRING(100) NOT NULL,
    Name STRING(255) NOT NULL,
    Description STRING(MAX),
    IsActive BOOL NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (ID);

CREATE TABLE InventoryLots (
    ID STRING(36) NOT NULL,
    ProductID STRING(36) NOT NULL,
    WarehouseID STRING(36), -- Soft Link crossing domains to internal/eam
    Quantity INT64 NOT NULL,
    CostBasis INT64 NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (ID);

-- Accelerated spatial fetching without restricting cross-table PKs
CREATE INDEX InventoryByProduct ON InventoryLots(ProductID);
