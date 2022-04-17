CREATE TABLE users (
   id serial NOT NULL,
   user_name     VARCHAR(1000) NOT NULL,
   pass     VARCHAR(1000) NOT NULL,
   first_name    VARCHAR(255) NOT NULL,
   last_name     VARCHAR(1000) NOT NULL,
   email     VARCHAR(1000) NOT NULL
);

CREATE TABLE device (
   id serial NOT NULL,
   name    VARCHAR(255) NOT NULL,
   description     VARCHAR(1000) NOT NULL,
   CONSTRAINT device_id PRIMARY KEY (id)
);

CREATE TABLE sensor_groups (
   id   serial       NOT NULL,
   group_name    VARCHAR(255) NOT NULL,
   CONSTRAINT sensor_groups_id PRIMARY KEY (id)
);

CREATE TABLE sensor (
   id   serial       NOT NULL ,
   name    VARCHAR(255) NOT NULL,
   description     VARCHAR(1000) NOT NULL,
   unit VARCHAR(15) NOT NULL,
   sensor_groups_id integer NOT NULL REFERENCES sensor_groups(id),
   CONSTRAINT sensor_id PRIMARY KEY (id)
);

CREATE TABLE device_sensor (
   device_id  integer REFERENCES device(id) ON UPDATE CASCADE ON DELETE CASCADE,
   sensor_id  integer REFERENCES sensor(id) ON UPDATE CASCADE ON DELETE CASCADE,
   CONSTRAINT device_sensor_pkey PRIMARY KEY (device_id, sensor_id)
);

INSERT INTO device(
	name, description)
	VALUES ('device_name', 'my laptop device');

 INSERT INTO sensor_groups(
	group_name)
	VALUES 
		('CPU_TEMP'),
		('CPU_USAGE'),
		('MEMORY_USAGE');

INSERT INTO sensor(
	name, description, unit, sensor_groups_id)
	VALUES 
		('cpuTemperature', 'Measures CPU temperature in provided unit', 'C', '1'),
		('cpuUsagePercentage', 'Measures CPU usage in percentages', '%', '2'),
		('cpuCores', 'Gets the number of CPU cores', 'count', '2'),
		('cpuFrequency', 'Measures CPU frequency in a provided unit', 'GHz', '2'),
		('memoryTotal', 'Measures memory total RAM', 'GigaBytes', '3'),
		('memoryAvailable', 'Gets the available RAM in a provided unit', 'Bytes', '3'),
		('memoryUsed', 'Gets the used RAM from the programs in a provided unit', 'Bytes', '3'),
		('memoryUsedPercentage', 'Used percentage RAM from the programs', '%', '3');

INSERT INTO device_sensor(
	device_id, sensor_id)
	VALUES 
		('1', '1'), 
		('1', '2'),
		('1', '3'),
		('1', '4'),
		('1', '5'),
		('1', '6'),
		('1', '7'),
		('1', '8');