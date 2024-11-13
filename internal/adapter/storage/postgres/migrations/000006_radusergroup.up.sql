--
-- Table structure for table 'radusergroup'
--
CREATE TABLE radusergroup (
                              id			serial PRIMARY KEY,
                              UserName		text NOT NULL DEFAULT '',
                              GroupName		text NOT NULL DEFAULT '',
                              priority		integer NOT NULL DEFAULT 0
);
create index radusergroup_UserName on radusergroup (UserName);
--
-- Use this index if you use case insensitive queries
--
-- create index radusergroup_UserName_lower on radusergroup (lower(UserName));