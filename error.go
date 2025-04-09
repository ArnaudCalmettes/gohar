package gohar

//go:generate gen-errfuncs $GOFILE

import (
	"errors"
	"fmt"
)

var (
	ErrBufferOverflow      = errors.New("buffer overflow")
	ErrNilBuffer           = errors.New("nil buffer")
	ErrInvalidBaseNote     = errors.New("invalid base note")
	ErrInvalidAlteration   = errors.New("invalid alteration")
	ErrUnknownScalePattern = errors.New("unknown scale pattern")
	ErrInvalidDegree       = errors.New("invalid degree")
)

func checkOutputBuffer[T any](buffer []T, capacity int) error {
	if buffer == nil {
		return ErrNilBuffer
	}
	if cap(buffer) < capacity {
		return fmt.Errorf(
			"%w: output slice has capacity %d < %d needed",
			ErrBufferOverflow, cap(buffer), capacity,
		)
	}
	return nil
}

func checkValidNote(note Note) error {
	var err error
	if note.Base < 'A' || note.Base > 'G' {
		err = errors.Join(err,
			wrapErrorf(ErrInvalidBaseNote, "'%c'", note.Base),
		)
	}
	if note.Alt < -2 || note.Alt > 2 {
		err = errors.Join(err,
			wrapErrorf(ErrInvalidAlteration, "%d", note.Alt),
		)
	}
	return err
}

func wrapErrorf(err error, msg string, args ...any) error {
	return fmt.Errorf("%w: %s", err, fmt.Sprintf(msg, args...))
}
