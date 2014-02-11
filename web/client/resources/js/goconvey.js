var convey = {

	// Configure the GoConvey web UI client here
	config: {
		// Install new themes by adding them here; the first one will be default
		themes: {
			"dark": {name: "Dark", filename: "dark.css"},
			"light": {name: "Light", filename: "light.css"}
		},

		// Path to the themes (end with forward-slash)
		themePath: "/resources/css/themes/"
	},



	//	*** Don't edit below here unless you're brave ***


	statuses: {				// contains some constants related to overall test status
		pass: { class: 'ok', text: "Pass" },
		fail: { class: 'fail', text: "Fail" },
		panic: { class: 'panic', text: "Panic" },
		buildfail: { class: 'buildfail', text: "Build Failure" }
	},
	intervals: {},			// intervals that execute periodically
	poller: new Poller(),	// the server poller
	status: "",				// what the server is currently doing
	overall: {},			// current overall results
	theme: "",				// current theme being used
	layout: {
		selClass: "sel",	// CSS class when an element is "selected"
		header: undefined,	// Container element of the header area (overall, controls)
		frame: undefined,	// Container element of the main body area (above footer)
		footer: undefined	// Container element of the footer (stuck to bottom)
	}
};


$(init);

function init()
{
	log("Welcome to GoConvey. Initializing UI...");
	convey.overall = emptyOverall();
	loadTheme();
	initPoller();
	wireup();
}

function emptyOverall()
{
	return {
		status: convey.statuses.pass,
		duration: 0,
		assertions: 0,
		passed: 0,
		panics: 0,
		failures: 0,
		skipped: 0,
		failedBuilds: 0
	};
}

function loadTheme(thmID)
{
	var defaultTheme = "dark";
	var linkTagId = "themeRef";

	if (!thmID)
		thmID = localStorage.getItem('theme');

	log("Initializing theme: " + thmID);

	if (!thmID || !convey.config.themes[thmID])
	{
		replacement = Object.keys(convey.config.themes)[0] || defaultTheme;
		log("WARNING: Could not find '" + thmID + "' theme; defaulting to '" + replacement + "'");
		thmID = replacement;
	}

	convey.theme = thmID;
	localStorage.setItem('theme', convey.theme);

	var linkTag = $('#'+linkTagId);
	var fullPath = convey.config.themePath
					+ convey.config.themes[convey.theme].filename;

	if (linkTag.length == 0)
	{
		$('head').append('<link rel="stylesheet" href="'
			+ fullPath + '" id="themeRef">');
	}
	else
		linkTag.attr('href', fullPath);
}

function initPoller()
{
	$(convey.poller).on('serverstarting', function(event)
	{
		log("Server is starting...");
		convey.status = "starting";
		$('#run-tests').addClass('spin-slowly disabled');
	});

	$(convey.poller).on('pollsuccess', function(event, data)
	{
		// These two if statements determine if the server is now busy
		// (and wasn't before) or is not busy (regardless of whether it was before)
		if ((!convey.status || convey.status == "idle")
				&& data.status && data.status != "idle")
			$('#run-tests').addClass('spin-slowly disabled');
		else if (convey.status != "idle" && data.status == "idle")
			$('#run-tests').removeClass('spin-slowly disabled');

		switch (data.status)
		{
			case "executing":
				$(convey.poller).trigger('serverexec', data);
				break;
			case "parsing":
				$(convey.poller).trigger('serverparsing', data);
				break;
			case "idle":
				$(convey.poller).trigger('serveridle', data);
				break;
		}

		convey.status = data.status;
	});

	$(convey.poller).on('serverexec', function(event, data)
	{
		log("Server status: executing");
	});

	$(convey.poller).on('serverparsing', function(event, data)
	{
		log("Server status: Parsing");
	});

	$(convey.poller).on('serveridle', function(event, data)
	{
		log("Server status: idle");
		// TODO: If execution just finished, get the latest...
	});

	convey.poller.start();
}

