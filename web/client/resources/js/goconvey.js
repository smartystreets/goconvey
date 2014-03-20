var convey = {

	// Configure the GoConvey web UI client here
	config: {
		// Install new themes by adding them here; the first one will be default
		themes: {
			"dark": { name: "Dark", filename: "dark.css", coverage: "hsla({{hue}}, 75%, 30%, .5)" },
			"light": { name: "Light", filename: "light.css", coverage: "hsla({{hue}}, 62%, 75%, 1)" }
		},

		// Path to the themes (end with forward-slash)
		themePath: "/resources/css/themes/"
	},



	//	*** Don't edit below here unless you're brave ***


	statuses: {				// contains some constants related to overall test status
		pass: { class: 'ok', text: "Pass" },	// class name must also be that in the favicon file name
		fail: { class: 'fail', text: "Fail" },
		panic: { class: 'panic', text: "Panic" },
		buildfail: { class: 'buildfail', text: "Build Failure" }
	},
	frameCounter: 0,		// gives each frame a unique ID
	maxHistory: 20,			// how many tests to keep in the history
	notif: undefined,		// the notification currently being displayed
	notifTimer: undefined,	// the timer that clears the notifications automatically
	poller: new Poller(),	// the server poller
	status: "",				// what the _server_ is currently doing (not overall test results)
	overallClass: "",		// class name of the "overall" status banner
	theme: "",				// theme currently being used
	packageStates: {},		// packages manually collapsed or expanded during this page's lifetime
	layout: {
		selClass: "sel",	// CSS class when an element is "selected"
		header: undefined,	// container element of the header area (overall, controls)
		frame: undefined,	// container element of the main body area (above footer)
		footer: undefined	// container element of the footer (stuck to bottom)
	},
	history: [],			// complete history of states (test results and aggregated data), including the current one
	moments: {},			// elements that display time relative to the current time, keyed by ID, with the moment() as a value
	intervals: {},			// ntervals that execute periodically
	intervalFuncs: {		// Functions executed by each interval in convey.intervals
		time: function()
		{
			var t = new Date();
			var h = zerofill(t.getHours(), 2);
			var m = zerofill(t.getMinutes(), 2);
			var s = zerofill(t.getSeconds(), 2);
			$('#time').text(h + ":" + m + ":" + s);
		},
		momentjs: function()
		{
			for (var id in convey.moments)
				$('#'+id).html(convey.moments[id].fromNow());
		}
	}
};


$(init);

$(window).load(function()
{
	// Things may shift after all the elements (images/fonts) are loaded
	reframe();
});

function init()
{
	log("Welcome to GoConvey!");
	log("Initializing interface");
	convey.overall = emptyOverall();
	loadTheme();
	$('body').show();
	initPoller();
	wireup();
}

function loadTheme(thmID)
{
	var defaultTheme = "dark";
	var linkTagId = "themeRef";

	if (!thmID)
		thmID = get('theme') || defaultTheme;

	log("Initializing theme: " + thmID);

	if (!convey.config.themes[thmID])
	{
		replacement = Object.keys(convey.config.themes)[0] || defaultTheme;
		log("WARNING: Could not find '" + thmID + "' theme; defaulting to '" + replacement + "'");
		thmID = replacement;
	}

	convey.theme = thmID;
	save('theme', convey.theme);

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

	colorizeCoverageBars();	// their color is set dynamically, so we have to use the theme's template
}

