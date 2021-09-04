# HTTP Web Server


### Endpoints description

    1. /measurement
        GET - gets all sensor measurement records from time-series database.
        POST - save measurements passed from current sensor groups specified in the web_hook_url from 
        sensor cli application.

    2. /sensor
        GET - gets all sensor information from relational database.
        PUT - updates senosor information.
        POST - add new sensor in relational database. Sensor should have id, deviceId, name, description, unit, measurements properties.
        DELETE - delete chosen sensor information.

    3. /device
        GET - gets device information by current name.
        PUT - updates device information.
        POST -  add new sensor in relational database. Device should have id, name, description, sensors properties.
        DELETE - delete chosen device information.

### Start postres database

1. Build the docker image for postgredb(./pkg/database)

```
    docker build -t <image>:<tag> <dockerFilePath>
```

2. Run PostresDB:
```
    docker run --publish 5432:5432  <image>
```

2. Start PgAdmin

```
   docker run --rm -p 5050:5050 thajeztah/pgadmin4
```

If the IP is not published and can be accessed from outside you should get the internal IP by running:

```
docker inspect <containerID>
```

- Access to pgAdmin from the browser http://localhost:5050 (The port is that specified with -p flag in docker run command)



https://docs.influxdata.com/influxdb/v2.0/reference/flux/stdlib/built-in/transformations/range/