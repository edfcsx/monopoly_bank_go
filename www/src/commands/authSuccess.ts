import { Commands, NetworkingMessage } from '../types'
import { PopUp } from '../popUp'
import { getLoginButton } from '../main'

export class AuthSuccessCommand extends Commands {
	private popup: PopUp = new PopUp()

 	public execute (serverMessage: NetworkingMessage) {
		this.popup.fire(
 			'Monopoly Bank',
 			`Bem vindo, ${serverMessage.args?.username}!`,
 			'success',
 			3000,
 		)

 		if (serverMessage.args) {
 			sessionStorage.setItem('auth', JSON.stringify({
 				username: serverMessage.args!.username,
 				password: serverMessage.args!.password,
 				token: serverMessage.player_hash
 			}))

 			getLoginButton()?.removeAttribute('disabled')

 			if (window.location.pathname != '/bank') {
 				setTimeout(() => {
 					window.location.href = '/bank'
 				}, 1000)
 			}
 		} else {
 			this.popup.fire(
 				'Monopoly Bank',
 				'Ocorreu um erro no sistema, por favor contate o desenvolvedor',
 				'error',
 				5000,
 			)
 		}
 	}
 }
