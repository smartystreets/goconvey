convey.zen.gen = {
	tab: "\t",
	template: '',
	isFunc: function(scope)
	{
		if (!scope.title || typeof scope.depth === 'undefined')
			return false;

		return scope.title.indexOf("Test") === 0 && scope.depth == 0;
	}
};

$(function()
{
	var maxRows = 15;
	var lastKeyWasEnter = false;
	var gen = convey.zen.gen;


	window.onbeforeunload = function()
	{
		if ($('#gen-input').val().length > 20)
			return "Wait! You still have test cases in the GoConvey generator!";
	};


	$('#gen-input').keydown(function(e)
	{
		var rows = parseInt($(this).attr('rows'));
		if (e.keyCode == 13)	// Enter
		{
			lastKeyWasEnter = true;
			$(this).attr('rows', Math.min(rows + 1, maxRows));
		}
		else
			lastKeyWasEnter = false;
	}).keyup(function(e)
	{
		if (lastKeyWasEnter)
			return;

		if (e.keyCode == 13)
		{
			var rows = parseInt($(this).attr('rows'));
			$(this).attr('rows', Math.min(rows + 1, maxRows));
		}
		else
		{
			var newlines = $(this).val().match(/\n/g) || [];
			$(this).attr('rows', Math.min(newlines.length + 1, maxRows));
		}

		generate($(this).val());
	});

	gen.template = $('#tpl-convey').text();

	// Inserts a sample
	$('#gen-sample').click(function()
	{
		if ($('#gen-input').val().length > 10)
		{
			if (!confirm("This will clear your story. Continue?"))
				return false;
		}
		$('#gen-input').val('TestWhatGoConveyCanDo\n\tThe first line can be your "go test" function name\n\tIndented lines are tests to be wrapped in Convey()\n\t\tYou can--and should--nest your statements like this\n\t\tYou can fill out the details later.\n\tJust type away!').keyup();
	});


	var betterTextArea = new EnhancedTextArea('gen-input');

	// Original from: http://potch.me/projects/textarea
	// with fixes by yours truly
	function EnhancedTextArea(id, tab) {
		var el = document.getElementById(id);
		tabText = tab ? tab : "\t";
		el.onkeydown = function(e) {
			if (e.keyCode == 9) {
				var ta = el;
				var val = ta.value;
				var ss = ta.selectionStart;
				var nv = val.substring(0,ss) + tabText + val.substring(ss);
				ta.value = nv;
				ta.selectionStart = ss + tabText.length;
				ta.selectionEnd = ss + tabText.length;
				suppress(e);
			}
			if (e.keyCode == 13) {
				var ta = el;
				var val = ta.value;
				var ss = ta.selectionStart;
				var bl = val.lastIndexOf("\n",ss-1);
				var line = val.substring(bl,ss);
				var lm = line.match(/^\s+/);
				var ns = lm ? lm[0].length-1 : 0;
				var nv = val.substring(0,ss) + "\n";
				for (var i=0; i<ns; i++)
					nv += tabText;
				ta.value = nv+val.substring(ss);
				ta.selectionStart = ss + ns + 1;
				ta.selectionEnd = ss + ns + 1;
				suppress(e);
			}
		};
	}


	Mark.pipes.recursivelyRender = function(val)
	{
		return !val || val.length == 0 ? "\n" : Mark.up(gen.template, val);
	}

	Mark.pipes.indent = function(val)
	{
		return new Array(val + 1).join("\t");
	}

	Mark.pipes.notTestFunc = function(scope)
	{
		return !convey.zen.gen.isFunc(scope);
	}

	Mark.pipes.safeFunc = function(val)
	{
		return val.replace(/[^a-z0-9_]/gi, '');
	}

});




function generate(input)
{
	var root = parseInput(input);
	$('#gen-output').text(Mark.up(convey.zen.gen.template, root.stories));
}

function parseInput(input)
{
	lines = input.split("\n");
	
	if (!lines)
		return;

	var root = {
		title: "(root)",
		stories: []
	};

	for (i in lines)
	{
		line = lines[i];
		lineText = $.trim(line);

		if (!lineText)
			continue;

		// Figure out how deep to put this story
		indent = line.match(new RegExp("^" + convey.zen.gen.tab + "+"));
		tabs = indent ? indent[0].length / convey.zen.gen.tab.length : 0;

		// Starting at root, traverse into the right spot in the arrays
		var curScope = root, prevScope = root;
		for (j = 0; j < tabs && curScope.stories.length > 0; j++)
		{
			curScope = curScope.stories[curScope.stories.length - 1];
			prevScope = curScope;
		}
		
		// Don't go crazy, though! (avoid excessive indentation)
		if (tabs > curScope.depth + 1)
			tabs = curScope.depth + 1;

		// Only top-level Convey() calls need the *testing.T object passed in
		var showT = convey.zen.gen.isFunc(prevScope)
					|| (!convey.zen.gen.isFunc(curScope)
							&& tabs == 0);

		// Save the story at this scope
		curScope.stories.push({
			title: lineText.replace(/"/g, "\\\""),		// escape quotes
			stories: [],
			depth: tabs,
			showT: showT
		});
	}

	return root;
}