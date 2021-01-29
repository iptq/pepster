import {Command, Reply} from ".";
import {Backend, Message} from "../Backend";
import {DiscordBackend, DiscordMessage} from "../DiscordBackend";
import Bot from "../Bot";

export class HelpCommand<B extends Backend, M extends Message<B>> implements Command<B, M> {
    name: "help";
    aliases: [];

    constructor(private bot: Bot<B, M>) {
    }

    handleMessage(message: M, args: string[]) {
        for (let command of commands) {
        }
    }
}

class HelpBaseReply {
}

function discordOutput<T extends new (...args: any[]) => {}>(Base: T) {
    return class DiscordExtended extends Base implements Reply<DiscordBackend, DiscordMessage, Command<DiscordBackend, DiscordMessage>> {
        constructor(...args: any[]) {
            super(args);
        }

        toReply() {
        }
    }
}

export const HelpReply = discordOutput(HelpBaseReply);