function wireup()
{
	log("Wireup");

	var themes = [];
	for (var k in convey.config.themes)
		themes.push({ id: k, name: convey.config.themes[k].name });
	$('#theme').html(render('tpl-theme-enum', themes));

	enumSel("theme", convey.theme);


	$('.enum#theme').on('click', 'li', function()
	{
		if (!$(this).hasClass(convey.layout.selClass))
			loadTheme($(this).data('theme'));
	});

	convey.layout.header = $('header').first();
	convey.layout.frame = $('.frame').first();
	convey.layout.footer = $('footer').last();

	updateWatchPath(true);	// true tells the server we're a new client

	// Updates the watched directory with the server and make sure it exists
	$('#path').change(function()
	{
		var tb = $(this);
		var newpath = encodeURIComponent($.trim(tb.val()));
		$.post('/watch?root='+newpath)
			.done(function() { tb.removeClass('error'); })
			.fail(function() { tb.addClass('error'); });
	});

	$('#run-tests').click(function()
	{
		var self = $(this);
		if (self.hasClass('spin-slowly') || self.hasClass('disabled'))
			return;
		$.get("/execute");
	});

	$('#play-pause').click(function()
	{
		$(this).toggleClass("throb " + convey.layout.selClass);
	});

	$('#toggle-notif').click(function()
	{
		$(this).toggleClass("fa-bell-o fa-bell " + convey.layout.selClass);

		localStorage.setItem('notifications', !notif());

		if (notif() && 'Notification' in window)
		{
			if (Notification.permission !== 'denied')
			{
				Notification.requestPermission(function(per)
				{
					if (!('permission' in Notification))
						Notification.permission = per;
				});
			}
		}
	});

	$('#show-history').click(function()
	{
		toggle($('.history'), $(this));
	});

	$('#show-settings').click(function()
	{
		toggle($('.settings'), $(this));
	});

	$('.controls li').tipsy({ live: true });

	$('.toggler').not('.narrow').prepend('<i class="fa fa-angle-up fa-lg"></i>');
	$('.toggler.narrow').prepend('<i class="fa fa-angle-down fa-lg"></i>');

	$('.toggler').not('.narrow').click(function()
	{
		var target = $('#' + $(this).data('toggle'));
		$('.fa-angle-down, .fa-angle-up', this).toggleClass('fa-angle-down fa-angle-up');
		target.toggleClass('hide');
	});

	$('.toggler.narrow').click(function()
	{
		var target = $('#' + $(this).data('toggle'));
		$('.fa-angle-down, .fa-angle-up', this).toggleClass('fa-angle-down fa-angle-up');
		target.toggleClass('hide-narrow show-narrow');
	});

	// Enumerations are horizontal lists where one item can be selected at a time
	$('.enum').on('click', 'li', enumSel);

	$(window).resize(reframe);
	reframe();
	latest();

	convey.intervals.time = setInterval(function()
	{
		var t = new Date();
		var h = zerofill(t.getUTCHours(), 2);
		var m = zerofill(t.getUTCMinutes(), 2);
		var s = zerofill(t.getUTCSeconds(), 2);
		//var ms = zerofill(t.getUTCMilliseconds(), 1);
		var ms = 0;
		$('#time').text(h + ":" + m + ":" + s + "." + ms);
	}, 1000);

	$('.story-line').click(function()
	{
		$('.story-line-sel').not(this).removeClass('story-line-sel');
		$(this).toggleClass('story-line-sel');
	});

	$('#stories').on('click', '.fa.ignore', function(event)
	{
		var pkg = $(this).closest('.pkg-story').data('pkg');
		if ($(this).hasClass('disabled'))
			return;
		else if ($(this).hasClass('unwatch'))
			$.get("/ignore", { path: pkg });
		else
			$.get("/reinstate", { path: pkg });
		$(this).toggleClass('watch')
			.toggleClass('unwatch')
			.toggleClass('fa-eye')
			.toggleClass('fa-eye-slash')
			.toggleClass('clr-red');
		return suppress(event);
	});

	$('#stories').on('click', '.story-pkg', function()
	{
		var pkg = $(this).data('pkg');
		$('tr.story-line.pkg-'+pkg).toggle();
		$('.fa-collapse-o, .fa-expand-o', this).toggleClass('fa-collapse-o fa-expand-o');
		return suppress(event);
	});

	// Temporary, for effect:
	setTimeout(function() { changeStatus(convey.statuses.panic) }, 2000);

	setTimeout(function() { changeStatus(convey.statuses.buildfail) }, 16000);

	setTimeout(function() { changeStatus(convey.statuses.fail) }, 25000);

	setTimeout(function() { changeStatus(convey.statuses.pass) }, 35000);
}













