# AndersenTest
This simple webservice grabs a bunch of children names from 
https://catalog.data.gov/dataset/most-popular-baby-names-by-sex-and-mothers-ethnic-group-new-york-city-8c742
and responds with a new JSON with 5 random child names in it.
# Launch
Just build it via
```
go build main.go
```
Then launch the outbut Build.exe (for windows) using command line
```
go main.exe
```
The server will be set up on **localhost:3000** and the response will be sent there.
