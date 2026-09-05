import bcrypt from "bcrypt";
import jwt from "jsonwebtoken";
import "dotenv/config";

import {
    deleteUserById,
    getAllUsersFromDB,
    getUserById,
    saveUser,
    userExistsByEmail,
} from "../repositories/user.repository.js";

import ApiError from "../utils/apiError.js";

export const createUser = async ({ name, email, password, role }) => {
    const existingUser = await userExistsByEmail(email);

    if (existingUser) {
        throw new ApiError(409, "Email already exists");
    }

    const passwordHash = await bcrypt.hash(password, 10);

    const user = await saveUser({
        name,
        email,
        passwordHash,
        role,
    });

    return user;
};

export const loginUser = async ({ email, password }) => {
    const user = await userExistsByEmail(email);

    if (!user) {
        throw new ApiError(401, "Invalid credentials");
    }

    const validPassword = await bcrypt.compare(
        password,
        user.passwordHash
    );

    if (!validPassword) {
        throw new ApiError(401, "Invalid credentials");
    }

    const payload = {
        id: user.id,
        email: user.email,
        role: user.role,
    };

    console.log(
    "USER SERVICE JWT SECRET:",
    JSON.stringify(process.env.JWT_SECRET)
    );

    const accessToken = jwt.sign(
        payload,
        process.env.JWT_SECRET,
        {
            expiresIn: process.env.JWT_SECRET_EXPIRE,
        }
    );

    const refreshToken = jwt.sign(
        payload,
        process.env.JWT_REFRESH,
        {
            expiresIn: process.env.JWT_REFRESH_EXPIRE,
        }
    );

    return {
        accessToken,
        refreshToken,
    };
};

export const getUser = async (id) => {
    const user = await getUserById(id);

    if (!user) {
        throw new ApiError(404, "User not found");
    }

    return user;
};

export const getAllUsers = async () => {
    return getAllUsersFromDB();
};

export const deleteUser = async (id) => {
    const user = await getUserById(id);

    if (!user) {
        throw new ApiError(404, "User not found");
    }

    await deleteUserById(id);
};