CREATE EXTENSION IF NOT EXISTS vector;

-- 1. Create the partitioned master table
CREATE TABLE coams_chunks (
    id UUID DEFAULT gen_random_uuid(),
    channel_id VARCHAR(50) NOT NULL,
    document_id UUID NOT NULL,
    content TEXT NOT NULL,
    embedding VECTOR(768),
    PRIMARY KEY (channel_id, id)
) PARTITION BY LIST (channel_id);

-- 2. Provision two physical test shards
CREATE TABLE coams_chunks_engineering PARTITION OF coams_chunks FOR VALUES IN ('engineering');
CREATE TABLE coams_chunks_hr PARTITION OF coams_chunks FOR VALUES IN ('hr');

-- 3. Create isolated HNSW indexes
CREATE INDEX idx_engineering_emb ON coams_chunks_engineering USING hnsw (embedding vector_cosine_ops);
CREATE INDEX idx_hr_emb ON coams_chunks_hr USING hnsw (embedding vector_cosine_ops);

-- 4. Seed test data
INSERT INTO coams_chunks (channel_id, document_id, content, embedding)
VALUES 
('engineering', gen_random_uuid(), 'Kubernetes cluster deployment specs', array_fill(0.1, ARRAY[768])::vector),
('hr', gen_random_uuid(), 'Top secret executive compensation plans', array_fill(0.9, ARRAY[768])::vector);
