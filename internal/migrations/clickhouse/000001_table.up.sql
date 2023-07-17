
CREATE TABLE IF NOT EXISTS events
(
 id         int,
campaign_id  int,
name String,
description String,
priority int,
removed Bool ,
created_at TIMESTAMP 
)ENGINE = Log;




