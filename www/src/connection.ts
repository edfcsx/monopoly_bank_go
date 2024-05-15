import { PopUp } from './popUp'
import { CommandsRequest, CommandsResponse } from './types'
import { createUniqueID } from './id'

export function __Connection(): Connection {
	return Connection.getInstance()
}

interface NetworkingMessage {
	command: CommandsRequest | CommandsResponse,
	args?: { [key:string]: any }
	[key: string]: any
}

class Connection {
	private socket: WebSocket | null = null;
	private is_open: boolean = false;
	private messages: any[] = []
	private messages_worker: any
	public commands: { [command: string]: any } = {}
	private args_repository: Map<string, { [key: string]: any }> = new Map()

	private static _instance: Connection | null = null

	constructor() {
		this.createSocket()
		this.createCommands()
		this.createWorker()
	}

	public static getInstance (): Connection {
		if (!Connection._instance) {
			Connection._instance = new Connection()
		}

		return Connection._instance
	}

	private createSocket (): void {
		this.socket = new WebSocket("ws://192.168.15.10:4444");
		this.is_open = false;

		this.socket.onopen = () => {
			this.is_open = true;
		}

		this.socket.onclose = () => {
			this.is_open = false;
		};

		this.socket.onerror = () => {
			this.is_open = false
			new PopUp().fire('Conexão', 'Não foi possível conectar ao servidor', 'error', 5000)
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
		// this.commands[CommandsResponse.AuthenticateSuccess] = new AuthenticateSuccess()
		// this.commands[CommandsResponse.AuthenticateFailed] = new AuthenticateFailed()
		// this.commands[CommandsRequest.SendProfile] = new ProfileCommand()
		// this.commands[CommandsResponse.GlobalMessage] = new GlobalMessage()
	}

	public isOpen (): boolean {
		return this.is_open;
	}

	public send (data: NetworkingMessage): void {
		if (this.isOpen()) {
			const unique_id = createUniqueID()

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
