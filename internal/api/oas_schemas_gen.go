// Code generated by ogen, DO NOT EDIT.

package api

import (
	"fmt"
	"time"
)

func (s *ErrorResponseStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

// Contains a message.
// Ref: #/components/schemas/AcceptanceResponse
type AcceptanceResponse struct {
	Message string `json:"message"`
}

// GetMessage returns the value of Message.
func (s *AcceptanceResponse) GetMessage() string {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *AcceptanceResponse) SetMessage(val string) {
	s.Message = val
}

type Access struct {
	Token string
}

// GetToken returns the value of Token.
func (s *Access) GetToken() string {
	return s.Token
}

// SetToken sets the value of Token.
func (s *Access) SetToken(val string) {
	s.Token = val
}

// Contains an error.
// Ref: #/components/schemas/ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetError returns the value of Error.
func (s *ErrorResponse) GetError() string {
	return s.Error
}

// SetError sets the value of Error.
func (s *ErrorResponse) SetError(val string) {
	s.Error = val
}

// ErrorResponseStatusCode wraps ErrorResponse with StatusCode.
type ErrorResponseStatusCode struct {
	StatusCode int
	Response   ErrorResponse
}

// GetStatusCode returns the value of StatusCode.
func (s *ErrorResponseStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *ErrorResponseStatusCode) GetResponse() ErrorResponse {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *ErrorResponseStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *ErrorResponseStatusCode) SetResponse(val ErrorResponse) {
	s.Response = val
}

// Contains a message as well as optional properties.
// Ref: #/components/schemas/MessageRequest
type MessageRequest struct {
	Message string `json:"message"`
}

// GetMessage returns the value of Message.
func (s *MessageRequest) GetMessage() string {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *MessageRequest) SetMessage(val string) {
	s.Message = val
}

// Contains a message as well as optional properties.
// Ref: #/components/schemas/MessageResponse
type MessageResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
	Version int32  `json:"version"`
}

// GetID returns the value of ID.
func (s *MessageResponse) GetID() int64 {
	return s.ID
}

// GetMessage returns the value of Message.
func (s *MessageResponse) GetMessage() string {
	return s.Message
}

// GetVersion returns the value of Version.
func (s *MessageResponse) GetVersion() int32 {
	return s.Version
}

// SetID sets the value of ID.
func (s *MessageResponse) SetID(val int64) {
	s.ID = val
}

// SetMessage sets the value of Message.
func (s *MessageResponse) SetMessage(val string) {
	s.Message = val
}

// SetVersion sets the value of Version.
func (s *MessageResponse) SetVersion(val int32) {
	s.Version = val
}

// Contains metadata.
// Ref: #/components/schemas/MessagesMetadataResponse
type MessagesMetadataResponse struct {
	CurrentPage  int32 `json:"current_page"`
	FirstPage    int32 `json:"first_page"`
	LastPage     int32 `json:"last_page"`
	PageSize     int32 `json:"page_size"`
	TotalRecords int64 `json:"total_records"`
}

// GetCurrentPage returns the value of CurrentPage.
func (s *MessagesMetadataResponse) GetCurrentPage() int32 {
	return s.CurrentPage
}

// GetFirstPage returns the value of FirstPage.
func (s *MessagesMetadataResponse) GetFirstPage() int32 {
	return s.FirstPage
}

// GetLastPage returns the value of LastPage.
func (s *MessagesMetadataResponse) GetLastPage() int32 {
	return s.LastPage
}

// GetPageSize returns the value of PageSize.
func (s *MessagesMetadataResponse) GetPageSize() int32 {
	return s.PageSize
}

// GetTotalRecords returns the value of TotalRecords.
func (s *MessagesMetadataResponse) GetTotalRecords() int64 {
	return s.TotalRecords
}

// SetCurrentPage sets the value of CurrentPage.
func (s *MessagesMetadataResponse) SetCurrentPage(val int32) {
	s.CurrentPage = val
}

// SetFirstPage sets the value of FirstPage.
func (s *MessagesMetadataResponse) SetFirstPage(val int32) {
	s.FirstPage = val
}

// SetLastPage sets the value of LastPage.
func (s *MessagesMetadataResponse) SetLastPage(val int32) {
	s.LastPage = val
}

// SetPageSize sets the value of PageSize.
func (s *MessagesMetadataResponse) SetPageSize(val int32) {
	s.PageSize = val
}

// SetTotalRecords sets the value of TotalRecords.
func (s *MessagesMetadataResponse) SetTotalRecords(val int64) {
	s.TotalRecords = val
}

// Contains messages and metadata objects.
// Ref: #/components/schemas/MessagesResponse
type MessagesResponse struct {
	Messages []MessageResponse        `json:"messages"`
	Metadata MessagesMetadataResponse `json:"metadata"`
}

// GetMessages returns the value of Messages.
func (s *MessagesResponse) GetMessages() []MessageResponse {
	return s.Messages
}

// GetMetadata returns the value of Metadata.
func (s *MessagesResponse) GetMetadata() MessagesMetadataResponse {
	return s.Metadata
}

// SetMessages sets the value of Messages.
func (s *MessagesResponse) SetMessages(val []MessageResponse) {
	s.Messages = val
}

// SetMetadata sets the value of Metadata.
func (s *MessagesResponse) SetMetadata(val MessagesMetadataResponse) {
	s.Metadata = val
}

// NewOptInt32 returns new OptInt32 with value set to v.
func NewOptInt32(v int32) OptInt32 {
	return OptInt32{
		Value: v,
		Set:   true,
	}
}

// OptInt32 is optional int32.
type OptInt32 struct {
	Value int32
	Set   bool
}

