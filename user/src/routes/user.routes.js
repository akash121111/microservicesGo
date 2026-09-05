import express from "express";

import asyncHandler from "../utils/asyncHandler.js";

import {
    createUserController,
    deleteUserController,
    getAllUserController,
    getUserByIdController,
    loginUserController,
} from "../controller/user.controller.js";

import { validate } from "../middleware/validate.js";

import {
    createUserValidation,
    loginUserValidation,
} from "../validation/user.validation.js";

const userRouter = express.Router();

userRouter.post(
    "/login",
    validate(loginUserValidation),
    asyncHandler(loginUserController)
);

userRouter.post(
    "/",
    validate(createUserValidation),
    asyncHandler(createUserController)
);

userRouter.get(
    "/",
    asyncHandler(getAllUserController)
);

userRouter.get(
    "/:id",
    asyncHandler(getUserByIdController)
);

userRouter.delete(
    "/:id",
    asyncHandler(deleteUserController)
);

export default userRouter;