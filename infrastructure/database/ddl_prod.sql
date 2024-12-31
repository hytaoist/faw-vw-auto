Version 1.0（24.12.31）
-------------------
-- test.faw_auth definition
-- Drop table
-- DROP TABLE test.faw_auth;
CREATE TABLE faw_auth (
	id serial4 NOT NULL,
	access_token text NOT NULL,
	token_type varchar NOT NULL,
	expires_in varchar NULL,
	create_at timestamptz NULL,
	CONSTRAINT faw_auth_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN prod.faw_auth.id IS '主键';
COMMENT ON COLUMN prod.faw_auth.access_token IS '授权Token';
COMMENT ON COLUMN prod.faw_auth.token_type IS 'Token类型';
COMMENT ON COLUMN prod.faw_auth.expires_in IS 'Token有效时间（毫秒级）';
COMMENT ON COLUMN prod.faw_auth.create_at IS '创建时间';
