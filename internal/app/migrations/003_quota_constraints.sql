alter table quota_limits drop constraint if exists quota_limits_user_id_api_key_id_model_window_key;
create unique index if not exists quota_limits_scope_idx on quota_limits (
  coalesce(user_id, '00000000-0000-0000-0000-000000000000'::uuid),
  coalesce(api_key_id, '00000000-0000-0000-0000-000000000000'::uuid),
  coalesce(model, ''),
  "window"
);
