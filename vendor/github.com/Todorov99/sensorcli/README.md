# Sensor CLI

 Sensor CLI is an application that gets the CPU temperature, CPU cores count, CPU frequency CPU usage in percentage, available RAM memory, used RAM memory from the running programmes, used RAM memory in percentage.

## Description

This application run on local OS and in a Docker container.
The application suports measurements from multiple sensor groups.

- Available groups:
   - CPU_TEMP
   - CPU_USAGE
   - MEMORY_USAGE

Sensor and device information is parsed from a config yaml file.

Sensor measurement data could be exported in JSON and YAML format.

### Installing

go build -o sensorcli ./

## Usage
```
  sensorcli cpu [flags]

```

## Flags:

```

      --delta_duration int     The period of which you will get your sensor data. (default 3)

      --format string          The data could be printed either JSON or YAML format. (default "JSON")

  -h, --help                   help for cpu

      --output_file string     Writing the output into CSV file.
      
      --sensor_group strings   There are three main sensor groups: CPU_TEMP, CPU_USAGE and MEMORY_USAGE

      --total_duration float   Terminating the whole program after specified duration (default 60)

      --web_hook_url string    Expose to current port.

```

## Use as go mod

export GOPRIVATE=github.com/Todorov99

go get -u github.com/Todorov99/sensorcli