-- +migrate Up
CREATE TABLE audit_logs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type smallint NOT NULL,
    entity_id uuid NOT NULL,
    action smallint NOT NULL,
    details jsonb,
    created_at timestamp with time zone DEFAULT now(),
    
    CONSTRAINT audit_logs_entity_type_check CHECK (entity_type BETWEEN 1 AND 4),
    CONSTRAINT audit_logs_action_check CHECK (action BETWEEN 1 AND 11)
);

CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at DESC);

-- +migrate Down
DROP TABLE audit_logs;