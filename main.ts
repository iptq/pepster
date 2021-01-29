import {config} from "dotenv";
config();

import Bot from "./src/Bot";
import {MatrixBackend} from "./src/MatrixBackend";
import {DiscordBackend} from "./src/DiscordBackend";

let bots = [];

let homeserverUrl = process.env.MATRIX_HOMESERVER;
let accessToken = process.env.MATRIX_ACCESS_TOKEN;
if (homeserverUrl && accessToken) {
    let matrixBackend = new MatrixBackend(
        process.env.MATRIX_HOMESERVER,
        process.env.MATRIX_ACCESS_TOKEN,
    );
    let matrixBot = new Bot(matrixBackend);
    bots.push(matrixBot.start());
}

let botToken = process.env.DISCORD_BOT_TOKEN;
if (botToken) {
    let discordBackend = new DiscordBackend(process.env.DISCORD_BOT_TOKEN);
    let discordBot = new Bot(discordBackend);
    bots.push(discordBot.start());
}

async function main() {
    await Promise.all(bots);
}

main();
