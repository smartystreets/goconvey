package examples

type Game struct {
	rolls     []int
	rollIndex int
}

func NewGame() *Game {
	game := Game{}
	game.rolls = make([]int, 21)
	return &game
}

func (self *Game) Roll(pins int) {
	self.rolls[self.rollIndex] = pins
	self.rollIndex++
}

func (self *Game) Score() int {
	sum, throw, frame := 0, 0, 0
	for ; frame < 10; frame++ {
		if self.isStrike(throw) {
			sum += self.strikeBonus(throw)
			throw += 1
		} else if self.isSpare(throw) {
			sum += self.spareBonus(throw)
			throw += 2
		} else {
			sum += self.currentFrame(throw)
			throw += 2
		}
	}
	return sum
}

func (self *Game) isStrike(throw int) bool {
	return self.rolls[throw] == 10
}
func (self *Game) isSpare(throw int) bool {
	return self.rolls[throw]+self.rolls[throw+1] == 10
}
func (self *Game) strikeBonus(throw int) int {
	return 10 + self.rolls[throw+1] + self.rolls[throw+2]
}
func (self *Game) spareBonus(throw int) int {
	return 10 + self.rolls[throw+2]
}
func (self *Game) currentFrame(throw int) int {
	return self.rolls[throw] + self.rolls[throw+1]
}
