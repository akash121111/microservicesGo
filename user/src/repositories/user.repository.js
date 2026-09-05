import prisma from "../lib/prisma.js";

const safeUserFields = {
    id: true,
    name: true,
    email: true,
    role: true,
    createdAt: true,
    updatedAt: true,
};

export const saveUser = async (user) => {
    return prisma.user.create({
        data: user,
        select: safeUserFields,
    });
};

export const getUserById = async (id) => {
    
    return prisma.user.findUnique({
        where: {
            id: id,
        },
        select: safeUserFields,
    });
};

export const getAllUsersFromDB = async () => {
    return prisma.user.findMany({
        select: safeUserFields,
    });
};

export const deleteUserById = async (id) => {
  
    return prisma.user.delete({
        where: {
            id: id,
        },
    });
};

export const userExistsByEmail = async (email) => {
    return prisma.user.findUnique({
        where: {
            email,
        },
    });
};