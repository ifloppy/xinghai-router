alter table groups drop constraint if exists groups_multiplier_check;
alter table groups add constraint groups_multiplier_check check (multiplier >= 0);
