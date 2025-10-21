package api

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type signupRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type signupResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

// TrinityCore-compatible SRP6 parameters and helper.
// g = 7, N = 0x894B... (WoW SRP6 modulus)
var (
	srpG = big.NewInt(7)
	srpN = func() *big.Int {
		// This is the standard WoW SRP6 prime used by TrinityCore.
		n, _ := new(big.Int).SetString("894B645E89E1535BBDAD5B8B290650530801B18EBFBF5E8FAB3C82872A3E9BB7", 16)
		return n
	}()
)

// computeSRP6 generates a random 32-byte salt and computes the verifier as v = g^x mod N,
// where x = SHA1(salt || SHA1(UPPER(username) ":" password)).
// Returns salt and verifier as big-endian byte slices, left-padded to the modulus size.
func computeSRP6(usernameUpper, password string) ([]byte, []byte, error) {
	// salt: 32 bytes random
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	// inner = SHA1(UPPER(username) + ":" + password)
	inner := sha1.Sum([]byte(usernameUpper + ":" + password))

	// x = SHA1(salt || inner)
	h := sha1.New()
	_, _ = h.Write(salt)
	_, _ = h.Write(inner[:])
	xBytes := h.Sum(nil)
	x := new(big.Int).SetBytes(xBytes)

	// v = g^x mod N
	v := new(big.Int).Exp(srpG, x, srpN)

	// left-pad v to modulus size (in bytes)
	nLen := (srpN.BitLen() + 7) / 8
	vBytes := v.Bytes()
	if len(vBytes) < nLen {
		padded := make([]byte, nLen)
		copy(padded[nLen-len(vBytes):], vBytes)
		vBytes = padded
	}

	// TrinityCore schema typically uses VARBINARY(32) for salt/verifier.
	// If modulus length exceeds 32, truncate to 32 least-significant bytes.
	if len(salt) > 32 {
		salt = salt[len(salt)-32:]
	}
	if len(vBytes) > 32 {
		vBytes = vBytes[len(vBytes)-32:]
	}

	return salt, vBytes, nil
}

func ensureUnique(db *sql.DB, usernameUpper, email string) error {
	var cnt int
	if err := db.QueryRow("SELECT COUNT(1) FROM account WHERE username = ?", usernameUpper).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("username already exists")
	}
	cnt = 0
	if err := db.QueryRow("SELECT COUNT(1) FROM account WHERE email = ?", email).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("email already exists")
	}
	return nil
}

func insertAccount(db *sql.DB, usernameUpper, email, password string) (int64, error) {
	salt, verifier, err := computeSRP6(usernameUpper, password)
	if err != nil {
		return 0, err
	}

	res, err := db.Exec(`
        INSERT INTO account (
            username, salt, verifier, email, reg_mail
        ) VALUES (
            ?, ?, ?, ?, ?
        )
    `, usernameUpper, salt, verifier, email, email)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// CreateAccount performs validation, uniqueness checks and DB insert. Useful for tests without HTTP.
func CreateAccount(db *sql.DB, username, email, password string) (int64, string, error) {
	usernameUpper := strings.ToUpper(strings.TrimSpace(username))
	if usernameUpper == "" {
		return 0, "", errors.New("invalid username")
	}
	if err := ensureUnique(db, usernameUpper, email); err != nil {
		return 0, "", err
	}
	id, err := insertAccount(db, usernameUpper, email, password)
	if err != nil {
		return 0, "", err
	}
	return id, usernameUpper, nil
}

func HandleSignup(c *gin.Context, db *sql.DB) {
	var body signupRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, usernameUpper, err := CreateAccount(db, body.Username, body.Email, body.Password)
	if err != nil {
		status := http.StatusInternalServerError
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, signupResponse{ID: id, Username: usernameUpper, Email: body.Email})
}
