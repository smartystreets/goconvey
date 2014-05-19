var convey = {
	overallClass: 'buildfail',		// initial overall status class
	statuses: {
		pass: { class: 'ok', text: "All Systems Go" },
		fail: { class: 'fail', text: "Initializing..." },
		panic: { class: 'panic', text: "Panic" },
		buildfail: { class: 'buildfail', text: "Loading..." }
	}
};

// Stop using the github path, gee whiz
if (window.location.host == "smartystreets.github.io")
	window.location = "http://goconvey.co";


// Wait until elements of page, including images, finish loading
$(window).load(function()
{
	setTimeout(function() { changeStatus(convey.statuses.fail); }, 500);
	setTimeout(function() { changeStatus(convey.statuses.pass); }, 2500);
});

// When document ready
$(function()
{
	// Sticky nav
	$('nav').waypoint('sticky');

	// Load theme as previously toggled
	if (window.localStorage)
	{
		var lastTheme = window.localStorage.getItem("theme");
		if (lastTheme)
			switchTheme(lastTheme);
	}

	// Toggle theme
	$('.toggle-theme').click(function()
	{
		if ($('#theme').attr('href') == "resources/css/dark.css")
			switchTheme("light");
		else
			switchTheme("dark");
	});

	// Carousels
	$('.carousel-switch > .fa').click(function()
	{
		$(this).siblings('.fa-circle').removeClass('fa-circle').addClass('fa-circle-o');
		$(this).addClass('fa-circle');
		var images = $(this).parents('.carousel').find('.carousel-images img');
		var image = images.eq($(this).index());
		images.hide();
		image.css('display', 'block');
	}).tipsy();
});


function switchTheme(theme)
{
	var otherTheme = theme == "light" ? "dark" : "light";

	$('#theme').attr('href', "resources/css/" + theme + ".css");
	$('.screenshot').each(function()
	{
		var src = $(this).attr('src');
		$(this).attr('src', src.replace("-" + otherTheme, "-" + theme));
	});

	if (window.localStorage)
		window.localStorage.setItem("theme", theme);
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
			$(this).effect("pulsate", {times: 1}, 100, function()
			{
				$(this).effect("pulsate", {times: 2}, 300, function()
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
		$('.overall').switchClass(convey.overallClass, newStatus.class, 750);

	convey.overallClass = newStatus.class;
}


// BEGIN GOOGLE ANALYTICS
(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
})(window,document,'script','//www.google-analytics.com/analytics.js','ga');

ga('create', 'UA-86578-21', 'goconvey.co');
ga('send', 'pageview');
// END GOOGLE ANALYTICS