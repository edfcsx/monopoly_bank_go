"use strict";
var _a;
function get_context() {
    return window.monopoly;
}
function create_unique_id() {
    return Math.random().toString(36).substr(2, 9);
}
class PopUp {
    fire(title, message, type, duration = 10000) {
        this.create_popup(title, message, type, duration);
    }
    create_popup(title, message, type, duration) {
        var _a, _b;
        const popup_container = document.getElementById('popup__container');
        if (!popup_container) {
            const container_element = `<div id="popup__container" class="popup__container"></div>`;
            document.getElementsByTagName('body')[0].insertAdjacentHTML('beforeend', container_element);
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
		`;
        (_a = document.getElementById('popup__container')) === null || _a === void 0 ? void 0 : _a.insertAdjacentHTML('beforeend', popup);
        (_b = document.getElementById(`close-${unique_id}`)) === null || _b === void 0 ? void 0 : _b.addEventListener('click', () => {
            const popup = document.getElementById(unique_id);
            popup.remove();
        });
        setTimeout(() => {
            const popup = document.getElementById(unique_id);
            popup === null || popup === void 0 ? void 0 : popup.remove();
        }, duration);
    }
}
var CommandsResponse;
(function (CommandsResponse) {
    CommandsResponse["AuthenticateFailed"] = "AuthenticateFailed";
    CommandsResponse["AuthenticateSuccess"] = "AuthenticateSuccess";
    CommandsResponse["Pong"] = "Pong";
    CommandsResponse["ProfileData"] = "ProfileData";
    CommandsResponse["TransferSuccess"] = "TransferSuccess";
    CommandsResponse["TransferFailed"] = "TransferFailed";
    CommandsResponse["TransferInsufficientFunds"] = "TransferInsufficientFunds";
    CommandsResponse["TransferReceived"] = "TransferReceived";
    CommandsResponse["BadRequest"] = "BadRequest";
    CommandsResponse["GlobalMessage"] = "GlobalMessage";
})(CommandsResponse || (CommandsResponse = {}));
var CommandsRequest;
(function (CommandsRequest) {
    CommandsRequest["Authenticate"] = "Authenticate";
    CommandsRequest["Ping"] = "Ping";
    CommandsRequest["SendProfile"] = "SendProfile";
    CommandsRequest["Transfer"] = "Transfer";
})(CommandsRequest || (CommandsRequest = {}));
class Commands {
}
class AuthenticateSuccess extends Commands {
    execute(serverMessage) {
        var _a, _b;
        console.log('MENSAGEM RECEBIDA DE AUTH:>', serverMessage);
        get_context().pop_up.fire('Monopoly Bank', `Bem vindo, ${(_a = serverMessage.args) === null || _a === void 0 ? void 0 : _a.username}!`, 'success', 3000);
        if (serverMessage.args) {
            sessionStorage.setItem('auth', JSON.stringify({
                username: serverMessage.args.username,
                password: serverMessage.args.password,
                token: serverMessage.player_hash
            }));
            (_b = getLoginButton()) === null || _b === void 0 ? void 0 : _b.removeAttribute('disabled');
            if (window.location.pathname != '/bank') {
                setTimeout(() => {
                    window.location.href = '/bank';
                }, 1000);
            }
        }
        else {
            get_context().pop_up.fire('Monopoly Bank', 'Ocorreu um erro no sistema, por favor contate o desenvolvedor', 'error', 5000);
        }
    }
}
class AuthenticateFailed extends Commands {
    execute(serverMessage) {
        var _a;
        get_context().pop_up.fire('Monopoly Bank', 'Usuário ou senha incorretos', 'error', 5000);
        sessionStorage.removeItem('auth');
        (_a = getLoginButton()) === null || _a === void 0 ? void 0 : _a.removeAttribute('disabled');
    }
}
class ProfileCommand extends Commands {
    execute(serverMessage) {
        console.log(serverMessage);
    }
}
class GlobalMessage extends Commands {
    execute(serverMessage) {
        get_context().pop_up.fire('Monopoly Bank', serverMessage.message, 'success', 5000);
    }
}
class Connection {
    constructor() {
        this.socket = null;
        this.is_open = false;
        this.messages = [];
        this.commands = {};
        this.args_repository = new Map();
        this.createSocket();
        this.createCommands();
        this.createWorker();
    }
    createSocket() {
        this.socket = new WebSocket("ws://192.168.15.8:4444");
        this.is_open = false;
        this.socket.onopen = () => {
            this.is_open = true;
        };
        this.socket.onclose = () => {
            this.is_open = false;
        };
        this.socket.onerror = () => {
            this.is_open = false;
            get_context().pop_up.fire('Conexão', 'Não foi possível conectar ao servidor', 'error', 5000);
        };
        this.socket.onmessage = (e) => {
            const [command, data] = String(e.data).split('|');
            const msg = data.length ? JSON.parse(data) : {};
            msg.command = command;
            if (msg.args_id && this.args_repository.has(msg.args_id)) {
                const args = this.args_repository.get(msg.args_id);
                this.args_repository.delete(msg.args_id);
                if (args) {
                    msg.args = args;
                }
            }
            this.messages.push(msg);
            if (!this.messages_worker) {
                this.createWorker();
            }
        };
    }
    openSocket() {
        var _a;
        if (this.is_open) {
            (_a = this.socket) === null || _a === void 0 ? void 0 : _a.close();
        }
        this.createSocket();
    }
    createWorker() {
        this.messages_worker = setInterval(() => {
            while (this.messages.length) {
                const message = this.messages.shift();
                console.log('mensagem recebida', message);
                if (message) {
                    const command = this.commands[`${message.command}`];
                    if (command) {
                        command.execute(message);
                    }
                }
            }
            if (!this.messages.length) {
                clearInterval(this.messages_worker);
                this.messages_worker = null;
            }
        }, 10);
    }
    createCommands() {
        this.commands[CommandsResponse.AuthenticateSuccess] = new AuthenticateSuccess();
        this.commands[CommandsResponse.AuthenticateFailed] = new AuthenticateFailed();
        this.commands[CommandsRequest.SendProfile] = new ProfileCommand();
        this.commands[CommandsResponse.GlobalMessage] = new GlobalMessage();
    }
    isOpen() {
        return this.is_open;
    }
    send(data) {
        var _a;
        if (this.isOpen()) {
            const unique_id = create_unique_id();
            if (data.args) {
                this.args_repository.set(unique_id, data.args);
                delete data.args;
            }
            const sendingData = Object.assign(Object.assign({}, data), { args_id: unique_id });
            if (sessionStorage.getItem('auth')) {
                const auth = JSON.parse(sessionStorage.getItem('auth'));
                sendingData.player_hash = auth.token;
            }
            // @ts-ignore
            delete sendingData.command;
            (_a = this.socket) === null || _a === void 0 ? void 0 : _a.send(`${data.command}|${JSON.stringify(sendingData)}`);
        }
        else {
            console.error('Connection is not open');
        }
    }
}
window.monopoly = {
    connection: new Connection(),
    pop_up: new PopUp()
};
// Check if the user is in bank page
if (window.location.pathname === '/bank') {
    // await connection to be open
    const interval = setInterval(() => {
        if (get_context().connection.isOpen()) {
            clearInterval(interval);
            get_context().connection.send({ command: CommandsRequest.SendProfile });
        }
    }, 100);
}
/*_______________________________________________________________________________
    this space reserved for login
_______________________________________________________________________________ */
function getLoginButton() {
    return document.getElementById('login-dispatch');
}
(_a = document.getElementById('login_form')) === null || _a === void 0 ? void 0 : _a.addEventListener('submit', (e) => {
    var _a;
    e.preventDefault();
    const formData = new FormData(e.target);
    const data = {
        username: formData.get('username'),
        password: formData.get('password')
    };
    get_context().connection.send(Object.assign({ command: CommandsRequest.Authenticate, args: Object.assign({}, data) }, data));
    (_a = getLoginButton()) === null || _a === void 0 ? void 0 : _a.setAttribute('disabled', 'true');
});
