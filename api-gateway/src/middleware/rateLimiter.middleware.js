import asyncHandler from "../utils/asyncHandler.js";
import { redisClient } from "../config/redis.js";

const rateLimiter = (route, limit, windowMs) => {
    return asyncHandler(async (req, res, next) => {
        const ip = req.ip;

        const key = `rateLimiter:${route}:${ip}`;

        const currentCount = await redisClient.incr(key);

        if (currentCount === 1) {
            await redisClient.expire(key, windowMs);
        }

        if (currentCount > limit) {
            return res.status(429).json({
                success: false,
                message: "Too many requests",
            });
        }

        next();
    });
};

export default rateLimiter;