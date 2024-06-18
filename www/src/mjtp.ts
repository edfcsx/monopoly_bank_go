export class MJTP {
	private resource: string;
	private version: string;
	private body: { [key: string]: any };

	constructor(resource: string, body: { [key: string]: any }) {
		this.resource = resource;
		this.version = "MJTP/1.0";
		this.body = body;
	}

	toString(): string {
		return `${this.resource} ${this.version} ${JSON.stringify(this.body)}\r\n\r\n`;
	}

	static parse(data: string): MJTP {
		// validate
		const parts = data.split(" ");

		if (parts.length < 3) {
			throw new Error("invalid message format");
		} else if (parts[1] !== "MJTP/1.0") {
			throw new Error("invalid version");
		} else if (parts[0].length === 0) {
			throw new Error("invalid resource");
		}

		let buffer: string = data
		let bufferPosition: number

		const mjtp = new MJTP('', {})

		bufferPosition = buffer.indexOf(' ')
		mjtp.resource = buffer.slice(0, bufferPosition)

		buffer = buffer.slice(bufferPosition + 1)
		bufferPosition = buffer.indexOf(' ')
		mjtp.version = buffer.slice(0, bufferPosition)

		buffer = buffer.slice(bufferPosition + 1)
		mjtp.body = JSON.parse(buffer.split('\r\n\r\n')[0])
		
		return mjtp
	}
}
