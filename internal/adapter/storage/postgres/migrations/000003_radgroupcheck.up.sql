--
-- Table structure for table 'radgroupcheck'
--
CREATE TABLE radgroupcheck (
                               id			serial PRIMARY KEY,
                               GroupName		text NOT NULL DEFAULT '',
                               Attribute		text NOT NULL DEFAULT '',
                               op			VARCHAR(2) NOT NULL DEFAULT '==',
                               Value			text NOT NULL DEFAULT ''
);
create index radgroupcheck_GroupName on radgroupcheck (GroupName,Attribute);