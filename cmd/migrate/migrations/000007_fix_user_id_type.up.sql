ALTER TABLE user_invitations
DROP COLUMN user_id,
ADD COLUMN user_id uuid NOT NULL;