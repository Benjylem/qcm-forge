package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Paramètres argon2id : minimum recommandé par l'OWASP. 19 Mo par hash
// plutôt que 64 Mo : sur un petit VPS, des connexions simultanées ne
// doivent pas épuiser la RAM. Les paramètres sont écrits dans chaque
// hash, on peut donc les augmenter plus tard sans casser l'existant.
const (
	argonMemory  = 19 * 1024 // en Kio
	argonTime    = 2         // nombre de passes
	argonThreads = 1
	saltLen      = 16 // octets
	keyLen       = 32 // octets
)

var ErrInvalidHash = errors.New("auth: hash au format invalide")

// HashPassword renvoie une chaîne au format standard (PHC) :
// $argon2id$v=19$m=19456,t=2,p=1$<sel base64>$<hash base64>
func HashPassword(password string) (string, error) {
	// Sel aléatoire, différent pour chaque utilisateur : deux mots de
	// passe identiques donnent deux hash différents.
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, keyLen)

	b64 := base64.RawStdEncoding
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonTime, argonThreads,
		b64.EncodeToString(salt), b64.EncodeToString(key)), nil
}

// CheckPassword recalcule le hash du mot de passe saisi avec le sel et
// les paramètres lus dans encoded, puis compare.
func CheckPassword(password, encoded string) (bool, error) {
	// "$argon2id$v=19$m=...$sel$hash" découpé sur "$" donne 6 morceaux,
	// le premier étant vide.
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, ErrInvalidHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false, ErrInvalidHash
	}

	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false, ErrInvalidHash
	}

	b64 := base64.RawStdEncoding
	salt, err := b64.DecodeString(parts[4])
	if err != nil {
		return false, ErrInvalidHash
	}
	want, err := b64.DecodeString(parts[5])
	if err != nil {
		return false, ErrInvalidHash
	}

	got := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(want)))

	// Comparaison en temps constant : un == classique s'arrête au premier
	// octet différent, et la durée de la réponse renseignerait un attaquant.
	return subtle.ConstantTimeCompare(got, want) == 1, nil
}
