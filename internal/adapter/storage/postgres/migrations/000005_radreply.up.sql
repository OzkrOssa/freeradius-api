--
-- Table structure for table 'radreply'
--
CREATE TABLE radreply (
    id			serial PRIMARY KEY,
    UserName		text NOT NULL DEFAULT '',
    Attribute		text NOT NULL DEFAULT '',
    op			VARCHAR(2) NOT NULL DEFAULT '=',
    Value			text NOT NULL DEFAULT ''
);
create index radreply_UserName on radreply (UserName,Attribute);
--
-- Use this index if you use case insensitive queries
--
-- create index radreply_UserName_lower on radreply (lower(UserName),Attribute);