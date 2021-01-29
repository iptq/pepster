import {Backend, Message} from "./Backend";
import {Command} from "./commands";

import HelpCommand from "./commands/help";

export default class Bot<B extends Backend<M>, M extends Message> {
    private prefixMap = new Map<string, string>();
    private startTime: Date;
    private commands: Command[] = [];

    constructor(private backend: B) {
        this.startTime = new Date();
        this.commands.push(new HelpCommand());

        this.backend.registerMessageHandler(message => {
            // ignore events that are sent before bot starts up
            let timestamp = message.sentAt();
            if (timestamp < this.startTime) return;

            let roomId = message.getRoomId();
            let prefix = this.getRoomPrefix(roomId);
            console.log("message", message);

            if (message.startsWithPrefix(prefix)) {
                let commandParts = message.getCommandParts(prefix);
                console.log("got a new command", commandParts);
            }
        });
    }

    getRoomPrefix(roomId: string): string {
        if (!this.prefixMap.has(roomId)) {
            // fetch this prefix using the database
            // TODO
            this.prefixMap.set(roomId, "!");
        }

        return this.prefixMap.get(roomId);
    }

    async start() {
        console.log("starting", this.backend.name);
        await this.backend.start();
    }
}
