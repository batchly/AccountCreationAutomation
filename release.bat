@ECHO OFF
CLS

ECHO Cleaning up bin directory
DEL ..\..\bin\accountcreationautomation.exe
DEL ..\..\bin\linux_amd64\accountcreationautomation
DEL ..\..\bin\darwin_amd64\accountcreationautomation


ECHO Starting Windows build
SET GOOS=windows
go install
ECHO Windows build complete

ECHO Starting Linux build
SET GOOS=linux
go install
ECHO Linux build complete

ECHO Starting Mac build
SET GOOS=darwin
SET GOARCH=amd64
go install
ECHO Mac build complete

ECHO Resetting build parameters
SET GOOS=windows

ECHO Build output in bin directory
