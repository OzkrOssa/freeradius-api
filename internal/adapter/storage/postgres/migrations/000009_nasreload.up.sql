--
-- Table structure for table 'nasreload'
--

CREATE TABLE IF NOT EXISTS nasreload (
    NASIPAddress		inet PRIMARY KEY,
    ReloadTime		timestamp with time zone NOT NULL
);