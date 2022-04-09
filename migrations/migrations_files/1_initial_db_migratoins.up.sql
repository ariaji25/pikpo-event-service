CREATE TABLE events(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    date TIMESTAMP NOT NULL,
    start_time VARCHAR(10) NOT NULL,
    end_time VARCHAR(10) NOT NULL,
    allday BOOLEAN DEFAULT false,
    event_type INTEGER DEFAULT 0,
    timezone VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP NULL
);