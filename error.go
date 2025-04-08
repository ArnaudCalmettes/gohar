package gohar

//go:generate gen-errfuncs $GOFILE

import (
	"errors"
	"fmt"
)

var (
	ErrBufferOverflow         = errors.New("buffer overflow")
	ErrNilBuffer              = errors.New("nil buffer")
	ErrInvalidBaseNote        = errors.New("invalid base note")
	ErrNonPrintableAlteration = errors.New("non-printable alteration")
	ErrUnknownScalePattern    = errors.New("unknown scale pattern")
	ErrInvalidDegree          = errors.New("invalid degree")
)

func CheckOutputBuffer[T any](buffer []T, capacity int) error {
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

func CheckNoteIsPrintable(note Note) error {
	var err error
	if note.Base < 'A' || note.Base > 'G' {
		err = errors.Join(err,
			fmt.Errorf("%w: '%c'", ErrInvalidBaseNote, note.Base),
		)
	}
	if note.Alt < -2 || note.Alt > 2 {
		err = errors.Join(err,
			fmt.Errorf("%w: %d", ErrNonPrintableAlteration, note.Alt),
		)
	}
	return err
}
