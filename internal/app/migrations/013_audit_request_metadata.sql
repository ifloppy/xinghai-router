alter table audit_logs
  add column if not exists client_ip text not null default '',
  add column if not exists forwarded_for text not null default '',
  add column if not exists user_agent text not null default '',
  add column if not exists browser text not null default '',
  add column if not exists browser_version text not null default '',
  add column if not exists operating_system text not null default '',
  add column if not exists operating_system_version text not null default '',
  add column if not exists device_type text not null default '',
  add column if not exists is_bot boolean not null default false,
  add column if not exists request_method text not null default '',
  add column if not exists request_path text not null default '',
  add column if not exists request_id text not null default '';

create index if not exists audit_logs_client_ip_idx on audit_logs(client_ip, created_at desc);
