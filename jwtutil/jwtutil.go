package jwtutil

import (
	"log"

	//	"github.com/gorilla/securecookie"
	"gopkg.in/gorilla/securecookie.v1"
)

// -----------------------------------------------------------------------------

// Flags is a package flags sample
// in form ready for use with github.com/jessevdk/go-flags
type Flags struct {
	AppKey   string `long:"psw_session_key" description:"Key to encode user session (default: random key reset on restart)"`
	BlockKey string `long:"psw_block_key" default:"T<8rYvXmgLBdND(YW}3QRcLwh4$4P5eq" description:"Key to encode session blocks (16,32 or 62 byte)"`
}

type Cryptor interface {
	Decode(name, value string, dst interface{}) error
	Encode(name string, value interface{}) (string, error)
}

// App is a package general type
type App struct {
	Log     *log.Logger
	Config  *Flags
	Cryptor Cryptor
}

// -----------------------------------------------------------------------------

// New creates mailer object
// Configuration should be set via functional options
func New(logger *log.Logger, cfg *Flags, options ...func(a *App) error) (*App, error) {
	a := App{Config: cfg, Log: logger}
	for _, option := range options {
		err := option(&a)
		if err != nil {
			return nil, err
		}
	}
	if a.Cryptor == nil {
		a.setDefaultCryptor()
	}
	return &a, nil
}

func (a *App) setDefaultCryptor() error {

	var hashKeyBytes = []byte(a.Config.AppKey)
	if a.Config.AppKey == "" {
		hashKeyBytes = securecookie.GenerateRandomKey(32)
		a.Log.Print("info: Random key generated. Sessions will be expired on restart")
	}
	var blockKeyBytes = []byte(a.Config.BlockKey) // "txVzHcURYJrK]UQ:d/YDmx97*Adwb;/%")

	var s = securecookie.New(hashKeyBytes, blockKeyBytes)
	s.SetSerializer(securecookie.JSONEncoder{})
	a.Cryptor = s
	return nil
}
