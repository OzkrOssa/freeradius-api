--
-- Table structure for table 'nas'
--
CREATE TABLE nas
(
    id          serial PRIMARY KEY,
    nasname     text NOT NULL,
    shortname   text NOT NULL,
    type        text NOT NULL DEFAULT 'other',
    ports       integer,
    secret      text NOT NULL,
    server      text,
    community   text,
    description text
);
create index nas_nasname on nas (nasname);
