import { createClient } from "redis";

export const redisClient = createClient({
    socket: {
        host: process.env.REDIS_HOST,
        port: Number(process.env.REDIS_PORT),
    },
});

redisClient.on("error", (error) => {
    console.error("Redis client error:", error);
});

redisClient.on("connect", () => {
    console.log("Redis connected");
});

redisClient.on("ready", () => {
    console.log("Redis ready");
});

redisClient.on("reconnecting", () => {
    console.log("Redis reconnecting");
});

redisClient.on("end", () => {
    console.log("Redis connection closed");
});

export const connectRedis = async () => {
    if (!redisClient.isOpen) {
        await redisClient.connect();
    }
};