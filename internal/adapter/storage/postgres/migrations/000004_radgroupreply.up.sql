--
-- Table structure for table 'radgroupreply'
--
CREATE TABLE radgroupreply (
    id			serial PRIMARY KEY,
    GroupName		text NOT NULL DEFAULT '',
    Attribute		text NOT NULL DEFAULT '',
    op			VARCHAR(2) NOT NULL DEFAULT '=',
    Value			text NOT NULL DEFAULT ''
);
create index radgroupreply_GroupName on radgroupreply (GroupName,Attribute);