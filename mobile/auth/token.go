package auth

var token string

func SetToken(value string) {
	token = value
}

func GetToken() string {
	return token
}

func ClearToken() {
	token = ""
}

func IsLoggedIn() bool {
	return token != ""
}
