/*
 * ATTENTION: The "eval" devtool has been used (maybe by default in mode: "development").
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
/******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./src/connection.ts":
/*!***************************!*\
  !*** ./src/connection.ts ***!
  \***************************/
/***/ ((__unused_webpack_module, exports, __webpack_require__) => {

eval("\nObject.defineProperty(exports, \"__esModule\", ({ value: true }));\nexports.__Connection = void 0;\nconst popUp_1 = __webpack_require__(/*! ./popUp */ \"./src/popUp.ts\");\nconst id_1 = __webpack_require__(/*! ./id */ \"./src/id.ts\");\nfunction __Connection() {\n    return Connection.getInstance();\n}\nexports.__Connection = __Connection;\nclass Connection {\n    constructor() {\n        this.socket = null;\n        this.is_open = false;\n        this.messages = [];\n        this.commands = {};\n        this.args_repository = new Map();\n        this.createSocket();\n        this.createCommands();\n        this.createWorker();\n    }\n    static getInstance() {\n        if (!Connection._instance) {\n            Connection._instance = new Connection();\n        }\n        return Connection._instance;\n    }\n    createSocket() {\n        this.socket = new WebSocket(\"ws://192.168.15.10:4444\");\n        this.is_open = false;\n        this.socket.onopen = () => {\n            this.is_open = true;\n        };\n        this.socket.onclose = () => {\n            this.is_open = false;\n        };\n        this.socket.onerror = () => {\n            this.is_open = false;\n            new popUp_1.PopUp().fire('Conexão', 'Não foi possível conectar ao servidor', 'error', 5000);\n        };\n        this.socket.onmessage = (e) => {\n            const [command, data] = String(e.data).split('|');\n            const msg = data.length ? JSON.parse(data) : {};\n            msg.command = command;\n            if (msg.args_id && this.args_repository.has(msg.args_id)) {\n                const args = this.args_repository.get(msg.args_id);\n                this.args_repository.delete(msg.args_id);\n                if (args) {\n                    msg.args = args;\n                }\n            }\n            this.messages.push(msg);\n            if (!this.messages_worker) {\n                this.createWorker();\n            }\n        };\n    }\n    openSocket() {\n        var _a;\n        if (this.is_open) {\n            (_a = this.socket) === null || _a === void 0 ? void 0 : _a.close();\n        }\n        this.createSocket();\n    }\n    createWorker() {\n        this.messages_worker = setInterval(() => {\n            while (this.messages.length) {\n                const message = this.messages.shift();\n                console.log('mensagem recebida', message);\n                if (message) {\n                    const command = this.commands[`${message.command}`];\n                    if (command) {\n                        command.execute(message);\n                    }\n                }\n            }\n            if (!this.messages.length) {\n                clearInterval(this.messages_worker);\n                this.messages_worker = null;\n            }\n        }, 10);\n    }\n    createCommands() {\n        // this.commands[CommandsResponse.AuthenticateSuccess] = new AuthenticateSuccess()\n        // this.commands[CommandsResponse.AuthenticateFailed] = new AuthenticateFailed()\n        // this.commands[CommandsRequest.SendProfile] = new ProfileCommand()\n        // this.commands[CommandsResponse.GlobalMessage] = new GlobalMessage()\n    }\n    isOpen() {\n        return this.is_open;\n    }\n    send(data) {\n        var _a;\n        if (this.isOpen()) {\n            const unique_id = (0, id_1.createUniqueID)();\n            if (data.args) {\n                this.args_repository.set(unique_id, data.args);\n                delete data.args;\n            }\n            const sendingData = Object.assign(Object.assign({}, data), { args_id: unique_id });\n            if (sessionStorage.getItem('auth')) {\n                const auth = JSON.parse(sessionStorage.getItem('auth'));\n                sendingData.player_hash = auth.token;\n            }\n            // @ts-ignore\n            delete sendingData.command;\n            (_a = this.socket) === null || _a === void 0 ? void 0 : _a.send(`${data.command}|${JSON.stringify(sendingData)}`);\n        }\n        else {\n            console.error('Connection is not open');\n        }\n    }\n}\nConnection._instance = null;\n\n\n//# sourceURL=webpack://www/./src/connection.ts?");

/***/ }),

/***/ "./src/id.ts":
/*!*******************!*\
  !*** ./src/id.ts ***!
  \*******************/
/***/ ((__unused_webpack_module, exports) => {

eval("\nObject.defineProperty(exports, \"__esModule\", ({ value: true }));\nexports.createUniqueID = void 0;\nfunction createUniqueID() {\n    return Math.random().toString(36).substr(2, 9);\n}\nexports.createUniqueID = createUniqueID;\n\n\n//# sourceURL=webpack://www/./src/id.ts?");

/***/ }),

/***/ "./src/main.ts":
/*!*********************!*\
  !*** ./src/main.ts ***!
  \*********************/
