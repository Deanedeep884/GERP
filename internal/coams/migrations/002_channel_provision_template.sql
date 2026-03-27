-- /Users/Duncan/GERP/internal/coams/migrations/002_channel_provision_template.sql
-- This template is executed dynamically by COAMS when a new channel/partition is mapped.
-- Replace {{CHANNEL_ID}} with the actual channel name (e.g., 'engineering', 'sales').

-- 1. Create the partition for Documents
CREATE TABLE coams_documents_{{CHANNEL_ID}} 
PARTITION OF coams_documents 
FOR VALUES IN ('{{CHANNEL_ID}}');

-- 2. Create the partition for Chunks
CREATE TABLE coams_chunks_{{CHANNEL_ID}} 
PARTITION OF coams_chunks 
FOR VALUES IN ('{{CHANNEL_ID}}');

-- 3. Create the partition for Links (Agent-Index)
CREATE TABLE coams_links_{{CHANNEL_ID}} 
PARTITION OF coams_links 
FOR VALUES IN ('{{CHANNEL_ID}}');

-- 4. Create the isolated HNSW Index for this specific channel vectors
CREATE INDEX idx_chunks_{{CHANNEL_ID}}_emb 
ON coams_chunks_{{CHANNEL_ID}} USING hnsw (embedding vector_cosine_ops);
