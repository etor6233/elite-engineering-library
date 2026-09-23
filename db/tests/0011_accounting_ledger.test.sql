begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','accounting-test','Accounting','Accounting');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','franchise','franchise','Franchise','franchisee');
insert into accounting.account values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','CASH','Cash','asset',true,clock_timestamp()),('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','REVENUE','Revenue','revenue',true,clock_timestamp());
insert into accounting.period(tenant_id,period_id,starts_on,ends_on,status,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','2026-08','2026-08-01','2026-09-01','open',1);
insert into accounting.journal(tenant_id,journal_id,organization_id,period_id,source_type,source_id,currency,posting_date,status,total_debit_minor_units,total_credit_minor_units,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','journal','franchise','2026-08','SALE','order-1','ARS','2026-08-15','draft',100,100,1);
insert into accounting.journal_line values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','journal',1,'CASH','cash',100,0),('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','journal',2,'REVENUE','sale',0,100);
do $$ begin
  begin update accounting.journal_line set debit_minor_units=1 where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29e00'; raise exception 'line mutation accepted'; exception when raise_exception then if sqlerrm <> 'immutable accounting history' then raise; end if; end;
  begin insert into accounting.journal_line values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','journal',3,'CASH','bad',1,1); raise exception 'dual-sided line accepted'; exception when check_violation then null; end;
end $$;
rollback;
