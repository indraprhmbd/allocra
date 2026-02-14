-- Insert default admin user
INSERT INTO users (name, email, role) VALUES
('Admin User', 'admin@allocra.local', 'admin'),
('Test User', 'user@allocra.local', 'user');

-- Insert sample rooms
INSERT INTO rooms (name, capacity) VALUES
('Meeting Room A', 10),
('Meeting Room B', 6),
('Conference Hall', 50);
