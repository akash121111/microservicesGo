import express, { json } from 'express';
import "dotenv/config"
import userRouter from './src/routes/user.routes.js';
import { errorHandler } from './src/middleware/errorHandler.js';

const app=express();
const PORT=process.env.PORT;

app.use(express.json())

app.get("/health",(req,res)=>{
    res.status(200).json({"sucess":true,"message":"ok"})
})
app.use("/api/v1/user",userRouter)

app.use(errorHandler)

app.listen(PORT,()=>{
    console.log("user server started at Port: ",PORT)
})