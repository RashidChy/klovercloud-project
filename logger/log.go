package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func Init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

}
