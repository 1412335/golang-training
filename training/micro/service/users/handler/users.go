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
	ErrMissingId         = errors.BadRequest("MISSING_ID", "Missing id")

	ErrConnectDB = errors.InternalServerError("CONNECT_DB", "Connecting to database failed")
	ErrNotFound  = errors.NotFound("NOT_FOUND", "User not found")

	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	ttlToken   = 24 * time.Hour
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
	Tokens    []Token
}

func (u *User) sanitize() *users.User {
	return &users.User{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}

type Token struct {
	Key       string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	UserID    string // foreign key
	// User      User
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
			logger.Errorf("Error connecting from db: %v", err)
			return ErrConnectDB
		}

		// token
		token := &Token{
			Key:       uuid.New().String(),
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(ttlToken),
		}
		if err := tx.Create(token).Error; err != nil {
			logger.Errorf("Error connecting from db: %v", err)
			return ErrConnectDB
		}

		rsp.User = user.sanitize()
		rsp.Token = token.Key

		return nil
	})
}

func (h *Users) Update(ctx context.Context, req *users.UpdateRequest, rsp *users.UpdateResponse) error {
	if len(req.Id) == 0 {
		return ErrMissingId
	}
	if req.FirstName != nil && len(req.FirstName.Value) == 0 {
		return ErrMissingFirstName
	}
	if req.LastName != nil && len(req.LastName.Value) == 0 {
		return ErrMissingLastName
	}
	if req.Email != nil && isValidEmail(req.Email.Value) == false {
		return ErrInvalidEmail
	}
	if req.Password != nil && isValidPassword(req.Password.Value) == false {
		return ErrInvalidPassword
	}

	// lookup user
	var user User
	if err := h.DB.Where(&User{ID: req.Id}).First(&user).Error; err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error connecting from db: %v", err)
		return ErrConnectDB
	}

	if req.FirstName != nil {
		user.FirstName = req.FirstName.Value
	}
	if req.LastName != nil {
		user.LastName = req.LastName.Value
	}
	if req.Email != nil {
		user.Email = strings.ToLower(req.Email.Value)
	}
	if req.Password != nil {
		// hash password
		pwdHashed, err := genHash(req.Password.Value)
		if err != nil {
			logger.Errorf("Hash password: %v", err)
			return errors.InternalServerError("HASH_PWD_FAILED", "Hash password failed")
		}
		user.Password = pwdHashed
	}

	if err := h.DB.Save(user).Error; err != nil && strings.Contains(err.Error(), "idx_users_email") {
		return ErrDuplicateEmail
	} else if err != nil {
		logger.Errorf("Error connecting from db: %v", err)
		return ErrConnectDB
	}

	rsp.User = user.sanitize()
	return nil
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
			logger.Errorf("Error connecting from db: %v", err)
			return ErrConnectDB
		}

		// compare password
		if err := compareHash(user.Password, req.Password); err != nil {
			return ErrIncorrectPassword
		}

		// generate token
		token := &Token{
			Key:       uuid.New().String(),
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(ttlToken),
		}
		if err := tx.Create(token).Error; err != nil {
			logger.Errorf("Error connecting from db: %v", err)
			return ErrConnectDB
		}

		rsp.User = user.sanitize()
		rsp.Token = token.Key

		return nil
	})
}

func (h *Users) List(ctx context.Context, req *users.ListRequest, rsp *users.ListResponse) error {
	var us []User
	if err := h.DB.Order("created_at desc").Find(&us).Error; err != nil {
		logger.Errorf("Error connecting from db: %v", err)
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
		return ErrMissingId
	}

	var us []User
	// db.Find(&users, []int{1,2,3})
	// db.Where([]int64{20, 21, 22}).Find(&users)
	if err := h.DB.Where("id IN ?", req.Ids).Find(&us).Error; err != nil {
		logger.Errorf("Error connecting from db: %v", err)
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
		logger.Errorf("Error connecting from db: %v", err)
		return ErrConnectDB
	}

	// map not order
	rsp.Users = make(map[string]*users.User, len(us))
	for _, user := range us {
		rsp.Users[user.Email] = user.sanitize()
	}
	return nil
}

func (h *Users) Delete(ctx context.Context, req *users.DeleteRequest, rsp *users.DeleteResponse) error {
	if len(req.Ids) == 0 {
		return ErrMissingId
	}

	return h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(req.Ids).Delete(&User{}).Error; err != nil {
			logger.Errorf("Error connecting from db: %v", err)
			return ErrConnectDB
		}
		return nil
	})
}
