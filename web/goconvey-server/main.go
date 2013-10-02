package main

import (
	"encoding/json"
	"fmt"
	"github.com/howeyc/fsnotify"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	go reactToChanges(watcher)

	working, _ := os.Getwd()
	watcher.Watch(working)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/latest", reportHandler)
	http.ListenAndServe(":8080", nil)
}

func reactToChanges(watcher *fsnotify.Watcher) {
	busy := false
	done := make(chan bool)
	ready := make(chan bool)

	for {
		select {
		case ev := <-watcher.Event:
			if strings.HasSuffix(ev.Name, ".go") && !busy {
				busy = true
				go test(done)
			}

		case err := <-watcher.Error:
			fmt.Println(err)

		case <-done:
			// TODO: rethink this delay?
			time.AfterFunc(1500*time.Millisecond, func() {
				ready <- true
			})

		case <-ready:
			busy = false
		}
	}
}

func test(done chan bool) {
	fmt.Println("Running tests...")

	// TODO: recurse into subdirectories and run tests...
	// oh yeah, and always check for new packages sprouting up,
	// or existing ones being removed...
	output, _ := exec.Command("go", "test", "-json").Output()
	result := parsePackageResult(string(output))

	serialized, _ := json.Marshal(result)
	// var buffer bytes.Buffer
	// json.Indent(&buffer, serialized, "", "  ")

	latest = string(serialized)
	done <- true
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, HOME_HTML)
}

func reportHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, latest)
}

var latest string
var HOME_HTML = `<!DOCTYPE html>
<html>

<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>Latest GoConvey Execution Report</title>

	<script type="text/javascript" src="https://ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
	<script type="text/javascript" src="https://raw.github.com/adammark/Markup.js/master/src/markup.min.js"></script>
	<script type="text/javascript">
        $(document).ready(function() {
            $("#output").hide();

            setInterval(function() {
                $.ajax({
                    url: "/latest",
                    complete: function(xhr, status) {
                    },
                    success: function(data) {
                    	console.log(data);

                        $("#output").hide();
                        $("#output").empty();
                        $("#output").html(data);
                        $("#output").slideDown();
                    },
                    async: false
                });
            }, 2000);
        });
    </script>

    <style>
	    /* http://meyerweb.com/eric/tools/css/reset/ v2.0 | 20110126 License: none (public domain) */

		html, body, div, span, applet, object, iframe, h1, h2, h3, h4, h5, h6, p, blockquote, pre, a, abbr, acronym, address, big, cite, code, del, dfn, em, img, ins, kbd, q, s, samp, small, strike, strong, sub, sup, tt, var, b, u, i, center, dl, dt, dd, ol, ul, li, fieldset, form, label, legend, table, caption, tbody, tfoot, thead, tr, th, td, article, aside, canvas, details, embed, figure, figcaption, footer, header, hgroup, menu, nav, output, ruby, section, summary, time, mark, audio, video {
		    margin: 0;
		    padding: 0;
		    border: 0;
		    font-size: 100%;
		    font: inherit;
		    vertical-align: baseline;
		}
		/* HTML5 display-role reset for older browsers */
		article, aside, details, figcaption, figure, footer, header, hgroup, menu, nav, section {
		    display: block;
		}
		body { line-height: 1; }
		ol, ul { list-style: none; }
		blockquote, q { quotes: none; }
		blockquote:before, blockquote:after, q:before, q:after {
		    content: '';
		    content: none;
		}
		table {
		    border-collapse: collapse;
		    border-spacing: 0;
		}

		/*  --------------------- Styles ====================== */

		html {
		    text-align: center;
		}

		body {
		    margin: 0 auto;
		}

		nav ul li { 
		    display: inline; 
		    margin-right: 10%;
		}

		.dots {
		    margin: auto 10%;
		    padding-bottom: 40px;
		}

		li {
		    display: block;
		    float: left;
		    height: 7px;
		    width: 16px;
		    font-size: 100%;
		}
		li.passed:before { content: "✔"; }
		li.failed:before { content: "✘"; }
		li.error:before  { content: "✘"; }

		div.story {
		    text-align: left;
		    margin-left: 10%;
		    font-family: Monaco;
		    line-height: 1.4em;
		    color: #586e75;
		}
		div.story .t2:before { content: "\00a0"; }
		div.story .t3:before { content: "\00a0\00a0"; }
		div.story .t4:before { content: "\00a0\00a0\00a0"; }
		div.story .t5:before { content: "\00a0\00a0\00a0\00a0"; }
		div.story .t6:before { content: "\00a0\00a0\00a0\00a0\00a0"; }
		div.story .t7:before { content: "\00a0\00a0\00a0\00a0\00a0\00a0"; }
		div.story .t8:before { content: "\00a0\00a0\00a0\00a0\00a0\00a0\00a0"; }

		h1 {
		    font-family: Monaco;
		    font-size: 15px;
		    padding: 5px;
		    margin: 40px 10%;
		    margin-left: 10%;
		    text-align: left;
		    border-radius: 5px;
		}
    </style>
</head>

<body>
	<div id="output"></div>

	<script id="stories" type="text/template">

	</script>
</body>

</html>`
