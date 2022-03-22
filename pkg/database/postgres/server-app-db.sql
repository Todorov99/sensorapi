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
	VALUES ('cpuTempCelsius', 'Measures CPU temp Celsius', 'C', '1', '1');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('cpuUsagePercent', 'Measures CPU usage percent', '%', '1', '2');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('cpuCoresCount', 'Measures CPU cores count', 'count', '1', '2');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('cpuFrequency', 'Measures CPU frequency', 'GHz', '1', '2');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('memoryTotal', 'Measures memory total', 'GigaBytes', '1', '3');

INSERT INTO sensor(
	name, description, unit, device_id, sensor_groups_id)
	VALUES ('memoryAvailableBytes', 'Measures memory available Bytes', 'Bytes', '1', '3');

INSERT INTO sensor(
	name, description, unit, device_id, sensor_groups_id)
	VALUES ('memoryUsedBytes', 'Measures memory used Bytes', 'Bytes', '1', '3');

INSERT INTO sensor(
	 name, description, unit, device_id, sensor_groups_id)
	VALUES ('memoryUsedPercent', 'Measures memory used percent', '%', '1', '3');    

