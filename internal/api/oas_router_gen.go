// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ogen-go/ogen/uri"
)

// ServeHTTP serves http request as defined by OpenAPI v3 specification,
// calling handler that matches the path or returning not found error.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	elem := r.URL.Path
	elemIsEscaped := false
	if rawPath := r.URL.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
			elemIsEscaped = strings.ContainsRune(elem, '%')
		}
	}
	if prefix := s.cfg.Prefix; len(prefix) > 0 {
		if strings.HasPrefix(elem, prefix) {
			// Cut prefix from the path.
			elem = strings.TrimPrefix(elem, prefix)
		} else {
			// Prefix doesn't match.
			s.notFound(w, r)
			return
		}
	}
	if len(elem) == 0 {
		s.notFound(w, r)
		return
	}
	args := [1]string{}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/v1/"
			if l := len("/v1/"); len(elem) >= l && elem[0:l] == "/v1/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'm': // Prefix: "messages"
				if l := len("messages"); len(elem) >= l && elem[0:l] == "messages" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch r.Method {
					case "GET":
						s.handleGetUserMessagesRequest([0]string{}, elemIsEscaped, w, r)
					case "POST":
						s.handleNewMessageRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "GET,POST")
					}

					return
				}
				switch elem[0] {
				case '/': // Prefix: "/"
					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "id"
					// Leaf parameter
					args[0] = elem
					elem = ""

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "DELETE":
							s.handleDeleteMessageRequest([1]string{
								args[0],
							}, elemIsEscaped, w, r)
						case "GET":
							s.handleGetMessageRequest([1]string{
								args[0],
							}, elemIsEscaped, w, r)
						case "PUT":
							s.handleUpdateMessageRequest([1]string{
								args[0],
							}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "DELETE,GET,PUT")
						}

						return
					}
				}
			case 't': // Prefix: "tokens/"
				if l := len("tokens/"); len(elem) >= l && elem[0:l] == "tokens/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'a': // Prefix: "ac"
					if l := len("ac"); len(elem) >= l && elem[0:l] == "ac" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'c': // Prefix: "cess"
						if l := len("cess"); len(elem) >= l && elem[0:l] == "cess" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleNewAccessTokenRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}
					case 't': // Prefix: "tivation"
						if l := len("tivation"); len(elem) >= l && elem[0:l] == "tivation" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "POST":
								s.handleNewActivationTokenRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}
					}
				case 'p': // Prefix: "password-reset"
					if l := len("password-reset"); len(elem) >= l && elem[0:l] == "password-reset" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleNewPasswordResetTokenRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}
				case 'r': // Prefix: "refresh"
					if l := len("refresh"); len(elem) >= l && elem[0:l] == "refresh" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleNewRefreshTokenRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}
				}
			case 'u': // Prefix: "users/"
				if l := len("users/"); len(elem) >= l && elem[0:l] == "users/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'a': // Prefix: "activate"
					if l := len("activate"); len(elem) >= l && elem[0:l] == "activate" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "PATCH":
							s.handleActivateUserRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "PATCH")
						}

						return
					}
				case 'r': // Prefix: "register"
					if l := len("register"); len(elem) >= l && elem[0:l] == "register" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "POST":
							s.handleNewUserRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "POST")
						}

						return
					}
				case 'u': // Prefix: "update-password"
					if l := len("update-password"); len(elem) >= l && elem[0:l] == "update-password" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "PATCH":
							s.handleUpdateUserPasswordRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "PATCH")
						}

						return
					}
				}
			}
		}
	}
	s.notFound(w, r)
}

// Route is route object.
type Route struct {
	name        string
	operationID string
	pathPattern string
	count       int
	args        [1]string
}

// Name returns ogen operation name.
//
// It is guaranteed to be unique and not empty.
func (r Route) Name() string {
	return r.name
}

// OperationID returns OpenAPI operationId.
func (r Route) OperationID() string {
	return r.operationID
}

// PathPattern returns OpenAPI path.
func (r Route) PathPattern() string {
	return r.pathPattern
}

// Args returns parsed arguments.
func (r Route) Args() []string {
	return r.args[:r.count]
}

// FindRoute finds Route for given method and path.
//
// Note: this method does not unescape path or handle reserved characters in path properly. Use FindPath instead.
func (s *Server) FindRoute(method, path string) (Route, bool) {
	return s.FindPath(method, &url.URL{Path: path})
}

