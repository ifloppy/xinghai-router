create table if not exists email_verification_codes (
	id uuid primary key default gen_random_uuid(),
	email text not null,
	code_hash text not null,
	purpose text not null default 'register',
	attempts integer not null default 0,
	expires_at timestamptz not null,
	consumed_at timestamptz,
	created_at timestamptz not null default now()
);

create index if not exists email_verification_codes_email_idx on email_verification_codes(email, purpose, created_at desc);