function initPoller()
{
	$(convey.poller).on('serverstarting', function(event)
	{
		log("Server is starting...");
		convey.status = "starting";
		showServerDown("Server starting");
		$('#run-tests').addClass('spin-slowly disabled');
	});

	$(convey.poller).on('pollsuccess', function(event, data)
	{
		if (convey.status != "starting")
			hideServerDown();

		// These two if statements determine if the server is now busy
		// (and wasn't before) or is not busy (regardless of whether it was before)
		if ((!convey.status || convey.status == "idle")
				&& data.status && data.status != "idle")
			$('#run-tests').addClass('spin-slowly disabled');
		else if (convey.status != "idle" && data.status == "idle")
		{
			$('#run-tests').removeClass('spin-slowly disabled');
			latest();	// TODO: Move this into server-is-idle handler?
		}

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

	$(convey.poller).on('pollfail', function(event, data)
	{
		log("Poll failed; server down");
		convey.status = "down";
		showServerDown("GoConvey server down");
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
		// TODO: If execution just finished, get the latest...?
	});

	convey.poller.start();
}

function wireup()
{
	log("Wireup");

	customMarkupPipes();

	var themes = [];
	for (var k in convey.config.themes)
		themes.push({ id: k, name: convey.config.themes[k].name });
	$('#theme').html(render('tpl-theme-enum', themes));

	enumSel("theme", convey.theme);

	loadSettingsFromStorage();

	$('#stories').on('click', '.toggle-all-pkg', function(event)
	{
		if ($(this).closest('.story-pkg').data('pkg-state') == "expanded")
			collapseAll();
		else
			expandAll();
		return suppress(event);
	});

	$('.enum#pkg-expand-collapse').on('click', 'li:not(.sel)', function()
	{
		var newSetting = $(this).data('pkg-expand-collapse');
		convey.packageStates = {};
		save('pkg-expand-collapse', newSetting);
		if (newSetting == "expanded")
			expandAll();
		else
			collapseAll();
	});

	$('.enum#theme').on('click', 'li:not(.sel)', function()
	{
		loadTheme($(this).data('theme'));
	});

	convey.layout.header = $('header').first();
	convey.layout.frame = $('.frame').first();
	convey.layout.footer = $('footer').last();

	updateWatchPath();

	$('#path').change(function()
	{
		// Updates the watched directory with the server and makes sure it exists
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

		save('notifications', !notif());

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

	$('#show-gen').click(function() {
		var writer = window.open("/composer.html");
		if (window.focus)
			writer.focus();
	});

	$('.controls li, .pkg-cover-name').tipsy({ live: true });

	$('.toggler').not('.narrow').prepend('<i class="fa fa-angle-up fa-lg"></i>');
	$('.toggler.narrow').prepend('<i class="fa fa-angle-down fa-lg"></i>');

	$('.toggler').not('.narrow').click(function()
	{
		var target = $('#' + $(this).data('toggle'));
		$('.fa-angle-down, .fa-angle-up', this).toggleClass('fa-angle-down fa-angle-up');
		target.toggle();
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

	convey.intervals.time = setInterval(convey.intervalFuncs.time, 1000);
	convey.intervals.momentjs = setInterval(convey.intervalFuncs.momentjs, 5000);
	convey.intervalFuncs.time();

	$('#stories').on('click', '.fa.ignore', function(event)
	{
		var pkg = $(this).data('pkg');
		if ($(this).hasClass('disabled'))
			return;
		else if ($(this).hasClass('unwatch'))
			$.get("/ignore", { paths: pkg });
		else
			$.get("/reinstate", { paths: pkg });
		$(this).toggleClass('watch unwatch fa-eye fa-eye-slash clr-red');
		return suppress(event);
	});

	$('#stories').on({
		mouseenter: function() { $('.toggle-all-pkg', this).stop().show('fast'); },
		mouseleave: function() { $('.toggle-all-pkg', this).stop().hide('fast'); }
	}, '.pkg-toggle-container');

	$('#stories').on('click', '.story-pkg', function(event)
	{
		togglePackage(this, true);
		return suppress(event);
	});

	$('#stories').on('click', '.story-line', function()
	{
		$('.story-line-sel').not(this).removeClass('story-line-sel');
		$(this).toggleClass('story-line-sel');
	});
}

function expandAll()
{
	$('.story-pkg').each(function() { expandPackage($(this).data('pkg')); });
}

function collapseAll()
{
	$('.story-pkg').each(function() { collapsePackage($(this).data('pkg')); });
}

function expandPackage(pkgId)
{
	var pkg = $('.story-pkg.pkg-'+pkgId);
	var rows = $('.story-line.pkg-'+pkgId);

	pkg.data('pkg-state', "expanded").addClass('expanded').removeClass('collapsed');

	$('.pkg-toggle', pkg)
		.addClass('fa-minus-square-o')
		.removeClass('fa-plus-square-o');

	rows.show();
}

function collapsePackage(pkgId)
{
	var pkg = $('.story-pkg.pkg-'+pkgId);
	var rows = $('.story-line.pkg-'+pkgId);

	pkg.data('pkg-state', "collapsed").addClass('collapsed').removeClass('expanded');

	$('.pkg-toggle', pkg)
		.addClass('fa-plus-square-o')
		.removeClass('fa-minus-square-o');

	rows.hide();
}

function togglePackage(storyPkgElem)
{
	var pkgId = $(storyPkgElem).data('pkg');
	if ($(storyPkgElem).data('pkg-state') == "expanded")
	{
		collapsePackage(pkgId);
		convey.packageStates[$(storyPkgElem).data('pkg-name')] = "collapsed";
	}
	else
	{
		expandPackage(pkgId);
		convey.packageStates[$(storyPkgElem).data('pkg-name')] = "expanded";
	}
}

function loadSettingsFromStorage()
{
	var pkgExpCollapse = get("pkg-expand-collapse");
	if (!pkgExpCollapse)
	{
		pkgExpCollapse = "expanded";
		save("pkg-expand-collapse", pkgExpCollapse);
	}
	enumSel("pkg-expand-collapse", pkgExpCollapse);

	if (notif())
		$('#toggle-notif').toggleClass("fa-bell-o fa-bell " + convey.layout.selClass);
}











function latest()
{
	log("Fetching latest test results");
	$.getJSON("/latest", process);
}

function process(data, status, jqxhr)
{
	console.log("Latest", data, status, jqxhr);	// TODO: temp

	if (!data || !data.Revision)
		return;

	if (data.Revision == current().results.Revision)
	{
		changeStatus(current().overall.status);	// re-assure that status is unchanged
		return;
	}



	convey.history.push(newFrame());





	current().results = data;

	updateWatchPath();

	// Remove all templated items from the DOM as we'll
	// replace them with new ones; also remove tipsy tooltips
	// that may have lingered around
	$('.templated, .tipsy').remove();

	var uniqueID = 0;

	var packages = {
		tested: [],
		coverage: {},
		nogofiles: [],
		notestfiles: [],
		notestfn: []
	};

	// Look for failures and panics through the packages->tests->stories...
	for (var i in data.Packages)
	{
		pkg = makeContext(data.Packages[i]);
		current().overall.duration += pkg.Elapsed;
		pkg._id = uniqueID++;

		if (pkg.Outcome == "build failure")
		{
			current().overall.failedBuilds++;
			current().failedBuilds.push(pkg);
			continue;
		}


		if (pkg.Outcome == "no go code")
			packages.nogofiles.push(pkg);
		else if (pkg.Outcome == "no test files")
			packages.notestfiles.push(pkg);
		else if (pkg.Outcome == "no test functions")
			packages.notestfn.push(pkg);
		else
		{
			packages.coverage[pkg.PackageName] = pkg.Coverage;
			packages.tested.push(pkg);
		}


		for (var j in pkg.TestResults)
		{
			test = makeContext(pkg.TestResults[j]);
			test._id = uniqueID++;
			test._pkgid = pkg._id;
			test._pkg = pkg.PackageName;

			if (test.Stories.length == 0)
			{
				// Here we've got ourselves a classic Go test,
				// not a GoConvey test that has stories and assertions
				// so we'll treat this whole test as a single assertion
				current().overall.assertions++;

				if (test.Error)
				{
					test._status = convey.statuses.panic;
					pkg._panicked++;
					test._panicked++;
					current().assertions.panicked.push(test);
				}
				else if (test.Passed === false)
				{
					test._status = convey.statuses.fail;
					pkg._failed++;
					test._failed++;
					current().assertions.failed.push(test);
				}
				else
				{
					test._status = convey.statuses.pass;
					pkg._passed++;
					test._passed++;
					current().assertions.passed.push(test);
				}
			}
			else
				test._status = convey.statuses.pass;

			var storyPath = [{ Depth: -1, Title: test.TestName }];	// Maintains the current assertion's story as we iterate

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
				story._pkgid = pkg._id;
				current().overall.assertions += story.Assertions.length;

				for (var l in story.Assertions)
				{
					var assertion = story.Assertions[l];
					assertion._id = uniqueID;
					assertion._pkg = pkg.PackageName;
					assertion._maxDepth = storyPath[storyPath.length - 1].Depth;
					$.extend(assertion._path = [], storyPath);

					if (assertion.Failure)
					{
						current().assertions.failed.push(assertion);
						pkg._failed++;
						test._failed++;
						story._failed++;
					}
					if (assertion.Error)
					{
						current().assertions.panicked.push(assertion);
						pkg._panicked++;
						test._panicked++;
						story._panicked++;
					}
					if (assertion.Skipped)
					{
						current().assertions.skipped.push(assertion);
						pkg._skipped++;
						test._skipped++;
						story._skipped++;
					}
					if (!assertion.Failure && !assertion.Error && !assertion.Skipped)
					{
						current().assertions.passed.push(assertion);
						pkg._passed++;
						test._passed++;
						story._passed++;
					}
				}

				assignStatus(story);
				uniqueID++;
			}
		}
	}

	current().overall.passed = current().assertions.passed.length;
	current().overall.panics = current().assertions.panicked.length;
	current().overall.failures = current().assertions.failed.length;
	current().overall.skipped = current().assertions.skipped.length;

	current().overall.duration = Math.round(current().overall.duration * 1000) / 1000;

	// Build failures trump panics,
	// Panics trump failures,
	// Failures trump pass.
	if (current().overall.failedBuilds)
		changeStatus(convey.statuses.buildfail);
	else if (current().overall.panics)
		changeStatus(convey.statuses.panic);
	else if (current().overall.failures)
		changeStatus(convey.statuses.fail);
	else
		changeStatus(convey.statuses.pass);

	// Save our organized package lists
	current().packages = packages;


	// Render... Render ALL THE THINGS! (all model/state modifications are DONE!)

	$('#coverage').html(render('tpl-coverage', packages.tested.sort(sortPackages)));
	$('#nogofiles').html(render('tpl-nogofiles', packages.nogofiles.sort(sortPackages)));
	$('#notestfiles').html(render('tpl-notestfiles', packages.notestfiles.sort(sortPackages)));
	$('#notestfn').html(render('tpl-notestfn', packages.notestfn.sort(sortPackages)));

	if (current().overall.failedBuilds)
	{
		$('.buildfailures').show();
		$('#buildfailures').html(render('tpl-buildfailures', current().failedBuilds));
	}
	else
		$('.buildfailures').hide();

	if (current().overall.panics)
	{
		$('.panics').show();
		$('#panics').html(render('tpl-panics', current().assertions.panicked));
	}
	else
		$('.panics').hide();


	if (current().overall.failures)
	{
		$('.failures').show();
		$('#failures').html(render('tpl-failures', current().assertions.failed));
		$(".failure").each(function() {
			$(this).prettyTextDiff();
		});
	}
	else
		$('.failures').hide();

	$('#stories').html(render('tpl-stories', packages.tested.sort(sortPackages)));

	var pkgDefaultView = get('pkg-expand-collapse');
	$('.story-pkg.expanded').each(function()
	{
		if (pkgDefaultView == "collapsed" && convey.packageStates[$(this).data('pkg-name')] != "expanded")
			collapsePackage($(this).data('pkg'));
	});

	redrawCoverageBars();

	$('#assert-count').html("<b>"+current().overall.assertions+"</b> assertion"
							+ (current().overall.assertions != 1 ? "s" : ""));
	$('#skip-count').html("<b>"+current().assertions.skipped.length + "</b> skipped");
	$('#fail-count').html("<b>"+current().assertions.failed.length + "</b> failed");
	$('#panic-count').html("<b>"+current().assertions.panicked.length + "</b> panicked");
	$('#duration').html("<b>"+current().overall.duration + "</b>s");

	$('#narrow-assert-count').html("<b>"+current().overall.assertions+"</b>");
	$('#narrow-skip-count').html("<b>"+current().assertions.skipped.length + "</b>");
	$('#narrow-fail-count').html("<b>"+current().assertions.failed.length + "</b>");
	$('#narrow-panic-count').html("<b>"+current().assertions.panicked.length + "</b>");

	convey.moments['last-test'] = moment();	// save timestamp when this test was executed

	// Render the frame in the history pane
	var framePiece = render('tpl-history', current());
	$('.history .container').prepend(framePiece);
	convey.moments['frame-'+current().id] = moment();
	if (convey.history.length > convey.maxHistory)
	{
		// Delete the oldest frame out of the history pane
		convey.history.splice(0, 1);
		$('.history .container .item').last().remove();
	}

	convey.intervalFuncs.momentjs();	// Shows "Last test ..." in footer immediately


	// TODO: An homage to Star Wars?
	//if (convey.overall.status == convey.statuses.pass && window.location.hash == "#anakin")
	//	$('body').append(render('tpl-ok-audio'));

		
	// Show notification, if enabled
	if (notif())
	{
		if (convey.notif)
		{
			clearTimeout(convey.notifTimer);
			convey.notif.close();
		}

		var notifText = notifSummary(current().overall)

		convey.notif = new Notification(notifText.title, {
			body: notifText.body,
			icon: $('.favicon').attr('href')
		});

		convey.notifTimer = setTimeout(function() { convey.notif.close(); }, 5000);
	}

	// Update title in title bar
	if (current().overall.passed == current().overall.assertions && current().overall.status.class == "ok")
		$('title').text("GoConvey (ALL PASS)");
	else
		$('title').text("GoConvey (" + current().overall.passed + "/" + current().overall.assertions + " pass)");

	
	// All done!
	$(convey).trigger('loaded');
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

		downPoller = setInterval(function()
		{
			// If the server is still down, do a ping to see
			// if it's up; pollSuccess() will do the rest.
			if (!up)
				$.get(endpoints.down).done(pollSuccess);
		}, intervalMs);

		$(self).trigger('pollfail', {
			exception: exception,
			message: message,
			jqxhr: jqxhr
		});
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

	var sameStatus = newStatus.class == convey.overallClass;

	// The CSS class .flash and the jQuery UI 'pulsate' effect don't play well together.
	// This series of callbacks does the flickering/pulsating as well as
	// enabling/disabling flashing in the proper order so that they don't overlap.
	// TODO: I suppose the pulsating could also be done with just CSS, maybe...?

	var times = sameStatus ? 3 : 2;
	var duration = sameStatus ? 500 : 300;

	$('.overall .status').removeClass('flash').effect("pulsate", {times: times}, duration, function()
	{
		$(this).text(newStatus.text);

		if (newStatus != convey.statuses.pass)	// only flicker extra when not currently passing
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

	if (!sameStatus)	// change the color
		$('.overall').switchClass(convey.overallClass, newStatus.class, 1000);

	current().overall.status = newStatus;
	convey.overallClass = newStatus.class;
	$('.favicon').attr('href', '/resources/ico/goconvey-'+newStatus.class+'.ico');
}

function updateWatchPath()
{
	$.get("/watch", function(data)
	{
		$('#path').val($.trim(data));
	});
}

function notifSummary(overall)
{
	var body = "(" + overall.duration + "s)\r\n";

	body += overall.passed + " passed, ";

	if (overall.failedBuilds)
		body += overall.failedBuilds + " build" + (overall.failedBuilds != 1 ? "s" : "") + " failed, ";
	if (overall.failures)
		body += overall.failures + " failed, ";
	if (overall.panics)
		body += overall.panics + " panicked, ";

	body += overall.skipped + " skipped";

	return {
		title: overall.status.text.toUpperCase(),
		body: body
	};
}

function redrawCoverageBars()
{
	$('.pkg-cover-bar').each(function()
	{
		var pkgName = $(this).data("pkg");
		var hue = $(this).data("width");
		var hueDiff = hue;

		if (convey.history.length > 1)
		{
			var oldHue = convey.history[convey.history.length - 2].packages.coverage[pkgName] || 0;
			$(this).width(oldHue + "%");
			hueDiff = hue - oldHue;
		}

		$(this).animate({
			width: "+=" + hueDiff + "%"
		}, 1250);
	});

	colorizeCoverageBars();
}

function colorizeCoverageBars()
{
	var colorTpl = convey.config.themes[convey.theme].coverage
					|| "hsla({{hue}}, 75%, 30%, .3)";	//default color template

	$('.pkg-cover-bar').each(function()
	{
		var hue = $(this).data("width");
		$(this).css({
			background: colorTpl.replace("{{hue}}", hue)
		});
	});
}

function render(templateID, context)
{
	var tpl = $('#' + templateID).text();
	return $($.trim(Mark.up(tpl, context)));
}

function reframe()
{
	var heightBelowHeader = $(window).height() - convey.layout.header.outerHeight();
	var middleHeight = heightBelowHeader - convey.layout.footer.outerHeight();
	convey.layout.frame.height(middleHeight);

	var pathWidth = $(window).width() - $('#logo').outerWidth() - $('#control-buttons').outerWidth() - 10;
	$('#path-container').width(pathWidth);

}

function notif()
{
	return get('notifications') === "true";	// stored as strings
}

function showServerDown(message)
{
	$('.server-down .message').text(message);
	$('.server-down').show();
	$('.server-not-down').hide();
	reframe();
}

function hideServerDown()
{
	$('.server-down').hide();
	$('.server-not-down').show();
	reframe();
}

function log(msg)
{
	var jqLog = $('#log');
	if (jqLog.length > 0)
	{
		var t = new Date();
		var h = zerofill(t.getHours(), 2);
		var m = zerofill(t.getMinutes(), 2);
		var s = zerofill(t.getSeconds(), 2);
		var ms = zerofill(t.getMilliseconds(), 3);
		date = h + ":" + m + ":" + s + "." + ms;

		$(jqLog).append(render('tpl-log-line', { time: date, msg: msg }));
		$(jqLog).parent('.col').scrollTop(jqLog[0].scrollHeight);
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

function sortPackages(a, b)
{
	// sorts packages ascending by only the last part of their name
	var aPkg = splitPathName(a.PackageName);
	var bPkg = splitPathName(b.PackageName);

	if (aPkg.length == 0 || bPkg.length == 0)
		return 0;

	var aName = aPkg.parts[aPkg.parts.length - 1].toLowerCase();
	var bName = bPkg.parts[bPkg.parts.length - 1].toLowerCase();

	if (aName < bName)
		return -1;
	else if (aName > bName)
		return 1;
	else
		return 0;

	/*
	Use to sort by entire package name:
	if (a.PackageName < b.PackageName) return -1;
	else if (a.PackageName > b.PackageName) return 1;
	else return 0;
	*/
}

function get(key)
{
	var val = localStorage.getItem(key);
	if (val && (val[0] == '[' || val[0] == '{'))
		return JSON.parse(val);
	else
		return val;
}

function save(key, val)
{
	if (typeof val === 'object' || typeof val === 'array')
		val = JSON.stringify(val);
	else if (typeof val === 'number' || typeof val === "boolean")
		val = "" + val;
	localStorage.setItem(key, val);
}

function splitPathName(str)
{
	var delim = str.indexOf('\\') > -1 ? '\\' : '/';
	return { delim: delim, parts: str.split(delim) };
}

function newFrame()
{
	return {
		results: {},					// response from server (with some of our own context info)
		packages: {},					// packages organized into statuses for convenience (like with coverage)
		overall: emptyOverall(),		// overall status info, compiled from server's response
		assertions: emptyAssertions(),	// lists of assertions, compiled from server's response
		failedBuilds: [],				// list of packages that failed to build
		timestamp: moment(),			// the timestamp of this "freeze-state"
		id: convey.frameCounter++		// unique ID for this frame
	};
}

function emptyOverall()
{
	return {
		status: {},
		duration: 0,
		assertions: 0,
		passed: 0,
		panics: 0,
		failures: 0,
		skipped: 0,
		failedBuilds: 0
	};
}

function emptyAssertions()
{
	return {
		passed: [],
		failed: [],
		panicked: [],
		skipped: []
	};
}

function makeContext(obj)
{
	obj._passed = 0;
	obj._failed = 0;
	obj._panicked = 0;
	obj._skipped = 0;
	obj._status = '';
	return obj;
}

function current()
{
	return convey.history[convey.history.length - 1] || newFrame();
}

function assignStatus(obj)
{
	if (obj._skipped)
		obj._status = 'skip';
	else if (obj.Outcome == "ignored")
		obj._status = convey.statuses.ignored;
	else if (obj._panicked)
		obj._status = convey.statuses.panic;
	else if (obj._failed || obj.Outcome == "failed")
		obj._status = convey.statuses.fail;
	else
		obj._status = convey.statuses.pass;
}

function customMarkupPipes()
{
	// MARKUP.JS
	// Custom pipes
	// TODO: Are all of these being used?
	Mark.pipes.relativePath = function(str)
	{
		basePath = new RegExp($('#path').val()+'[\\/]', 'gi');
		return str.replace(basePath, '');
	};
	Mark.pipes.showhtml = function(str)
	{
		return str.replace(/</g, "&lt;").replace(/>/g, "&gt;");
	};
	Mark.pipes.nothing = function(str)
	{
		return str == "no test files" || str == "no test functions" || str == "no go code"
	};
	Mark.pipes.boldPkgName = function(str)
	{
		var pkg = splitPathName(str);
		pkg.parts[0] = '<span class="not-pkg-name">' + pkg.parts[0];
		pkg.parts[pkg.parts.length - 1] = "</span><b>" + pkg.parts[pkg.parts.length - 1] + "</b>";
		return pkg.parts.join(pkg.delim);
	};
	Mark.pipes.chopEnd = function(str, n)
	{
		return str.length > n ? "..." + str.substr(str.length - n) : str;
	};
	Mark.pipes.needsDiff = function(test)
	{
		return !!test.Failure && (test.Expected != "" || test.Actual != "");
	};
	Mark.pipes.coveragePct = function(str)
	{
		// Expected input: 75% to be represented as: "75.0"
		var num = parseInt(str);	// we only need int precision
		if (num < 0)
			return "0";
		else if (num <= 5)
			return "5px";	// Still shows low coverage
		else if (num > 100)
			str = "100";
		return str;
	};
	Mark.pipes.coverageDisplay = function(str)
	{
		var num = parseFloat(str);
		return num < 0 ? "" : num + "% coverage";
	};
	Mark.pipes.coverageReportName = function(str)
	{
		return str.replace(/\//g, "-");
	};
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