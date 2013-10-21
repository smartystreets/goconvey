var convey = {
	speed: 'fast',
	statuses: {
		pass: 'ok',
		fail: 'fail',
		panic: 'panic',
		failedBuild: 'buildfail',
		skip: 'skip'
	},
	regex: {
		expected: /Expected:?\s+'?/,
		actual: /'?\s+(Actual:)\s+'?/,
		actualEnd: /$|('?\s+(\((Should|but))|$)/
	},
	lastScrollPos: 0,
	assertions: emptyAssertions(),
	failedBuilds: [],
	overall: emptyOverall(),
	zen: {},
	revisionHash: ""
};


$(function()
{
	// Focus on first textbox
	if ($('input').first().val() == "")
		$('input').first().focus();

	$('#banners').waypoint('sticky');

	// Markup.js custom pipes
	Mark.pipes.relativePath = function(str)
	{
		basePath = new RegExp($('#path').val(), 'g');
		return str.replace(basePath, '');
	}
	Mark.pipes.showhtml = function(str)
	{
		return str.replace(/</g, "&lt;").replace(/</g, "&gt;");
	}

	Mark.pipes.nothing = function(str)
	{
		return str == "no test files" || str == "no test functions" || str == "no go code"
	}


	// Show code generator if link is clicked
	$('#show-gen').click(function()
	{
		$('#generator').fadeIn(convey.speed, function() {
			$('#gen-input').focus();
			generate($('#gen-input').val());
		});
	});

	// Smooth scroll within page (props to css-tricks.com)
	$('body').on('click', 'a[href^=#]:not([href=#])', function()
	{
		if (location.pathname.replace(/^\//,'') == this.pathname.replace(/^\//,'') 
			|| location.hostname == this.hostname)
		{
			var target = $(this.hash);
			target = target.length ? target : $('[name=' + this.hash.slice(1) +']');
			   if (target.length) {
				 $('html,body').animate({
					 scrollTop: target.offset().top - 150
				}, 400);
				return false;
			}
		}
	});

	// Hide code generator (or any 'zen'-like window)
	$('.zen-close').click(function() {
		$('.zen').fadeOut(convey.speed);
	});

	// Immediately get test results and on every interval, too
	updatePath(update);
	var poller = setInterval(update, 1500);

	function updatePath(callback)
	{
		$.get('/watch', function(data) {
			$('#path').val($.trim(data));
			if (typeof callback === 'function')
				callback();
		});
	}

	$('#path').change(function() {
		var self = $(this)
		$.post('/watch?root='+encodeURIComponent($.trim($(this).val())))
			.done(function() {
				self.css('color', '');
			})
			.fail(function() {
				self.css('color', '#DD0000');
			});
	});

	function update()
	{
		// Save this so we can revert to the same place we were before the update
		convey.lastScrollPos = $(document).scrollTop();

		$.getJSON("/latest", function(data, status, jqxhr)
		{
			$('#server-down').slideUp(convey.speed);

			if (data.Revision == convey.revisionHash)
				return;

			convey.revisionHash = data.Revision;

			$('#spinner').show();

			updatePath();

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
				// Remove them from the DOM as we'll put new ones back in
				$('.templated').remove();

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

						if (test.Stories.length == 0)
						{
							// Here we've got ourselves a classic Go test,
							// not a GoConvey test that has stories and assertions
							// so we'll treat this whole test as a single assertion
							convey.overall.assertions ++;
							test._id = uniqueID;

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
								test._parsedExpected = parseExpected(test.Message);
								test._parsedActual = parseActual(test.Message);
								test._parsed = test._parsedExpected != "" && test._parsedActual != "";
								convey.assertions.failed.push(test);
							}
							else
							{
								test._status = convey.statuses.pass;
								pkg._passed ++;
								test._passed ++;
								convey.assertions.passed.push(test);
							}

							uniqueID ++;
						}
						else
							test._status = convey.statuses.pass;

						for (var k in test.Stories)
						{
							var story = makeContext(test.Stories[k]);

							story._id = uniqueID;
							convey.overall.assertions += story.Assertions.length;

							for (var l in story.Assertions)
							{
								var assertion = story.Assertions[l];
								assertion._id = uniqueID;

								if (assertion.Failure)
								{
									assertion._parsedExpected = parseExpected(assertion.Failure);
									assertion._parsedActual = parseActual(assertion.Failure);
									assertion._parsed = assertion._parsedExpected != "" && assertion._parsedActual != "";
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
					$('#left-sidebar').append(render('tpl-builds-shortcuts', convey.failedBuilds));
					$('#contents').append(render('tpl-builds', convey.failedBuilds));
				}
				if (convey.overall.panics > 0)
				{
					console.log(convey.overall);
					$('#left-sidebar').append(render('tpl-panic-shortcuts', convey.assertions.panicked));
					$('#contents').append(render('tpl-panics', convey.assertions.panicked));
				}
				if (convey.overall.failures > 0)
				{
					$('#left-sidebar').append(render('tpl-failure-shortcuts', convey.assertions.failed));
					$('#contents').append(render('tpl-failures', convey.assertions.failed));
				}

				// Show stories
				$('#contents').append(render('tpl-stories', data));


				// Compute diffs
				$(".failure").each(function() {
					$(this).prettyTextDiff();
				});


				// Finally, show all the results at once, which appear below the banner,
				// and hide the loading spinner, and update the title

				$('#loading').hide();
				
				var cleanedStatus = $.trim($('.overall .summary').text())
									.replace(/\n+\s*|\s-\s/g, ', ')
									.replace(/\s+|\t|-/ig, ' ');
				$('title').text("GoConvey: " + cleanedStatus);
				
				$(this).fadeIn(function()
				{
					// Loading is finished
					$('#spinner').hide();

					// Scroll to same position as before (doesn't account for different-sized content)
					$(document).scrollTop(convey.lastScrollPos);

					// Remove the height attribute which smoothed out the transition
					$('html,body').css('height', '');
				});
			});
		}).fail(function(jqxhr, message, error)
		{
			// If the server is still just starting up... faux that.
			if (jqxhr.responseText == "" && message == "parsererror")
				message = "starting";

			$('#server-down').remove();

			$('#banners').prepend(render('tpl-server-down', {
				jqxhr: jqxhr,
				message: message,
				error: error
			}));
		});
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

	function assignStatus(obj)
	{
		if (obj._skipped)
			obj._status = 'skip';
		else if (obj._panicked)
			obj._status = convey.statuses.panic;
		else if (obj._failed || obj.Outcome == "failed")
			obj._status = convey.statuses.fail;
		else
			obj._status = convey.statuses.pass;
	}

	function render(templateID, context)
	{
		var tpl = $.trim($('#' + templateID).text());
		return $($.trim(Mark.up(tpl, context)));
	}

	function parseExpected(str)
	{
		return stringBetween(str, convey.regex.expected, convey.regex.actual);
	}

	function parseActual(str)
	{
		return stringBetween(str, convey.regex.actual, convey.regex.actualEnd);
	}

	function stringBetween(str, startExpr, endExpr)
	{
		var startMatch = str.match(startExpr);
		
		if (!startMatch)
			return "";

		var start = startMatch.index + startMatch[0].length;

		var endMatch = str.substr(start).match(endExpr);

		if (!endMatch)
			return "";
		
		var end = start + endMatch.index;

		return str.substring(start, end);
	}
});


function emptyOverall()
{
	return {
		status: 'ok',
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