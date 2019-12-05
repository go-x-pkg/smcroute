package smcroute

const (
	nullCharacter       byte   = byte(0x00)
	nullCharacterString string = "\x00"

	responseBufferSize int = 255

	defaultSocketPath string = "/run/smcroute.sock"
	defaultNetwork    string = "unix"
)