function latest()
{
	log("Fetching latest test results");
	$.getJSON("/latest", process);
}

function process(data, status, jqxhr)
{
	console.log("Latest", data, status, jqxhr);
	/*if (!data || !data.Revision)
		return showServerDown(jqxhr, "starting");
	else
		$('#server-down').slideUp(convey.speed);

	if (data.Revision == convey.revisionHash)
		return;

	convey.revisionHash = data.Revision;
	convey.payload = data;

	updateWatchPath();

	// Empty out the data from the last update
	convey.overall = emptyOverall();
	convey.assertions = emptyAssertions();
	convey.failedBuilds = [];

	// Force page height to help smooth out the transition
	$('html,body').css('height', $(document).outerHeight());

	// Remove existing/old test results
	$('.overall').slideUp(convey.speed);
	$('#results').fadeOut(convey.speed, function()
	{
		// Remove all templated items from the DOM as we'll
		// replace them with new ones; also remove tipsy tooltips
		// that may have lingered around
		$('.templated, .tipsy').remove();

		var uniqueID = 0;

		// Look for failures and panics through the packages->tests->stories...
		for (var i in data.Packages)
		{
			pkg = makeContext(data.Packages[i]);
			convey.overall.duration += pkg.Elapsed;
			pkg._id = uniqueID++;

			if (pkg.Outcome == "build failure")
			{
				convey.overall.failedBuilds ++;
				convey.failedBuilds.push(pkg);
				continue;
			}

			for (var j in pkg.TestResults)
			{
				test = makeContext(pkg.TestResults[j]);
				test._id = uniqueID;
				uniqueID ++;

				if (test.Stories.length == 0)
				{
					// Here we've got ourselves a classic Go test,
					// not a GoConvey test that has stories and assertions
					// so we'll treat this whole test as a single assertion
					convey.overall.assertions ++;

					if (test.Error)
					{
						test._status = convey.statuses.panic;
						pkg._panicked ++;
						test._panicked ++;
						convey.assertions.panicked.push(test);
					}
					else if (test.Passed === false)
					{
						test._status = convey.statuses.fail;
						pkg._failed ++;
						test._failed ++;
						convey.assertions.failed.push(test);
					}
					else
					{
						test._status = convey.statuses.pass;
						pkg._passed ++;
						test._passed ++;
						convey.assertions.passed.push(test);
					}
				}
				else
					test._status = convey.statuses.pass;

				var storyPath = [{ Depth: -1, Title: test.TestName }];	// Will maintain the current assertion's path

				for (var k in test.Stories)
				{
					var story = makeContext(test.Stories[k]);

					// Establish the current story path so we can report the context
					// of failures and panicks more conveniently at the top of the page
					if (storyPath.length > 0)
						for (var x = storyPath[storyPath.length - 1].Depth; x >= test.Stories[k].Depth; x--)
							storyPath.pop();
					
					storyPath.push({ Depth: test.Stories[k].Depth, Title: test.Stories[k].Title });

					story._id = uniqueID;
					convey.overall.assertions += story.Assertions.length;

					for (var l in story.Assertions)
					{
						var assertion = story.Assertions[l];
						assertion._id = uniqueID;
						$.extend(assertion._path = [], storyPath);

						if (assertion.Failure)
						{
							convey.assertions.failed.push(assertion);
							pkg._failed ++;
							test._failed ++;
							story._failed ++;
						}
						if (assertion.Error)
						{
							convey.assertions.panicked.push(assertion);
							pkg._panicked ++;
							test._panicked ++;
							story._panicked ++;
						}
						if (assertion.Skipped)
						{
							convey.assertions.skipped.push(assertion);
							pkg._skipped ++;
							test._skipped ++;
							story._skipped ++;
						}
						if (!assertion.Failure && !assertion.Error && !assertion.Skipped)
						{
							convey.assertions.passed.push(assertion);
							pkg._passed ++;
							test._passed ++;
							story._passed ++;
						}
					}

					assignStatus(story);
					uniqueID ++;
				}
			}
		}

		convey.overall.passed = convey.assertions.passed.length;
		convey.overall.panics = convey.assertions.panicked.length;
		convey.overall.failures = convey.assertions.failed.length;
		convey.overall.skipped = convey.assertions.skipped.length;

		convey.overall.duration = Math.round(convey.overall.duration * 1000) / 1000;

		// Build failures trump panics,
		// Panics trump failures,
		// Failures trump passing.
		if (convey.overall.failedBuilds)
			convey.overall.status = convey.statuses.failedBuild;
		else if (convey.overall.panics)
			convey.overall.status = convey.statuses.panic;
		else if (convey.overall.failures)
			convey.overall.status = convey.statuses.fail;

		// Show the overall status: passed, failed, or panicked
		if (convey.overall.status == convey.statuses.pass)
			$('#banners').append(render('tpl-overall-ok', convey.overall));
		else if (convey.overall.status == convey.statuses.fail)
			$('#banners').append(render('tpl-overall-fail', convey.overall));
		else if (convey.overall.status == convey.statuses.panic)
			$('#banners').append(render('tpl-overall-panic', convey.overall));
		else
			$('#banners').append(render('tpl-overall-buildfail', convey.overall));

		// Show overall status
		$('.overall').slideDown();
		$('.favicon').attr('href', '/ico/goconvey-'+convey.overall.status+'.ico');

		// Show shortucts and builds/failures/panics details
		if (convey.overall.failedBuilds > 0)
		{
			$('#right-sidebar').append(render('tpl-builds-shortcuts', convey.failedBuilds));
			$('#contents').append(render('tpl-builds', convey.failedBuilds));
		}
		if (convey.overall.panics > 0)
		{
			$('#right-sidebar').append(render('tpl-panic-shortcuts', convey.assertions.panicked));
			$('#contents').append(render('tpl-panics', convey.assertions.panicked));
		}
		if (convey.overall.failures > 0)
		{
			$('#right-sidebar').append(render('tpl-failure-shortcuts', convey.assertions.failed));
			$('#contents').append(render('tpl-failures', convey.assertions.failed));
		}

		// Show stories
		$('#contents').append(render('tpl-stories', data));

		// Show shortcut links to packages
		$('#left-sidebar').append(render('tpl-packages', data.Packages.sort(sortPackages)));

		// Compute diffs
		$(".failure").each(function() {
			$(this).prettyTextDiff();
		});


		// Finally, show all the results at once, which appear below the banner,
		// and hide the loading spinner, and update the title

		$('#loading').hide();
		
		var cleanSummary = $.trim($('.overall .summary').text())
							.replace(/\n+\s*|\s-\s/g, ', ')
							.replace(/\s+|\t|-/ig, ' ');
		$('title').text("GoConvey: " + cleanSummary);

		// An homage to Star Wars
		if (convey.overall.status == convey.statuses.pass && window.location.hash == "#anakin")
			$('body').append(render('tpl-ok-audio'));
		
		if (notif())
		{
			if (convey.notif)
				convey.notif.close();

			var cleanStatus = $.trim($('.overall:visible .status').text()).toUpperCase();

			convey.notif = new Notification(cleanStatus, {
				body: cleanSummary,
				icon: $('.favicon').attr('href')
			});

			setTimeout(function() { convey.notif.close(); }, 3500);
		}

		$(this).fadeIn(function()
		{
			// Loading is finished
			doneExecuting();

			// Scroll to same position as before (doesn't account for different-sized content)
			$(document).scrollTop(convey.lastScrollPos);

			if ($('.stuck .overall').is(':visible'))
				bannerClickToTop(true);	// make the banner across the top clickable again

			// Remove the height attribute which smoothed out the transition
			$('html,body').css('height', '');
		});
	});
	*/
}

















