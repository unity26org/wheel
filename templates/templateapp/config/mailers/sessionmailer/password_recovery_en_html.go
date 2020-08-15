package sessionmailer

var PasswordRecoveryEnPath = []string{"config", "mailers", "session", "password_recovery.en.html"}

var PasswordRecoveryEnContent = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	</head>
	<body>
		<p>Hello {{ .UserFirstName }}!</p>

		<p>Someone has requested a link to change your password. You can do this through the link below.</p>

		<p><a href="{{ .LinkToPasswordRecovery }}">Change my password</a></p>

		<p>If you didn't request this, please ignore this email.</p>
		<p>Your password won't change until you access the link above and create a new one.</p>
	</body>
</html>`
