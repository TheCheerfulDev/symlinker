package entity

type Symlink struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type Symlinks struct {
	Links []Symlink `yaml:"symlinks"`
}
