# github.com/seanflannery10/core

Routes for core API

## Routes

<details>
<summary>`/debug/vars`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
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
- [Authenticate.func1]()
- **/v1/messages**
	- [RequireAuthenticatedUser]()
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
- [Authenticate.func1]()
- **/v1/messages**
	- [RequireAuthenticatedUser]()
	- **/{id}**
		- **/**
			- _PATCH_
				- [UpdateMessageHandler.func1]()
			- _DELETE_
				- [DeleteMessageHandler.func1]()
			- _GET_
				- [GetMessageHandler.func1]()

</details>
<details>
<summary>`/v1/tokens/activation`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [Authenticate.func1]()
- **/v1/tokens**
	- **/activation**
		- _PUT_
			- [CreateTokenActivationHandler.func1]()

</details>
<details>
<summary>`/v1/tokens/authentication`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [Authenticate.func1]()
- **/v1/tokens**
	- **/authentication**
		- _POST_
			- [CreateTokenAuthHandler.func1]()

</details>
<details>
<summary>`/v1/tokens/password-reset`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [Authenticate.func1]()
- **/v1/tokens**
	- **/password-reset**
		- _PUT_
			- [CreateTokenPasswordResetHandler.func1]()

</details>
<details>
<summary>`/v1/users/activate`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
- [Authenticate.func1]()
- **/v1/users**
	- **/activate**
		- _PUT_
			- [ActivateUserHandler.func1]()

</details>
<details>
<summary>`/v1/users/register`</summary>

- [StartSpan.func1]()
- [Metrics]()
- [RecoverPanic]()
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
- [Authenticate.func1]()
- **/v1/users**
	- **/update-password**
		- _PUT_
			- [UpdateUserPasswordHandler.func1]()

</details>

Total # of routes: 9
