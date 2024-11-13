--
-- Table structure for table 'radpostauth'
--
CREATE TABLE radpostauth (
    id			bigserial PRIMARY KEY,
    username		text NOT NULL,
    pass			text,
    reply			text,
    CalledStationId		text,
    CallingStationId	text,
    authdate		timestamp with time zone NOT NULL default now(),
    Class			text
);
