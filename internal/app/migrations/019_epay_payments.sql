create table if not exists payment_orders (
  id uuid primary key,
  order_no text not null unique,
  user_id uuid not null references users(id) on delete cascade,
  provider text not null check (provider in ('epay')),
  payment_type text not null check (payment_type in ('alipay','wxpay','qqpay')),
  amount numeric(20,2) not null check (amount > 0),
  status text not null default 'pending' check (status in ('pending','paid','failed','expired')),
  provider_trade_no text,
  paid_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists payment_orders_user_idx on payment_orders(user_id, created_at desc);
create unique index if not exists payment_orders_provider_trade_idx on payment_orders(provider, provider_trade_no) where provider_trade_no is not null;
