import {
    createUser,
    deleteUser,
    getAllUsers,
    getUser,
    loginUser,
} from "../services/user.service.js";

export const createUserController = async (req, res) => {
    const user = await createUser(req.body);

    res.status(201).json({
        success: true,
        message: "User created successfully",
        data: user,
    });
};

export const loginUserController = async (req, res) => {
    const token = await loginUser(req.body);

    res.status(200).json({
        success: true,
        message: "Login successful",
        data: token,
    });
};

export const getUserByIdController = async (req, res) => {
    const user = await getUser(req.params.id);

    res.status(200).json({
        success: true,
        message: "User fetched successfully",
        data: user,
    });
};

export const getAllUserController = async (req, res) => {
    const users = await getAllUsers();

    res.status(200).json({
        success: true,
        message: "Users fetched successfully",
        data: users,
    });
};

export const deleteUserController = async (req, res) => {
    await deleteUser(req.params.id);

    res.status(204).send();
};