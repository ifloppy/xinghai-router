alter table site_settings add column if not exists auto_disable_failed_channels boolean not null default false;
