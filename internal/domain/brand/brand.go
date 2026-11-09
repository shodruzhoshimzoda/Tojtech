package brand

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	MaxNameLength        = 100
	MaxSlugLength        = 100
	MaxDescriptionLength = 500
	// A logo may be sent as a data URL (the admin panel can upload a file).
	// It is stored as text and returned in every brand list, so keep it small.
	MaxLogoDataURLLength = 700_000
)

var (
	ErrEmptyBrandName     = errors.New("brand name cannot be empty")
	ErrBrandNameTooLong   = errors.New("brand name must be at most 100 characters")
	ErrEmptySlug          = errors.New("brand slug cannot be empty")
	ErrInvalidSlug        = errors.New("brand slug may contain only lowercase Latin letters, digits and single hyphens (for example: asus-rog)")
	ErrLongDescription    = errors.New("brand description must be at most 500 characters")
	ErrInvalidLogoURL     = errors.New("logo must be a link starting with http:// or https:// or an uploaded image")
	ErrLogoTooLarge       = errors.New("logo image is too large, choose a smaller one")
	ErrBrandNotFound      = errors.New("brand not found")
	ErrBrandAlreadyExists = errors.New("a brand with this name or slug already exists")
)

// Brand - HP, ASUS, MSI ... It is not tied to one category: the same brand can
// have graphics cards, RAM and laptops. A product points to its brand (optional).
type Brand struct {
	ID          int64     `db:"id" json:"-"`
	UUID        uuid.UUID `db:"uuid" json:"uuid"`
	Name        string    `db:"name" json:"name"`
	Slug        string    `db:"slug" json:"slug"`
	LogoURL     string    `db:"logo_url" json:"logo_url"`
	Description string    `db:"description" json:"description"`
	// ProductCount is filled only by list queries, it is not a column of the brands table.
	ProductCount int       `db:"-" json:"product_count"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Normalize cleans the input: removes extra spaces, makes the slug lowercase.
// If the slug is empty it is made from the name ("ASUS ROG" -> "asus-rog").
func (b *Brand) Normalize() {
	b.Name = strings.TrimSpace(b.Name)
	b.Slug = strings.ToLower(strings.TrimSpace(b.Slug))
	b.LogoURL = strings.TrimSpace(b.LogoURL)
	b.Description = strings.TrimSpace(b.Description)

	if b.Slug == "" {
		b.Slug = Slugify(b.Name)
	}
}

// Validate checks the brand. Call Normalize first.
func (b *Brand) Validate() error {
	if b.Name == "" {
		return ErrEmptyBrandName
	}
	if utf8.RuneCountInString(b.Name) > MaxNameLength {
		return ErrBrandNameTooLong
	}
	if b.Slug == "" {
		return ErrEmptySlug
	}
	if len(b.Slug) > MaxSlugLength || !slugPattern.MatchString(b.Slug) {
		return ErrInvalidSlug
	}
	if utf8.RuneCountInString(b.Description) > MaxDescriptionLength {
		return ErrLongDescription
	}
	return validateLogoURL(b.LogoURL)
}

// The logo is optional: an empty value means "no logo".
func validateLogoURL(raw string) error {
	if raw == "" {
		return nil
	}
	if strings.HasPrefix(raw, "data:image/") {
		if len(raw) > MaxLogoDataURLLength {
			return ErrLogoTooLarge
		}
		return nil
	}
	parsed, err := url.ParseRequestURI(raw)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return ErrInvalidLogoURL
	}
	return nil
}

// Slugify turns a name into a url-friendly slug: lowercase Latin letters and digits
// separated by single hyphens. Other characters (spaces, symbols, Cyrillic...) become separators.
// It may return "" (for a name with no Latin letters or digits) - then the caller must ask for a slug.
func Slugify(name string) string {
	var sb strings.Builder
	pendingHyphen := false

	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			if pendingHyphen && sb.Len() > 0 {
				sb.WriteByte('-')
			}
			pendingHyphen = false
			sb.WriteRune(r)
		} else {
			pendingHyphen = true
		}
	}
	return sb.String()
}