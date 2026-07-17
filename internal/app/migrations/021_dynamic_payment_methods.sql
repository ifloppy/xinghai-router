alter table payment_orders drop constraint if exists payment_orders_payment_type_check;

alter table payment_settings drop column if exists payment_types;

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
