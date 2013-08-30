GoConvey - BDD in Go
====================






RoadMap:
--------

TODO:

	- Reporting
	- Presenters
		- Story
			- Markdown (lists, errors are block quotes)
		- Dot
	- Randomized execution of stories (including resets)
	- Idle (tests re-run at every save)


Would be awesome:

	- Output Story presentation to HTML file (https://github.com/russross/blackfriday)
	- Create http endpoint that serves the html output
	- make http endpoint poll for updates reload report (collapse all but failed and erred stuff)
	- clicking on filename/line-number in web report shows that file as a web page w/ problem line highlighted