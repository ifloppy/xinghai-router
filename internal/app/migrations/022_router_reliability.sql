alter table site_settings add column if not exists retry_count integer not null default 3 check (retry_count between 0 and 10);
alter table site_settings add column if not exists retry_status_codes text not null default '100-199,300-407,409-503,505-523,525-599';
alter table site_settings add column if not exists health_check_mode text not null default 'off' check (health_check_mode in ('off','scheduled_all','passive_recovery'));
alter table site_settings add column if not exists health_check_interval_minutes integer not null default 5 check (health_check_interval_minutes between 1 and 1440);
alter table site_settings add column if not exists health_check_auto_recover boolean not null default true;
alter table site_settings add column if not exists health_check_channel_ids text not null default '';
alter table site_settings add column if not exists auto_disable_on_test_failure boolean not null default false;
alter table site_settings add column if not exists auto_disable_slow_seconds integer not null default 0 check (auto_disable_slow_seconds between 0 and 600);
alter table site_settings add column if not exists auto_disable_status_codes text not null default '401,429,503';
alter table site_settings add column if not exists auto_disable_keywords text not null default 'Your credit balance is too low
This organization has been disabled.
You exceeded your current quota
Permission denied
The security token included in the request is invalid
Operation not allowed
Your account is not authorized
订阅额度不足或未配置订阅
所有账号暂时不可用
已达到 Token Plan 用量上限
Weekly usage limit reached.
5-hour usage limit reached
Invalid token
Too Many Requests
You have exceeded the monthly usage quota
You have exceeded the weekly usage quota. It will reset at ';

alter table channels add column if not exists auto_disabled boolean not null default false;
alter table channels add column if not exists disabled_reason text not null default '';
