import express from "express";

import { userProxy } from "../proxy/user.proxy.js";
import { orderProxy } from "../proxy/order.proxy.js";
import { productProxy } from "../proxy/product.proxy.js";

import rateLimiter from "../middleware/rateLimiter.middleware.js";
import auth from "../middleware/auth.middleware.js";
import authorize from "../middleware/authorize.middleware.js";

const proxyRouter = express.Router();


// =====================================================
// USER ROUTES
// =====================================================

// ---------- Public User Routes ----------

const publicUserRoutes = express.Router();

publicUserRoutes.post(
    "/login",
    rateLimiter("/user/login", 5, 60),
    userProxy
);

publicUserRoutes.post(
    "/",
    rateLimiter("/user", 5, 60),
    userProxy
);


// ---------- Protected User Routes ----------

const protectedUserRoutes = express.Router();

protectedUserRoutes.use(auth);

protectedUserRoutes.get(
    "/",
    rateLimiter("/user", 30, 60),
    userProxy
);

protectedUserRoutes.get(
    "/:id",
    rateLimiter("/user/:id", 30, 60),
    userProxy
);


// ---------- Admin User Routes ----------

const adminUserRoutes = express.Router();

adminUserRoutes.use(auth);
adminUserRoutes.use(authorize("ADMIN"));

adminUserRoutes.delete(
    "/:id",
    rateLimiter("/user/:id", 10, 60),
    userProxy
);


// Mount User Routes

proxyRouter.use("/user", publicUserRoutes);
proxyRouter.use("/user", protectedUserRoutes);
proxyRouter.use("/user", adminUserRoutes);


// =====================================================
// ORDER ROUTES
// =====================================================

// All order routes require authentication

const protectedOrderRoutes = express.Router();

protectedOrderRoutes.use(auth);

protectedOrderRoutes.post(
    "/",
    rateLimiter("/order", 30, 60),
    orderProxy
);

protectedOrderRoutes.post(
    "/:id",
    rateLimiter("/order", 30, 60),
    orderProxy
);
protectedOrderRoutes.patch(
    "/my",
    rateLimiter("/order", 30, 60),
    productProxy
);

protectedOrderRoutes.post(
    "/cancel/:orderId",
    rateLimiter("/order", 30, 60),
    orderProxy
);
protectedOrderRoutes.use(orderProxy);


proxyRouter.use(
    "/order",
    protectedOrderRoutes
);


// =====================================================
// PRODUCT ROUTES
// =====================================================

// ---------- Public Product Routes ----------

const publicProductRoutes = express.Router();

publicProductRoutes.get(
    "/",
    rateLimiter("/product", 90, 60),
    productProxy
);

publicProductRoutes.get(
    "/:id",
    rateLimiter("/product/:id", 90, 60),
    productProxy
);

// ---------- protected Product Routes for release stock or book stock ----------
const protectedProductRoutes=express.Router();
protectedProductRoutes.use(auth)

protectedProductRoutes.patch(
    "/:id/stocks",
    rateLimiter("/product/:id", 30, 60),
    productProxy
);

protectedProductRoutes.patch(
    "/:id/relstocks",
    rateLimiter("/product/:id", 30, 60),
    productProxy
);


// ---------- Admin Product Routes ----------


const adminProductRoutes = express.Router();

adminProductRoutes.use(auth);
adminProductRoutes.use(authorize("ADMIN"));

adminProductRoutes.post(
    "/",
    rateLimiter("/product", 30, 60),
    productProxy
);

adminProductRoutes.put(
    "/:id",
    rateLimiter("/product/:id", 30, 60),
    productProxy
);

adminProductRoutes.patch(
    "/:id",
    rateLimiter("/product/:id", 30, 60),
    productProxy
);

adminProductRoutes.delete(
    "/:id",
    rateLimiter("/product/:id", 30, 60),
    productProxy
);


// Mount Product Routes

proxyRouter.use(
    "/product",
    publicProductRoutes
);

proxyRouter.use(
    "/product",
    protectedProductRoutes
);

proxyRouter.use(
    "/product",
    adminProductRoutes
);



export default proxyRouter;