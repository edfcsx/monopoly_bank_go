import { Commands, NetworkingMessage } from '../types'

export class ProfileCommand extends Commands {
 	public execute (serverMessage: NetworkingMessage) {
 		console.log('profile command:>', serverMessage)
 	}
}
