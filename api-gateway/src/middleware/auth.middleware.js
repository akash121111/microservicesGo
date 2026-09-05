import ApiError from "../utils/apiError.js";
import jwt from "jsonwebtoken";
import "dotenv/config";

const auth = (req, res, next) => {
    const authorization = req.get("Authorization");

    if (!authorization) {
        return next(
            new ApiError(401, "Authorization header missing")
        );
    }

    const [type, token] = authorization.split(" ");


    if (type !== "Bearer" || !token) {
        return next(
            new ApiError(401, "Invalid authorization header")
        );
    }

    try {

        const decoded = jwt.verify(
            token,
            process.env.JWT_SECRET
        );
        req.user = decoded;
        

        next();
    } catch (error) {

        if (error.name === "TokenExpiredError") {
            return next(
                new ApiError(401, "Access token expired")
            );
        }

        if (error.name === "JsonWebTokenError") {
            return next(
                new ApiError(401, "Invalid access token")
            );
        }

        return next(error);
    }
};

export default auth;