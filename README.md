# Webb forum för försvarsmakten

## Installation
Installera mysql genom XAMPP eller liknande program, sedan skapar du en databas.
Default namnet på databasen är webbforum.

Installera Golang och klona detta projekt in i %GOPATH%/src.
Läs mer om %GOPATH% på https://golang.org/doc/install och google.

Sen krävs det att dep är installerat för att enklast installera libraries/dependencies, https://github.com/golang/dep.
Enklast är att köra kommandot:
```
https://github.com/golang/dep
```

När dep är installerat så kan man ladda ner alla dependencies/libraries genom att köra kommandot:
```
dep ensure
```
I mappen som detta projekt är klonat i.

## Kör programmet
När alla dependencies är nedladdat så är det enklast att köra kommandot
```
go run cmd/webbforum/webbforum.go
```
För att starta webbforumet.

Om du skapade en databas med default namnet och inte har ändrat någon port/login så ska allt funka.
Om du har ändrat något så ska en fil som heter `config.json` ha skapats där du kan ändra olika värden.