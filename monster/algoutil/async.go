package algoutil

import "github.com/sirupsen/logrus"

func pcall(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("aync/pcall: error=%v", err)
		}
	}()
	fn()
}
func Run(fn func()) {
	go pcall(fn)
}
