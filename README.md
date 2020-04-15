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

Navigate into the root directory of the project. Build an executable with

```go build```

The result of this puts an executable file called ```httpmonitorutil``` into your root directory.


### Running

You can then run this in the normal Unix fashion, passing it a required filename argument which serves as the input to the program:

```./httpmonitorutil <csv_file_path_ and_name> ```
  
The file must have a CSV format with headers on the first row.  The program should now be scanning
the file and producing output to stdout.

You can pass an optional flag as in:

```.httpmonitorutil <csv_file_path_and_name> --alertThreshold``` to customize the threshold at which the alerts fire.  Details can be found if you type ```./httpmonitorutil```


### Comments 

I modeled the log monitoring application as a pipeline that could take a varying number of stages. My goal was to build a generic pipeline with enough abstraction that it could be easily extended to other domains and to facilitate adding processing facility and output streams to the current application. I chose this abstraction because Go's channel and goroutine constructs make this relatively simple to model and in some ways makes it easier to reason about the program, since real-time streaming domain is a natural fit for this model.  

Tradeoffs for the genericity and asynchronous nature of the program are that there are some additional inefficiencies in the code.  More than a few interfaces to keep track of and a decent amount of type assertions and sometimes copying data around to fit in the correct type container was required given Go's lack of generics. I believe there's enough code reuse to justify this decision however.  Then there is the added complexity of managing and closing channels so that all the goroutines can make progress. The tradeoff was accepted to increase the performance of the application.

### Improvements / Future Work

Given my reflections about the design, I believe it is fairly robust and performant as it stands. 

The csv reader uses buffered io (golang bufio) under the covers, so that most read calls should not be too expensive, as most of them won't invoke the OS's read file method. Only one row is read into memory at a time and is sent downstream for other goroutines to do work so both scaling the size of the file and computing the results should be efficient.  If file size becomes unwieldy or if the file contents change such that the reading is no longer efficient then a reader can be implemented built on top of bufio directly and it can be fine tuned for the specific needs.

If the row processing facility becomes constrained (perhaps the number of columns in the csv become very big) we can scale up the number of worker goroutines that are reading from the csv source channel.  The current pipeline would support this with not too much work.  We would need a fan-out stage for the log reader and then multiple goroutines can read from it and do work. We would need some concurrency constructs and a way to merge the results at the end.  Similarly, if the metrics computation requirements become very intensive we can scale the goroutines that perform the computation for each stage to run them in parallel and then merge the result.  We need a good understanding of concurrency in Go but the current platform supports the basic structure.   

Error handling is one area where the current application needs work.  Currently, the csv reader can write errors to an error channel but nobody is listening on this channel.  The downstream processes should receive the errors and handle them accordingly.  In the simplest case this could signal that the application should stop but potentially each stage could handle them differently,and/or the main goroutine could be the ultimate arbiter of whether the application should continue.

Explicit cancellation and better channel management is another area of improvement.  There is currently no way for the main goroutine to tell the others that they should stop.  This is an important feature to have for increased versatility and reliabilty.  Currently most goroutine creation is encapsulated in the pipeline's Start function.  This is clean but does not give the caller the power to stop them.  A caller of a goroutine should be able to stop it.

### Time spent

I spent about 4.5 days of coding on this project plus a little time for design.
