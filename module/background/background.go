package background

import (
	"github.com/robfig/cron"
)

var c *cron.Cron

// SetBackground -
func SetBackground() {
	c = cron.New()
	c.AddFunc("0 * * * * *", readFacebookPage)
	c.Start()
}
