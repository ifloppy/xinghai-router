alter table channels add column if not exists weight integer not null default 100 check (weight > 0);
alter table channels add column if not exists failure_count integer not null default 0;
alter table channels add column if not exists cooldown_until timestamptz;
alter table channels add column if not exists last_checked_at timestamptz;
alter table channels add column if not exists last_error text;

create table if not exists model_routes (
  id uuid primary key,
  public_model text not null,
  upstream_model text not null,
  channel_id uuid references channels(id) on delete cascade,
  priority integer not null default 100,
  weight integer not null default 100 check (weight > 0),
  enabled boolean not null default true,
  created_at timestamptz not null default now(),
  unique(public_model, channel_id)
);
create index if not exists model_routes_public_model_idx on model_routes(public_model, priority);

create table if not exists pricing_rules (
  id uuid primary key,
  model text not null unique,
  input_per_million numeric(20,8) not null default 0,
  cached_input_per_million numeric(20,8) not null default 0,
  output_per_million numeric(20,8) not null default 0,
  multiplier numeric(12,6) not null default 1 check (multiplier > 0),
  enabled boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists user_wallets (
  user_id uuid primary key references users(id) on delete cascade,
  balance numeric(20,8) not null default 0 check (balance >= 0),
  reserved numeric(20,8) not null default 0 check (reserved >= 0),
  updated_at timestamptz not null default now()
);
create table if not exists wallet_ledger (
  id uuid primary key,
  user_id uuid not null references users(id) on delete cascade,
  amount numeric(20,8) not null,
  balance_after numeric(20,8) not null,
  kind text not null check (kind in ('topup','reservation','charge','release','refund','adjustment')),
  request_id text,
  note text,
  created_at timestamptz not null default now()
);
create index if not exists wallet_ledger_user_idx on wallet_ledger(user_id, created_at desc);

create table if not exists usage_records (
  id uuid primary key,
  request_id text not null unique references request_logs(request_id) on delete cascade,
  user_id uuid not null references users(id) on delete cascade,
  api_key_id uuid references api_keys(id) on delete set null,
  model text not null,
  prompt_tokens integer not null default 0,
  cached_prompt_tokens integer not null default 0,
  completion_tokens integer not null default 0,
  cost numeric(20,8) not null default 0,
  status text not null default 'settled',
  created_at timestamptz not null default now()
);
create index if not exists usage_records_user_idx on usage_records(user_id, created_at desc);

create table if not exists quota_limits (
  id uuid primary key,
  user_id uuid references users(id) on delete cascade,
  api_key_id uuid references api_keys(id) on delete cascade,
  model text,
  "window" text not null check ("window" in ('minute','day','month')),
  max_requests bigint,
  max_tokens bigint,
  unique(user_id, api_key_id, model, "window")
);

create table if not exists audit_logs (
  id uuid primary key,
  action text not null,
  actor text not null,
  entity_type text not null,
  entity_id text,
  details jsonb not null default '{}',
  created_at timestamptz not null default now()
);
create index if not exists audit_logs_created_idx on audit_logs(created_at desc);

insert into user_wallets(user_id) select id from users on conflict do nothing;
