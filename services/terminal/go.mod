module terminal

go 1.19

replace (
	github.com/Inoi-K/Find-Me/pkg v0.0.0 => ../../pkg
	github.com/Inoi-K/Find-Me/services/recommendations v0.0.0 => ../recommendations
)

require (
	github.com/Inoi-K/Find-Me/pkg v0.0.0
	github.com/Inoi-K/Find-Me/services/recommendations v0.0.0
)
