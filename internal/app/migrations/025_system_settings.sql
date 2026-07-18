alter table site_settings
	add column if not exists geetest_captcha_id text not null default '',
	add column if not exists geetest_captcha_key_encrypted text not null default '',
	add column if not exists smtp_host text not null default '',
	add column if not exists smtp_port text not null default '465',
	add column if not exists smtp_username text not null default '',
	add column if not exists smtp_password_encrypted text not null default '',
	add column if not exists smtp_from text not null default '';
