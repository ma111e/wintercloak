package engines

import "strings"

type Engines []EngineMeta

type EngineMeta struct {
	Name             string     `json:"name"`
	ShortDescription string     `json:"shortDescription"`
	LongDescription  string     `json:"longDescription"`
	Language         string     `json:"language"`
	IconifyIcon      string     `json:"iconifyIcon"` // https://icon-sets.iconify.design/?query=
	Args             EngineArgs `json:"args"`
}

var (
	Available            = Engines{}
	NameToFactoryMapping = map[string]func(args EngineArgs) (Engine, error){}
)

func (av *Engines) RegisterNewEngine(name string, short string, long string, lang string, iconifyIcon string, args EngineArgs, factory func(args EngineArgs) (Engine, error)) {
	if long == "" {
		long = short
	}

	NameToFactoryMapping[strings.ToLower(name)] = factory

	*av = append(*av, EngineMeta{
		Name:             name,
		ShortDescription: short,
		LongDescription:  long,
		Language:         lang,
		IconifyIcon:      iconifyIcon,
		Args:             args,
	})
}

func (av *Engines) List() []EngineMeta {
	return *av
}

func (av *Engines) GetArgsFor(engineName string) EngineArgs {
	for _, eng := range *av {
		if strings.ToLower(eng.Name) == engineName {
			return eng.Args
		}
	}

	return nil
}
