import {Backend, Message} from "../Backend";

export interface Command<B extends Backend, M extends Message<B>> {
    /** Name of the command */
    name: string;

    /** Aliases of the command */
    aliases: string[];

    /** Handle message */
    handleMessage(message: M, args: string[]): Reply<B, M, this>;
}

export interface Reply<B extends Backend, M extends Message<B>, C extends Command<B, M>> {
    toReply();
}
