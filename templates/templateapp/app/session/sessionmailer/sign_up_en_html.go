package sessionmailer

var SignUpEnPath = []string{"app", "session", "mailer", "sign_up.en.html"}

var SignUpEnContent = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	</head>
	<body>
		<p>Welcome {{ .UserFirstName }},</p>

		<p>Thank you for registering for <strong>{{ .AppName }}</strong>!</p>

		<p>Anything you need from us, please let us know.</p>

		<p>Best,</p>

		<p>Name, Job Position<br>email@example.com</p>
	</body>
</html>`