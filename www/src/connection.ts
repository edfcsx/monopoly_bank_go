import { PopUp } from './popUp'
import { CommandsRequest, CommandsResponse, NetworkingMessage } from './types'
import { createUniqueID } from './id'
import { Commands } from './commands'
import {MJTP} from "./mjtp";

export function __Connection(): Connection {
	return Connection.getInstance()
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
		const playerHash = sessionStorage.getItem('player_hash')
		this.socket = new WebSocket(`ws://192.168.15.9:4444?player_hash=${playerHash}`)
		this.is_open = false;

		this.socket.onopen = () => {
			this.is_open = true;

			console.log('Conexão aberta')
			console.log('authenticando')

			const auth = new MJTP('/authenticate', { id: playerHash })
			this.socket?.send(auth.toString())

			setTimeout(() => {
				console.log('requesting profile')
				const mjtp = new MJTP('/status', {})
				this.socket?.send(mjtp.toString())
			}, 1000)
			//
			// setInterval(() => {
			// 	const mjtp = new MJTP('/profile', {})
			// 	this.socket?.send(mjtp.toString())
			// }, 500)
		}

		this.socket.onclose = () => {
			this.is_open = false;
		};

		this.socket.onerror = () => {
			this.is_open = false
			new PopUp().fire('Conexão', 'Sessão expirada!', 'error', 5000)

			setTimeout(() => {
				sessionStorage.clear()
				window.location.href = '/'
			}, 2000)
		};

		this.socket.onmessage = (e) => {
			const mjtp = MJTP.parse(e.data)
			console.log('MENSAGEM RECEBIDA:>', mjtp)
			// const [command, data] = String(e.data).split('|')
			// const msg: NetworkingMessage = data.length ? JSON.parse(data) : {} as NetworkingMessage
			//
			// msg.command = command as CommandsResponse | CommandsRequest
			//
			// if (msg.args_id && this.args_repository.has(msg.args_id)) {
			// 	const args = this.args_repository.get(msg.args_id)
			// 	this.args_repository.delete(msg.args_id)
			//
			// 	if (args) {
			// 		msg.args = args
			// 	}
			// }
			//
			// this.messages.push(msg)
			//
			// if (!this.messages_worker) {
			// 	this.createWorker()
			// }
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
		this.commands[CommandsResponse.AuthenticateSuccess] = new Commands.AuthSuccessCommand()
		this.commands[CommandsResponse.AuthenticateFailed] = new Commands.AuthFailedCommand()
		this.commands[CommandsResponse.GlobalMessage] = new Commands.GlobalMessageCommand()
		this.commands[CommandsResponse.ProfileData] = new Commands.ProfileCommand()
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
