-- Add the 'tags' column (an array of VARCHAR with max length 100)
ALTER TABLE posts
ADD COLUMN tags VARCHAR(100)[];

-- Add 'updated_at' column if it does not exist (timestamp with time zone, default to current timestamp)
ALTER TABLE posts
ADD COLUMN updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW();
