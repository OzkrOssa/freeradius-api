--
-- Table structure for table 'radcheck'
--
CREATE TABLE radcheck (
                          id			serial PRIMARY KEY,
                          UserName		text NOT NULL DEFAULT '',
                          Attribute		text NOT NULL DEFAULT '',
                          op			VARCHAR(2) NOT NULL DEFAULT '==',
                          Value			text NOT NULL DEFAULT ''
);
create index radcheck_UserName on radcheck (UserName,Attribute);
--
-- Use this index if you use case insensitive queries
--
-- create index radcheck_UserName_lower on radcheck (lower(UserName),Attribute);