CREATE TABLE IF NOT EXISTS cats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name varchar(50) NOT NULL,
    description varchar(255),
    location varchar(255) NOT NULL,
    photo_path VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    last_seen timestamp(0) with time zone NOT NULL DEFAULT NOW()
)