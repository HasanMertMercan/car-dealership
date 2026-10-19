CREATE TABLE IF NOT EXISTS car (
    id UUID PRIMARY KEY,
    brand TEXT NOT NULL,
    model TEXT NOT NULL,
    year INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_car_brand ON car(brand);
CREATE INDEX IF NOT EXISTS idx_car_model ON car(model);