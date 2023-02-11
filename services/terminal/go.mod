module github.com/Inoi-K/Find-Me/services/terminal

go 1.19

replace (
	github.com/Inoi-K/Find-Me/pkg v0.0.0 => ../../pkg
	github.com/Inoi-K/Find-Me/services/rengine v0.0.0 => ./../rengine
)

require (
	github.com/Inoi-K/Find-Me/pkg v0.0.0
	github.com/Inoi-K/Find-Me/services/rengine v0.0.0
)
