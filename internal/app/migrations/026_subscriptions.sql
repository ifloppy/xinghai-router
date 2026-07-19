create table if not exists subscription_plans (
  id uuid primary key,
  name text not null,
  description text not null default '',
  price numeric(20,2) not null default 0 check (price >= 0),
  currency text not null default 'CNY',
  billing_period text not null check (billing_period in ('month','year')),
  credit_amount numeric(20,8) not null default 0 check (credit_amount >= 0),
  group_id uuid references groups(id) on delete set null,
  model_whitelist text[] not null default '{}',
  max_requests_per_period bigint,
  max_tokens_per_period bigint,
  sort_order integer not null default 0,
  enabled boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create unique index if not exists subscription_plans_name_idx on subscription_plans(name);
create index if not exists subscription_plans_enabled_sort_idx on subscription_plans(enabled, sort_order, name);

create table if not exists user_subscriptions (
  id uuid primary key,
  user_id uuid not null references users(id) on delete cascade,
  plan_id uuid not null references subscription_plans(id) on delete cascade,
  status text not null default 'pending' check (status in ('pending','active','expired','cancelled')),
  current_period_start timestamptz,
  current_period_end timestamptz,
  auto_renew boolean not null default false,
  cancelled_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists user_subscriptions_user_idx on user_subscriptions(user_id, created_at desc);
create index if not exists user_subscriptions_active_idx on user_subscriptions(user_id, status, current_period_end);

create table if not exists subscription_orders (
  id uuid primary key,
  order_no text not null unique,
  subscription_id uuid not null references user_subscriptions(id) on delete cascade,
  user_id uuid not null references users(id) on delete cascade,
  plan_id uuid not null references subscription_plans(id) on delete cascade,
  provider text not null default 'epay' check (provider in ('epay')),
  payment_type text not null,
  amount numeric(20,2) not null check (amount > 0),
  status text not null default 'pending' check (status in ('pending','paid','failed','expired')),
  provider_trade_no text,
  period_kind text not null default 'new' check (period_kind in ('new','renewal')),
  paid_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists subscription_orders_user_idx on subscription_orders(user_id, created_at desc);
create index if not exists subscription_orders_subscription_idx on subscription_orders(subscription_id, created_at desc);
create unique index if not exists subscription_orders_provider_trade_idx on subscription_orders(provider, provider_trade_no) where provider_trade_no is not null;
