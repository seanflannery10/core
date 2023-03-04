# github.com/seanflannery10/core

Routes for core API

## Routes

<details>
<summary>`/debug/vars`</summary>

- [Metrics]()
- [RecoverPanic]()
- [Middleware.func1]()
- [SetQueriesCtx.func1]()
- [SetMailerCtx.func1]()
- [Authenticate]()
- **/debug/vars**
	- _GET_
		- [ttp.Handler.ServeHTTP-fm]()

</details>
<details>
<summary>`/v1/messages`</summary>

- [Metrics]()
- [RecoverPanic]()
- [Middleware.func1]()
- [SetQueriesCtx.func1]()
- [SetMailerCtx.func1]()
- [Authenticate]()
- **/v1/messages**
	- [RequireAuthenticatedUser]()
	- **/**
		- _GET_
			- [GetMessagesUserHandler]()
		- _POST_
			- [CreateMessageHandler]()

</details>
<details>
<summary>`/v1/messages/{id}`</summary>

- [Metrics]()
- [RecoverPanic]()
- [Middleware.func1]()
- [SetQueriesCtx.func1]()
- [SetMailerCtx.func1]()
- [Authenticate]()
- **/v1/messages**
	- [RequireAuthenticatedUser]()
	- **/{id}**
		- **/**
			- _GET_
				- [GetMessageHandler]()
			- _PATCH_
				- [UpdateMessageHandler]()
			- _DELETE_
				- [DeleteMessageHandler]()

</details>
<details>
<summary>`/v1/tokens/activation`</summary>

- [Metrics]()
- [RecoverPanic]()
- [Middleware.func1]()
- [SetQueriesCtx.func1]()
- [SetMailerCtx.func1]()
- [Authenticate]()
- **/v1/tokens**
	- **/activation**
		- _PUT_
			- [CreateTokenActivationHandler]()

</details>
<details>
<summary>`/v1/tokens/authentication`</summary>

- [Metrics]()
- [RecoverPanic]()
- [Middleware.func1]()
- [SetQueriesCtx.func1]()
- [SetMailerCtx.func1]()
- [Authenticate]()
- **/v1/tokens**
	- **/authentication**
		- _POST_
			- [CreateTokenAuthHandler]()

</details>
<details>
<summary>`/v1/tokens/password-reset`</summary>

- [Metrics]()
- [RecoverPanic]()
- [Middleware.func1]()
- [SetQueriesCtx.func1]()
- [SetMailerCtx.func1]()
- [Authenticate]()
- **/v1/tokens**
	- **/password-reset**
		- _PUT_
			- [CreateTokenPasswordResetHandler]()

</details>
<details>
<summary>`/v1/users/activate`</summary>

- [Metrics]()
- [RecoverPanic]()
- [Middleware.func1]()
- [SetQueriesCtx.func1]()
- [SetMailerCtx.func1]()
- [Authenticate]()
- **/v1/users**
	- **/activate**
		- _PUT_
			- [ActivateUserHandler]()

</details>
<details>
<summary>`/v1/users/register`</summary>

- [Metrics]()
- [RecoverPanic]()
- [Middleware.func1]()
- [SetQueriesCtx.func1]()
- [SetMailerCtx.func1]()
- [Authenticate]()
- **/v1/users**
	- **/register**
		- _POST_
			- [CreateUserHandler]()

</details>
<details>
<summary>`/v1/users/update-password`</summary>

- [Metrics]()
- [RecoverPanic]()
- [Middleware.func1]()
- [SetQueriesCtx.func1]()
- [SetMailerCtx.func1]()
- [Authenticate]()
- **/v1/users**
	- **/update-password**
		- _PUT_
			- [UpdateUserPasswordHandler]()

</details>

Total # of routes: 9
