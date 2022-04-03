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
   device_id integer NOT NULL REFERENCES device(id),
   sensor_groups_id integer NOT NULL REFERENCES sensor_groups(id),
   CONSTRAINT sensor_id PRIMARY KEY (id)
);

INSERT INTO device(
	name, description)
	VALUES ('device_name', 'my laptop device');

 INSERT INTO sensor_groups(
	group_name)
	VALUES ('CPU_TEMP');

INSERT INTO sensor_groups(
	group_name)
	VALUES ('CPU_USAGE');

INSERT INTO sensor_groups(
	group_name)
	VALUES ('MEMORY_USAGE');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('cpuTemperature', 'Measures CPU temperature in provided unit', 'C', '1', '1');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('cpuUsagePercentage', 'Measures CPU usage in percentages', '%', '1', '2');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('cpuCores', 'Gets the number of CPU cores', 'count', '1', '2');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('cpuFrequency', 'Measures CPU frequency in a provided unit', 'GHz', '1', '2');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('memoryTotal', 'Measures memory total RAM', 'GigaBytes', '1', '3');

INSERT INTO sensor(
	name, description, unit, device_id, sensor_groups_id)
	VALUES ('memoryAvailable', 'Gets the available RAM in a provided unit', 'Bytes', '1', '3');

INSERT INTO sensor(
	name, description, unit, device_id, sensor_groups_id)
	VALUES ('memoryUsed', 'Gets the used RAM from the programs in a provided unit', 'Bytes', '1', '3');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('memoryUsedPercentage', 'Used percentage RAM from the programs', '%', '1', '3');    

