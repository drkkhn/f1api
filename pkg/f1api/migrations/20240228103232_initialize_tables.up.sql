CREATE TABLE IF NOT EXISTS teams (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    car text NOT NULL
);

CREATE TABLE IF NOT EXISTS racers (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    firstName text NOT NULL,
    lastName text NOT NULL,
    teamId bigserial REFERENCES teams(id)
);