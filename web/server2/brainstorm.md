
--------------------------------

- scan for all file system items
- filter on depth
- identfy all folders/packages (by grouping items and profiles?)
- flag all ignored packages (UI)
- parse all watched profiles
- flag the disabled packages (profiles)
- calculate checksums of items in watched packages
- send all packages with appropriate flags to executor (and consequently to the server)

---------------------------------------------------------------------------------------


------------------------------------

- Scan returns a chan FileSystemItem
- We consume the channel once, filtering on depth, creating maps that reveal the folder structure (one map) and the profiles (another map). Maps are keyed by containing folder path.
- maps are scanned and flagged (ignored, profile parsed to discover disabled packages) and checksummed

------------------------------------------------------------------------------------------------------