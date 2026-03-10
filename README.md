# Go TUI Evaluator

A sandbox application to evaluate different Go TUI libraries (Bubbletea, Tview, Vaxis, PTerm) displaying a shared mocked dashboard service.

## Usage

Build the application:
```bash
go build -o myapp main.go
```

Run the different frontends:
```bash
./myapp bubbletea
./myapp tview
./myapp vaxis
./myapp pterm
```

## Architecture
- `internal/dashboard/` holds the decoupled mock data service.
- `internal/ui/` contains isolated UI implementations for each framework.
- The UI implementations adhere strictly to the shared `MetricsProvider` interface.
