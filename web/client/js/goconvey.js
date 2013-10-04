var convey = {
	speed: 'fast',
	statuses: {
		pass: 'ok',
		fail: 'fail',
		panic: 'panic',
		skip: 'skip'
	},
	assertions: emptyAssertions(),
	overall: emptyOverall(),
	zen: {},
	revisionHash: ""
};

$(function()
{
	// Focus on first textbox
	if ($('input').first().val() == "")
		$('input').first().focus();

	// Show code generator
	$('#show-gen').click(function()
	{
		$('#generator').fadeIn(convey.speed, function() {
			$('#gen-input').focus();
			generate($('#gen-input').val());
		});
	});

	// Hide code generator (or any 'zen'-like window)
	$('.zen-close').click(function() {
		$('.zen').fadeOut(convey.speed);
	});

	// Immediately get test results and on every interval, too
	update();
	var poller = setInterval(update, 1500);

	function updatePath()
	{
		$.get('/watch', function(data) {
			$('#path').val($.trim(data));
		});
	}
	
	updatePath();

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
		$.getJSON("/latest", function(data, status, jqxhr)
		{
			if (data.Revision == convey.revisionHash)
				return;

			convey.revisionHash = data.Revision;

			$('#spinner').show();

			updatePath();

			// Empty out the data from the last update
			convey.overall = emptyOverall();
			convey.assertions = emptyAssertions();

			// Remove existing/old test results
			$('.overall').slideUp(convey.speed);

			$('#results').fadeOut(convey.speed, function()
			{
				// Remove them from the DOM as we'll put new ones back in
				$('.templated').remove();

				// Look for failures and panics through the packages->tests->stories...
				for (var i in data.Packages)
				{
					pkg = makeContext(data.Packages[i]);
					convey.overall.duration += pkg.Elapsed;

					for (var j in pkg.TestResults)
					{
						test = makeContext(pkg.TestResults[j]);

						if (test.Stories.length == 0)
						{
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

						for (var k in test.Stories)
						{
							var story = makeContext(test.Stories[k]);

							convey.overall.assertions += story.Assertions.length;

							for (var l in story.Assertions)
							{
								var assertion = story.Assertions[l];
								if (assertion.Failure)
								{
									assertion._parsedExpected = parseExpected(assertion.Failure);
									assertion._parsedActual = parseActual(assertion.Failure);

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
						}
					}
				}
				console.log("PANICKED", convey.assertions.panicked);
				convey.overall.passed = convey.assertions.passed.length;
				convey.overall.panics = convey.assertions.panicked.length;
				convey.overall.failures = convey.assertions.failed.length;
				convey.overall.skipped = convey.assertions.skipped.length;

				// Panics trump failures overall
				if (convey.overall.panics)
					convey.overall.status = convey.statuses.panic;
				else if (convey.overall.failures)
					convey.overall.status = convey.statuses.fail;

				// Show the overall status: passed, failed, or panicked
				if (convey.overall.status == convey.statuses.pass)
					$('header').after(render('tpl-overall-ok', convey.overall));
				else if (convey.overall.status == convey.statuses.fail)
					$('header').after(render('tpl-overall-fail', convey.overall));
				else
					$('header').after(render('tpl-overall-panic', convey.overall));

				// Show overall status
				$('.overall').slideDown();

				// Show shortucts and failures/panics details
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

				// Finally, show all the results at once, which appear below the banner
				// and hide the loading spinner
				$('#loading').hide();
				$(this).fadeIn(function() {
					$('#spinner').hide();
				});
			});
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
		else if (obj._failed || obj.Passed === false)
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
		return stringBetween(str, "Expected: '", "'\nActual: '");
	}

	function parseActual(str)
	{
		return stringBetween(str, "'\nActual: '", "'\n(Should");
	}

	function stringBetween(str, startStr, endStr)
	{
		var start = str.indexOf(startStr);
		
		if (start < 0)
			return "";
		
		start += startStr.length;

		var end = str.indexOf(endStr, start);

		if (end < 0)
			return "";

		return str.substring(start, end);
	}
});


Mark.pipes.relativePath = function(str)
{
	basePath = new RegExp($('#path').val(), 'g');
	return str.replace(basePath, '');
}


function emptyOverall()
{
	return {
		status: 'ok',
		duration: 0,
		assertions: 0,
		passed: 0,
		panics: 0,
		failures: 0,
		skipped: 0
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