// IsSet returns true if OptInt32 was set.
func (o OptInt32) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptInt32) Reset() {
	var v int32
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptInt32) SetTo(v int32) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptInt32) Get() (v int32, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptInt32) Or(d int32) int32 {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

type Refresh struct {
	APIKey string
}

// GetAPIKey returns the value of APIKey.
func (s *Refresh) GetAPIKey() string {
	return s.APIKey
}

// SetAPIKey sets the value of APIKey.
func (s *Refresh) SetAPIKey(val string) {
	s.APIKey = val
}

// Contains a plaintext token as well as optional properties.
// Ref: #/components/schemas/TokenRequest
type TokenRequest struct {
	Token string `json:"token"`
}

// GetToken returns the value of Token.
func (s *TokenRequest) GetToken() string {
	return s.Token
}

// SetToken sets the value of Token.
func (s *TokenRequest) SetToken(val string) {
	s.Token = val
}

// Contains a plaintext token as well as optional properties.
// Ref: #/components/schemas/TokenResponse
type TokenResponse struct {
	Scope  string    `json:"scope"`
	Expiry time.Time `json:"expiry"`
	Token  string    `json:"token"`
}

// GetScope returns the value of Scope.
func (s *TokenResponse) GetScope() string {
	return s.Scope
}

// GetExpiry returns the value of Expiry.
func (s *TokenResponse) GetExpiry() time.Time {
	return s.Expiry
}

// GetToken returns the value of Token.
func (s *TokenResponse) GetToken() string {
	return s.Token
}

// SetScope sets the value of Scope.
func (s *TokenResponse) SetScope(val string) {
	s.Scope = val
}

// SetExpiry sets the value of Expiry.
func (s *TokenResponse) SetExpiry(val time.Time) {
	s.Expiry = val
}

// SetToken sets the value of Token.
func (s *TokenResponse) SetToken(val string) {
	s.Token = val
}

// TokenResponseHeaders wraps TokenResponse with response headers.
type TokenResponseHeaders struct {
	SetCookie OptString
	Response  TokenResponse
}

// GetSetCookie returns the value of SetCookie.
func (s *TokenResponseHeaders) GetSetCookie() OptString {
	return s.SetCookie
}

// GetResponse returns the value of Response.
func (s *TokenResponseHeaders) GetResponse() TokenResponse {
	return s.Response
}

// SetSetCookie sets the value of SetCookie.
func (s *TokenResponseHeaders) SetSetCookie(val OptString) {
	s.SetCookie = val
}

// SetResponse sets the value of Response.
func (s *TokenResponseHeaders) SetResponse(val TokenResponse) {
	s.Response = val
}

// Contains a password and token object.
// Ref: #/components/schemas/UpdateUserPasswordRequest
type UpdateUserPasswordRequest struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

// GetPassword returns the value of Password.
func (s *UpdateUserPasswordRequest) GetPassword() string {
	return s.Password
}

// GetToken returns the value of Token.
func (s *UpdateUserPasswordRequest) GetToken() string {
	return s.Token
}

// SetPassword sets the value of Password.
func (s *UpdateUserPasswordRequest) SetPassword(val string) {
	s.Password = val
}

// SetToken sets the value of Token.
func (s *UpdateUserPasswordRequest) SetToken(val string) {
	s.Token = val
}

// Contains an email address.
// Ref: #/components/schemas/UserEmailRequest
type UserEmailRequest struct {
	Email string `json:"email"`
}

// GetEmail returns the value of Email.
func (s *UserEmailRequest) GetEmail() string {
	return s.Email
}

// SetEmail sets the value of Email.
func (s *UserEmailRequest) SetEmail(val string) {
	s.Email = val
}

// Contains an email address and password.
// Ref: #/components/schemas/UserLoginRequest
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetEmail returns the value of Email.
func (s *UserLoginRequest) GetEmail() string {
	return s.Email
}

// GetPassword returns the value of Password.
func (s *UserLoginRequest) GetPassword() string {
	return s.Password
}

// SetEmail sets the value of Email.
func (s *UserLoginRequest) SetEmail(val string) {
	s.Email = val
}

// SetPassword sets the value of Password.
func (s *UserLoginRequest) SetPassword(val string) {
	s.Password = val
}

// Contains a username, email and password.
// Ref: #/components/schemas/UserRequest
type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetName returns the value of Name.
func (s *UserRequest) GetName() string {
	return s.Name
}

// GetEmail returns the value of Email.
func (s *UserRequest) GetEmail() string {
	return s.Email
}

// GetPassword returns the value of Password.
func (s *UserRequest) GetPassword() string {
	return s.Password
}

// SetName sets the value of Name.
func (s *UserRequest) SetName(val string) {
	s.Name = val
}

// SetEmail sets the value of Email.
func (s *UserRequest) SetEmail(val string) {
	s.Email = val
}

// SetPassword sets the value of Password.
func (s *UserRequest) SetPassword(val string) {
	s.Password = val
}

// Contains a username, email and password.
// Ref: #/components/schemas/UserResponse
type UserResponse struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Version int32  `json:"version"`
}

// GetName returns the value of Name.
func (s *UserResponse) GetName() string {
	return s.Name
}

// GetEmail returns the value of Email.
func (s *UserResponse) GetEmail() string {
	return s.Email
}

// GetVersion returns the value of Version.
func (s *UserResponse) GetVersion() int32 {
	return s.Version
}

// SetName sets the value of Name.
func (s *UserResponse) SetName(val string) {
	s.Name = val
}

// SetEmail sets the value of Email.
func (s *UserResponse) SetEmail(val string) {
	s.Email = val
}

// SetVersion sets the value of Version.
func (s *UserResponse) SetVersion(val int32) {
	s.Version = val
}
