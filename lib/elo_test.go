package lib

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEloRanking(t *testing.T) {
	Convey("After a battle", t, func() {
		originalScore := 1500.0
		expected := Expected(originalScore, originalScore)
		Convey("Winner has a higher score", func() {
			winnerScore := Win(originalScore, expected)
			So(winnerScore, ShouldBeGreaterThan, originalScore)
		})
		Convey("Loser has a lower score", func() {
			loserScore := Loss(originalScore, expected)
			So(loserScore, ShouldBeLessThan, originalScore)
		})
	})
}
