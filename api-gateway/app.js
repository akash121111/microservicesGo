import express, { json } from 'express';
import { errorHandler } from './src/middleware/errorHandler.middleware.js';
import "dotenv/config"
import proxyRouter from './src/routes/proxy.routes.js';
import { connectRedis } from './src/config/redis.js';
import correlationId from './src/middleware/correlationId.js';



const app=express();
const PORT=process.env.PORT;
app.get("/health",(req,res)=>{
    res.status(200).json({"sucess":true,"message":"ok"})
})
app.use(correlationId)
app.use((req, res, next) => {
    console.log("🔥 APP REQUEST:", req.method, req.originalUrl,req.headers);
    next();
});
app.use("/api/v1",proxyRouter)

app.use(errorHandler)

app.listen(PORT,async()=>{
    await connectRedis()
    console.log("api gateway started at Port: ",PORT)
})