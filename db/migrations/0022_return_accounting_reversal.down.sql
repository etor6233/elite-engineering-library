begin;
drop index if exists accounting.return_effect_posting_source_idx;
drop trigger if exists return_effect_posting_immutable on accounting.return_effect_posting;
drop function if exists accounting.prevent_return_effect_posting_mutation();
drop table if exists accounting.return_effect_posting;
commit;
