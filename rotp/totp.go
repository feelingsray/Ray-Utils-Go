package rotp

import "time"

type TOTP struct {
	*HOTP
	x int64
}

func NewTOTP(secret []byte, digits, x int) (t *TOTP) {
	t = new(TOTP)
	t.HOTP = NewHOTP(secret, digits)
	t.x = int64(x)
	return t
}

func (t TOTP) At(timestamp int64) string {
	counter := uint64(timestamp / t.x)
	return t.HOTP.At(counter)
}

func (t TOTP) Now() string {
	now := time.Now()
	timestamp := now.Unix()
	return t.At(timestamp)
}

func (t TOTP) Time(dt time.Time) string {
	dts := time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), 0, 0, time.Local)
	return t.At(dts.Unix())
}

func (t TOTP) Verify(code string) bool {
	return t.Now() == code
}

func (t TOTP) VerifyWithTime(code string, dt time.Time) bool {
	return t.Time(dt) == code
}
