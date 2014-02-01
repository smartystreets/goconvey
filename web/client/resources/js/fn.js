var convey = convey || {};
convey.fn = {};

convey.fn.emptyOverall = function() {
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
};