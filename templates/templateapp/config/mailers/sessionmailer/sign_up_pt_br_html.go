package sessionmailer

var SignUpPtBrPath = []string{"config", "mailers", "session", "sign_up.pt-BR.html"}

var SignUpPtBrContent = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	</head>
	<body>
		<p>Bem-vindo {{ .UserFirstName }},</p>

		<p>Obrigado por se cadastrar em <strong>{{ .AppName }}</strong>!</p>

		<p>Qualquer coisa que você precisar, por favor, entre em contato conosco.</p>

		<p>Atenciosamente,</p>

		<p>Nome, Cargo<br>email@example.com</p>
	</body>
</html>`
