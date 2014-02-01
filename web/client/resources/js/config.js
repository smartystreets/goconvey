var convey = convey || {};

convey.config = {
	statuses: {
		pass: {
			class: 'ok',
			text: "PASS"
		},
		fail: {
			class: 'fail',
			text: "FAIL"
		},
		panic: {
			class: 'panic',
			text: "PANIC"
		},
		buildfail: {
			class: 'buildfail',
			text: "BUILD FAILED"
		}
	},
};