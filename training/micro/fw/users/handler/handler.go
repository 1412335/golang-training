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

	pb "fw/users/proto"
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
	ErrMissingToken      = errors.BadRequest("MISSING_TOKEN", "Missing token")

	ErrConnectDB = errors.InternalServerError("CONNECT_DB", "Connecting to database failed")
	ErrNotFound  = errors.NotFound("NOT_FOUND", "User not found")

	ErrTokenGenerated = errors.InternalServerError("TOKEN_GEN_FAILED", "Generate token failed")
	ErrTokenInvalid   = errors.Unauthorized("TOKEN_INVALID", "Token invalid")
	ErrTokenNotFound  = errors.BadRequest("TOKEN_NOT_FOUND", "Token not found")
	ErrTokenExpired   = errors.Unauthorized("TOKEN_EXPIRE", "Token expired")

	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	// ttlToken   = 24 * time.Hour
)

func isValidEmail(email string) bool {
	if len(email) == 0 || len(email) > 255 {
		return false
	}
	return emailRegex.MatchString(email)
}

func isValidPassword(password string) bool {
	return len(password) >= 8
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
	ValidFrom time.Time
	ValidTo   time.Time
	Active    bool
	Password  string
	Name      string
	Email     string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) sanitize() *pb.User {
	return &pb.User{}
}

type usersHandler struct {
	DB         *gorm.DB
	JWTManager *JWTManager
}

func NewUsersHandler(db *gorm.DB, jwtManager *JWTManager) *usersHandler {
	return &usersHandler{
		DB:         db,
		JWTManager: jwtManager,
	}
}

func (h *usersHandler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate request
	if len(req.FirstName) == 0 {
		return ErrMissingFirstName
	}
	if len(req.LastName) == 0 {
		return ErrMissingLastName
	}
	if !isValidEmail(req.Email) {
		return ErrInvalidEmail
	}
	if !isValidPassword(req.Password) {
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

		// create token
		token, err := h.JWTManager.Generate(user)
		if err != nil {
			logger.Errorf("Error gen token: %v", err)
			return ErrTokenGenerated
		}

		rsp.User = user.sanitize()
		rsp.Token = token

		return nil
	})
}

func (h *usersHandler) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	if len(req.Id) == 0 {
		return ErrMissingId
	}
	if req.FirstName != nil && len(req.FirstName.Value) == 0 {
		return ErrMissingFirstName
	}
	if req.LastName != nil && len(req.LastName.Value) == 0 {
		return ErrMissingLastName
	}
	if req.Email != nil && !isValidEmail(req.Email.Value) {
		return ErrInvalidEmail
	}
	if req.Password != nil && !isValidPassword(req.Password.Value) {
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

func (h *usersHandler) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
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
		token, err := h.JWTManager.Generate(&user)
		if err != nil {
			logger.Errorf("Error gen token: %v", err)
			return ErrTokenGenerated
		}

		rsp.User = user.sanitize()
		rsp.Token = token

		return nil
	})
}

func (h *usersHandler) Logout(ctx context.Context, req *pb.LogoutRequest, rsp *pb.LogoutResponse) error {
	if len(req.Id) == 0 {
		return ErrMissingId
	}
	return h.DB.Transaction(func(tx *gorm.DB) error {
		// lookup user
		// var user User
		// if err := tx.Where(&User{ID: req.Id}).Preload("Tokens").First(&user).Error; err == gorm.ErrRecordNotFound {
		// 	return ErrNotFound
		// } else if err != nil {
		// 	logger.Errorf("Error connecting from db: %v", err)
		// 	return ErrConnectDB
		// }

		return nil
	})
}

func (h *usersHandler) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	if len(req.Token) == 0 {
		return ErrMissingToken
	}
	return h.DB.Transaction(func(tx *gorm.DB) error {
		// verrify token
		claims, err := h.JWTManager.Verify(req.Token)
		if err != nil {
			return ErrTokenInvalid
		}

		rsp.User = claims.User.sanitize()
		return nil
	})
}

func (h *usersHandler) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	var us []User
	if err := h.DB.Order("created_at desc").Find(&us).Error; err != nil {
		logger.Errorf("Error connecting from db: %v", err)
		return ErrConnectDB
	}

	rsp.Users = make([]*pb.User, len(us))
	for i, user := range us {
		rsp.Users[i] = user.sanitize()
	}
	return nil
}

func (h *usersHandler) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
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
	rsp.Users = make(map[string]*pb.User, len(us))
	for _, user := range us {
		rsp.Users[user.ID] = user.sanitize()
	}
	return nil
}

func (h *usersHandler) ReadByEmail(ctx context.Context, req *pb.ReadByEmailRequest, rsp *pb.ReadByEmailResponse) error {
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
	rsp.Users = make(map[string]*pb.User, len(us))
	for _, user := range us {
		rsp.Users[user.Email] = user.sanitize()
	}
	return nil
}

func (h *usersHandler) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
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
