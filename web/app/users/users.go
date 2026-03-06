package users

import (
	"math/rand"

	"github.com/gin-gonic/gin"
)

func GETUsers(ctx *gin.Context) {
	username := GenerateUsername()
	ctx.JSON(200, gin.H{
		"username": username,
	})
}

func POSTUsers(ctx *gin.Context) {
	username := GenerateUsername()
	ctx.JSON(201, gin.H{
		"result":   "user_created",
		"username": username,
	})
}

var (
	adjectives = []string{
		"silent", "swift", "dark", "bright", "cool", "wild", "calm", "bold",
		"epic", "lazy", "sneaky", "fuzzy", "mighty", "crazy", "lucky", "happy",
		"clever", "fierce", "gentle", "shiny", "rusty", "dusty", "golden", "silver",
	}

	nouns = []string{
		"wolf", "fox", "hawk", "bear", "tiger", "panda", "eagle", "shark",
		"ninja", "coder", "gamer", "rider", "hunter", "ranger", "wizard", "knight",
		"pixel", "storm", "blaze", "frost", "shadow", "thunder", "comet", "nova",
	}

	patterns = []func() string{
		// adjective + noun (e.g. "silentfox")
		func() string {
			return pick(adjectives) + pick(nouns)
		},
		// adjective + underscore + noun (e.g. "swift_hawk")
		func() string {
			return pick(adjectives) + "_" + pick(nouns)
		},
		// noun + underscore + noun (e.g. "fox_hunter")
		func() string {
			return pick(nouns) + "_" + pick(nouns)
		},
		// x + adjective + noun (e.g. "xepicbear")
		func() string {
			return "x" + pick(adjectives) + pick(nouns)
		},
		// adjective + adjective + noun (e.g. "darkcalmwolf")
		func() string {
			return pick(adjectives) + pick(adjectives) + pick(nouns)
		},
	}
)

func pick(list []string) string {
	return list[rand.Intn(len(list))]
}

// GenerateUsername returns a random realistic-looking username without numbers.
func GenerateUsername() string {
	pattern := patterns[rand.Intn(len(patterns))]
	return pattern()
}
