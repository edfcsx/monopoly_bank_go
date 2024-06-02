import { __Connection } from './connection'
import { CommandsRequest } from './types'
import { Worker } from './worker'
import { PopUp } from "./popUp";

// Check if the user is in bank page
if (window.location.pathname === '/bank') {
	if (!sessionStorage.getItem('player_hash')) {
		sessionStorage.clear()
		window.location.href = '/'
	}

	__Connection()
	// Worker.Make((w: any) => {
	// 	if (__Connection().isOpen()) {
	// 		Worker.Clear(w)
	// 		__Connection().send({ command: CommandsRequest.SendProfile })
	// 	}
	// })
}

export function getLoginButton(): HTMLElement | null {
	return document.getElementById('login-dispatch')
}

document.getElementById('login_form')?.addEventListener('submit', (e) => {
	e.preventDefault()
	const formData = new FormData(e.target as HTMLFormElement)

	// getLoginButton()?.setAttribute('disabled', 'true')

	fetch( 'http://192.168.15.10:7600/login',{
		body: JSON.stringify({ username: formData.get('username'), password: formData.get('password') }),
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		}
	}).then(async (response) => {
		if (response.status === 200) {
			const data = await response.json()
			sessionStorage.setItem('player_hash', data.player_hash)
			window.location.href = '/bank'
		} else {
			new PopUp().fire('Login', 'Usuário ou senha inválidos', 'error', 5000)
			getLoginButton()?.removeAttribute('disabled')
		}
	})
})
