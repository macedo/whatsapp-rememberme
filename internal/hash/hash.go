package hash

import "golang.org/x/crypto/bcrypt"

func NewEncryptor(pepper string) Encryptor {
	return Encryptor{[]byte(pepper)}
}

type Encryptor struct {
	pepper []byte
}

func (e *Encryptor) Compare(hashed_input, input string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashed_input),
		append([]byte(input), e.pepper...))

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false
		}

		panic(err)
	}

	return true
}

func (e *Encryptor) Digest(input string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		append([]byte(input), e.pepper...),
		bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (e *Encryptor) MustDigest(input string) string {
	digest, err := e.Digest(input)
	if err != nil {
		panic(`hash: Digest(` + input + `): ` + err.Error())
	}

	return digest
}
