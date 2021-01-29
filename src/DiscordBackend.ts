import {Client, Channel, Message as DiscordJsMessage, APIMessage} from "discord.js";
import {parse as parseShell} from "shell-quote";

import {MessageHandlerFunc, Backend, Message, Reply} from "./Backend";

export class DiscordBackend implements Backend {
    name = "discord";

    private client: Client;

    constructor(private botToken: string) {
        this.client = new Client();
    }

    registerMessageHandler(func: MessageHandlerFunc<this>) {
        this.client.on("message", event => {
            let message = new DiscordMessage(event);
            func(message);
        });
    }

    async start() {
        await this.client.login(this.botToken);
    }

    async sendReply(reply: Reply<this>) {
    }
}

export class DiscordMessage implements Message<DiscordBackend> {
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

export class DiscordReply implements Reply<DiscordBackend> {
    constructor(target: Channel, message: APIMessage) {
    }
}