function Poller(config)
{
	// CONFIGURABLE
	var endpoints = {
		up: "/status/poll",		// url to poll when the server is up
		down: "/status"			// url to poll at regular intervals when the server is down
	};
	var timeout =  60000 * 2;	// how many ms between polling attempts
	var intervalMs = 1000;		// ms between polls when the server is down

	// INTERNAL STATE
	var up = true;				// whether or not we can connect to the server
	var req;					// the pending ajax request
	var downPoller;				// the setInterval for polling when the server is down
	var self = this;

	if (typeof config === 'object')
	{
		if (typeof config.endpoints === 'object')
		{
			endpoints.up = config.endpoints.up;
			endpoints.down = config.endpoints.down;
		}
		if (config.timeout)
			timeout = config.timeout;
		if (config.interval)
			intervalMs = config.interval;
	}

	$(self).on('pollstart', function(event, data) {
		log("Started poller");
	}).on('pollstop', function(event, data) {
		log("Stopped poller");
	});


	this.start = function()
	{
		if (req)
			return false;
		doPoll();
		$(self).trigger('pollstart', {url: endpoints.up, timeout: timeout});
		return true;
	};

	this.stop = function()
	{
		if (!req)
			return false;
		req.abort();
		req = undefined;
		stopped = true;
		stopDownPoller();
		$(self).trigger('pollstop', {});
		return true;
	};

	this.setTimeout = function(tmout)
	{
		timeout = tmout;	// takes effect at next poll
	};

	this.isUp = function()
	{
		return up;
	};

	function doPoll()
	{
		req = $.ajax({
			url: endpoints.up + "?timeout=" + timeout,
			timeout: timeout
		}).done(pollSuccess).fail(pollFailed);
	}

	function pollSuccess(data, message, jqxhr)
	{
		stopDownPoller();
		doPoll();
		
		var wasUp = up;
		up = true;
		status = data;

		var arg = {
			status: status,
			data: data,
			jqxhr: jqxhr
		};

		if (!wasUp)
			$(convey.poller).trigger('serverstarting', arg);
		else
			$(self).trigger('pollsuccess', arg);
	}

	function pollFailed(jqxhr, message, exception)
	{
		if (message == "timeout")
		{
			log("Poller timeout; re-polling...", req);
			doPoll();	// in our case, timeout actually means no activity; poll again
			return;
		}

		up = false;

		log("Poll failed; server down");

		downPoller = setInterval(function()
		{
			// If the server is still down, do a ping to see
			// if it's up; pollSuccess() will do the rest.
			if (!up)
				$.get(endpoints.down).done(pollSuccess);
		}, intervalMs);
	}

	function stopDownPoller()
	{
		if (!downPoller)
			return;
		clearInterval(downPoller);
		downPoller = undefined;
	}
}























