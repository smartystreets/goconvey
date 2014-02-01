var convey = {

	// Configure the GoConvey web UI client here
	config: {

		// Install new themes by adding them here; the first one will be default
		themes: {
			"dark": {name: "Dark Theme", filename: "dark.css"},
			"light": {name: "Light Theme", filename: "light.css"}
		},

		// Path to the themes (end with forward-slash)
		themePath: "/resources/css/themes/",

		// Whether to enable debug output on the console
		debug: true
	},



	//	*** Don't edit below here unless you're brave ***


	statuses: {				// contains some constants related to overall test status
		pass: { class: 'ok', text: "PASS" },
		fail: { class: 'fail', text: "FAIL" },
		panic: { class: 'panic', text: "PANIC" },
		buildfail: { class: 'buildfail', text: "BUILD FAILED" }
	},
	intervals: {},			// intervals that execute periodically
	overall: {},			// current overall results
	theme: "",				// current theme being used
	selClass: "sel"			// CSS class when an element is "selected"
};



init();

function init()
{
	convey.overall = emptyOverall();
	loadTheme();
	$(wireup);
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

	if (!thmID || !convey.config.themes[thmID])
	{
		replacement = Object.keys(convey.config.themes)[0] || defaultTheme;
		log("WARNING", "Could not find '" + thmID + "' theme; defaulting to '" + replacement + "'");
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

	//enumSel("theme", convey.theme);
}

function wireup()
{
	var themes = [];
	for (var k in convey.config.themes)
		themes.push({ id: k, name: convey.config.themes[k].name });
	$('#theme').html(render('tpl-theme-enum', themes));
	
	enumSel("theme", convey.theme);

	$('.enum#theme').on('click', 'li', function()
	{
		if (!$(this).hasClass(convey.selClass))
			loadTheme($(this).data('theme'));
	});

	$('#run-tests').click(function()
	{
		var self = $(this);

		if (self.hasClass('spin-slowly'))
			return;	// Tests already running (TODO: better detection; maybe a state variable)
		
		// Add the CSS class with the animation
		self.addClass('spin-slowly');
		
		// TODO: This will spin while tests are executing, until it finishes
		setTimeout(function() { self.removeClass('spin-slowly'); }, 3350);
	});

	$('#play-pause').click(function()
	{
		$(this).toggleClass("throb " + convey.selClass);
	});

	$('#toggle-notif').click(function()
	{
		$(this).toggleClass("fa-bell-o fa-bell " + convey.selClass);
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

	$('.toggler').click(function()
	{
		var toggle = $('#' + $(this).data('toggle'));
		$('.fa-angle-down, .fa-angle-up', this).toggleClass('fa-angle-down fa-angle-up');
		toggle.toggleClass('hide-narrow show-narrow');
		//$('.' + togId).toggle();
	});

	// Enumerations are lists where one item can be selected at a time
	$('.enum').on('click', 'li', enumSel);

/*
	convey.intervals.time = setInterval(function()
	{
		var t = new Date();
		var h = zerofill(t.getUTCHours(), 2);
		var m = zerofill(t.getUTCMinutes(), 2);
		var s = zerofill(t.getUTCSeconds(), 2);
		var ms = zerofill(t.getUTCMilliseconds(), 3);
		$('#time').text(h + ":" + m + ":" + s + "." + ms);
	}, 71);
*/

	// Temporary, for effect:
	setTimeout(function() { changeStatus(convey.statuses.panic) }, 2000);

	setTimeout(function() { changeStatus(convey.statuses.buildfail) }, 16000);

	setTimeout(function() { changeStatus(convey.statuses.fail) }, 25000);

	setTimeout(function() { changeStatus(convey.statuses.pass) }, 35000);
}

function enumSel(id, val)
{
	if (typeof id === "string" && typeof val === "string")
	{
		$('.enum#'+id+' > li').each(function()
		{
			if ($(this).data(id) == val)
			{
				$(this).addClass(convey.selClass).siblings().removeClass(convey.selClass);
				return false;
			}
		});
	}
	else
		$(this).addClass(convey.selClass).siblings().removeClass(convey.selClass);
}

function toggle(jqelem, switchelem)
{
	var speed = 250;
	var transition = 'easeInOutQuart';
	var containerSel = '.container';

	if (!jqelem.is(':visible'))
	{
		$(containerSel, jqelem).css('opacity', 0);
		jqelem.slideDown(speed, transition, function()
		{
			if (switchelem)
				switchelem.toggleClass(convey.selClass);
			$(containerSel, jqelem).stop().animate({
				opacity: 1
			}, speed);
		});
	}
	else
	{
		$(containerSel, jqelem).stop().animate({
			opacity: 0
		}, speed, function()
		{
			if (switchelem)
				switchelem.toggleClass(convey.selClass);
			jqelem.slideUp(speed, transition);
		});
	}
}

function changeStatus(newStatus)
{
	if (!newStatus || !newStatus.class || !newStatus.text)
		newStatus = convey.statuses.pass;

	// The .flash CSS class and the pulsate effect don't play well together.
	// This series of callbacks does the flickering/pulsating as well as
	// enabling/disabling flashing in the proper order.

	$('.overall .status').removeClass('flash').effect("pulsate", {times: 4}, 650, function()
	{
		 $(this).text(newStatus.text);

		if (newStatus != convey.statuses.pass)
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
	});

	$('.overall').switchClass(convey.overall.status.class, newStatus.class, 1500);

	convey.overall.status = newStatus;
}

function render(templateID, context)
{
	var tpl = $('#' + templateID).text();
	return $($.trim(Mark.up(tpl, context)));
}

function zerofill(val, count)
{
	// Cheers to http://stackoverflow.com/a/9744576/1048862
	var pad = new Array(1 + count).join('0');
	return (pad + val).slice(-pad.length);
}

function log(type, msg)
{
	if (convey.config.debug)
		console.log(type+":", msg);
}