/*_______________________________________________________________________________
    this space reserved for context and global classes
_______________________________________________________________________________ */
interface Context {
    connection: Connection
    pop_up: PopUp
}

type json = string;

function get_context(): Context {
    return (window as any).monopoly;
}

function create_unique_id(): string {
    return Math.random().toString(36).substr(2, 9);
}

class PopUp {
	public fire(title: string, message: string, type: 'error' | 'success', duration: number = 10000) {
		this.create_popup(title, message, type, duration);
	}

	private create_popup (title: string, message: string, type: 'error' | 'success', duration: number) {
		const popup_container = document.getElementById('popup__container')

		if (!popup_container) {
			const container_element = `<div id="popup__container" class="popup__container"></div>`
			document.getElementsByTagName('body')[0].insertAdjacentHTML('beforeend', container_element)
		}

		const unique_id = create_unique_id();

		const popup = `
			<div id="${unique_id}" class="popup popup-${type}">
				<div class="popup__header">
					<span class="popup__header__title">${title}</span>
					<span id="close-${unique_id}" class="popup__header__close">X</span>
				</div>
							
				<div class="popup__body">
					<span>${message}</span>
				</div>
			</div>
		`

		document.getElementById('popup__container')?.insertAdjacentHTML('beforeend', popup)

		document.getElementById(`close-${unique_id}`)?.addEventListener('click', () => {
		const popup = document.getElementById(unique_id)
				popup!.remove()
		})

		setTimeout(() => {
			const popup = document.getElementById(unique_id)
				popup?.remove()
		}, duration)
	}
}

enum CommandsResponse {
    AuthenticateFailed         = "AuthenticateFailed",
    AuthenticateSuccess        = "AuthenticateSuccess",
    Pong                       = "Pong",
    ProfileData                = "ProfileData",
    TransferSuccess            = "TransferSuccess",
    TransferFailed             = "TransferFailed",
    TransferInsufficientFunds  = "TransferInsufficientFunds",
    TransferReceived           = "TransferReceived",
    BadRequest                 = "BadRequest",
    GlobalMessage              = "GlobalMessage"
}

enum CommandsRequest {
    Authenticate               = "Authenticate",
    Ping                       = "Ping",
    SendProfile                = "SendProfile",
    Transfer                   = "Transfer"
}

interface NetworkingMessage {
    command: CommandsRequest | CommandsResponse,
    args?: { [key:string]: any }
    [key: string]: any
}

abstract class Commands {
	abstract execute (serverMessage: NetworkingMessage): void;
}

class AuthenticateSuccess extends Commands {
	public execute (serverMessage: NetworkingMessage) {
		console.log('MENSAGEM RECEBIDA DE AUTH:>', serverMessage)
		get_context().pop_up.fire(
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
			get_context().pop_up.fire(
				'Monopoly Bank',
				'Ocorreu um erro no sistema, por favor contate o desenvolvedor',
				'error',
				5000,
			)
		}
	}
}

class AuthenticateFailed extends Commands {
	public execute (serverMessage: NetworkingMessage) {
		get_context().pop_up.fire(
			'Monopoly Bank',
			'Usuário ou senha incorretos',
			'error',
			5000,
		)

		sessionStorage.removeItem('auth')
		getLoginButton()?.removeAttribute('disabled')
	}
}

class ProfileCommand extends Commands {
	public execute (serverMessage: NetworkingMessage) {
		console.log(serverMessage)
	}
}

class GlobalMessage extends Commands {
	public execute (serverMessage: NetworkingMessage) {
		get_context().pop_up.fire(
			'Monopoly Bank',
			serverMessage.message,
			'success',
			5000,
		)
	}
}

class Connection {
	private socket: WebSocket | null = null;
	private is_open: boolean = false;
	private messages: NetworkingMessage[] = []
	private messages_worker: any
	public commands: { [command: string]: Commands } = {}
	private args_repository: Map<string, { [key: string]: any }> = new Map()

	constructor() {
		this.createSocket()
		this.createCommands()
		this.createWorker()
	}

	private createSocket (): void {
		this.socket = new WebSocket("ws://192.168.15.8:4444");
		this.is_open = false;

		this.socket.onopen = () => {
			this.is_open = true;
		}

		this.socket.onclose = () => {
			this.is_open = false;
		};

		this.socket.onerror = () => {
			this.is_open = false
			get_context().pop_up.fire('Conexão', 'Não foi possível conectar ao servidor', 'error', 5000)
		};

		this.socket.onmessage = (e) => {
			const [command, data] = String(e.data).split('|')
			const msg: NetworkingMessage = data.length ? JSON.parse(data) : {} as NetworkingMessage

			msg.command = command as CommandsResponse | CommandsRequest

			if (msg.args_id && this.args_repository.has(msg.args_id)) {
				const args = this.args_repository.get(msg.args_id)
				this.args_repository.delete(msg.args_id)

				if (args) {
					msg.args = args
				}
			}

			this.messages.push(msg)

			if (!this.messages_worker) {
					this.createWorker()
			}
		}
	}

	public openSocket (): void {
		if (this.is_open) {
				this.socket?.close()
		}

		this.createSocket()
	}

	private createWorker (): void {
		this.messages_worker = setInterval(() => {
			while (this.messages.length) {
				const message = this.messages.shift()
				console.log('mensagem recebida', message)

				if (message) {
					const command = this.commands[`${message.command}`]

					if (command) {
						command.execute(message)
					}
				}
			}

			if (!this.messages.length) {
				clearInterval(this.messages_worker)
				this.messages_worker = null
			}
		}, 10)
	}

	private createCommands () {
		this.commands[CommandsResponse.AuthenticateSuccess] = new AuthenticateSuccess()
		this.commands[CommandsResponse.AuthenticateFailed] = new AuthenticateFailed()
		this.commands[CommandsRequest.SendProfile] = new ProfileCommand()
		this.commands[CommandsResponse.GlobalMessage] = new GlobalMessage()
	}

	public isOpen (): boolean {
		return this.is_open;
	}

	public send (data: NetworkingMessage): void {
		if (this.isOpen()) {
			const unique_id = create_unique_id()

			if (data.args) {
					this.args_repository.set(unique_id, data.args)
					delete data.args
			}

			const sendingData: {[key: string]: any} = { ...data, args_id: unique_id }

			if (sessionStorage.getItem('auth')) {
				const auth = JSON.parse(sessionStorage.getItem('auth') as string)
				sendingData.player_hash = auth.token
			}

			// @ts-ignore
			delete sendingData.command
			this.socket?.send(`${data.command}|${JSON.stringify(sendingData)}`)
		} else {
				console.error('Connection is not open')
		}
	}
}

(window as any).monopoly = {
	connection: new Connection(),
	pop_up: new PopUp()
} as Context;

// Check if the user is in bank page
if (window.location.pathname === '/bank') {
	// await connection to be open
	const interval = setInterval(() => {
		if (get_context().connection.isOpen()) {
			clearInterval(interval)
			get_context().connection.send({ command: CommandsRequest.SendProfile })
		}
	}, 100)
}

/*_______________________________________________________________________________
	this space reserved for login
_______________________________________________________________________________ */

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

	get_context().connection.send({
		command: CommandsRequest.Authenticate,
		args: { ...data },
		...data,
	})

	getLoginButton()?.setAttribute('disabled', 'true')
})
