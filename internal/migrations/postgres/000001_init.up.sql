CREATE TABLE IF NOT EXISTS campaings
(
    id         serial  primary key,
  name text not null unique
);

create index if not exists index_id_campaings on campaings (id);

INSERT INTO campaings (name)
    VALUES ('запись 1');

CREATE TABLE IF NOT EXISTS items
(
 id         serial primary key,
campaign_id  int references campaings(id) on delete cascade not null,
name text not null,
description text default '',
priority int,
removed boolean default false,
created_at TIMESTAMPTZ default now()
);

CREATE OR REPLACE FUNCTION items_priority_max() RETURNS int
    AS 'SELECT COALESCE(max(priority),0)+1 FROM items'
    LANGUAGE 'sql';

 alter table items alter column priority set default items_priority_max();

create index if not exists index_id_items on items (id);
create index if not exists index_campaign_id_items on items (campaign_id);
create index if not exists index_name_items on items (name);