ALTER TABLE user_invitations
DROP COLUMN user_id,
ADD COLUMN user_id bigint NOT NULL