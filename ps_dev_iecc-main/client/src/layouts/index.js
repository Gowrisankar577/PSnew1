import {
  BrowserRouter as Router,
  Route,
  Routes,
  useLocation,
} from "react-router-dom";
import AppSidebar from "./sidebar";
import AuthLogin from "../auth/login";
import { Suspense, useEffect, useState } from "react";
import PrivateRoute from "../auth/PrivateRoute";
import { apiGetRequest } from "../utils/api";
import { useAuth } from "../auth/AuthContext";
import routes from "./routes";
import { appBase, imagesBase, portalName } from "../utils/settings";
import BorderLinearProgress from "../components/progress";
import AuthLogout from "../auth/logout";

function AppContent() {
  const location = useLocation();
  const [userData, setUserData] = useState("loading");
  const { isAuthenticated } = useAuth();

  const [loading, setLoading] = useState(true);
  const [menus, setMenus] = useState([]);
  const [selectedMenu, setSelectedMenu] = useState(1);

  const currentUrl = location.pathname;
  const [sidebarState, setSideBarState] = useState(false);
  const handleSideBar = () => {
    setSideBarState(!sidebarState);
  };

  const fetchMenus = async () => {
    const response = await apiGetRequest("/resources");
    if (response.success) {
      setMenus(response.data.resources);
      setUserData({
        user_id: response.data.user_id,
        user_name: response.data.user_name,
        id: response.data.id,
      });
    }
    setLoading(false);
  };

  const updateSelectedMenu = () => {
    if (!menus) return;

    const basePath = "/" + location.pathname.split("/")[1];
    const menu = menus.find((menu) => basePath === menu.path);

    if (menu) {
      setSelectedMenu(menu.id);
    }
  };

  useEffect(() => {
    if (
      !isAuthenticated &&
      location.pathname !== appBase + "/auth/login" &&
      location.pathname !== appBase + "/tiny" &&
      location.pathname !== appBase + "/parent/leave-approval"
    ) {
      window.location = appBase + "/auth/login";
      return;
    }

    updateSelectedMenu();
  }, [location.pathname, menus]);

  useEffect(() => {
    if (
      ![
        appBase + "/auth/login",
        appBase + "/auth/logout",
        appBase + "/error-404",
        appBase + "/tiny",
        appBase + "/parent/leave-approval",
      ].includes(location.pathname)
    ) {
      fetchMenus();
    } else {
      setLoading(false);
    }
  }, []);

  if (loading) {
    return (
      <div style={{ width: "100%", marginTop: -7 }}>
        <BorderLinearProgress />
      </div>
    );
  }

  const showSidebar =
    ![
      appBase + "/error-404",
      appBase + "/auth/logout",
      appBase + "/auth/login",
    ].includes(currentUrl) 

  return (
    <div className="w-screen h-screen flex">
      {showSidebar && (
        <AppSidebar
          open={sidebarState}
          handleSideBar={handleSideBar}
          menu={menus || []}
          selectedmenu={selectedMenu}
          onSelectMenu={setSelectedMenu}
        />
      )}

      <div className="flex-1 flex flex-col h-screen w-full overflow-auto">
        {showSidebar && (
          <div className="w-full bg-white p-3  h-18 flex items-center justify-between drop-shadow-sm	z-20">
            <div className="flex gap-5 items-center">
              <i
                onClick={handleSideBar}
                className={`bx text-3xl ${
                  sidebarState ? "bx-x" : "bx-menu"
                } cursor-pointer block sm:hidden`}
              ></i>
              <h3 className="font-medium text-lg">{portalName}</h3>
            </div>
            <div className="flex gap-7 items-center bg-background py-1 px-5 rounded-md cursor-pointer">
              <img
                className="w-[35px] h-[35px] rounded-full object-cover"
                src={`${imagesBase}/user/images/${userData.user_id}.jpg`}
              />
              <div className="sm:flex hidden flex-col">
                <h2 className="text-[13px] font-medium">{userData.user_id}</h2>
                <h2 className="text-[16px] font-semibold">
                  {userData.user_name}
                </h2>
              </div>
            </div>
          </div>
        )}

        <div className="flex-1 w-full overflow-auto">
          <Suspense
            fallback={
              <div style={{ width: "100%", marginTop: -7 }}>
                <BorderLinearProgress />
              </div>
            }
          >
            <Routes>
              <Route path={appBase + "/auth/login"} element={<AuthLogin />} />
              <Route path={appBase + "/auth/logout"} element={<AuthLogout />} />

              {menus.map((menu) => {
                const Component = routes[menu.element];

                if (!Component) {
                  return null;
                }

                return (
                  <Route
                    key={menu.path}
                    path={menu.path}
                    element={
                      <PrivateRoute>
                        <Component />
                      </PrivateRoute>
                    }
                  />
                );
              })}
              <Route
                path="*"
                element={
                  <h1 className="p-3 text-xl text-center">
                    Request URL 404 | Not Found
                  </h1>
                }
              />
            </Routes>
          </Suspense>
        </div>
      </div>
    </div>
  );
}

function App() {
  return (
    <Router>
      <AppContent />
    </Router>
  );
}

export default App;