function enumSel(id, val)
{
	if (typeof id === "string" && typeof val === "string")
	{
		$('.enum#'+id+' > li').each(function()
		{
			if ($(this).data(id) == val)
			{
				$(this).addClass(convey.layout.selClass).siblings().removeClass(convey.layout.selClass);
				return false;
			}
		});
	}
	else
		$(this).addClass(convey.layout.selClass).siblings().removeClass(convey.layout.selClass);
}

function toggle(jqelem, switchelem)
{
	var speed = 250;
	var transition = 'easeInOutQuart';
	var containerSel = '.container';

	if (!jqelem.is(':visible'))
	{
		$(containerSel, jqelem).css('opacity', 0);
		jqelem.stop().slideDown(speed, transition, function()
		{
			if (switchelem)
				switchelem.toggleClass(convey.layout.selClass);
			$(containerSel, jqelem).stop().animate({
				opacity: 1
			}, speed);
			reframe();
		});
	}
	else
	{
		$(containerSel, jqelem).stop().animate({
			opacity: 0
		}, speed, function()
		{
			if (switchelem)
				switchelem.toggleClass(convey.layout.selClass);
			jqelem.stop().slideUp(speed, transition, function() { reframe(); });
		});
	}
}

function changeStatus(newStatus)
{
	if (!newStatus || !newStatus.class || !newStatus.text)
		newStatus = convey.statuses.pass;

	// The CSS class .flash and the jQuery UI 'pulsate' effect don't play well together.
	// This series of callbacks does the flickering/pulsating as well as
	// enabling/disabling flashing in the proper order so that they don't overlap.
	// TODO: I suppose the pulsating could also be done with just CSS...

	$('.overall .status').removeClass('flash').effect("pulsate", {times: 2}, 300, function()
	{
		 $(this).text(newStatus.text);

		if (newStatus != convey.statuses.pass)
		{
			$(this).effect("pulsate", {times: 2}, 300, function()
			{
				$(this).effect("pulsate", {times: 3}, 500, function()
				{
					if (newStatus == convey.statuses.panic
							|| newStatus == convey.statuses.buildfail)
						$(this).addClass('flash');
					else
						$(this).removeClass('flash');
				});
			});
		}
	});

	$('.overall').switchClass(convey.overall.status.class, newStatus.class, 750);
	convey.overall.status = newStatus;
}

