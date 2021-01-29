import {Client, Message as DiscordJsMessage} from "discord.js";
import {parse as parseShell} from "shell-quote";

import {Backend, Message} from "./Backend";

export class DiscordBackend implements Backend<DiscordMessage> {
    name = "discord";

    private client: Client;

    constructor(private botToken: string) {
        this.client = new Client();
    }

    registerMessageHandler(func) {
        this.client.on("message", event => {
            let message = new DiscordMessage(event);
            func(message);
        });
    }

    async start() {
        await this.client.login(this.botToken);
    }
}

export class DiscordMessage implements Message {
    constructor(private message: DiscordJsMessage) {
    }

    getRoomId(): string {
        return this.message.channel.id;
    }

    sentAt(): Date {
        return this.message.createdAt;
    }

    startsWithPrefix(prefix: string): boolean {
        return this.message.content.startsWith(prefix);
    }

    getCommandParts(prefix: string): string[] {
        let body = this.message.content;
        if (body.startsWith(prefix)) {
            body = body.substr(prefix.length);
        }

        return parseShell(body).map(entry => entry.toString());
    }
}
