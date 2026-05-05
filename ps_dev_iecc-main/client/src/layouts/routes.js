import { lazy } from "react";

const routes = {
  AppDashboard: lazy(() => import("../pages/dashboard")),
};

export default routes;
