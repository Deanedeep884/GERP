-- /Users/Duncan/GERP/internal/coams/migrations/001_coams_schema.sql

-- 1. Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- 2. Physical Partitions (Channels)
CREATE TABLE coams_documents (
    id UUID DEFAULT gen_random_uuid(),
    tenant_id VARCHAR(255) NOT NULL,
    channel_id VARCHAR(50) NOT NULL, -- e.g., 'engineering', 'hr'
    title VARCHAR(255) NOT NULL,
    raw_markdown TEXT NOT NULL,
    
    -- Verbose Metadata (Managed by COAMS, tied to IAM workspace identities)
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    version INT DEFAULT 1,
    
    PRIMARY KEY (channel_id, id)
) PARTITION BY LIST (channel_id);

-- 3. Vector Storage (Partitioned to match Documents)
CREATE TABLE coams_chunks (
    id UUID DEFAULT gen_random_uuid(),
    channel_id VARCHAR(50) NOT NULL,
    document_id UUID NOT NULL,
    header_path TEXT NOT NULL,
    content TEXT NOT NULL,
    tokens INT NOT NULL DEFAULT 0,
    embedding VECTOR(768),
    
    PRIMARY KEY (channel_id, id),
    FOREIGN KEY (channel_id, document_id) REFERENCES coams_documents(channel_id, id) ON DELETE CASCADE
) PARTITION BY LIST (channel_id);

-- 4. Outbound Edges / Agent-Index Graph (Partitioned)
CREATE TABLE coams_links (
    id UUID DEFAULT gen_random_uuid(),
    channel_id VARCHAR(50) NOT NULL,
    source_document_id UUID NOT NULL,
    target_document_id UUID, -- NULL if external
    is_external BOOLEAN DEFAULT false,
    external_url TEXT,
    
    PRIMARY KEY (channel_id, id),
    FOREIGN KEY (channel_id, source_document_id) REFERENCES coams_documents(channel_id, id) ON DELETE CASCADE
) PARTITION BY LIST (channel_id);