// FindPath finds Route for given method and URL.
func (s *Server) FindPath(method string, u *url.URL) (r Route, _ bool) {
	var (
		elem = u.Path
		args = r.args
	)
	if rawPath := u.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
		}
		defer func() {
			for i, arg := range r.args[:r.count] {
				if unescaped, err := url.PathUnescape(arg); err == nil {
					r.args[i] = unescaped
				}
			}
		}()
	}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/v1/"
			if l := len("/v1/"); len(elem) >= l && elem[0:l] == "/v1/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'm': // Prefix: "messages"
				if l := len("messages"); len(elem) >= l && elem[0:l] == "messages" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch method {
					case "GET":
						r.name = "GetUserMessages"
						r.operationID = "GetUserMessages"
						r.pathPattern = "/v1/messages"
						r.args = args
						r.count = 0
						return r, true
					case "POST":
						r.name = "NewMessage"
						r.operationID = "NewMessage"
						r.pathPattern = "/v1/messages"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}
				switch elem[0] {
				case '/': // Prefix: "/"
					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "id"
					// Leaf parameter
					args[0] = elem
					elem = ""

					if len(elem) == 0 {
						switch method {
						case "DELETE":
							// Leaf: DeleteMessage
							r.name = "DeleteMessage"
							r.operationID = "DeleteMessage"
							r.pathPattern = "/v1/messages/{id}"
							r.args = args
							r.count = 1
							return r, true
						case "GET":
							// Leaf: GetMessage
							r.name = "GetMessage"
							r.operationID = "GetMessage"
							r.pathPattern = "/v1/messages/{id}"
							r.args = args
							r.count = 1
							return r, true
						case "PUT":
							// Leaf: UpdateMessage
							r.name = "UpdateMessage"
							r.operationID = "UpdateMessage"
							r.pathPattern = "/v1/messages/{id}"
							r.args = args
							r.count = 1
							return r, true
						default:
							return
						}
					}
				}
			case 't': // Prefix: "tokens/"
				if l := len("tokens/"); len(elem) >= l && elem[0:l] == "tokens/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'a': // Prefix: "ac"
					if l := len("ac"); len(elem) >= l && elem[0:l] == "ac" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'c': // Prefix: "cess"
						if l := len("cess"); len(elem) >= l && elem[0:l] == "cess" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "POST":
								// Leaf: NewAccessToken
								r.name = "NewAccessToken"
								r.operationID = "NewAccessToken"
								r.pathPattern = "/v1/tokens/access"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
					case 't': // Prefix: "tivation"
						if l := len("tivation"); len(elem) >= l && elem[0:l] == "tivation" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "POST":
								// Leaf: NewActivationToken
								r.name = "NewActivationToken"
								r.operationID = "NewActivationToken"
								r.pathPattern = "/v1/tokens/activation"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
					}
				case 'p': // Prefix: "password-reset"
					if l := len("password-reset"); len(elem) >= l && elem[0:l] == "password-reset" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: NewPasswordResetToken
							r.name = "NewPasswordResetToken"
							r.operationID = "NewPasswordResetToken"
							r.pathPattern = "/v1/tokens/password-reset"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				case 'r': // Prefix: "refresh"
					if l := len("refresh"); len(elem) >= l && elem[0:l] == "refresh" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: NewRefreshToken
							r.name = "NewRefreshToken"
							r.operationID = "NewRefreshToken"
							r.pathPattern = "/v1/tokens/refresh"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				}
			case 'u': // Prefix: "users/"
				if l := len("users/"); len(elem) >= l && elem[0:l] == "users/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'a': // Prefix: "activate"
					if l := len("activate"); len(elem) >= l && elem[0:l] == "activate" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "PATCH":
							// Leaf: ActivateUser
							r.name = "ActivateUser"
							r.operationID = "ActivateUser"
							r.pathPattern = "/v1/users/activate"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				case 'r': // Prefix: "register"
					if l := len("register"); len(elem) >= l && elem[0:l] == "register" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "POST":
							// Leaf: NewUser
							r.name = "NewUser"
							r.operationID = "NewUser"
							r.pathPattern = "/v1/users/register"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				case 'u': // Prefix: "update-password"
					if l := len("update-password"); len(elem) >= l && elem[0:l] == "update-password" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "PATCH":
							// Leaf: UpdateUserPassword
							r.name = "UpdateUserPassword"
							r.operationID = "UpdateUserPassword"
							r.pathPattern = "/v1/users/update-password"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				}
			}
		}
	}
	return r, false
}
