-- internal/lms/schema.ddl
-- Cloud Spanner DDL (Strict DDD Isolation)

CREATE TABLE Courses (
    ID STRING(36) NOT NULL,
    Title STRING(255) NOT NULL,
    IsMandatory BOOL NOT NULL,
    ValidForDays INT64 NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (ID);

CREATE TABLE Enrollments (
    CourseID STRING(36) NOT NULL,
    ID STRING(36) NOT NULL,
    EmployeeID STRING(36) NOT NULL, -- Soft Link crossing domains to internal/hcm
    Status STRING(50) NOT NULL,
    Score INT64,
    CompletedAt TIMESTAMP,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (CourseID, ID),
  INTERLEAVE IN PARENT Courses ON DELETE CASCADE;

CREATE TABLE Certifications (
    ID STRING(36) NOT NULL,
    EmployeeID STRING(36) NOT NULL, -- Soft Link to HCM
    CourseID STRING(36) NOT NULL,   -- Standalone correlation to Course
    IssuedAt TIMESTAMP NOT NULL,
    ExpiresAt TIMESTAMP NOT NULL,
    Revoked BOOL NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (ID);

CREATE INDEX CertificationsByEmployee ON Certifications(EmployeeID);
