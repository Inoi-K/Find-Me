package consts

const (
	MainSphereCoefficient  = 1.0
	OtherSphereCoefficient = .2
)

// TODO find datasets of tags for spheres!
const (
	// WORK

	GameDev  = "gamedev"
	Python   = "python"
	CSharp   = "#"
	Clojure  = "clojure"
	Java     = "java"
	CPP      = "c++"
	Golang   = "golang"
	Ruby     = "ruby"
	Rust     = "rust"
	PHP      = "php"
	Solidity = "solidity"

	Swift  = "swift"
	Kotlin = "kotlin"

	HTML = "html"
	CSS  = "css"
	JS   = "javascript"
	TS   = "typescript"

	Scala   = "scala"
	Haskell = "haskell"
	Erlang  = "erlang"
	Elixir  = "elixir"

	// LOVE

	Walking = "walking"

	Dogs      = "dogs"
	Cats      = "cats"
	Elephants = "elephants"
	// Example of confusion between panda animal and pandas python library
	Pandas = "pandas"

	// CHILL

	Dota = "dota"
	// Space handling via regexp required: cs:go = cs go = csgo
	CS = "cs"
)

var (
	Synonyms = map[string]string{
		// WORK
		GameDev:            GameDev,
		"game development": GameDev,

		Python:    Python,
		"python2": Python,
		"python3": Python,
		"py":      Python,
		"pip":     Python,

		CSharp:   CSharp,
		"csharp": CSharp,

		Clojure: Clojure,

		Java: Java,

		CPP:   CPP,
		"cpp": CPP,

		Golang: Golang,
		// can be easily misinterpreted
		"go": Golang,

		Ruby:    Ruby,
		"rails": Ruby,

		Rust: Rust,

		PHP: PHP,

		Solidity: Solidity,

		Swift: Swift,

		Kotlin: Kotlin,

		HTML: HTML,

		CSS: CSS,

		JS:   JS,
		"js": JS,

		TS:   TS,
		"ts": TS,

		Scala: Scala,

		Haskell: Haskell,

		Erlang: Erlang,

		Elixir: Elixir,

		// LOVE

		// TODO Handle a lot of single/plural repetitions
		Walking: Walking,
		//"go": Walking,
		"walk": Walking,

		Dogs:  Dogs,
		"dog": Dogs,

		Cats:  Cats,
		"cat": Cats,

		Elephants:  Elephants,
		"elephant": Elephants,

		Pandas:  Pandas,
		"panda": Pandas,

		// CHILL
		Dota:    Dota,
		"dota2": Dota,

		CS:      CS,
		"csgo":  CS,
		"cs:go": CS,
		"cs go": CS,
	}
)
