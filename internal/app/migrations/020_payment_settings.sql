create table if not exists payment_settings (
  provider text primary key check (provider in ('epay')),
  enabled boolean not null default false,
  base_url text not null default '',
  merchant_id text not null default '',
  merchant_key_encrypted text not null default '',
  public_base_url text not null default '',
  updated_at timestamptz not null default now()
);

insert into payment_settings(provider) values('epay') on conflict do nothing;

create table if not exists payment_methods (
  id uuid primary key,
  provider text not null references payment_settings(provider) on delete cascade,
  code text not null,
  name text not null,
  enabled boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique(provider, code)
);

create index if not exists payment_methods_provider_idx on payment_methods(provider, enabled, created_at);
