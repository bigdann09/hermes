CREATE TYPE notification_type AS ENUM ('system', 'social');

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    type notification_type NOT NULL,
    metadata JSONB NOT NULL,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_type ON notifications(type);
CREATE INDEX idx_read_at ON notifications(read_at);
CREATE INDEX idx_user_id ON notifications(user_id);