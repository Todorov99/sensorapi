FROM postgres:12.3

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=Password321
ENV POSTGRES_PORT=5432
ENV POSTGRES_DB=sensorCLI

#VOLUME "var/lib/postgresdb" 
COPY server-app-db.sql /docker-entrypoint-initdb.d/001_data.sql

