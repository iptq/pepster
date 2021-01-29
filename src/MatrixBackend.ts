import {MatrixClient, AutojoinRoomsMixin} from "matrix-bot-sdk";
import {parse as parseShell} from "shell-quote";

import {Backend, Message} from "./Backend";

export class MatrixBackend implements Backend {
    name = "matrix";

    private client: MatrixClient;

    constructor(homeserverUrl: string, accessToken: string) {
        this.client = new MatrixClient(homeserverUrl, accessToken);
        AutojoinRoomsMixin.setupOnClient(this.client);
    }

    registerMessageHandler(func) {
        this.client.on("room.message", (roomId, event) => {
            let message = new MatrixMessage(roomId, event);
            func(message);
        });
    }

    async start() {
        await this.client.start();
    }
}

export class MatrixMessage implements Message<MatrixBackend> {
    constructor(private roomId: string, private event: any) { }

    getRoomId() {
        return this.roomId;
    }

    sentAt(): Date {
        return new Date(this.event.origin_server_ts);
    }

    startsWithPrefix(prefix: string): boolean {
        let msgtype = this.event.content.msgtype;
        if (msgtype === "m.text") {
            let body: string = this.event.content.body;
            return body.startsWith(prefix);
        }
        return false;
    }

    getCommandParts(prefix: string): string[] {
        let msgtype = this.event.content.msgtype;
        if (msgtype === "m.text") {
            let body: string = this.event.content.body;
            if (body.startsWith(prefix)) {
                body = body.substr(prefix.length);
            }

            let command = parseShell(body).map(entry => entry.toString());
            return command;
        }
        return [];
    }
}
