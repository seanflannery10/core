// Code generated by ogen, DO NOT EDIT.

package api

import (
	"time"
)

// SetFake set fake values.
func (s *AcceptanceResponse) SetFake() {
	{
		{
			s.Message = "string"
		}
	}
}

// SetFake set fake values.
func (s *ErrorResponse) SetFake() {
	{
		{
			s.Error = "string"
		}
	}
}

// SetFake set fake values.
func (s *MessageRequest) SetFake() {
	{
		{
			s.Message = "string"
		}
	}
}

// SetFake set fake values.
func (s *MessageResponse) SetFake() {
	{
		{
			s.ID = int64(0)
		}
	}
	{
		{
			s.Message = "string"
		}
	}
	{
		{
			s.Version = int32(0)
		}
	}
}

// SetFake set fake values.
func (s *MessagesMetadataResponse) SetFake() {
	{
		{
			s.CurrentPage = int32(0)
		}
	}
	{
		{
			s.FirstPage = int32(0)
		}
	}
	{
		{
			s.LastPage = int32(0)
		}
	}
	{
		{
			s.PageSize = int32(0)
		}
	}
	{
		{
			s.TotalRecords = int64(0)
		}
	}
}

// SetFake set fake values.
func (s *MessagesResponse) SetFake() {
	{
		{
			s.Messages = nil
			for i := 0; i < 0; i++ {
				var elem MessageResponse
				{
					elem.SetFake()
				}
				s.Messages = append(s.Messages, elem)
			}
		}
	}
	{
		{
			s.Metadata.SetFake()
		}
	}
}

// SetFake set fake values.
func (s *TokenRequest) SetFake() {
	{
		{
			s.Token = "string"
		}
	}
}

// SetFake set fake values.
func (s *TokenResponse) SetFake() {
	{
		{
			s.Scope = "string"
		}
	}
	{
		{
			s.Expiry = time.Now()
		}
	}
	{
		{
			s.Token = "string"
		}
	}
}

// SetFake set fake values.
func (s *UpdateUserPasswordRequest) SetFake() {
	{
		{
			s.Password = "string"
		}
	}
	{
		{
			s.Token = "string"
		}
	}
}

// SetFake set fake values.
func (s *UserEmailRequest) SetFake() {
	{
		{
			s.Email = "string"
		}
	}
}

// SetFake set fake values.
func (s *UserLoginRequest) SetFake() {
	{
		{
			s.Email = "string"
		}
	}
	{
		{
			s.Password = "string"
		}
	}
}

// SetFake set fake values.
func (s *UserRequest) SetFake() {
	{
		{
			s.Name = "string"
		}
	}
	{
		{
			s.Email = "string"
		}
	}
	{
		{
			s.Password = "string"
		}
	}
}

// SetFake set fake values.
func (s *UserResponse) SetFake() {
	{
		{
			s.Name = "string"
		}
	}
	{
		{
			s.Email = "string"
		}
	}
	{
		{
			s.Version = int32(0)
		}
	}
}
