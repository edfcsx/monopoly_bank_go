import { NetworkingMessage, Commands } from '../types'
import { PopUp } from '../popUp'
import { getLoginButton } from '../main'

export class AuthFailedCommand extends Commands {
	private popup: PopUp = new PopUp()

	public execute (serverMessage: NetworkingMessage) {
		this.popup.fire(
			'Monopoly Bank',
			'Usuário ou senha incorretos',
			'error',
			5000,
		)

		sessionStorage.removeItem('auth')
		getLoginButton()?.removeAttribute('disabled')
	}
}
