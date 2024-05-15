import { createUniqueID } from './id'

export class PopUp {
	public fire(title: string, message: string, type: 'error' | 'success', duration: number = 10000) {
		this.createPopUp(title, message, type, duration)
	}

	private createPopUp (title: string, message: string, type: 'error' | 'success', duration: number) {
		const popupContainer = document.getElementById('popup__container')

		if (!popupContainer) {
			const container_element = `<div id="popup__container" class="popup__container"></div>`
			document.getElementsByTagName('body')[0]?.insertAdjacentHTML('beforeend', container_element)
		}

		const unique_id = createUniqueID();

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
