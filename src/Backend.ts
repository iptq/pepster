export interface Backend<M extends Message> {
    /** The unique name for this backend */
    name: string;

    /** Registers the default message handler */
    registerMessageHandler(func: (_: M) => void): void;

    /** Starts the server */
    start(): Promise<void>;
}

export interface Message {
    /** Returns a unique string identifying this room. */
    getRoomId(): string;

    /** Returns the timestamp when this message was sent */
    sentAt(): Date;

    /** Checks if the message content starts with the given prefix. */
    startsWithPrefix(prefix: string): boolean;

    /** Get an array consisting of the command and the arguments */
    getCommandParts(prefix: string): string[];
}
