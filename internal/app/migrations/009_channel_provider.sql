alter table channels add column if not exists provider text not null default 'openai';
alter table channels drop constraint if exists channels_provider_check;
alter table channels add constraint channels_provider_check check (provider in ('openai','ollama','kimi','opencode_go','anthropic'));
