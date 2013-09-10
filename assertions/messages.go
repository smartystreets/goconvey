package assertions

const ( // equality
	shouldHaveBeenEqual         = "Expected '%v'\nto equal '%v'\n(but it didn't)!"
	shouldNotHaveBeenEqual      = "Expected     '%v'\nto NOT equal '%v'\n(but it did)!"
	shouldHaveResembled         = "Expected '%v'\nto resemble '%v'\n(but it didn't)!"
	shouldNotHaveResembled      = "Expected        '%v'\nto NOT resemble '%v'\n(but it did)!"
	shouldBePointers            = "Both arguments should be pointers "
	shouldHaveBeenNonNilPointer = shouldBePointers + "(the %s was %s)!"
	shouldHavePointedTo         = "Expected '%v' (address: '%v') and '%v' (address: '%v') to be the same address (but their weren't)!"
	shouldNotHavePointedTo      = "Expected '%v' and '%v' to be different references (but they matched: '%v')!"
	shouldHaveBeenNil           = "Expected '%v' to be nil (but it wasn't)!"
	shouldNotHaveBeenNil        = "Expected '%v' to NOT be nil (but it was)!"
	shouldHaveBeenTrue          = "Expected 'true' (not '%v')!"
	shouldHaveBeenFalse         = "Expected 'false' (not '%v')!"
)

const ( // quantity comparisons
	shouldHaveBeenGreater            = "Expected '%v' to be greater than '%v' (but it wasn't)!"
	shouldHaveBeenGreaterOrEqual     = "Expected '%v' to be greater than or equal to '%v' (but it wasn't)!"
	shouldHaveBeenLess               = "Expected '%v' to be less than '%v' (but it wasn't)!"
	shouldHaveBeenLessOrEqual        = "Expected '%v' to be less than or equal to '%v' (but it wasn't)!"
	shouldHaveBeenBetween            = "Expected '%v' to be between '%v' and '%v' (but it wasn't)!"
	shouldNotHaveBeenBetween         = "Expected '%v' NOT to be between '%v' and '%v' (but it was)!"
	shouldHaveDifferentUpperAndLower = "The lower and upper bounds must be different values (they were both '%v')."
	shouldHaveBeenBetweenOrEqual     = "Expected '%v' to be between '%v' and '%v' or equal to one of them (but it wasn't)!"
	shouldNotHaveBeenBetweenOrEqual  = "Expected '%v' NOT to be between '%v' and '%v' or equal to one of them (but it was)!"
)

const ( // collections
	shouldHaveContained                 = "Expected the container (%v) to contain: '%v' (but it didn't)!"
	shouldNotHaveContained              = "Expected the container (%v) NOT to contain: '%v' (but it did)!"
	shouldHaveBeenIn                    = "Expected '%v' to be in the container (%v, but it wasn't)!"
	shouldNotHaveBeenIn                 = "Expected '%v' NOT to be in the container (%v, but it was)!"
	shouldHaveBeenAValidCollection      = "You must provide a valid container (was %v)!"
	shouldHaveProvidedCollectionMembers = "This assertion requires at least 1 comparison value (you provided 0)."
)

const ( // strings
	shouldHaveStartedWith    = "Expected '%v' to start with: \n         '%v' (but it didn't)!"
	shouldNotHaveStartedWith = "Expected '%v' NOT to start with: \n         '%v' (but it did)!"
	shouldHaveEndedWith      = "Expected '%v' to end with: \n         '%v' (but it didn't)!"
	shouldNotHaveEndedWith   = "Expected '%v' NOT to end with: \n         '%v' (but it didn't)!"
)

const ( // panics
	shouldHavePanickedWith = "Expected func() to panic with '%v' (but it panicked with '%v')!"
	shouldHavePanicked     = "Expected func() to panic with '%v' (but it didn't panic at all)!"
)

const ( // type checking
	shouldHaveBeenA    = "Expected '%v' to be a '%v' (but was a '%v')!"
	shouldNotHaveBeenA = "Expected '%v to NOT be a '%v' (but it was)!"
)