/***/ ((__unused_webpack_module, exports, __webpack_require__) => {

eval("\nvar _a;\nObject.defineProperty(exports, \"__esModule\", ({ value: true }));\nconst connection_1 = __webpack_require__(/*! ./connection */ \"./src/connection.ts\");\nconst types_1 = __webpack_require__(/*! ./types */ \"./src/types.ts\");\nconst worker_1 = __webpack_require__(/*! ./worker */ \"./src/worker.ts\");\n(0, connection_1.__Connection)();\n// Check if the user is in bank page\nif (window.location.pathname === '/bank') {\n    worker_1.Worker.Make((w) => {\n        if ((0, connection_1.__Connection)().isOpen()) {\n            worker_1.Worker.Clear(w);\n            (0, connection_1.__Connection)().send({ command: types_1.CommandsRequest.SendProfile });\n        }\n    });\n}\nfunction getLoginButton() {\n    return document.getElementById('login-dispatch');\n}\n(_a = document.getElementById('login_form')) === null || _a === void 0 ? void 0 : _a.addEventListener('submit', (e) => {\n    var _a;\n    e.preventDefault();\n    const formData = new FormData(e.target);\n    const data = {\n        username: formData.get('username'),\n        password: formData.get('password')\n    };\n    (0, connection_1.__Connection)().send(Object.assign({ command: types_1.CommandsRequest.Authenticate, args: Object.assign({}, data) }, data));\n    (_a = getLoginButton()) === null || _a === void 0 ? void 0 : _a.setAttribute('disabled', 'true');\n});\n// abstract class Commands {\n// \tabstract execute (serverMessage: NetworkingMessage): void;\n// }\n//\n// class AuthenticateSuccess extends Commands {\n// \tpublic execute (serverMessage: NetworkingMessage) {\n// \t\tconsole.log('MENSAGEM RECEBIDA DE AUTH:>', serverMessage)\n// \t\tget_context().pop_up.fire(\n// \t\t\t'Monopoly Bank',\n// \t\t\t`Bem vindo, ${serverMessage.args?.username}!`,\n// \t\t\t'success',\n// \t\t\t3000,\n// \t\t)\n//\n// \t\tif (serverMessage.args) {\n// \t\t\tsessionStorage.setItem('auth', JSON.stringify({\n// \t\t\t\tusername: serverMessage.args!.username,\n// \t\t\t\tpassword: serverMessage.args!.password,\n// \t\t\t\ttoken: serverMessage.player_hash\n// \t\t\t}))\n//\n// \t\t\tgetLoginButton()?.removeAttribute('disabled')\n//\n// \t\t\tif (window.location.pathname != '/bank') {\n// \t\t\t\tsetTimeout(() => {\n// \t\t\t\t\twindow.location.href = '/bank'\n// \t\t\t\t}, 1000)\n// \t\t\t}\n// \t\t} else {\n// \t\t\tget_context().pop_up.fire(\n// \t\t\t\t'Monopoly Bank',\n// \t\t\t\t'Ocorreu um erro no sistema, por favor contate o desenvolvedor',\n// \t\t\t\t'error',\n// \t\t\t\t5000,\n// \t\t\t)\n// \t\t}\n// \t}\n// }\n//\n// class AuthenticateFailed extends Commands {\n// \tpublic execute (serverMessage: NetworkingMessage) {\n// \t\tget_context().pop_up.fire(\n// \t\t\t'Monopoly Bank',\n// \t\t\t'Usuário ou senha incorretos',\n// \t\t\t'error',\n// \t\t\t5000,\n// \t\t)\n//\n// \t\tsessionStorage.removeItem('auth')\n// \t\tgetLoginButton()?.removeAttribute('disabled')\n// \t}\n// }\n//\n// class ProfileCommand extends Commands {\n// \tpublic execute (serverMessage: NetworkingMessage) {\n// \t\tconsole.log(serverMessage)\n// \t}\n// }\n//\n// class GlobalMessage extends Commands {\n// \tpublic execute (serverMessage: NetworkingMessage) {\n// \t\tget_context().pop_up.fire(\n// \t\t\t'Monopoly Bank',\n// \t\t\tserverMessage.message,\n// \t\t\t'success',\n// \t\t\t5000,\n// \t\t)\n// \t}\n// }\n/*_______________________________________________________________________________\n    this space reserved for login\n_______________________________________________________________________________ */\n\n\n//# sourceURL=webpack://www/./src/main.ts?");

/***/ }),

/***/ "./src/popUp.ts":
/*!**********************!*\
  !*** ./src/popUp.ts ***!
  \**********************/
