# Crylog

<img align="right" src="https://img.shields.io/badge/License-MIT-blue.svg">

### Summary
 <a href="https://www.linkedin.com/in/sven-schaper/" style="text-align: right"><img align="right" src="https://www.lime-anchor.com/img/connect.png" height="250"></a><br>

Using **Crylog for Golang** you can add log entries with different severety in your golang project. You can define log entries for ERROR, COMMUNICATE, REPORT, INFO,  WARN and DEBUG. After you enriched your golang project with advanced log entries, you can put a config file next to your main.go file, where you can specify the loglevel you want to enable. It will follow the direction described above, so if you chose Info inside your config file it will print out Error, Communicate, Report and Info. Each level will be printed in a specific color to the log, so you are able to differentiate them visually. You can set the loglevel per package.

### Set Up

Download the project
```bash
go get https://github.com/svenschaper/crylog
```

Import the package into your go File
```golang
import "github.com/svenschaper/crylog"
```


### How to use?

Create a config.yml and put it next to your main.go file 
```
log.level: DEBUG                        //Default log level for all packages
log.level.yourPackageName: INFO         //Specific log level for the package yourPackageName
log.level.main: ERROR                   //Specific log level for the package main
```

Add the following to your go files
```golang
var logger logging.Logger

func init() {
	logger = logging.GeneralInitLogger("yourPackageName")
}


func yourFunction(){
    logger.Info("This is a green info message")
	logger.Error("This is a red error message")
	logger.Debug("This is a blue debug message")
	logger.Warn("This is a purple warn message")
	logger.Report("This is a cyan report message")
    logger.Communicate("This is yellow communication message")
    

	logger.InfoWithC("This is a green info message", "This part of the message should be encrypted")
}

```





<br>
<br>
<br>
<br>
<br>


### Need custom development?

<a href="https://lime-anchor.com"><img align="right" src="https://www.lime-anchor.com/img/gint.png" height="50"></a>

**Cloud Computing**
* Infrastructure automation based on ansible
* Custom solutions based on AWS
* Custom solutions based on Google Cloud
* Custom solutions based on Azure Cloud

**Integration**
* API Management and Design
* Lightning fast Microservices based on Golang



# License

MIT licensed. In short -> Have fun with it!
