BEGIN;

CREATE TABLE tracks (
    id SERIAL PRIMARY KEY,
    artist VARCHAR(256) NOT NULL,
    name VARCHAR(256) NOT NULL,
    price INT NOT NULL CHECK (price > 0)
);

CREATE TABLE playback_logs (
    id SERIAL PRIMARY KEY,
    track_id INT NOT NULL, 
    played_at TIMESTAMP NOT NULL DEFAULT NOW(),
    amount_paid INT NOT NULL CHECK (amount_paid > 0),

    CONSTRAINT fk_track 
        FOREIGN KEY (track_id) 
        REFERENCES tracks (id) 
        ON DELETE CASCADE
);

CREATE INDEX idx_playback_logs_track_id ON playback_logs (track_id);

INSERT INTO tracks (artist, name, price) VALUES
('The Gophers', 'Static Typing Blues', 99),
('The Gophers', 'Concurrency in C#', 149),
('Duck Type', 'Implicit Interface', 199),
('SQL Masters', 'The Trailing Comma', 49),
('The Compilers', 'Panic at the Handler', 299);

INSERT INTO playback_logs (track_id, amount_paid, played_at) VALUES
(1, 99, NOW() - interval '1 hour'), (1, 99, NOW() - interval '2 hours'), (1, 99, NOW() - interval '3 hours'), (1, 99, NOW() - interval '4 hours'), (1, 99, NOW() - interval '5 hours'),
(2, 149, NOW() - interval '1 hour'), (2, 149, NOW() - interval '2 hours'), (2, 149, NOW() - interval '3 hours'), (2, 149, NOW() - interval '4 hours'), (2, 149, NOW() - interval '5 hours'),
(3, 199, NOW() - interval '1 hour'), (3, 199, NOW() - interval '2 hours'), (3, 199, NOW() - interval '3 hours'), (3, 199, NOW() - interval '4 hours'), (3, 199, NOW() - interval '5 hours'),
(4, 49, NOW() - interval '1 hour'), (4, 49, NOW() - interval '2 hours'), (4, 49, NOW() - interval '3 hours'), (4, 49, NOW() - interval '4 hours'), (4, 49, NOW() - interval '5 hours'),
(5, 299, NOW() - interval '1 hour'), (5, 299, NOW() - interval '2 hours'), (5, 299, NOW() - interval '3 hours'), (5, 299, NOW() - interval '4 hours'), (5, 299, NOW() - interval '5 hours'),
(1, 99, NOW()), (2, 149, NOW()), (3, 199, NOW()), (4, 49, NOW()), (5, 299, NOW());

COMMIT;