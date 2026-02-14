-- Migration: Extend rooms with type and status
ALTER TABLE rooms ADD COLUMN type VARCHAR(20) DEFAULT 'shared';
ALTER TABLE rooms ADD COLUMN status VARCHAR(20) DEFAULT 'online';

-- Update existing data
UPDATE rooms SET type = 'shared', status = 'online';
