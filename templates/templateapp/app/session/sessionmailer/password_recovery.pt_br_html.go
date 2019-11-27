package sessionmailer

var PasswordRecoveryPtBrPath = []string{"app", "session", "mailer", "password_recovery.pt-BR.html"}

var PasswordRecoveryPtBrContent = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	</head>
	<body>
		<p>Olá {{ .UserFirstName }}!</p>

		<p>Alguém solicitou o link para alterar sua senha. Você pode fazer isso através do link abaixo.</p>

		<p><a href="{{ .LinkToPasswordRecovery }}">Alterar senha</a></p>

		<p>Se não foi você que solicitou, por favor, apenas ignore este email.</p>
		<p>Sua senha não será alterada até você acessar o link acima e criar uma nova.</p>
	</body>
</html>`
