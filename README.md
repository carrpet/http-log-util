# HTTPMonitorUtil

A metrics monitor and alerting CLI application that streams http log entries and displays
statistics about the logs and alerts on the volume of requests seen. 

## Getting Started

Download the project into your desired directory. The project was delivered in a compressed
archive using ```tar``` and named ```pcdatadog.tar```.
Unpack it by navigating to the directory where the .tar file
lives and run

```tar -zxvf pcdatadog.tar -C <path_to_extract_to>```

All remaining instructions assume that the commands are run from the root directory
(ie. directory where it was extracted to) of the project.  

### Prerequisites

You should have a very current version of Go. The app was developed and tested against go version 1.14.1.

## Running the Tests
All the tests live in the cmd package which is one directory below main in the /cmd directory.
To run all the tests in that directory navigate into the cmd directory:

```cd cmd```

then execute
```go test```

All tests should be executed and pass.

### Installing

You should first run the automated tests before installing.  See above section.

Navigate into the root directory of the project.  

