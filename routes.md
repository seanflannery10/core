# github.com/seanflannery10/core

Routes for core API

## Routes

<details>
<summary>`/debug/vars`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [o-chi/cors.(*Cors).Handler-fm]()
- [Authenticate.func1]()
- **/debug/vars**
	- _GET_
		- [ttp.Handler.ServeHTTP-fm]()

</details>
<details>
<summary>`/v1/messages`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [o-chi/cors.(*Cors).Handler-fm]()
- [Authenticate.func1]()
- **/v1/messages**
	- [RequireAuthenticatedUser.func1]()
	- **/**
		- _GET_
			- [GetMessagesUserHandler.func1]()
		- _POST_
			- [CreateMessageHandler.func1]()

</details>
<details>
<summary>`/v1/messages/{id}`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [o-chi/cors.(*Cors).Handler-fm]()
- [Authenticate.func1]()
- **/v1/messages**
	- [RequireAuthenticatedUser.func1]()
	- **/{id}**
		- **/**
			- _GET_
				- [GetMessageHandler.func1]()
			- _PUT_
				- [UpdateMessageHandler.func1]()
			- _DELETE_
				- [DeleteMessageHandler.func1]()

</details>
<details>
<summary>`/v1/tokens/access`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [o-chi/cors.(*Cors).Handler-fm]()
- [Authenticate.func1]()
- **/v1/tokens**
	- **/access**
		- _POST_
			- [CreateTokenAccessHandler.func1]()

</details>
<details>
<summary>`/v1/tokens/activation`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [o-chi/cors.(*Cors).Handler-fm]()
- [Authenticate.func1]()
- **/v1/tokens**
	- **/activation**
		- _POST_
			- [CreateTokenActivationHandler.func1]()

</details>
<details>
<summary>`/v1/tokens/password-reset`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [o-chi/cors.(*Cors).Handler-fm]()
- [Authenticate.func1]()
- **/v1/tokens**
	- **/password-reset**
		- _POST_
			- [CreateTokenPasswordResetHandler.func1]()

</details>
<details>
<summary>`/v1/tokens/refresh`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [o-chi/cors.(*Cors).Handler-fm]()
- [Authenticate.func1]()
- **/v1/tokens**
	- **/refresh**
		- _POST_
			- [CreateTokenRefreshHandler.func1]()

</details>
<details>
<summary>`/v1/users/activate`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [o-chi/cors.(*Cors).Handler-fm]()
- [Authenticate.func1]()
- **/v1/users**
	- **/activate**
		- _PATCH_
			- [ActivateUserHandler.func1]()

</details>
<details>
<summary>`/v1/users/register`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [o-chi/cors.(*Cors).Handler-fm]()
- [Authenticate.func1]()
- **/v1/users**
	- **/register**
		- _POST_
			- [CreateUserHandler.func1]()

</details>
<details>
<summary>`/v1/users/update-password`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [o-chi/cors.(*Cors).Handler-fm]()
- [Authenticate.func1]()
- **/v1/users**
	- **/update-password**
		- _PATCH_
			- [UpdateUserPasswordHandler.func1]()

</details>

Total # of routes: 10
