import { lazy } from "react";

const routes = {
  AppDashboard: lazy(() => import("../pages/dashboard")),
  AppCommunityDashboard: lazy(() => import("../pages/community")),
};

export default routes;
