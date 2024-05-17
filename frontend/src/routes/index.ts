import { lazy } from "solid-js";

const route = [
    {
        path: "/",
        component: lazy(() => import("../App"))
    },
    {
        path: "/home",
        component: lazy(() => import("../pages/Home"))
    }
];

export default route;
