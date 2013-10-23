High-level structs:
===================

    Scanner
      - scans for changes to .go files
      - update watcher with creation and deletion of new packages

[X] Watcher
      - initiates test executor
      - adjusts watched packages (adjust, ignore, etc...)
      - knows which packages are active

    Server
      - instructs watcher to adjust/ignore watched packages
      - provides status of executor
      - can trigger an adhoc test run from the executor
      - provides latest results


    Mid-level structs:
    ==================

    Executor
      - shares/receives a list of packages to test from watcher
      - runs tests across 'x' goroutines
      - passes output to parser for parsing
      - provides aggregated result to server

[X]   Parser
        - receives test output from executor
        - parses each package and aggregates a result
