import { Commands, NetworkingMessage } from '../types'
import {PopUp} from "../popUp";

export class GlobalMessageCommand extends Commands {
	private popup: PopUp = new PopUp()

 	public execute (serverMessage: NetworkingMessage) {
 		this.popup.fire(
 			'Monopoly Bank',
 			serverMessage.message,
 			'success',
 			5000,
 		)
 	}
 }