function updateWatchPath(newClient)
{
	var tb = $('#path')[0];
	var endpoint = "/watch";
	if (newClient)
		endpoint += "?newclient=1";

	$.get(endpoint, function(data)
	{
		$(tb).val($.trim(data));
	});
}

function coverageToColor(percent)
{
	// This converts a number between 0 and 360
	// to an HSL (not RGB) value appropriate for
	// displaying a basic coverage bar behind text.
	// It works for any value between 0 to 360,
	// but the hue at 120 happens to be about green,
	// and 0 is red, between is yellow; just what we want.
	var hue = percent * 1.2;
	return "hsl(" + hue + ", 100%, 75%)";
}

function render(templateID, context)
{
	var tpl = $('#' + templateID).text();
	return $($.trim(Mark.up(tpl, context)));
}

function reframe()
{
	var heightBelowHeader = $(window).height() - convey.layout.header.outerHeight();
	var middleHeight = heightBelowHeader - convey.layout.footer.outerHeight() - 1;	// -1 for borders
	convey.layout.frame.height(middleHeight);
}

function notif()
{
	return localStorage.getItem('notifications') === "true";	// stored as strings
}

function log(msg)
{
	var logElem = $('#log')[0];
	if (logElem)
	{
		var t = new Date();
		var h = zerofill(t.getUTCHours(), 2);
		var m = zerofill(t.getUTCMinutes(), 2);
		var s = zerofill(t.getUTCSeconds(), 2);
		var ms = zerofill(t.getUTCMilliseconds(), 3);
		date = h + ":" + m + ":" + s + "." + ms;

		$(logElem).append(render('tpl-log-line', { time: date, msg: msg }));
		$(logElem).scrollTop(logElem.scrollHeight);
	}
	else
		console.log(msg);
}

function zerofill(val, count)
{
	// Cheers to http://stackoverflow.com/a/9744576/1048862
	var pad = new Array(1 + count).join('0');
	return (pad + val).slice(-pad.length);
}

function suppress(event)
{
	if (!event)
		return false;
	if (event.preventDefault)
		event.preventDefault();
	if (event.stopPropagation)
		event.stopPropagation();
	event.cancelBubble = true;
	return false;
}