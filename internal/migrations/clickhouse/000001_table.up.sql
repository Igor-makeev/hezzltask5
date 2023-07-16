
CREATE TABLE IF NOT EXISTS events
(
 id         int,
campaign_id  int,
name text,
description text,
priority int,
removed UInt8 ,
created_at TIMESTAMP 
)ENGINE = Log;



