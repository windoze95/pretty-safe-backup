package setup

type Setup struct {
	Name        string
	Description string
	Source      string
	Excludes    []string
	//      Destination string
}

func (s Setup) submittable() bool {
	return s.Name != "" && s.Source != ""
}

var (
	answers *Setup = &Setup{}
)

func Build() {
	mainMenu()
}
