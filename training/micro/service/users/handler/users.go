package handler

import (
	"context"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	users "users/proto"
)

var (
	ErrMissingFirstName  = errors.BadRequest("MISSING_FIRST_NAME", "Missing first name")
	ErrMissingLastName   = errors.BadRequest("MISSING_LAST_NAME", "Missing last name")
	ErrMissingEmail      = errors.BadRequest("MISSING_EMAIL", "Missing email")
	ErrDuplicateEmail    = errors.BadRequest("DUPLICATE_EMAIL", "A user with this email address already exists")
	ErrInvalidEmail      = errors.BadRequest("INVALID_EMAIL", "The email provided is invalid")
	ErrInvalidPassword   = errors.BadRequest("INVALID_PASSWORD", "Password must be at least 8 characters long")
	ErrIncorrectPassword = errors.Unauthorized("INCORRECT_PASSWORD", "Password wrong")
	ErrMissingIds        = errors.BadRequest("MISSING_IDS", "Missing ids")

	ErrConnectDB = errors.InternalServerError("CONNECT_DB", "Connecting to database failed")
	ErrNotFound  = errors.NotFound("NOT_FOUND", "User not found")

	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func isValidEmail(email string) bool {
	if len(email) == 0 || len(email) > 255 {
		return false
	}
	return emailRegex.MatchString(email)
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func genHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func compareHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time
}

func (u *User) sanitize() *users.User {
	return &users.User{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}

type Users struct {
	DB *gorm.DB
}

func (h *Users) Create(ctx context.Context, req *users.CreateRequest, rsp *users.CreateResponse) error {
	// validate request
	if len(req.FirstName) == 0 {
		return ErrMissingFirstName
	}
	if len(req.LastName) == 0 {
		return ErrMissingLastName
	}
	if isValidEmail(req.Email) == false {
		return ErrInvalidEmail
	}
	if isValidPassword(req.Password) == false {
		return ErrInvalidPassword
	}

	// hash password
	pwdHashed, err := genHash(req.Password)
	if err != nil {
		logger.Errorf("Hash password: %v", err)
		return errors.InternalServerError("HASH_PWD_FAILED", "Hash password failed")
	}

	// create
	return h.DB.Transaction(func(tx *gorm.DB) error {
		user := &User{
			ID:        uuid.New().String(),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     strings.ToLower(req.Email),
			Password:  pwdHashed,
		}
		if err := tx.Create(user).Error; err != nil && strings.Contains(err.Error(), "idx_users_email") {
			return ErrDuplicateEmail
		} else if err != nil {
			logger.Errorf("create user: %v", err)
			return ErrConnectDB
		}

		rsp.User = user.sanitize()

		return nil
	})
}

func (h *Users) Login(ctx context.Context, req *users.LoginRequest, rsp *users.LoginResponse) error {
	// validate request
	if len(req.Email) == 0 {
		return ErrMissingEmail
	}
	if len(req.Password) == 0 {
		return ErrInvalidPassword
	}

	return h.DB.Transaction(func(tx *gorm.DB) error {
		// lookup user
		var user User
		if err := tx.Where(&User{Email: strings.ToLower(req.Email)}).First(&user).Error; err == gorm.ErrRecordNotFound {
			return ErrNotFound
		} else if err != nil {
			logger.Errorf("Error reading from db: %v", err)
			return ErrConnectDB
		}

		// compare password
		if err := compareHash(user.Password, req.Password); err != nil {
			return ErrIncorrectPassword
		}

		rsp.User = user.sanitize()

		return nil
	})
}

func (h *Users) List(ctx context.Context, req *users.ListRequest, rsp *users.ListResponse) error {
	var us []User
	if err := h.DB.Order("created_at desc").Find(&us).Error; err != nil {
		logger.Errorf("Error reading from db: %v", err)
		return ErrConnectDB
	}

	rsp.Users = make([]*users.User, len(us))
	for i, user := range us {
		rsp.Users[i] = user.sanitize()
	}
	return nil
}

func (h *Users) Read(ctx context.Context, req *users.ReadRequest, rsp *users.ReadResponse) error {
	if len(req.Ids) == 0 {
		return ErrMissingIds
	}

	var us []User
	// db.Find(&users, []int{1,2,3})
	// db.Where([]int64{20, 21, 22}).Find(&users)
	if err := h.DB.Where("id IN ?", req.Ids).Find(&us).Error; err != nil {
		logger.Errorf("Error reading from db: %v", err)
		return ErrConnectDB
	}

	// map not order
	rsp.Users = make(map[string]*users.User, len(us))
	for _, user := range us {
		rsp.Users[user.ID] = user.sanitize()
	}
	return nil
}

func (h *Users) ReadByEmail(ctx context.Context, req *users.ReadByEmailRequest, rsp *users.ReadByEmailResponse) error {
	if len(req.Emails) == 0 {
		return ErrMissingEmail
	}

	emails := make([]string, len(req.Emails))
	for i, email := range req.Emails {
		emails[i] = strings.ToLower(email)
	}

	var us []User
	// db.Find(&users, []int{1,2,3})
	// db.Where([]int64{20, 21, 22}).Find(&users)
	if err := h.DB.Where("lower(email) IN ?", emails).Find(&us).Error; err != nil {
		logger.Errorf("Error reading from db: %v", err)
		return ErrConnectDB
	}

	// map not order
	rsp.Users = make(map[string]*users.User, len(us))
	for _, user := range us {
		rsp.Users[user.Email] = user.sanitize()
	}
	return nil
}
