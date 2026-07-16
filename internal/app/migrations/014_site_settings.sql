create table if not exists site_settings (
    id boolean primary key default true check (id),
    name text not null default 'Xinghai Router',
    icon_url text not null default '',
    updated_at timestamptz not null default now()
);
insert into site_settings(id) values(true) on conflict (id) do nothing;
