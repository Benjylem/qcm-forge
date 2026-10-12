package auth

import (
	"strings"
	"testing"
)

func TestHashThenCheck(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(hash, "$argon2id$v=19$m=19456,t=2,p=1$") {
		t.Errorf("format inattendu : %s", hash)
	}

	ok, err := CheckPassword("correct horse battery staple", hash)
	if err != nil || !ok {
		t.Errorf("le bon mot de passe doit passer (ok=%v, err=%v)", ok, err)
	}

	ok, err = CheckPassword("mauvais", hash)
	if err != nil || ok {
		t.Errorf("un mauvais mot de passe doit échouer (ok=%v, err=%v)", ok, err)
	}
}

// Même mot de passe, deux hash différents : preuve que le sel est aléatoire.
func TestSaltIsRandom(t *testing.T) {
	a, _ := HashPassword("identique")
	b, _ := HashPassword("identique")
	if a == b {
		t.Error("deux hash du même mot de passe ne doivent pas être égaux")
	}
}

func TestCheckRejectsMalformedHash(t *testing.T) {
	for _, bad := range []string{
		"",
		"motdepasseenclair",
		"$2a$12$unhashbcrypt",
		"$argon2id$v=19$m=19456,t=2,p=1$pasdubase64!$xxx",
	} {
		if _, err := CheckPassword("x", bad); err != ErrInvalidHash {
			t.Errorf("%q : attendu ErrInvalidHash, obtenu %v", bad, err)
		}
	}
}
