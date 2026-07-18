alter table users
	add column if not exists leaderboard_opt_in boolean not null default true,
	add column if not exists leaderboard_mask_name boolean not null default true;
