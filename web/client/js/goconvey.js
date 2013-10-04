var convey = {
	statuses: {
		pass: 'ok',
		fail: 'fail',
		panic: 'panic',
		skip: 'skip'
	},
	revisionHash: "",
	speed: 'fast'
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

	// Immediately get test results
	update();

	var poller = setInterval(update, 1500);

	function update()
	{
		$.getJSON("oldschool-panic.json", function(data, status, jqxhr)
		{
			if (data.Revision == convey.revisionHash)
				return;

			convey.revisionHash = data.Revision;

			$('#spinner').show();

			// Remove existing/old test results
			$('.overall').slideUp(convey.speed);
			$('#results').fadeOut(convey.speed, function()
			{
				// Remove them from the DOM as we'll put new ones back in
				$('.templated').remove();

				// Prepare to begin iterating the new results
				var overallStatus = convey.statuses.pass;
				var ctx_overall = {
					duration: 0,
					assertions: 0,
					passed: 0,
					panics: 0,
					failures: 0,
					skipped: 0
				};
				var asserts_passed = [];
				var asserts_failed = [];
				var asserts_panicked = [];
				var asserts_skipped = [];

				// Look for failures and panics through the packages->tests->stories...
				for (var i in data.Packages)
				{
					pkg = makeContext(data.Packages[i]);
					ctx_overall.duration += pkg.Elapsed;

					for (var j in pkg.TestResults)
					{
						test = makeContext(pkg.TestResults[j]);

						if (test.Stories.length == 0)
						{
							ctx_overall.assertions ++;

							if (test.Error)
							{
								test._status = convey.statuses.panic;
								pkg._panicked ++;
								test._panicked ++;
								asserts_panicked.push(test);
							}
							else if (test.Passed === false)
							{
								test._status = convey.statuses.fail;
								pkg._failed ++;
								test._failed ++;
								asserts_failed.push(test);
							}
							else
							{
								test._status = convey.statuses.pass;
								pkg._passed ++;
								test._passed ++;
								asserts_passed.push(test);
							}
						}
						else
							test._status = convey.statuses.pass;

						for (var k in test.Stories)
						{
							var story = makeContext(test.Stories[k]);

							ctx_overall.assertions += story.Assertions.length;

							for (var l in story.Assertions)
							{
								var assertion = story.Assertions[l];
								if (assertion.Failure)
								{
									assertion._parsedExpected = parseExpected(assertion.Failure);
									assertion._parsedActual = parseActual(assertion.Failure);

									asserts_failed.push(assertion);
									pkg._failed ++;
									test._failed ++;
									story._failed ++;
								}
								if (assertion.Error)
								{
									asserts_panicked.push(assertion);
									pkg._panicked ++;
									test._panicked ++;
									story._panicked ++;
								}
								if (assertion.Skipped)
								{
									asserts_skipped.push(assertion);
									pkg._skipped ++;
									test._skipped ++;
									story._skipped ++;
								}
								if (!assertion.Failure && !assertion.Error && !assertion.Skipped)
								{
									asserts_passed.push(assertion);
									pkg._passed ++;
									test._passed ++;
									story._passed ++;
								}
							}

							assignStatus(story);
						}
					}
				}

				ctx_overall.passed = asserts_passed.length;
				ctx_overall.panics = asserts_panicked.length;
				ctx_overall.failures = asserts_failed.length;
				ctx_overall.skipped = asserts_skipped.length;

				// Panics trump failures overall
				if (ctx_overall.panics)
					overallStatus = convey.statuses.panic;
				else if (ctx_overall.failures)
					overallStatus = convey.statuses.fail;

				// Show the overall status: passed, failed, or panicked
				if (overallStatus == 'pass')
					$('header').after(render('tpl-overall-ok', ctx_overall));
				else if (overallStatus == convey.statuses.fail)
					$('header').after(render('tpl-overall-fail', ctx_overall));
				else
					$('header').after(render('tpl-overall-panic', ctx_overall));

				// Show overall status
				$('.overall').slideDown();

				// Show shortucts and failures/panics details
				if (ctx_overall.panics > 0)
				{
					$('#left-sidebar').append(render('tpl-panic-shortcuts', asserts_panicked));
					$('#contents').append(render('tpl-panics', asserts_panicked));
				}
				if (ctx_overall.failures > 0)
				{
					$('#left-sidebar').append(render('tpl-failure-shortcuts', asserts_failed));
					$('#contents').append(render('tpl-failures', asserts_failed));
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