/***/ ((__unused_webpack_module, exports, __webpack_require__) => {

eval("\nObject.defineProperty(exports, \"__esModule\", ({ value: true }));\nexports.PopUp = void 0;\nconst id_1 = __webpack_require__(/*! ./id */ \"./src/id.ts\");\nclass PopUp {\n    fire(title, message, type, duration = 10000) {\n        this.createPopUp(title, message, type, duration);\n    }\n    createPopUp(title, message, type, duration) {\n        var _a, _b, _c;\n        const popupContainer = document.getElementById('popup__container');\n        if (!popupContainer) {\n            const container_element = `<div id=\"popup__container\" class=\"popup__container\"></div>`;\n            (_a = document.getElementsByTagName('body')[0]) === null || _a === void 0 ? void 0 : _a.insertAdjacentHTML('beforeend', container_element);\n        }\n        const unique_id = (0, id_1.createUniqueID)();\n        const popup = `\r\n\t\t\t<div id=\"${unique_id}\" class=\"popup popup-${type}\">\r\n\t\t\t\t<div class=\"popup__header\">\r\n\t\t\t\t\t<span class=\"popup__header__title\">${title}</span>\r\n\t\t\t\t\t<span id=\"close-${unique_id}\" class=\"popup__header__close\">X</span>\r\n\t\t\t\t</div>\r\n\t\t\t\t\t\t\t\r\n\t\t\t\t<div class=\"popup__body\">\r\n\t\t\t\t\t<span>${message}</span>\r\n\t\t\t\t</div>\r\n\t\t\t</div>\r\n\t\t`;\n        (_b = document.getElementById('popup__container')) === null || _b === void 0 ? void 0 : _b.insertAdjacentHTML('beforeend', popup);\n        (_c = document.getElementById(`close-${unique_id}`)) === null || _c === void 0 ? void 0 : _c.addEventListener('click', () => {\n            const popup = document.getElementById(unique_id);\n            popup.remove();\n        });\n        setTimeout(() => {\n            const popup = document.getElementById(unique_id);\n            popup === null || popup === void 0 ? void 0 : popup.remove();\n        }, duration);\n    }\n}\nexports.PopUp = PopUp;\n\n\n//# sourceURL=webpack://www/./src/popUp.ts?");

/***/ }),

/***/ "./src/types.ts":
/*!**********************!*\
  !*** ./src/types.ts ***!
  \**********************/
/***/ ((__unused_webpack_module, exports) => {

eval("\nObject.defineProperty(exports, \"__esModule\", ({ value: true }));\nexports.CommandsRequest = exports.CommandsResponse = void 0;\nvar CommandsResponse;\n(function (CommandsResponse) {\n    CommandsResponse[\"AuthenticateFailed\"] = \"AuthenticateFailed\";\n    CommandsResponse[\"AuthenticateSuccess\"] = \"AuthenticateSuccess\";\n    CommandsResponse[\"Pong\"] = \"Pong\";\n    CommandsResponse[\"ProfileData\"] = \"ProfileData\";\n    CommandsResponse[\"TransferSuccess\"] = \"TransferSuccess\";\n    CommandsResponse[\"TransferFailed\"] = \"TransferFailed\";\n    CommandsResponse[\"TransferInsufficientFunds\"] = \"TransferInsufficientFunds\";\n    CommandsResponse[\"TransferReceived\"] = \"TransferReceived\";\n    CommandsResponse[\"BadRequest\"] = \"BadRequest\";\n    CommandsResponse[\"GlobalMessage\"] = \"GlobalMessage\";\n})(CommandsResponse || (exports.CommandsResponse = CommandsResponse = {}));\nvar CommandsRequest;\n(function (CommandsRequest) {\n    CommandsRequest[\"Authenticate\"] = \"Authenticate\";\n    CommandsRequest[\"Ping\"] = \"Ping\";\n    CommandsRequest[\"SendProfile\"] = \"SendProfile\";\n    CommandsRequest[\"Transfer\"] = \"Transfer\";\n})(CommandsRequest || (exports.CommandsRequest = CommandsRequest = {}));\n\n\n//# sourceURL=webpack://www/./src/types.ts?");

/***/ }),

/***/ "./src/worker.ts":
/*!***********************!*\
  !*** ./src/worker.ts ***!
  \***********************/
/***/ ((__unused_webpack_module, exports) => {

eval("\nObject.defineProperty(exports, \"__esModule\", ({ value: true }));\nexports.Worker = void 0;\nvar Worker;\n(function (Worker) {\n    function Make(callback, errCallback = () => { }, timeout = 5000) {\n        let t = 0;\n        const worker = setInterval(() => {\n            if (t > timeout) {\n                clearInterval(worker);\n                errCallback();\n                return;\n            }\n            callback(worker);\n            t += 100;\n        }, 100);\n    }\n    Worker.Make = Make;\n    function MakeHighResolution(callback, errCallback = () => { }, timeout = 1000) {\n        let t = 0;\n        const worker = setInterval(() => {\n            if (t > timeout) {\n                clearInterval(worker);\n                errCallback();\n                return;\n            }\n            callback(worker);\n            t += 10;\n        }, 10);\n    }\n    Worker.MakeHighResolution = MakeHighResolution;\n    function Clear(worker) {\n        clearInterval(worker);\n    }\n    Worker.Clear = Clear;\n})(Worker || (exports.Worker = Worker = {}));\n\n\n//# sourceURL=webpack://www/./src/worker.ts?");

/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module can't be inlined because the eval devtool is used.
/******/ 	var __webpack_exports__ = __webpack_require__("./src/main.ts");
/******/ 	
/******/ })()
;