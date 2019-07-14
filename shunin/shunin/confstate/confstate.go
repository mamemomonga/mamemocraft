package confstate

import (
	cs "github.com/mamemomonga/go-confstate"
)

type Configs struct {
	Rcon       CRcon   `yaml:"rcon"`
	Shunin     CShunin `yaml:"mastodon"`
}

type CRcon struct {
	Host    string `yaml:"host"`
	Password  string `yaml:"password"`
}

type CShunin struct {
	Server     string `yaml:"server"`
	Email      string `yaml:"email"`
	Password   string `yaml:"password"`
}

type States struct {
}

// Load Initalize and load
func Load(cf string) error {
	cs.ConfigsFile        = cf        // ConfigsFile Filename
	cs.DefaultBaseDirType = cs.DBTBin // Basedir is offset from executable binary
	cs.OffsetFromBin      = ".."      // Base directory offset from executable binary
	cs.DefaultConfigsFile = "etc/configs.yaml" // Configs File default filename
	cs.DefaultStatesFile  = "etc/states.json"  // States File default filename
	cs.Debug = true                   // Debug mode

	// Initalize Configs
	cs.Configs = &Configs{
		Rcon: CRcon{},
		Shunin: CShunin{},
	}
	// Initalize States
	cs.States = &States{
	}

	if err := cs.LoadConfigs(); err != nil {
		return err
	}

	if err := cs.LoadStates(); err != nil {
		return err
	}
	return nil
}

// C Accessor for Configs
func C() *Configs {
	return cs.Configs.(*Configs)
}

// S Accessor for States
func S() *States {
	return cs.States.(*States)
}

// GetDir Accessor for GetDir(Not required)
func GetDir(p string) (string, error) {
	return cs.GetDir(p)
}

// Save States file
func Save() error {
	return cs.SaveStates()
}


