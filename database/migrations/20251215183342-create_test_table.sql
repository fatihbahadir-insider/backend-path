-- +migrate Up
CREATE TABLE "tests" (
    guid uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    title varchar,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

-- +migrate Down
DROP TABLE "tests";

