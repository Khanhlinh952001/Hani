package users

import (
	"errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUserService(user *User, password string) error {
	if _, err := repoGetUserByEmail(user.Email); err == nil {
		return errors.New("email already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	if user.Status == 0 {
		user.Status = 1
	}

	if err := repoCreateUser(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("email already taken")
		}
		return err
	}
	return nil
}

func GetAllUsersService() ([]User, error) {
	return repoGetAllUsers()
}

func GetUserByEmailService(email string) (*User, error) {
	return repoGetUserByEmail(email)
}

func AuthenticateService(email, password string) (*User, error) {
	user, err := repoGetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

func GetUserByIDService(id string) (*User, error) {
	parsed, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return repoGetUserByID(parsed)
}

func UpdateUserService(id string, data *User, password string) error {
	parsed, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	existing, err := repoGetUserByID(parsed)
	if err != nil {
		return err
	}

	if data.Name != "" {
		existing.Name = data.Name
	}
	if data.Email != "" {
		existing.Email = data.Email
	}
	if data.PhoneNumber != "" {
		existing.PhoneNumber = data.PhoneNumber
	}
	if data.Provider != "" {
		existing.Provider = data.Provider
	}
	if data.ProviderId != "" {
		existing.ProviderId = data.ProviderId
	}
	if data.Avatar != "" {
		existing.Avatar = data.Avatar
	}
	if data.Gender != "" {
		existing.Gender = data.Gender
	}
	if data.AiProfileID != nil && *data.AiProfileID != "" {
		existing.AiProfileID = data.AiProfileID
	}
	if data.SelectedCharacterID != "" {
		existing.SelectedCharacterID = data.SelectedCharacterID
	}
	if data.Level != 0 {
		existing.Level = data.Level
	}
	if data.Address != "" {
		existing.Address = data.Address
	}
	if data.Status != 0 {
		existing.Status = data.Status
	}
	if data.Role != 0 {
		existing.Role = data.Role
	}
	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		existing.Password = string(hashed)
	}

	return repoUpdateUser(parsed, existing)
}

func FindOrCreateClerkUser(clerkID, email, name string) (*User, error) {
	if existing, err := repoGetUserByProviderId("clerk", clerkID); err == nil {
		if name != "" && existing.Name != name {
			existing.Name = name
			_ = repoUpdateUser(existing.ID, existing)
		}
		return existing, nil
	}

	if email == "" {
		email = clerkID + "@clerk.hani.local"
	}
	if name == "" {
		name = "친구"
	}

	// Placeholder password — Clerk users never authenticate locally.
	placeholder, err := bcrypt.GenerateFromPassword([]byte(clerkID+"-clerk"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:       name,
		Email:      email,
		Password:   string(placeholder),
		Provider:   "clerk",
		ProviderId: clerkID,
		Status:     1,
	}
	if err := repoCreateUser(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return repoGetUserByProviderId("clerk", clerkID)
		}
		return nil, err
	}
	return user, nil
}

func SetAiProfileService(userID int, profileID string, companionGender string) error {
	existing, err := repoGetUserByID(userID)
	if err != nil {
		return err
	}
	id := profileID
	existing.AiProfileID = &id
	if companionGender == "male" || companionGender == "female" {
		// no-op on user gender; profile stores companion gender
	}
	return repoUpdateUser(userID, existing)
}

func SetSelectedCharacterService(userID int, characterID string) error {
	existing, err := repoGetUserByID(userID)
	if err != nil {
		return err
	}
	existing.SelectedCharacterID = characterID
	return repoUpdateUser(userID, existing)
}

func DeleteUserService(id string) error {
	parsed, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return repoDeleteUser(parsed)
}
