create table users (
  id uuid primary key,
  email text not null unique,
  name text not null,
  role text not null default 'user' check (role in ('admin','operator','user')),
  enabled boolean not null default true,
  created_at timestamptz not null default now()
);
create table api_keys (
  id uuid primary key,
  user_id uuid not null references users(id) on delete cascade,
  name text not null,
  key_prefix text not null,
  secret_hash text not null unique,
  expires_at timestamptz,
  revoked_at timestamptz,
  last_used_at timestamptz,
  created_at timestamptz not null default now()
);
create table channels (
  id uuid primary key,
  name text not null unique,
  base_url text not null,
  api_key text not null,
  models jsonb not null default '[]',
  enabled boolean not null default true,
  priority integer not null default 100,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
create table request_logs (
  id uuid primary key,
  request_id text not null unique,
  user_id uuid references users(id) on delete set null,
  api_key_id uuid references api_keys(id) on delete set null,
  channel_id uuid references channels(id) on delete set null,
  model text not null,
  status_code integer not null,
  prompt_tokens integer,
  completion_tokens integer,
  total_tokens integer,
  duration_ms integer not null,
  error_code text,
  created_at timestamptz not null default now()
);
create index request_logs_created_at_idx on request_logs(created_at desc);
create index request_logs_user_id_idx on request_logs(user_id, created_at desc);
