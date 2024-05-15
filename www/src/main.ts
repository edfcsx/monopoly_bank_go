import { __Connection } from "./connection"
import { CommandsRequest } from "./types"
import { Worker } from "./worker"

__Connection()

// Check if the user is in bank page
if (window.location.pathname === '/bank') {
	Worker.Make((w: any) => {
		if (__Connection().isOpen()) {
			Worker.Clear(w)
			__Connection().send({ command: CommandsRequest.SendProfile })
		}
	})
}

function getLoginButton() {
	return document.getElementById('login-dispatch')
}

document.getElementById('login_form')?.addEventListener('submit', (e) => {
	e.preventDefault()
	const formData = new FormData(e.target as HTMLFormElement)

	const data = {
		username: formData.get('username'),
		password: formData.get('password')
	}

	__Connection().send({
		command: CommandsRequest.Authenticate,
		args: { ...data },
		...data,
	})

	getLoginButton()?.setAttribute('disabled', 'true')
})

// abstract class Commands {
// 	abstract execute (serverMessage: NetworkingMessage): void;
// }
//
// class AuthenticateSuccess extends Commands {
// 	public execute (serverMessage: NetworkingMessage) {
// 		console.log('MENSAGEM RECEBIDA DE AUTH:>', serverMessage)
// 		get_context().pop_up.fire(
// 			'Monopoly Bank',
// 			`Bem vindo, ${serverMessage.args?.username}!`,
// 			'success',
// 			3000,
// 		)
//
// 		if (serverMessage.args) {
// 			sessionStorage.setItem('auth', JSON.stringify({
// 				username: serverMessage.args!.username,
// 				password: serverMessage.args!.password,
// 				token: serverMessage.player_hash
// 			}))
//
// 			getLoginButton()?.removeAttribute('disabled')
//
// 			if (window.location.pathname != '/bank') {
// 				setTimeout(() => {
// 					window.location.href = '/bank'
// 				}, 1000)
// 			}
// 		} else {
// 			get_context().pop_up.fire(
// 				'Monopoly Bank',
// 				'Ocorreu um erro no sistema, por favor contate o desenvolvedor',
// 				'error',
// 				5000,
// 			)
// 		}
// 	}
// }
//
// class AuthenticateFailed extends Commands {
// 	public execute (serverMessage: NetworkingMessage) {
// 		get_context().pop_up.fire(
// 			'Monopoly Bank',
// 			'Usuário ou senha incorretos',
// 			'error',
// 			5000,
// 		)
//
// 		sessionStorage.removeItem('auth')
// 		getLoginButton()?.removeAttribute('disabled')
// 	}
// }
//
// class ProfileCommand extends Commands {
// 	public execute (serverMessage: NetworkingMessage) {
// 		console.log(serverMessage)
// 	}
// }
//
// class GlobalMessage extends Commands {
// 	public execute (serverMessage: NetworkingMessage) {
// 		get_context().pop_up.fire(
// 			'Monopoly Bank',
// 			serverMessage.message,
// 			'success',
// 			5000,
// 		)
// 	}
// }

/*_______________________________________________________________________________
	this space reserved for login
_______________________________________________________________________________ */


