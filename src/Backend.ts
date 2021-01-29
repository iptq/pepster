export type MessageHandlerFunc<B extends Backend> = (_: Message<B>) => void;

export interface Backend {
    /** The unique name for this backend */
    name: string;

    /** Registers the default message handler */
    registerMessageHandler(func: MessageHandlerFunc<this>): void;

    /** Starts the server */
    start(): Promise<void>;

    /** Send a reply */
    sendReply(reply: Reply<this>): Promise<void>;
}

export interface Message<B extends Backend> {
    /** Returns a unique string identifying this room. */
    getRoomId(): string;

    /** Returns a unique string identifying the author. */
    getAuthorId(): string;

    /** Returns the timestamp when this message was sent */
    sentAt(): Date;

    /** Checks if the message content starts with the given prefix. */
    startsWithPrefix(prefix: string): boolean;

    /** Get an array consisting of the command and the arguments */
    getCommandParts(prefix: string): string[];
}

export interface Reply<B extends Backend> {
}
