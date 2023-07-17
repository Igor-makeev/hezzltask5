
CREATE TABLE IF NOT EXISTS events(
id         int,
campaign_id  int,
name String,
description String,
priority int,
removed Bool ,
EventTime TIMESTAMP,
INDEX index_id_events id TYPE minmax GRANULARITY 3,
INDEX index_campaign_id_events campaign_id TYPE minmax GRANULARITY 3,
INDEX index_name_events name TYPE bloom_filter GRANULARITY 3
)ENGINE = MergeTree
ORDER BY EventTime;


