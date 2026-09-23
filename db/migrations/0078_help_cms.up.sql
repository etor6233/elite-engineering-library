begin;
create schema help;
create table help.article(
 tenant_id uuid not null,
 article_id text not null check(length(article_id)between 1 and 128),
 organization_id text not null,
 locale text not null check(locale in('es','en')),
 category text not null check(category~'^[a-z][a-z0-9_-]{0,63}$'),
 current_version bigint not null check(current_version>0),
 primary key(tenant_id,article_id),
 foreign key(tenant_id,organization_id)references org.organization(tenant_id,organization_id)
);
create table help.revision(
 tenant_id uuid not null,
 article_id text not null,
 version bigint not null check(version>0),
 command_id text not null check(length(command_id)between 1 and 128),
 action text not null check(action in('create','update','publish','archive')),
 actor_subject text not null check(length(actor_subject)between 1 and 256),
 request_sha256 text not null check(request_sha256~'^[a-f0-9]{64}$'),
 article_sha256 text not null check(article_sha256~'^[a-f0-9]{64}$'),
 snapshot jsonb not null check(jsonb_typeof(snapshot)='object'),
 title text not null check(octet_length(title)between 1 and 200),
 body text not null check(octet_length(body)between 1 and 16384),
 state text not null check(state in('draft','published','archived')),
 primary key(tenant_id,article_id,version),
 unique(tenant_id,command_id),
 foreign key(tenant_id,article_id)references help.article(tenant_id,article_id)
);
alter table help.article add constraint help_current_revision_fk foreign key(tenant_id,article_id,current_version)references help.revision(tenant_id,article_id,version)deferrable initially deferred;
create function help.guard_revision()returns trigger language plpgsql as $$begin
 if tg_op<>'INSERT'then raise exception 'help revision is immutable';end if;
 if new.snapshot->>'id' is distinct from new.article_id or new.snapshot->>'version' is distinct from new.version::text
 or new.snapshot->>'state' is distinct from new.state or new.snapshot->>'title' is distinct from new.title or new.snapshot->>'body' is distinct from new.body
 or not exists(select 1 from help.article h where h.tenant_id=new.tenant_id and h.article_id=new.article_id
 and h.organization_id=new.snapshot->>'organization_id'and h.locale=new.snapshot->>'locale'and h.category=new.snapshot->>'category')
 then raise exception 'help snapshot differs from head';end if;return new;
end$$;
create trigger help_revision_guard before insert or update or delete on help.revision for each row execute function help.guard_revision();
create function help.guard_head()returns trigger language plpgsql as $$begin
 if tg_op='DELETE'then raise exception 'help head evidence cannot be deleted';end if;
 if row(new.tenant_id,new.article_id,new.organization_id,new.locale,new.category)is distinct from row(old.tenant_id,old.article_id,old.organization_id,old.locale,old.category)
 or new.current_version<>old.current_version+1 then raise exception 'help head scope or version changed';end if;return new;
end$$;
create trigger help_head_guard before update or delete on help.article for each row execute function help.guard_head();
create index help_scope_page_idx on help.article(tenant_id,organization_id,locale,article_id);
commit;
