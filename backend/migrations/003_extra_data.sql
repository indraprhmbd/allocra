-- Add more sample rooms (nodes)
INSERT INTO rooms (name, capacity) VALUES
('NODE-AX-06', 128),
('NODE-AX-07', 256),
('NODE-AX-08', 64),
('NODE-AX-09', 512),
('NODE-AX-10', 1024)
ON CONFLICT (name) DO NOTHING;

-- Add some sample allocations for testing
INSERT INTO bookings (room_id, user_id, start_time, end_time, status)
SELECT 
    id, 
    1, 
    NOW() + (interval '1 hour' * (id % 5)), 
    NOW() + (interval '1 hour' * (id % 5)) + interval '2 hours',
    'approved'
FROM rooms
WHERE name LIKE 'NODE-AX-%'
LIMIT 5;
