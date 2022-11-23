package utilities

// This funcion is used to check whether the logged user is acting on his own account or not
func CheckSelfAction(loggedUser string, user string) bool {
	if loggedUser == user {
		return true
	}
	return false
}
