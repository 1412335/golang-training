# Notes Service

This is the Notes service

Generated with

```
micro new notes
```

## Usage

Generate the proto code

```
make proto
```

Run the service

```
micro run .
micro update .
micro logs -f notes
```

## Test
command
```
micro notes list
```

go test
```
go test --run TestNotes_UpdateStream ./handler
```