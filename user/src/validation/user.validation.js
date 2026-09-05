import Joi from "joi";

export const createUserValidation=Joi.object({
    name:Joi.string().min(3).required(),
    email:Joi.string().required(),
    password:Joi.string().min(8).required(),
    role:Joi.string().default('USER')
})

export const loginUserValidation=Joi.object({
    email:Joi.string().required(),
    password:Joi.string().min(8).required(),
})